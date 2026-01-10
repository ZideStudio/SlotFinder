package server

import (
	"app/commons/guard"
	"app/pkg/account"
	"app/pkg/auth"
	"app/pkg/availability"
	"app/pkg/event"
	"app/pkg/health"
	"app/pkg/provider"
	"app/pkg/signin"
	"app/pkg/slot"
	"app/pkg/sse"

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
		// Account routes
		accountGroup := v1.Group("/account")
		{
			accountRouter := account.NewAccountController(nil)

			accountGroup.POST("", accountRouter.Create)
			accountGroup.GET("/me", guard.AuthCheck(nil), accountRouter.GetMe)
			accountGroup.PATCH("", guard.AuthCheck(&guard.AuthCheckParams{RequireAuthentication: true, RequireValidatedAccount: false}), accountRouter.Update)
			accountGroup.PATCH("/avatar", guard.AuthCheck(nil), guard.MaxUploadSizeMiddleware(10<<20), accountRouter.UploadAvatar)
			accountGroup.POST("/forgot-password", accountRouter.ForgotPassword)
			accountGroup.POST("/reset-password", accountRouter.ResetPassword)
		}

		// Auth routes
		authGroup := v1.Group("/auth")
		{
			signinRouter := signin.NewSigninController(nil)
			providerRouter := provider.NewProviderController(nil)
			authRouter := auth.NewAuthController()

			authGroup.POST("/signin", signinRouter.Signin)

			authGroup.Use(guard.AuthCheck(&guard.AuthCheckParams{RequireAuthentication: false, RequireValidatedAccount: true})).GET("/:provider/url", providerRouter.ProviderUrl)
			authGroup.GET("/:provider/callback", providerRouter.ProviderCallback)

			authGroup.Use(guard.AuthCheck(nil)).GET("/status", authRouter.Status)
			authGroup.Use(guard.AuthCheck(nil)).POST("/logout", authRouter.Logout)

		}

		// Availability routes
		availabilityGroup := v1.Group("/availabilities")
		availabilityRouter := availability.NewAvailabilityController(nil)
		{
			availabilityGroup.DELETE("/:availabilityId", guard.AuthCheck(nil), availabilityRouter.Delete)
			availabilityGroup.PATCH("/:availabilityId", guard.AuthCheck(nil), availabilityRouter.Update)
		}

		// Event routes
		eventGroup := v1.Group("/events")
		{
			eventRouter := event.NewEventController(nil)

			eventGroup.GET("", guard.AuthCheck(nil), eventRouter.GetUserEvents)
			eventGroup.POST("", guard.AuthCheck(nil), eventRouter.Create)
			eventGroup.PATCH("/:eventId", guard.AuthCheck(nil), eventRouter.Update)

			specificEventGroup := eventGroup.Group("/:eventId")
			{
				specificEventGroup.GET("", guard.AuthCheck(&guard.AuthCheckParams{RequireAuthentication: false, RequireValidatedAccount: true}), eventRouter.GetEvent)
				specificEventGroup.POST("/join", guard.AuthCheck(nil), eventRouter.JoinEvent)
				specificEventGroup.PATCH("/profile", guard.AuthCheck(nil), eventRouter.UpdateProfile)
			}

			// Availability routes
			{
				eventGroup.POST("/:eventId/availability", guard.AuthCheck(nil), availabilityRouter.Create)
			}

			// SSE routes
			{
				sseRouter := sse.NewSSEController(nil)
				eventGroup.GET("/:eventId/sse", guard.AuthCheck(nil), sseRouter.Connect)
			}

		}

		// Slot routes
		slotGroup := v1.Group("/slots")
		{
			slotRouter := slot.NewSlotController(nil)

			slotGroup.POST("/:slotId/confirm", guard.AuthCheck(nil), slotRouter.ConfirmSlot)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
