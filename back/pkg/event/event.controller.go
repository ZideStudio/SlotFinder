package event

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Router /api/v1/events [post]
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

// @Summary Update an event
// @Tags Event
// @Accept json
// @Produce json
// @Param eventId path string true "Event Id"
// @Param data body EventUpdateDto true "Event parameters"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/events/{eventId} [patch]
func (ctl *EventController) Update(c *gin.Context) {
	var data EventUpdateDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	idUuid, err := uuid.Parse(c.Param("eventId"))
	if err != nil {
		helpers.HandleJSONResponse(c, nil, constants.ERR_EVENT_NOT_FOUND.Err)
		return
	}

	err = ctl.eventService.Update(idUuid, &data, user)

	helpers.HandleJSONResponse(c, nil, err)
}

// @Summary Get user events
// @Tags Event
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} EventResponse
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/events [get]
func (ctl *EventController) GetUserEvents(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	events, err := ctl.eventService.GetUserEvents(user)
	helpers.HandleJSONResponse(c, events, err)
}

// @Summary Get event
// @Tags Event
// @Param eventId path string true "Event Id"
// @Accept json
// @Produce json
// @Success 200 {object} EventResponse
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/events/{eventId} [get]
func (ctl *EventController) GetEvent(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	idUuid, err := uuid.Parse(c.Param("eventId"))
	if err != nil {
		helpers.HandleJSONResponse(c, nil, constants.ERR_EVENT_NOT_FOUND.Err)
		return
	}

	events, err := ctl.eventService.GetEvent(idUuid, user)
	helpers.HandleJSONResponse(c, events, err)
}

// @Summary Join event
// @Tags Event
// @Param eventId path string true "Event Id"
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/events/{eventId}/join [post]
func (ctl *EventController) JoinEvent(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	idUuid, err := uuid.Parse(c.Param("eventId"))
	if err != nil {
		helpers.HandleJSONResponse(c, nil, constants.ERR_EVENT_NOT_FOUND.Err)
		return
	}

	events, err := ctl.eventService.JoinEvent(idUuid, user)
	helpers.HandleJSONResponse(c, events, err)
}

// @Summary Update event profile
// @Tags Event
// @Accept json
// @Produce json
// @Param eventId path string true "Event Id"
// @Param data body EventProfileDto true "Event profile parameters"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/events/{eventId}/profile [patch]
func (ctl *EventController) UpdateProfile(c *gin.Context) {
	var data EventProfileDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	idUuid, err := uuid.Parse(c.Param("eventId"))
	if err != nil {
		helpers.HandleJSONResponse(c, nil, constants.ERR_EVENT_NOT_FOUND.Err)
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	err = ctl.eventService.UpdateProfile(&data, idUuid, user)

	helpers.HandleJSONResponse(c, nil, err)
}
