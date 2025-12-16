package main

import (
	"fmt"
	"time"
	"vasvault/internal/repositories"
	"vasvault/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment")
	}
	db, err := repositories.Connect()
	if err != nil {
		panic("Failed to connect to database")
	}

	r := gin.Default()
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.InitRoutes(r, db.DB)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic("Failed to run server")
	}
}
