package server

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/pkg/firebase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(firebaseApp *firebase.App) *gin.Engine {
	r := gin.Default()

	// Handlers initialization
	userHandler := handlers.NewUserHandler(firebaseApp)

	// Public routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(firebaseApp))
	{
		api.GET("/profile", userHandler.GetProfile)
	}

	return r
}
