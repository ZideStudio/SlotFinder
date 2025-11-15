package availability

import (
	"app/commons/guard"
	"app/commons/helpers"

	"github.com/gin-gonic/gin"
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

// @Summary Create an availability
// @Tags Availability
// @Accept json
// @Produce json
// @Param data body AvailabilityCreateDto true "Availability parameters"
// @Security BearerAuth
// @Success 200 {object} model.Availability
// @Failure 400 {object} helpers.ApiError
// @Router /v1/availability [post]
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

	availability, err := ctl.availabilityService.Create(&data, user)

	helpers.HandleJSONResponse(c, availability, err)
}
