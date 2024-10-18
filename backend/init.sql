CREATE TABLE IF NOT EXISTS users (
    id bytea PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    username VARCHAR(255),
    role VARCHAR(50)
);

-- Create an index on the email column for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Create an index on the username column for faster lookups
CREATE INDEX idx_users_username ON users(username);
