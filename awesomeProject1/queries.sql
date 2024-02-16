-- queries.sql
-- name: CreateUser :one
INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING id;

-- name: FindUserByPhoneNumber :one
SELECT * FROM users WHERE phone_number = $1;

-- name: UpdateOTP :exec
UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE id = $3;

-- name: VerifyOTP :one
SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1 AND otp = $2 AND otp_expiration_time > NOW());
