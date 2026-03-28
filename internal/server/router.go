package server

import (
	"database/sql"
	"net/http"
	"time"

	auth_delivery "backend/internal/auth/delivery/http"
	auth_repo "backend/internal/auth/repository"
	auth_service "backend/internal/auth/service"
	auth_usecase "backend/internal/auth/usecase"
	chat_delivery "backend/internal/chat/delivery/http"
	chat_repo "backend/internal/chat/repository"
	chat_usecase "backend/internal/chat/usecase"
	chat_ws "backend/internal/chat/ws"
	"backend/internal/config"
	driver_delivery "backend/internal/driver/delivery/http"
	driver_repo "backend/internal/driver/repository"
	driver_usecase "backend/internal/driver/usecase"
	"backend/internal/middleware"
	passenger_delivery "backend/internal/passenger/delivery/http"
	passenger_repo "backend/internal/passenger/repository"
	passenger_usecase "backend/internal/passenger/usecase"
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func NewRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(logger.GinLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityMiddleware())

	// Rate Limiting
	limiter := middleware.NewIPRateLimiter(rate.Every(time.Minute), 100)
	r.Use(middleware.RateLimitMiddleware(limiter))

	// WebSocket Hub
	hub := chat_ws.NewHub()

	// Repositories
	userRepo := auth_repo.NewSQLUserRepository(db)
	driverRepo := driver_repo.NewSQLDriverRepository(db)
	rideRepo := passenger_repo.NewSQLRideRepository(db)
	interactionRepo := passenger_repo.NewSQLInteractionRepository(db)
	chatRepo := chat_repo.NewSQLChatRepository(db)

	// JWT Service (Infrastructure)
	jwtService := auth_service.NewJWTService(cfg.Auth.JWTSecret)

	// Usecases
	// Auth
	signupUC := auth_usecase.NewSignupUsecase(userRepo, driverRepo, jwtService)
	loginUC := auth_usecase.NewLoginUsecase(userRepo, jwtService)

	// Passenger
	requestRideUC := passenger_usecase.NewRequestRideUsecase(rideRepo)
	getRideUC := passenger_usecase.NewGetRideUsecase(rideRepo)
	getPassengerHistoryUC := passenger_usecase.NewGetPassengerHistoryUsecase(rideRepo)
	addSavedPlaceUC := passenger_usecase.NewAddSavedPlaceUsecase(interactionRepo)
	getSavedPlacesUC := passenger_usecase.NewGetSavedPlacesUsecase(interactionRepo)

	// Driver
	acceptRideUC := driver_usecase.NewAcceptRideUsecase(rideRepo, driverRepo)
	updateRideStatusUC := driver_usecase.NewUpdateRideStatusUsecase(rideRepo, driverRepo)
	updateLocationUC := driver_usecase.NewUpdateLocationUsecase(driverRepo)
	getNearbyRidesUC := driver_usecase.NewGetNearbyRidesUsecase(rideRepo)

	// Chat
	sendMessageUC := chat_usecase.NewSendMessageUsecase(chatRepo, hub)
	getChatHistoryUC := chat_usecase.NewGetChatHistoryUsecase(chatRepo)

	// Controllers (Adapters)
	authC := auth_delivery.NewAuthController(signupUC, loginUC)
	passengerC := passenger_delivery.NewPassengerController(
		requestRideUC,
		getRideUC,
		getPassengerHistoryUC,
		addSavedPlaceUC,
		getSavedPlacesUC,
	)
	driverC := driver_delivery.NewDriverController(
		acceptRideUC,
		updateRideStatusUC,
		updateLocationUC,
		getNearbyRidesUC,
	)
	chatC := chat_delivery.NewChatController(sendMessageUC, getChatHistoryUC, hub)

	// Routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authC.Signup)
		authGroup.POST("/login", authC.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtService))
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// Passenger
		passengerGroup := api.Group("/passenger")
		{
			passengerGroup.POST("/ride/request", passengerC.Request)
			passengerGroup.GET("/ride/:id", passengerC.GetRide)
			passengerGroup.GET("/history", passengerC.GetHistory)
			passengerGroup.POST("/saved-places", passengerC.AddSavedPlace)
			passengerGroup.GET("/saved-places", passengerC.GetSavedPlaces)
		}

		// Driver
		driverGroup := api.Group("/driver")
		{
			driverGroup.POST("/ride/:id/accept", driverC.Accept)
			driverGroup.PATCH("/ride/:id/status", driverC.UpdateStatus)
			driverGroup.POST("/location", driverC.UpdateLocation)
			driverGroup.GET("/rides/nearby", driverC.GetNearby)
		}

		// Chat
		chatGroup := api.Group("/chat")
		{
			chatGroup.POST("/:id/send", chatC.SendMessage)
			chatGroup.GET("/:id/history", chatC.GetHistory)
			chatGroup.GET("/:id/ws", chatC.HandleWebSocket)
		}
	}

	return r
}
