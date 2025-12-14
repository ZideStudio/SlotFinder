package sse

import (
	"app/commons/guard"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEController struct {
	sseService *SSEService
}

func NewSSEController(service *SSEService) *SSEController {
	if service != nil {
		return &SSEController{
			sseService: service,
		}
	}

	return &SSEController{
		sseService: GetSSEService(),
	}
}

// Connect handles SSE connection for event updates
// @Summary Connect to SSE for event updates
// @Description Establishes a Server-Sent Events connection to receive real-time updates for a specific event
// @Tags SSE
// @Param eventId path string true "Event ID"
// @Success 200 {string} string "SSE connection established"
// @Failure 400 {object} map[string]string "Invalid event ID"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Router /v1/events/{eventId}/sse [get]
func (ctrl *SSEController) Connect(c *gin.Context) {
	eventIdStr := c.Param("eventId")
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userId := user.Id

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	ctrl.sseService.HandleSSEConnection(c, userId, eventId)
}
