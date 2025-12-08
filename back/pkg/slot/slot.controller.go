package slot

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SlotController struct {
	slotService *SlotService
}

func NewSlotController(ctl *SlotController) *SlotController {
	if ctl != nil {
		return ctl
	}

	return &SlotController{
		slotService: NewSlotService(nil),
	}
}

// @Summary Confirm a slot
// @Tags Slot
// @Param slotId path string true "Slot Id"
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body ConfirmSlotDto true "Confirm Slot parameters"
// @Success 200 {object} model.Slot
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/slots/{slotId}/confirm [post]
func (ctl *SlotController) ConfirmSlot(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	var data ConfirmSlotDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	slotId, err := uuid.Parse(c.Param("slotId"))
	if err != nil {
		helpers.HandleJSONResponse(c, nil, constants.ERR_SLOT_NOT_FOUND.Err)
		return
	}

	slots, err := ctl.slotService.ConfirmSlot(data, slotId, user.Id)
	helpers.HandleJSONResponse(c, slots, err)
}
