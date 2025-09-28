package guard

import (
	"app/config"
	"errors"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UnsignedResponse struct {
	Message any `json:"message"`
}

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
		err = errors.New("token invalid")
		return claims, err
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		err = errors.New("token expired")
		return claims, err
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

			c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
				Message: err.Error(),
			})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}

		claims.Jwt = jwt

		c.Set("user", claims)

		c.Next()
	}
}
