package routes

import (
	"vasvault/internal/handlers"
	"vasvault/internal/middleware"
	"vasvault/internal/repositories"
	"vasvault/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// API v1 routes
	apiV1 := r.Group("/api/v1")
	{
		// Public routes
		apiV1.POST("/login", userHandler.Login)
		apiV1.POST("/register", userHandler.Register)
		apiV1.POST("/refresh", userHandler.Refresh)

		// Protected routes
		protected := apiV1.Group("")
		protected.Use(middleware.GinBearerAuth())
		{
			protected.GET("/me", userHandler.Me)
			protected.PUT("/profile", userHandler.UpdateProfile)
		}
	}
}
