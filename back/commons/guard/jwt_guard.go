package guard

import (
	"app/commons/constants"
	"app/commons/helpers"
	"app/commons/lib"
	"app/config"
	"errors"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Id       uuid.UUID `json:"id"`
	Username *string   `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

func GetUserClaims(c *gin.Context, user **Claims) error {
	userClaims, ok := c.Keys["user"]
	if !ok || userClaims == nil {
		return nil
	}
	userParsed, ok := userClaims.(*Claims)
	if !ok {
		return errors.New("user claims type assertion failed")
	}

	*user = userParsed

	return nil
}

func ParseToken(jwtToken string) (*Claims, error) {
	claims := &Claims{}

	config := config.GetConfig()

	f, err := os.ReadFile(config.Auth.PublicPemPath)
	if err != nil {
		return claims, err
	}

	_, err = jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (any, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(f))
	})
	if err != nil {
		return claims, constants.ERR_TOKEN_INVALID.Err
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		return claims, constants.ERR_TOKEN_EXPIRED.Err
	}

	return claims, err
}

type AuthCheckParams struct {
	RequireAuthentication bool
	RequireUsername       bool
}

// Set params to nil to enable all checks
func AuthCheck(params *AuthCheckParams) gin.HandlerFunc {
	if params == nil {
		params = &AuthCheckParams{true, true}
	}

	return func(c *gin.Context) {
		jwt, err := c.Cookie("access_token")
		if err != nil {
			if !params.RequireAuthentication { // validate auth
				c.Next()
				return
			}
			if errors.Is(err, http.ErrNoCookie) {
				err = constants.ERR_NOT_AUTHENTICATED.Err
			}

			helpers.HandleJSONResponse(c, nil, err)
			return
		}

		claims, err := ParseToken(jwt)
		if err != nil {
			if err.Error() == "token expired" {
				lib.SetAccessTokenCookie(c, "", -1)
			}

			helpers.HandleJSONResponse(c, nil, err)
			return
		}

		if params.RequireUsername && claims.Username == nil {
			helpers.HandleJSONResponse(c, nil, constants.ERR_USERNAME_MISSING.Err)
			return
		}

		c.Set("user", claims)

		c.Next()
	}
}
