package slot

import (
	"app/commons/guard"
	model "app/db/models"
	"app/pkg/mail"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSlotTestContext(t *testing.T, method string, body []byte, slotId string, user *guard.Claims) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "slotId", Value: slotId}}
	if user != nil {
		c.Set("user", user)
	}
	return c, recorder
}

func TestNewSlotController_ReusesProvidedInstance(t *testing.T) {
	existing := &SlotController{}
	assert.Same(t, existing, NewSlotController(existing))
}

func TestNewSlotController_Nil_BuildsDefault(t *testing.T) {
	ctl := NewSlotController(nil)
	assert.NotNil(t, ctl.slotService)
}

func TestSlotController_ConfirmSlot_InvalidClaimsType(t *testing.T) {
	ctl := &SlotController{slotService: newTestSlotService(t)}
	body, _ := json.Marshal(ConfirmSlotDto{StartsAt: time.Now(), EndsAt: time.Now().Add(time.Hour)})
	c, recorder := newSlotTestContext(t, http.MethodPost, body, uuid.New().String(), nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.ConfirmSlot(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestSlotController_RemoveValidatedSlot_InvalidClaimsType(t *testing.T) {
	ctl := &SlotController{slotService: newTestSlotService(t)}
	c, recorder := newSlotTestContext(t, http.MethodDelete, nil, uuid.New().String(), nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.RemoveValidatedSlot(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestSlotController_ConfirmSlot_InvalidSlotId(t *testing.T) {
	ctl := &SlotController{slotService: newTestSlotService(t)}
	body, _ := json.Marshal(ConfirmSlotDto{StartsAt: time.Now(), EndsAt: time.Now().Add(time.Hour)})
	c, recorder := newSlotTestContext(t, http.MethodPost, body, "not-a-uuid", &guard.Claims{Id: uuid.New()})

	ctl.ConfirmSlot(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestSlotController_ConfirmSlot_InvalidBody(t *testing.T) {
	ctl := &SlotController{slotService: newTestSlotService(t)}
	c, recorder := newSlotTestContext(t, http.MethodPost, []byte(`not-json`), uuid.New().String(), &guard.Claims{Id: uuid.New()})

	ctl.ConfirmSlot(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSlotController_ConfirmSlot_Success(t *testing.T) {
	original := mail.SmtpSendFunc
	called := make(chan struct{}, 1)
	mail.SmtpSendFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		called <- struct{}{}
		return nil
	}
	t.Cleanup(func() { mail.SmtpSendFunc = original })

	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)
	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute)}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	ctl := &SlotController{slotService: s}
	body, _ := json.Marshal(ConfirmSlotDto{StartsAt: slotEntity.StartsAt, EndsAt: slotEntity.EndsAt})
	c, recorder := newSlotTestContext(t, http.MethodPost, body, slotEntity.Id.String(), &guard.Claims{Id: owner.Id})

	ctl.ConfirmSlot(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatal("expected the confirmation email goroutine to run")
	}
}

func TestSlotController_RemoveValidatedSlot_InvalidSlotId(t *testing.T) {
	ctl := &SlotController{slotService: newTestSlotService(t)}
	c, recorder := newSlotTestContext(t, http.MethodDelete, nil, "not-a-uuid", &guard.Claims{Id: uuid.New()})

	ctl.RemoveValidatedSlot(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestSlotController_RemoveValidatedSlot_Success(t *testing.T) {
	original := mail.SmtpSendFunc
	called := make(chan struct{}, 1)
	mail.SmtpSendFunc = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		called <- struct{}{}
		return nil
	}
	t.Cleanup(func() { mail.SmtpSendFunc = original })

	s := newTestSlotService(t)
	owner := createTestAccount(t)
	event := createTestEvent(t, s, owner.Id)
	slotEntity := model.Slot{Id: uuid.New(), EventId: event.Id, StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(30 * time.Minute), IsValidated: true}
	require.NoError(t, s.slotRepository.Create(&slotEntity))

	ctl := &SlotController{slotService: s}
	c, recorder := newSlotTestContext(t, http.MethodDelete, nil, slotEntity.Id.String(), &guard.Claims{Id: owner.Id})

	ctl.RemoveValidatedSlot(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatal("expected the cancellation email goroutine to run")
	}
}
