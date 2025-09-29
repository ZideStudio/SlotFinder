package guard

import (
	"app/commons/constants"
	"app/commons/helpers"
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
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Jwt      string    `json:"-"`
	jwt.RegisteredClaims
}

func GetUserClaims(c *gin.Context, user **Claims) error {
	userClaims, ok := c.Keys["user"]
	if !ok || userClaims == nil {
		return nil
	}
	userParsed, ok := userClaims.(*Claims)
	if !ok {
		return errors.New("User claims type assertion failed")
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

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return claims, constants.ERR_TOKEN_EXPIRED.Err
	}

	return claims, err
}

func AuthCheck(requireAuthentication bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := c.Cookie("access_token")
		if err != nil {
			if !requireAuthentication { // validate auth
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
				c.SetCookie(
					"access_token",            // name
					"",                        // value
					-1,                        // max age in seconds
					"/",                       // path
					config.GetConfig().Domain, // domain
					true,                      // secure
					true,                      // httpOnly
				)
			}

			helpers.HandleJSONResponse(c, nil, err)
			return
		}

		claims.Jwt = jwt

		c.Set("user", claims)

		c.Next()
	}
}
