/*Question-01*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

func main() {
	config := getConfig()
	initDB(config)

	defer db.Close()

	router := gin.Default()

	router.POST("/users", createUser)
	router.POST("/generate-otp", generateOTP)
	router.POST("/verify-otp", verifyOTP)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

func getConfig() Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func initDB(config Config) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	pool, err := pgxpool.Connect(nil, connStr)
	if err != nil {
		log.Fatal(err)
	}

	db = pool
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uuid.New()

	createdUser, err := createUserDB(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func createUserDB(user User) (User, error) {
	var createdUser User
	err := db.QueryRow(context.TODO(), createUserQuery, user.ID, user.Username, user.Phone, user.OTP).
		Scan(&createdUser.ID, &createdUser.Username, &createdUser.Phone, &createdUser.OTP, &createdUser.OTPExp)

	return createdUser, err
}

func generateOTP(c *gin.Context) {
	var phoneRequest struct {
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&phoneRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByPhone(phoneRequest.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	otpSecret := generateRandomOTP()
	updatedUser, err := updateOTP(user.ID, otpSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating OTP"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func verifyOTP(c *gin.Context) {
	var verificationRequest struct {
		UserID uuid.UUID `json:"userId"`
	}

	if err := c.ShouldBindJSON(&verificationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verifiedUser, err := verifyOTPDB(verificationRequest.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying OTP"})
		return
	}

	c.JSON(http.StatusOK, verifiedUser)
}

func verifyOTPDB(userID uuid.UUID) (User, error) {
	var verifiedUser User
	err := db.QueryRow(context.TODO(), verifyOTPQuery, userID).
		Scan(&verifiedUser.ID, &verifiedUser.Username, &verifiedUser.Phone, &verifiedUser.OTP, &verifiedUser.OTPExp)

	return verifiedUser, err
}

func updateOTP(userID uuid.UUID, otp string) (User, error) {
	var updatedUser User
	err := db.QueryRow(context.TODO(), updateOTPQuery, userID, otp).
		Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Phone, &updatedUser.OTP, &updatedUser.OTPExp)

	return updatedUser, err
}

func getUserByPhone(email string) (User, error) {
	var user User
	err := db.QueryRow(context.TODO(), getUserByPhoneQuery, email).
		Scan(&user.ID, &user.Username, &user.Phone, &user.OTP, &user.OTPExp)

	return user, err
}

// Add User struct definition
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Phone    string    `json:"phone"`
	OTP      string    `json:"otp"`
	OTPExp   bool      `json:"otpExpTime"`
}

const (
	createUserQuery      = "CreateUser"
	verifyOTPQuery       = "VerifyOTP"
	updateOTPQuery       = "UpdateOTP"
	getUserByPhoneQuery  = "GetUserByPhone"
	generateRandomOTPLen = 6
)

func generateRandomOTP() string {
	// Implement logic to generate a random OTP of a specific length (e.g., 6 digits)
	// You can use a library or implement it manually
	// For simplicity, a placeholder function is used here
	return "123456"
}
