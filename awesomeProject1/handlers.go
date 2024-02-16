package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin_rest_postgresql/internal/database"
)

type UserHandler struct {
	DB *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func createUser(db *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user User
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Check if phone number already exists
        existingUser, err := db.FindUserByPhoneNumber(context.Background(), user.PhoneNumber)
        if err != nil && err != sql.ErrNoRows {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if existingUser != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
            return
        }

        // Create new user
        id, err := db.CreateUser(context.Background(), user.Name, user.PhoneNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"id": id})
    }
}

func generateOTP(db *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req OTPRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Find user by phone number
        user, err := db.FindUserByPhoneNumber(context.Background(), req.PhoneNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if user == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        // Generate OTP
        otp := generateRandomOTP()
        expirationTime := time.Now().Add(time.Minute)

        // Update user with OTP
        _, err = db.UpdateOTP(context.Background(), otp, expirationTime, user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"otp": otp})
    }
}

func verifyOTP(db *pgxpool.Pool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req OTPRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Check if OTP is correct and not expired
        exists, err := db.VerifyOTP(context.Background(), req.PhoneNumber, req.OTP)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if !exists {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
    }
}

// Helper function to generate random 4-digit OTP
func generateRandomOTP() string {
    rand.Seed(time.Now().UnixNano())
    return strconv.Itoa(rand.Intn(10000))
}
