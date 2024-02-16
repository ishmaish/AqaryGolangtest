package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func NewDBPool(connString string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &DB{pool}, nil
}

// CreateUser creates a new user in the database.
func (db *DB) CreateUser(ctx context.Context, name, phoneNumber string) (int, error) {
    var userID int
    err := db.QueryRow(ctx, "INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING id", name, phoneNumber).Scan(&userID)
    if err != nil {
        return 0, err
    }
    return userID, nil
}

// FindUserByPhoneNumber retrieves a user from the database by phone number.
func (db *DB) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
    var user User
    err := db.QueryRow(ctx, "SELECT * FROM users WHERE phone_number = $1", phoneNumber).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.OTP, &user.OTPExpirationTime)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil // User not found
        }
        return nil, err
    }
    return &user, nil
}

// UpdateOTP updates the OTP and its expiration time for a user.
func (db *DB) UpdateOTP(ctx context.Context, otp string, expirationTime time.Time, userID int) error {
    _, err := db.Exec(ctx, "UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE id = $3", otp, expirationTime, userID)
    if err != nil {
        return err
    }
    return nil
}

// VerifyOTP verifies if the provided OTP is valid for a user.
func (db *DB) VerifyOTP(ctx context.Context, phoneNumber, otp string) (bool, error) {
    var exists bool
    err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1 AND otp = $2 AND otp_expiration_time > NOW())", phoneNumber, otp).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}
