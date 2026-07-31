package availability

import (
	"app/commons/guard"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAvailabilityTestContext(t *testing.T, method string, body []byte, eventId, availabilityId string, user *guard.Claims) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	params := gin.Params{}
	if eventId != "" {
		params = append(params, gin.Param{Key: "eventId", Value: eventId})
	}
	if availabilityId != "" {
		params = append(params, gin.Param{Key: "availabilityId", Value: availabilityId})
	}
	c.Params = params
	if user != nil {
		c.Set("user", user)
	}
	return c, recorder
}

func TestNewAvailabilityController_ReusesProvidedInstance(t *testing.T) {
	existing := &AvailabilityController{}
	assert.Same(t, existing, NewAvailabilityController(existing))
}

func TestAvailabilityController_Create_InvalidBody(t *testing.T) {
	ctl := &AvailabilityController{availabilityService: newTestAvailabilityService(t)}
	c, recorder := newAvailabilityTestContext(t, http.MethodPost, []byte(`not-json`), uuid.New().String(), "", &guard.Claims{Id: uuid.New()})

	ctl.Create(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// NOTE: no "unauthenticated" controller test here — Create/Update/Delete
// don't nil-check `user` after guard.GetUserClaims before passing it into
// the service, so an unauthenticated call panics on a nil pointer deref.
// This is unreachable in production because every route wires
// guard.AuthCheck(nil) (RequireAuthentication: true) in front of these
// handlers (see server/router.go), which rejects the request with 401
// before the controller ever runs. Simulating "no claims set" directly
// against the controller (bypassing that middleware) would just crash the
// test process, so it's intentionally not covered here.

func TestAvailabilityController_Create_InvalidEventId(t *testing.T) {
	ctl := &AvailabilityController{availabilityService: newTestAvailabilityService(t)}
	body, _ := json.Marshal(AvailabilityCreateDto{StartsAt: alignedNow(t), EndsAt: alignedNow(t).Add(time.Hour)})
	c, recorder := newAvailabilityTestContext(t, http.MethodPost, body, "not-a-uuid", "", &guard.Claims{Id: uuid.New()})

	ctl.Create(c)

	// getEventIdParam maps a bad UUID to ERR_EVENT_NOT_FOUND, a registered
	// custom error with its own (404) status.
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestAvailabilityController_Create_Success(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	ctl := &AvailabilityController{availabilityService: s}

	body, _ := json.Marshal(AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)})
	c, recorder := newAvailabilityTestContext(t, http.MethodPost, body, event.Id.String(), "", &guard.Claims{Id: owner})

	ctl.Create(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAvailabilityController_Update_InvalidAvailabilityId(t *testing.T) {
	ctl := &AvailabilityController{availabilityService: newTestAvailabilityService(t)}
	c, recorder := newAvailabilityTestContext(t, http.MethodPatch, []byte(`{}`), "", "not-a-uuid", &guard.Claims{Id: uuid.New()})

	ctl.Update(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestAvailabilityController_Update_Success(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	ctl := &AvailabilityController{availabilityService: s}
	newEnd := event.StartsAt.Add(90 * time.Minute)
	body, _ := json.Marshal(AvailabilityUpdateDto{EndsAt: &newEnd})
	c, recorder := newAvailabilityTestContext(t, http.MethodPatch, body, "", created.Id.String(), claims)

	ctl.Update(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAvailabilityController_Delete_InvalidAvailabilityId(t *testing.T) {
	ctl := &AvailabilityController{availabilityService: newTestAvailabilityService(t)}
	c, recorder := newAvailabilityTestContext(t, http.MethodDelete, nil, "", "not-a-uuid", &guard.Claims{Id: uuid.New()})

	ctl.Delete(c)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestAvailabilityController_Delete_Success(t *testing.T) {
	s := newTestAvailabilityService(t)
	owner := createTestUser(t)
	event := createEventWithAccess(t, s.eventRepository, owner)
	claims := &guard.Claims{Id: owner}

	created, err := s.Create(&AvailabilityCreateDto{StartsAt: event.StartsAt, EndsAt: event.StartsAt.Add(time.Hour)}, event.Id, claims)
	require.NoError(t, err)

	ctl := &AvailabilityController{availabilityService: s}
	c, recorder := newAvailabilityTestContext(t, http.MethodDelete, nil, "", created.Id.String(), claims)

	ctl.Delete(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
