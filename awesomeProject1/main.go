/*Question-01*/
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin_rest_postgresql/internal/handlers"
	"github.com/yourusername/gin_rest_postgresql/internal/database"
	"log"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Initialize database connection pool
	pool, err := database.NewDBPool("postgresql://username:password@localhost:5432/otp")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Handlers
	userHandler := handlers.NewUserHandler(pool)
	router.POST("/api/users", userHandler.CreateUser)
	router.POST("/api/users/generateotp", userHandler.GenerateOTP)
	router.POST("/api/users/verifyotp", userHandler.VerifyOTP)

	// Start the server
	router.Run(":8080")
}
