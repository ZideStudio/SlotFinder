package sse

import (
	model "app/db/models"
	"app/db/repository"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEClient struct {
	Id      string
	userId  uuid.UUID
	EventId uuid.UUID
	Channel chan []byte
	Context context.Context
	Cancel  context.CancelFunc
}

type SSEService struct {
	clients         map[string]*SSEClient         // clientID -> client
	clientsByEvent  map[uuid.UUID]map[string]bool // eventID -> set of clientIDs
	mutex           sync.RWMutex
	eventRepository *repository.EventRepository
	slotRepository  *repository.SlotRepository
}

type SlotUpdateMessage []model.Slot

var sseServiceInstance *SSEService
var sseServiceOnce sync.Once

// GetSSEService returns the singleton SSE service instance
func GetSSEService() *SSEService {
	sseServiceOnce.Do(func() {
		sseServiceInstance = &SSEService{
			clients:        make(map[string]*SSEClient),
			clientsByEvent: make(map[uuid.UUID]map[string]bool),
			mutex:          sync.RWMutex{},
			slotRepository: &repository.SlotRepository{},
		}
	})
	return sseServiceInstance
}

// NewSSEService creates a new SSE service instance
func NewSSEService() *SSEService {
	return &SSEService{
		clients:        make(map[string]*SSEClient),
		clientsByEvent: make(map[uuid.UUID]map[string]bool),
		mutex:          sync.RWMutex{},
		slotRepository: &repository.SlotRepository{},
	}
}

// AddClient adds a new SSE client
func (s *SSEService) AddClient(clientID string, userId uuid.UUID, eventId uuid.UUID, ctx context.Context) *SSEClient {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientCtx, cancel := context.WithCancel(ctx)

	client := &SSEClient{
		Id:      clientID,
		userId:  userId,
		EventId: eventId,
		Channel: make(chan []byte, 10), // Buffer of 10 messages
		Context: clientCtx,
		Cancel:  cancel,
	}

	s.clients[clientID] = client

	// Add to event index
	if s.clientsByEvent[eventId] == nil {
		s.clientsByEvent[eventId] = make(map[string]bool)
	}
	s.clientsByEvent[eventId][clientID] = true

	return client
}

// RemoveClient removes an SSE client
func (s *SSEService) RemoveClient(clientID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[clientID]; exists {
		client.Cancel()
		close(client.Channel)
		delete(s.clients, clientID)

		// Remove from event index
		if eventClients, exists := s.clientsByEvent[client.EventId]; exists {
			delete(eventClients, clientID)
			// Clean up empty event entries
			if len(eventClients) == 0 {
				delete(s.clientsByEvent, client.EventId)
			}
		}
	}
}

// BroadcastSlotsUpdate sends slot updates to all participants of an event
func (s *SSEService) BroadcastSlotsUpdate(eventId uuid.UUID, slots []model.Slot) {
	s.mutex.RLock()

	message := SlotUpdateMessage(slots)

	messageBytes, err := json.Marshal(message)
	if err != nil {
		s.mutex.RUnlock()
		return
	}

	var connectedClients []string
	var clientsToRemove []string
	var sentCount int

	// O(1) lookup of clients for this event
	if eventClients, exists := s.clientsByEvent[eventId]; exists {
		for clientID := range eventClients {
			if client, exists := s.clients[clientID]; exists {
				select {
				case client.Channel <- messageBytes:
					connectedClients = append(connectedClients, clientID)
					sentCount++
				case <-client.Context.Done():
					// Collect clients to remove instead of removing immediately
					clientsToRemove = append(clientsToRemove, clientID)
				default:
					// Channel full, skip this client
				}
			} else {
				// Client doesn't exist, mark for removal from event index
				clientsToRemove = append(clientsToRemove, clientID)
			}
		}
	}

	s.mutex.RUnlock()

	// Remove disconnected clients after releasing the read lock
	for _, clientID := range clientsToRemove {
		s.RemoveClient(clientID)
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
	clientID := fmt.Sprintf("%s-%s-%d", userId.String(), eventId.String(), time.Now().UnixNano())

	// Add client to SSE service
	client := s.AddClient(clientID, userId, eventId, c.Request.Context())
	defer s.RemoveClient(clientID)

	// Check if user has access to the event
	var event model.Event
	err := s.eventRepository.FindOneById(eventId, &event)
	if err != nil || !event.HasUserAccess(&userId) {
		c.JSON(403, gin.H{"error": "Access denied to event"})
		return
	}

	// Send current event slots on connection
	var currentSlots []model.Slot
	if err := s.slotRepository.FindByEventId(eventId, &currentSlots); err != nil {
		currentSlots = []model.Slot{} // Fallback to empty array if error
	}

	initialMessage := SlotUpdateMessage(currentSlots)
	initialBytes, _ := json.Marshal(initialMessage)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(initialBytes))
	c.Writer.Flush()

	// Listen for messages and client disconnect
	for {
		select {
		case message := <-client.Channel:
			fmt.Fprintf(c.Writer, "data: %s\n\n", string(message))
			c.Writer.Flush()
		case <-client.Context.Done():
			return
		case <-c.Request.Context().Done():
			return
		}
	}
}
