package sse

import (
	model "app/db/models"
	"app/db/repository"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestService() *SSEService {
	return &SSEService{
		clients:        make(map[string]*SSEClient),
		clientsByEvent: make(map[uuid.UUID]map[string]bool),
	}
}

func TestNewSSEService(t *testing.T) {
	s := NewSSEService()
	assert.NotNil(t, s)
	assert.NotNil(t, s.clients)
	assert.NotNil(t, s.clientsByEvent)
}

func TestGetSSEService_Singleton(t *testing.T) {
	s1 := GetSSEService()
	s2 := GetSSEService()
	assert.Same(t, s1, s2)
}

func TestAddClient(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	userId := uuid.New()

	client := s.AddClient("client-1", userId, eventId, context.Background())
	assert.NotNil(t, client)
	assert.Equal(t, "client-1", client.Id)
	assert.Equal(t, 1, s.GetConnectedClientsCount(eventId))
}

func TestRemoveClient(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	client := s.AddClient("client-1", uuid.New(), eventId, context.Background())

	s.RemoveClient("client-1")

	assert.Equal(t, 0, s.GetConnectedClientsCount(eventId))
	_, stillOpen := <-client.Channel
	assert.False(t, stillOpen, "channel should be closed")
}

func TestRemoveClient_Unknown_NoOp(t *testing.T) {
	s := newTestService()
	// Should not panic when removing a client that was never added.
	s.RemoveClient("unknown")
}

func TestRemoveClient_MissingEventIndex_NoOp(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	ctx, cancel := context.WithCancel(context.Background())
	client := &SSEClient{Id: "client-1", EventId: eventId, Channel: make(chan []byte, 1), Context: ctx, Cancel: cancel}
	// Inserted directly into s.clients without a matching s.clientsByEvent
	// entry, exercising the "missing index" early-return branch instead of
	// the normal cleanup path (which AddClient's bookkeeping always avoids).
	s.clients["client-1"] = client

	s.RemoveClient("client-1")

	_, exists := s.clients["client-1"]
	assert.False(t, exists)
}

func TestRemoveClient_LeavesOtherClientsOfSameEvent(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	s.AddClient("client-1", uuid.New(), eventId, context.Background())
	s.AddClient("client-2", uuid.New(), eventId, context.Background())

	s.RemoveClient("client-1")

	assert.Equal(t, 1, s.GetConnectedClientsCount(eventId))
}

func TestGetConnectedClientsCount_NoClients(t *testing.T) {
	s := newTestService()
	assert.Equal(t, 0, s.GetConnectedClientsCount(uuid.New()))
}

func TestBroadcastSlotsUpdate_DeliversToClient(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	client := s.AddClient("client-1", uuid.New(), eventId, context.Background())

	s.BroadcastSlotsUpdate(eventId, nil)

	select {
	case msg := <-client.Channel:
		assert.Equal(t, "null", string(msg))
	case <-time.After(time.Second):
		t.Fatal("expected a message on the client channel")
	}
}

func TestBroadcastSlotsUpdate_MarshalFails_NoOp(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	s.AddClient("client-1", uuid.New(), eventId, context.Background())

	// time.Time.MarshalJSON errors for years outside [0,9999], forcing
	// json.Marshal(message) to fail before any client delivery is attempted.
	invalidSlot := model.Slot{StartsAt: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)}

	// Should not panic; the marshal error is logged and the call returns early.
	s.BroadcastSlotsUpdate(eventId, []model.Slot{invalidSlot})
}

func TestBroadcastSlotsUpdate_NoClientsForEvent_NoOp(t *testing.T) {
	s := newTestService()
	// Should be a no-op, not panic, when nobody is listening for this event.
	s.BroadcastSlotsUpdate(uuid.New(), nil)
}

func TestBroadcastSlotsUpdate_RemovesDisconnectedClients(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	client := s.AddClient("client-1", uuid.New(), eventId, context.Background())

	// Fill the buffered channel to capacity so the `client.Channel <- msg` select
	// case can never be ready, forcing the cancelled-context branch deterministically.
	for i := 0; i < defaultChannelBuffer; i++ {
		client.Channel <- []byte("filler")
	}
	client.Cancel()

	s.BroadcastSlotsUpdate(eventId, nil)

	assert.Equal(t, 0, s.GetConnectedClientsCount(eventId))
}

