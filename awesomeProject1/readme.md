SET up postgres and upload this code 

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    otp VARCHAR(4),
    otp_expiration_time TIMESTAMP
);
