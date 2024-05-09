CREATE TABLE IF NOT EXISTS Users (
                       User_id SERIAL PRIMARY KEY,
                       Social_club_id VARCHAR(255) NOT NULL,
                       Username VARCHAR(255) NOT NULL,
                       Password_Hash VARCHAR(255) NOT NULL,
                       Email VARCHAR(255) UNIQUE NOT NULL,
                       Created_At TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       Last_Login TIMESTAMP
);