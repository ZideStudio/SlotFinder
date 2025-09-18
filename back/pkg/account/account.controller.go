package account

import (
	"app/commons/guard"
	"app/commons/helpers"
	"errors"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService *AccountService
}

func NewAccountController(ctl *AccountController) *AccountController {
	if ctl != nil {
		return ctl
	}

	return &AccountController{
		accountService: NewAccountService(nil),
	}
}

// @Summary Create an account
// @Description Create a new account with the provided parameters.
// @Tags Account
// @Accept json
// @Produce json
// @Param data body AccountCreateDto true "Account parameters"
// @Success 200 {object} AccountCreateResponseDto
// @Failure 400 {object} helpers.ApiError
// @Router /v1/account [post]
func (ctl *AccountController) Create(c *gin.Context) {
	var data AccountCreateDto
	if err := helpers.ShouldBindJSON(c, &data); err != nil {
		return
	}

	account, err := ctl.accountService.Create(&data)

	helpers.ResponseJSON(c, account, err)
}

// @Summary Get My Account
// @Description Get the account information of the current user.
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} model.Account
// @Failure 400 {object} helpers.ApiError
// @Router /v1/account/me [get]
// @Security BearerAuth
func (ctl *AccountController) GetMe(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.ResponseJSON(c, nil, err)
		return
	}
	if user == nil {
		helpers.ResponseJSON(c, nil, errors.New("user not found"))
		return
	}

	account, err := ctl.accountService.GetMe(user.Id)

	helpers.ResponseJSON(c, account, err)
}

// @Summary Update my account
// @Description Update own account
// @Tags Account
// @Accept json
// @Produce json
// @Param data body AccountUpdateDto true "Account parameters"
// @Success 200 {object} model.Account
// @Failure 400 {object} helpers.ApiError
// @Router /v1/account [patch]
// @Security BearerAuth
func (ctl *AccountController) Update(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.ResponseJSON(c, nil, err)
		return
	}
	if user == nil {
		helpers.ResponseJSON(c, nil, errors.New("user not found"))
		return
	}

	var data AccountUpdateDto
	if err := helpers.ShouldBindJSON(c, &data); err != nil {
		return
	}

	account, err := ctl.accountService.Update(&data, user.Id)

	helpers.ResponseJSON(c, account, err)
}
