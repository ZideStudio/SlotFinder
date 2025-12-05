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

// extracts and validates the availability parameter from the URL path.
func (ctl *AvailabilityController) getAvailabilityIdParam(c *gin.Context) (availabilityIdUuid uuid.UUID, err error) {
	availability := c.Param("availabilityId")
	if availability == "" {
		return availabilityIdUuid, constants.ERR_AVAILABILITY_NOT_FOUND.Err
	}

	availabilityIdUuid, err = uuid.Parse(availability)
	if err != nil || availabilityIdUuid == uuid.Nil {
		return availabilityIdUuid, constants.ERR_AVAILABILITY_NOT_FOUND.Err
	}

	return availabilityIdUuid, nil
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
// @Router /api/v1/events/{eventId}/availability [post]
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

// @Summary Delete an availability
// @Tags Availability
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param availabilityId path string true "Availability ID"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/events/{eventId}/availability/{availabilityId} [delete]
func (ctl *AvailabilityController) Delete(c *gin.Context) {
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

	availabilityId, err := ctl.getAvailabilityIdParam(c)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	err = ctl.availabilityService.Delete(eventId, availabilityId, user)

	helpers.HandleJSONResponse(c, nil, err)
}
