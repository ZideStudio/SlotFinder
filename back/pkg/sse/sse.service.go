package sse

import (
	model "app/db/models"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SSEClient struct {
	ID      string
	userId  uuid.UUID
	EventId uuid.UUID
	Channel chan []byte
	Context context.Context
	Cancel  context.CancelFunc
}

type SSEService struct {
	clients map[string]*SSEClient
	mutex   sync.RWMutex
}

type SlotUpdateMessage []model.Slot

var sseServiceInstance *SSEService
var sseServiceOnce sync.Once

// GetSSEService returns the singleton SSE service instance
func GetSSEService() *SSEService {
	sseServiceOnce.Do(func() {
		sseServiceInstance = &SSEService{
			clients: make(map[string]*SSEClient),
			mutex:   sync.RWMutex{},
		}
	})
	return sseServiceInstance
}

// NewSSEService creates a new SSE service instance
func NewSSEService() *SSEService {
	return &SSEService{
		clients: make(map[string]*SSEClient),
		mutex:   sync.RWMutex{},
	}
}

// AddClient adds a new SSE client
func (s *SSEService) AddClient(clientID string, userId uuid.UUID, eventId uuid.UUID, ctx context.Context) *SSEClient {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientCtx, cancel := context.WithCancel(ctx)

	client := &SSEClient{
		ID:      clientID,
		userId:  userId,
		EventId: eventId,
		Channel: make(chan []byte, 10), // Buffer of 10 messages
		Context: clientCtx,
		Cancel:  cancel,
	}

	s.clients[clientID] = client
	log.Debug().Str("clientID", clientID).Str("userId", userId.String()).Str("eventId", eventId.String()).Msg("SSE client connected")

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
		log.Debug().Str("clientID", clientID).Msg("SSE client disconnected")
	}
}

// BroadcastSlotsUpdate sends slot updates to all participants of an event
func (s *SSEService) BroadcastSlotsUpdate(eventId uuid.UUID, slots []model.Slot) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	log.Debug().Str("eventId", eventId.String()).Int("slotsCount", len(slots)).Msg("Starting to broadcast slot update")

	message := SlotUpdateMessage(slots)

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal slot update message")
		return
	}

	log.Debug().Str("messageJSON", string(messageBytes)).Msg("Marshaled SSE message")

	var connectedClients []string
	var sentCount int
	var totalClients int

	for clientID, client := range s.clients {
		totalClients++
		if client.EventId == eventId {
			log.Debug().Str("clientID", clientID).Str("userId", client.userId.String()).Msg("Sending message to client")
			select {
			case client.Channel <- messageBytes:
				connectedClients = append(connectedClients, clientID)
				sentCount++
				log.Debug().Str("clientID", clientID).Msg("Message sent successfully to client")
			case <-client.Context.Done():
				// Client context is done, remove it
				log.Warn().Str("clientID", clientID).Msg("Client context done, removing client")
				go s.RemoveClient(clientID)
			default:
				log.Warn().Str("clientID", clientID).Msg("Client channel full, skipping message")
			}
		} else {
			log.Debug().Str("clientID", clientID).Str("clientEventID", client.EventId.String()).Str("targetEventID", eventId.String()).Msg("Client not for this event, skipping")
		}
	}

	log.Debug().Str("eventId", eventId.String()).Strs("clients", connectedClients).Int("sentCount", sentCount).Int("totalClients", totalClients).Msg("Broadcasted slot update to connected clients")
}

// GetConnectedClientsCount returns the number of connected clients for an event
func (s *SSEService) GetConnectedClientsCount(eventId uuid.UUID) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	count := 0
	for _, client := range s.clients {
		if client.EventId == eventId {
			count++
		}
	}
	return count
}

// HandleSSEConnection handles the SSE connection endpoint
func (s *SSEService) HandleSSEConnection(c *gin.Context, userId uuid.UUID, eventId uuid.UUID) {
	clientID := fmt.Sprintf("%s-%s-%d", userId.String(), eventId.String(), time.Now().UnixNano())

	// Add client to SSE service
	client := s.AddClient(clientID, userId, eventId, c.Request.Context())
	defer s.RemoveClient(clientID)

	log.Debug().Str("clientID", clientID).Str("userId", userId.String()).Str("eventId", eventId.String()).Msg("SSE client connection established")

	// Listen for messages and client disconnect
	for {
		select {
		case message := <-client.Channel:
			log.Debug().Str("clientID", clientID).Str("message", string(message)).Msg("Sending SSE message to client")

			fmt.Fprintf(c.Writer, "data: %s\n\n", string(message))
			c.Writer.Flush()

			log.Debug().Str("clientID", clientID).Msg("SSE message sent and flushed")
		case <-client.Context.Done():
			log.Debug().Str("clientID", clientID).Msg("SSE client context done")
			return
		case <-c.Request.Context().Done():
			log.Debug().Str("clientID", clientID).Msg("SSE connection closed by client")
			return
		}
	}
}
