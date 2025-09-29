package server

import (
	"app/commons/guard"
	"app/pkg/account"
	"app/pkg/auth"
	"app/pkg/event"
	"app/pkg/health"
	"app/pkg/provider"
	"app/pkg/signin"

	_ "app/docs"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(new(guard.CorsGuard).CorsCheck())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	healthRouter := new(health.HealthController)

	router.GET("/readyz", healthRouter.Ready)
	router.GET("/healthz", healthRouter.Status)

	v1 := router.Group("v1")
	{
		accountGroup := v1.Group("account")
		{
			accountRouter := account.NewAccountController(nil)

			accountGroup.POST("", accountRouter.Create)
			accountGroup.Use(guard.AuthCheck(true)).GET("/me", accountRouter.GetMe)
			accountGroup.Use(guard.AuthCheck(true)).PATCH("", accountRouter.Update)
		}

		authGroup := v1.Group("auth")
		{
			signinRouter := signin.NewSigninController(nil)
			providerRouter := provider.NewProviderController(nil)
			authRouter := auth.NewAuthController()

			authGroup.POST("signin", signinRouter.Signin)

			authGroup.Use(guard.AuthCheck(false)).GET("/providers/url", providerRouter.ProvidersUrl)
			authGroup.Use(guard.AuthCheck(false)).GET("/:provider/url", providerRouter.ProviderUrl)
			authGroup.GET("/:provider/callback", providerRouter.ProviderCallback)

			authGroup.Use(guard.AuthCheck(true)).GET("/status", authRouter.Status)

		}

		eventGroup := v1.Group("event")
		{
			eventRouter := event.NewEventController(nil)

			eventGroup.Use(guard.AuthCheck(true)).POST("", eventRouter.Create)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
