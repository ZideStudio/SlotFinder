package account

import (
	"app/commons/constants"
	"app/commons/guard"
	"app/commons/helpers"
	"app/commons/lib"
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"slices"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService *AccountService
	avatarService  *AvatarService
}

func NewAccountController(ctl *AccountController) *AccountController {
	if ctl != nil {
		return ctl
	}

	return &AccountController{
		accountService: NewAccountService(nil),
		avatarService:  NewAvatarService(nil),
	}
}

// @Summary Create an account
// @Description Create a new account with the provided parameters.
// @Tags Account
// @Accept json
// @Produce json
// @Param data body AccountCreateDto true "Account parameters"
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/account [post]
func (ctl *AccountController) Create(c *gin.Context) {
	var data AccountCreateDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	accessToken, err := ctl.accountService.Create(&data)
	if err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}

	lib.SetAccessTokenCookie(c, accessToken, 0)

	helpers.HandleJSONResponse(c, nil, err)
}

// @Summary Get My Account
// @Description Get the account information of the current user.
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} model.Account
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/account/me [get]
// @security AccessTokenCookie
func (ctl *AccountController) GetMe(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}
	if user == nil {
		helpers.HandleJSONResponse(c, nil, errors.New("user not found"))
		return
	}

	account, err := ctl.accountService.GetMe(user.Id)

	helpers.HandleJSONResponse(c, account, err)
}

// @Summary Update my account
// @Description Update own account
// @Tags Account
// @Accept json
// @Produce json
// @Param data body AccountUpdateDto true "Account parameters"
// @Success 200 {object} model.Account
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/account [patch]
// @security AccessTokenCookie
func (ctl *AccountController) Update(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}
	if user == nil {
		helpers.HandleJSONResponse(c, nil, errors.New("user not found"))
		return
	}

	var data AccountUpdateDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	account, accessToken, err := ctl.accountService.Update(&data, user.Id)
	if accessToken != nil {
		lib.SetAccessTokenCookie(c, *accessToken, 0)
	}

	helpers.HandleJSONResponse(c, account, err)
}

// @Summary Upload Avatar
// @Description UploadAvatar the avatar image of the current user
// @Tags Account
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Avatar image file"
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Failure 500 {object} helpers.ApiError
// @Router /api/v1/account/avatar [patch]
// @security AccessTokenCookie
func (ctl *AccountController) UploadAvatar(c *gin.Context) {
	var user *guard.Claims
	if err := guard.GetUserClaims(c, &user); err != nil {
		helpers.HandleJSONResponse(c, nil, err)
		return
	}
	if user == nil {
		helpers.HandleJSONResponse(c, nil, errors.New("user not found"))
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "missing image"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to open image"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read image"})
		return
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid image file"})
		return
	}
	if slices.Contains(constants.ALLOWED_PICTURE_FORMATS, constants.PictureFormat(format)) == false {
		c.JSON(400, gin.H{"error": "unsupported format"})
		return
	}

	err = ctl.avatarService.UploadUserAvatar(imageBytes, user.Id)

	helpers.HandleJSONResponse(c, nil, err)
}

// @Summary Forgot Password
// @Description Send a password reset email to the user
// @Tags Account
// @Accept json
// @Produce json
// @Param data body ForgotPasswordDto true "Email for password reset"
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/account/forgot-password [post]
func (ctl *AccountController) ForgotPassword(c *gin.Context) {
	var data ForgotPasswordDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	err := ctl.accountService.ForgotPassword(&data)
	helpers.HandleJSONResponse(c, nil, err)
}

// @Summary Reset Password
// @Description Reset password using reset token
// @Tags Account
// @Accept json
// @Produce json
// @Param data body ResetPasswordDto true "Reset token and new password"
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /api/v1/account/reset-password [post]
func (ctl *AccountController) ResetPassword(c *gin.Context) {
	var data ResetPasswordDto
	if err := helpers.SetHttpContextBody(c, &data); err != nil {
		return
	}

	err := ctl.accountService.ResetPassword(&data)
	helpers.HandleJSONResponse(c, nil, err)
}
