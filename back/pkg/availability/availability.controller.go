package availability

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AvailabilityController struct {
	availabilityService *AvailabilityService
}

func NewAvailabilityController(ctl *AvailabilityController) *AvailabilityController {
	if ctl != nil {
		return ctl
	}

	return &AvailabilityController{
		availabilityService: NewAvailabilityService(nil),
	}
}

// extracts and validates the eventId parameter from the URL path.
func (ctl *AvailabilityController) getEventIdParam(c *gin.Context) (eventIdUuid uuid.UUID, err error) {
	eventId := c.Param("eventId")
	if eventId == "" {
		return eventIdUuid, constants.ERR_EVENT_NOT_FOUND.Err
	}

	eventIdUuid, err = uuid.Parse(eventId)
	if err != nil || eventIdUuid == uuid.Nil {
		return eventIdUuid, constants.ERR_EVENT_NOT_FOUND.Err
	}

	return eventIdUuid, nil
}

// @Summary Create an availability
// @Tags Availability
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param data body AvailabilityCreateDto true "Availability parameters"
// @Security BearerAuth
// @Success 200 {object} model.Availability
// @Failure 400 {object} helpers.ApiError
// @Router /v1/events/{eventId}/availability [post]
func (ctl *AvailabilityController) Create(c *gin.Context) {
	var data AvailabilityCreateDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	eventId, err := ctl.getEventIdParam(c)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	availability, err := ctl.availabilityService.Create(&data, eventId, user)

	helpers.HandleJSONResponse(c, availability, err)
}