func TestBroadcastSlotsUpdate_SkipsClientMissingFromRegistry(t *testing.T) {
	s := newTestService()
	eventId := uuid.New()
	s.clientsByEvent[eventId] = map[string]bool{"ghost-client": true}

	// Exercises the `else` branch (client id present in the event index but
	// missing from s.clients). Must not panic. Note RemoveClient is itself a
	// no-op for ids absent from s.clients, so the stale index entry survives
	// this call — that's existing behavior, not something this test changes.
	s.BroadcastSlotsUpdate(eventId, nil)

	assert.Equal(t, 1, s.GetConnectedClientsCount(eventId))
}

func TestHandleSSEConnection_EventNotFound(t *testing.T) {
	s := newTestService()
	s.eventRepository = repository.NewEventRepository(testDB(t))
	s.slotRepository = repository.NewSlotRepository(testDB(t))

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	s.HandleSSEConnection(c, uuid.New(), uuid.New())

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestHandleSSEConnection_AccessDenied(t *testing.T) {
	db := testDB(t)
	s := newTestService()
	s.eventRepository = repository.NewEventRepository(db)
	s.slotRepository = repository.NewSlotRepository(db)

	owner := model.Account{Id: uuid.New()}
	require.NoError(t, db.Create(&owner).Error)
	event := model.Event{Id: uuid.New(), Name: "Event", Duration: 30, OwnerId: owner.Id, Status: "IN_DECISION"}
	require.NoError(t, db.Create(&event).Error)
	require.NoError(t, db.Create(&model.AccountEvent{AccountId: owner.Id, EventId: event.Id}).Error)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	// A user who is not a participant of the event should be denied access.
	s.HandleSSEConnection(c, uuid.New(), event.Id)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestHandleSSEConnection_StreamsInitialMessageThenCloses(t *testing.T) {
	db := testDB(t)
	s := newTestService()
	s.eventRepository = repository.NewEventRepository(db)
	s.slotRepository = repository.NewSlotRepository(db)

	userId := uuid.New()
	account := model.Account{Id: userId}
	require.NoError(t, db.Create(&account).Error)
	event := model.Event{Id: uuid.New(), Name: "Event", Duration: 30, OwnerId: userId, Status: "IN_DECISION"}
	require.NoError(t, db.Create(&event).Error)
	require.NoError(t, db.Create(&model.AccountEvent{AccountId: userId, EventId: event.Id}).Error)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/sse/:eventId", func(c *gin.Context) {
		s.HandleSSEConnection(c, userId, event.Id)
	})
	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(server.URL + "/sse/" + event.Id.String())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	buf := make([]byte, 64)
	n, err := resp.Body.Read(buf)
	require.NoError(t, err)
	assert.Contains(t, string(buf[:n]), "data: ")
}

// testDB opens a throwaway in-memory sqlite DB migrated for the models this
// package's repositories touch (Event/Slot/Account/Availability/AccountEvent).
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&model.Account{}, &model.Event{}, &model.Slot{}, &model.Availability{}, &model.AccountEvent{}))
	return database
}

// closedRepoDB returns a gorm.DB whose underlying connection is already
// closed, so any query through it fails immediately — used to force a
// single repository's calls to fail while the rest of the service keeps
// using a working DB.
func closedRepoDB(t *testing.T) *gorm.DB {
	t.Helper()
	database := testDB(t)
	sqlDB, err := database.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
	return database
}

func TestHandleSSEConnection_SlotsLookupFails_FallsBackToEmpty(t *testing.T) {
	db := testDB(t)
	s := newTestService()
	s.eventRepository = repository.NewEventRepository(db)
	s.slotRepository = repository.NewSlotRepository(closedRepoDB(t))

	userId := uuid.New()
	account := model.Account{Id: userId}
	require.NoError(t, db.Create(&account).Error)
	event := model.Event{Id: uuid.New(), Name: "Event", Duration: 30, OwnerId: userId, Status: "IN_DECISION"}
	require.NoError(t, db.Create(&event).Error)
	require.NoError(t, db.Create(&model.AccountEvent{AccountId: userId, EventId: event.Id}).Error)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/sse/:eventId", func(c *gin.Context) {
		s.HandleSSEConnection(c, userId, event.Id)
	})
	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(server.URL + "/sse/" + event.Id.String())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	buf := make([]byte, 64)
	n, err := resp.Body.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "data: []\n\n", string(buf[:n]))
}
