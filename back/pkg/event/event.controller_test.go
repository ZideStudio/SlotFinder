package event

import (
	"app/commons/guard"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newEventTestContext(t *testing.T, method, target string, body []byte, eventId string, user *guard.Claims) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, target, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	if eventId != "" {
		c.Params = gin.Params{{Key: "eventId", Value: eventId}}
	}
	if user != nil {
		c.Set("user", user)
	}
	return c, recorder
}

func TestNewEventController_ReusesProvidedInstance(t *testing.T) {
	existing := &EventController{}
	assert.Same(t, existing, NewEventController(existing))
}

func TestEventController_Create_InvalidClaimsType(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	body, _ := json.Marshal(validEventCreateDto())
	c, recorder := newEventTestContext(t, http.MethodPost, "/", body, "", nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.Create(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEventController_Create_InvalidBody(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodPost, "/", []byte(`not-json`), "", &guard.Claims{Id: uuid.New()})

	ctl.Create(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestEventController_Create_Success(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	username := "ctrl-owner"
	body, _ := json.Marshal(validEventCreateDto())
	c, recorder := newEventTestContext(t, http.MethodPost, "/", body, "", &guard.Claims{Id: uuid.New(), Username: &username})

	ctl.Create(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestEventController_Update_InvalidClaimsType(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodPatch, "/", []byte(`{}`), uuid.New().String(), nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.Update(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEventController_Update_InvalidEventId(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodPatch, "/", []byte(`{}`), "not-a-uuid", &guard.Claims{Id: uuid.New()})

	ctl.Update(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestEventController_Update_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	ctl := &EventController{eventService: s}

	name := "Renamed Event"
	body, _ := json.Marshal(EventUpdateDto{Name: &name})
	c, recorder := newEventTestContext(t, http.MethodPatch, "/", body, event.Id.String(), &guard.Claims{Id: owner})

	ctl.Update(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestEventController_GetUserEvents_InvalidClaimsType(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodGet, "/?page=1&limit=10", nil, "", nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.GetUserEvents(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEventController_GetUserEvents_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	createTestEventForOwner(t, s, owner)
	ctl := &EventController{eventService: s}

	c, recorder := newEventTestContext(t, http.MethodGet, "/?page=1&limit=10", nil, "", &guard.Claims{Id: owner})

	ctl.GetUserEvents(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestEventController_GetUserEvents_InvalidPagination(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodGet, "/?page=0", nil, "", &guard.Claims{Id: uuid.New()})

	ctl.GetUserEvents(c)

	assert.NotEqual(t, http.StatusOK, recorder.Code)
}

func TestEventController_GetEventSummary_InvalidEventId(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodGet, "/", nil, "not-a-uuid", nil)

	ctl.GetEventSummary(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestEventController_GetEventSummary_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	ctl := &EventController{eventService: s}

	c, recorder := newEventTestContext(t, http.MethodGet, "/", nil, event.Id.String(), nil)

	ctl.GetEventSummary(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestEventController_GetEvent_InvalidClaimsType(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodGet, "/", nil, uuid.New().String(), nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.GetEvent(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEventController_GetEvent_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	ctl := &EventController{eventService: s}

	c, recorder := newEventTestContext(t, http.MethodGet, "/", nil, event.Id.String(), &guard.Claims{Id: owner})

	ctl.GetEvent(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestEventController_JoinEvent_InvalidClaimsType(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	c, recorder := newEventTestContext(t, http.MethodPost, "/", nil, uuid.New().String(), nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.JoinEvent(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEventController_JoinEvent_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	ctl := &EventController{eventService: s}

	newMember := uuid.New()
	c, recorder := newEventTestContext(t, http.MethodPost, "/", nil, event.Id.String(), &guard.Claims{Id: newMember})

	ctl.JoinEvent(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestEventController_UpdateProfile_InvalidClaimsType(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	body, _ := json.Marshal(EventProfileDto{Color: "#FFFFFF"})
	c, recorder := newEventTestContext(t, http.MethodPatch, "/", body, uuid.New().String(), nil)
	c.Set("user", "not-a-claims-pointer")

	ctl.UpdateProfile(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestEventController_UpdateProfile_InvalidEventId(t *testing.T) {
	ctl := &EventController{eventService: newTestEventService(t)}
	body, _ := json.Marshal(EventProfileDto{Color: "#FFFFFF"})
	c, recorder := newEventTestContext(t, http.MethodPatch, "/", body, "not-a-uuid", &guard.Claims{Id: uuid.New()})

	ctl.UpdateProfile(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestEventController_UpdateProfile_Success(t *testing.T) {
	s := newTestEventService(t)
	owner := createTestOwner(t)
	event := createTestEventForOwner(t, s, owner)
	ctl := &EventController{eventService: s}

	body, _ := json.Marshal(EventProfileDto{Color: "#AABBCC"})
	c, recorder := newEventTestContext(t, http.MethodPatch, "/", body, event.Id.String(), &guard.Claims{Id: owner})

	ctl.UpdateProfile(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
