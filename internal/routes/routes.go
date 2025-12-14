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

	fileRepo := repositories.NewFileRepository(db)
	fileService := services.NewFileService(fileRepo, "./uploads")
	fileHandler := handlers.NewFileHandler(fileService)

	// Category module
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

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

			// Category endpoints
			protected.POST("/categories", categoryHandler.Create)
			protected.GET("/categories", categoryHandler.List)
			protected.GET("/categories/:id", categoryHandler.Detail)
			protected.PUT("/categories/:id", categoryHandler.Update)
			protected.DELETE("/categories/:id", categoryHandler.Delete)

			protected.POST("/files", fileHandler.Upload)
			protected.GET("/files", fileHandler.ListMyFiles)
			protected.GET("/files/:id", fileHandler.GetByID)
			protected.DELETE("/files/:id", fileHandler.Delete)
		}
	}
}
