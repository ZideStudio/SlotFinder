package sse

import (
	model "app/db/models"
	"app/db/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SSEClient struct {
	Id      string
	UserId  uuid.UUID
	EventId uuid.UUID
	Channel chan []byte
	Context context.Context
	Cancel  context.CancelFunc
}

type SSEService struct {
	clients         map[string]*SSEClient         // clientId -> client
	clientsByEvent  map[uuid.UUID]map[string]bool // eventId -> set of clientIds
	mutex           sync.RWMutex
	eventRepository *repository.EventRepository
	slotRepository  *repository.SlotRepository
}

type SlotUpdateMessage []model.Slot

const (
	defaultChannelBuffer = 10 // Buffer size for SSE client channels
)

var sseServiceInstance *SSEService
var sseServiceOnce sync.Once

// GetSSEService returns the singleton SSE service instance
func GetSSEService() *SSEService {
	sseServiceOnce.Do(func() {
		sseServiceInstance = &SSEService{
			clients:         make(map[string]*SSEClient),
			clientsByEvent:  make(map[uuid.UUID]map[string]bool),
			eventRepository: &repository.EventRepository{},
			slotRepository:  &repository.SlotRepository{},
		}
	})
	return sseServiceInstance
}

// NewSSEService creates a new SSE service instance
func NewSSEService() *SSEService {
	return &SSEService{
		clients:         make(map[string]*SSEClient),
		clientsByEvent:  make(map[uuid.UUID]map[string]bool),
		eventRepository: &repository.EventRepository{},
		slotRepository:  &repository.SlotRepository{},
	}
}

// AddClient adds a new SSE client
func (s *SSEService) AddClient(clientId string, userId uuid.UUID, eventId uuid.UUID, ctx context.Context) *SSEClient {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientCtx, cancel := context.WithCancel(ctx)

	client := &SSEClient{
		Id:      clientId,
		UserId:  userId,
		EventId: eventId,
		Channel: make(chan []byte, defaultChannelBuffer),
		Context: clientCtx,
		Cancel:  cancel,
	}

	s.clients[clientId] = client

	// Add to event index
	if s.clientsByEvent[eventId] == nil {
		s.clientsByEvent[eventId] = make(map[string]bool)
	}
	s.clientsByEvent[eventId][clientId] = true

	return client
}

// RemoveClient removes an SSE client
func (s *SSEService) RemoveClient(clientId string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[clientId]; exists {
		client.Cancel()
		close(client.Channel)
		delete(s.clients, clientId)

		// Remove from event index
		eventClients, exists := s.clientsByEvent[client.EventId]
		if !exists {
			return
		}

		delete(eventClients, clientId)
		// Clean up empty event entries
		if len(eventClients) != 0 {
			return
		}

		delete(s.clientsByEvent, client.EventId)
	}
}

// BroadcastSlotsUpdate sends slot updates to all participants of an event
func (s *SSEService) BroadcastSlotsUpdate(eventId uuid.UUID, slots []model.Slot) {
	s.mutex.RLock()

	message := SlotUpdateMessage(slots)

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to marshal SSE message")
		s.mutex.RUnlock()
		return
	}

	var clientsToRemove []string
	var sentCount int
	if eventClients, exists := s.clientsByEvent[eventId]; exists {
		for clientId := range eventClients {
			if client, exists := s.clients[clientId]; exists {
				select {
				case client.Channel <- messageBytes:
					sentCount++
				case <-client.Context.Done():
					// Collect clients to remove instead of removing immediately
					clientsToRemove = append(clientsToRemove, clientId)
				}
			} else {
				// Client doesn't exist, mark for removal from event index
				clientsToRemove = append(clientsToRemove, clientId)
			}
		}
	}

	s.mutex.RUnlock()

	// Remove disconnected clients after releasing the read lock
	for _, clientId := range clientsToRemove {
		s.RemoveClient(clientId)
	}
}

// GetConnectedClientsCount returns the number of connected clients for an event
func (s *SSEService) GetConnectedClientsCount(eventId uuid.UUID) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if eventClients, exists := s.clientsByEvent[eventId]; exists {
		return len(eventClients)
	}
	return 0
}

// HandleSSEConnection handles the SSE connection endpoint
func (s *SSEService) HandleSSEConnection(c *gin.Context, userId uuid.UUID, eventId uuid.UUID) {
	clientId := fmt.Sprintf("%s-%s-%d", userId.String(), eventId.String(), time.Now().UnixNano())

	// Add client to SSE service
	client := s.AddClient(clientId, userId, eventId, c.Request.Context())
	defer s.RemoveClient(clientId)

	// Check if user has access to the event
	var event model.Event
	err := s.eventRepository.FindOneById(eventId, &event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if !event.HasUserAccess(&userId) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to event"})
		return
	}

	// Send current event slots on connection
	var currentSlots []model.Slot
	if err := s.slotRepository.FindByEventId(eventId, &currentSlots); err != nil {
		currentSlots = []model.Slot{} // Fallback to empty array if error
	}

	initialMessage := SlotUpdateMessage(currentSlots)
	initialBytes, err := json.Marshal(initialMessage)
	if err != nil {
		log.Error().Err(err).Str("eventId", eventId.String()).Msg("Failed to marshal initial SSE message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Send initial message with error handling
	if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", string(initialBytes)); err != nil {
		return
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()

	// Listen for messages and client disconnect
	for {
		select {
		case message := <-client.Channel:
			if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", string(message)); err != nil {
				log.Error().Err(err).Str("clientId", clientId).Msg("Failed to send SSE message to client")
				return
			}

			flusher.Flush()
		case <-client.Context.Done():
			return
		case <-c.Request.Context().Done():
			return
		}
	}
}
