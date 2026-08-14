package sse

import (
	"app/commons/guard"
	"app/db/repository"
	"app/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newSSETestContext(eventId string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{{Key: "eventId", Value: eventId}}
	return c, recorder
}

func TestNewSSEController_ReusesProvidedInstance(t *testing.T) {
	svc := newTestService()
	ctrl := NewSSEController(svc)
	assert.Same(t, svc, ctrl.sseService)
}

func TestNewSSEController_Nil_FallsBackToSingleton(t *testing.T) {
	ctrl := NewSSEController(nil)
	assert.Same(t, GetSSEService(), ctrl.sseService)
}

func TestConnect_InvalidEventId(t *testing.T) {
	ctrl := &SSEController{sseService: newTestService()}
	c, recorder := newSSETestContext("not-a-uuid")

	ctrl.Connect(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestConnect_Unauthenticated_NoClaims(t *testing.T) {
	ctrl := &SSEController{sseService: newTestService()}
	c, recorder := newSSETestContext(uuid.New().String())

	ctrl.Connect(c)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestConnect_Unauthenticated_InvalidClaimsType(t *testing.T) {
	ctrl := &SSEController{sseService: newTestService()}
	c, recorder := newSSETestContext(uuid.New().String())
	c.Set("user", "not-a-claims-pointer")

	ctrl.Connect(c)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestConnect_DelegatesToService(t *testing.T) {
	db := testutils.TestDB(t)
	svc := newTestService()
	svc.eventRepository = repository.NewEventRepository(db)
	svc.slotRepository = repository.NewSlotRepository(db)
	ctrl := &SSEController{sseService: svc}

	c, recorder := newSSETestContext(uuid.New().String())
	c.Set("user", &guard.Claims{Id: uuid.New()})

	ctrl.Connect(c)

	// The service can't find the (nonexistent) event, so it responds 404 —
	// this proves Connect() successfully parsed the params/claims and
	// delegated into HandleSSEConnection.
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "text/event-stream", c.Writer.Header().Get("Content-Type"))
}
