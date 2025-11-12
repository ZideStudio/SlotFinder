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

	v1 := router.Group("/v1")
	{
		accountGroup := v1.Group("/account")
		{
			accountRouter := account.NewAccountController(nil)

			accountGroup.POST("", accountRouter.Create)
			accountGroup.GET("/me", guard.AuthCheck(nil), accountRouter.GetMe)
			accountGroup.PATCH("", guard.AuthCheck(&guard.AuthCheckParams{RequireAuthentication: true, RequireUsername: false}), accountRouter.Update)
		}

		authGroup := v1.Group("/auth")
		{
			signinRouter := signin.NewSigninController(nil)
			providerRouter := provider.NewProviderController(nil)
			authRouter := auth.NewAuthController()

			authGroup.POST("/signin", signinRouter.Signin)

			authGroup.Use(guard.AuthCheck(&guard.AuthCheckParams{RequireAuthentication: false, RequireUsername: true})).GET("/:provider/url", providerRouter.ProviderUrl)
			authGroup.GET("/:provider/callback", providerRouter.ProviderCallback)

			authGroup.Use(guard.AuthCheck(nil)).GET("/status", authRouter.Status)
			authGroup.Use(guard.AuthCheck(nil)).POST("/logout", authRouter.Logout)

		}

		eventGroup := v1.Group("/event")
		{
			eventRouter := event.NewEventController(nil)

			eventGroup.GET("", guard.AuthCheck(nil), eventRouter.GetUserEvents)
			eventGroup.GET("/:id", guard.AuthCheck(&guard.AuthCheckParams{RequireAuthentication: false, RequireUsername: true}), eventRouter.GetEvent)
			eventGroup.POST("/:id/join", guard.AuthCheck(nil), eventRouter.JoinEvent)
			eventGroup.POST("", guard.AuthCheck(nil), eventRouter.Create)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
