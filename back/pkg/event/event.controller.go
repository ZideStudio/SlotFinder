package event

import (
	"app/commons/guard"
	"app/commons/helpers"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService *EventService
}

func NewEventController(ctl *EventController) *EventController {
	if ctl != nil {
		return ctl
	}

	return &EventController{
		eventService: NewEventService(nil),
	}
}

// @Summary Create an event
// @Tags Event
// @Accept json
// @Produce json
// @Param data body EventCreateDto true "Event parameters"
// @Security BearerAuth
// @Success 200 {object} model.Event
// @Failure 400 {object} helpers.ApiError
// @Router /v1/event [post]
func (ctl *EventController) Create(c *gin.Context) {
	var data EventCreateDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	event, err := ctl.eventService.Create(&data, user)

	helpers.HandleJSONResponse(c, event, err)
}

// @Summary Get user events
// @Tags Event
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} EventResponse
// @Failure 400 {object} helpers.ApiError
// @Router /v1/event [get]
func (ctl *EventController) GetUserEvents(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	events, err := ctl.eventService.GetUserEvents(user)
	helpers.HandleJSONResponse(c, events, err)
}
