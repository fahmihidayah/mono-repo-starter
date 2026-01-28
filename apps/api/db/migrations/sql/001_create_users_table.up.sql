-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    reset_password_token VARCHAR(255),
    reset_password_token_expiry BIGINT DEFAULT 0,
    hashed_password TEXT NOT NULL,
    hash_salt VARCHAR(255) NOT NULL,
    hash_iterations INTEGER DEFAULT 10000,
    login_attempts INTEGER DEFAULT 0,
    lock_until BIGINT DEFAULT 0,
    is_super_user BOOLEAN DEFAULT FALSE,
    verification_code VARCHAR(255),
    verification_hash VARCHAR(255),
    verification_token_expiry BIGINT DEFAULT 0,
    verification_kind VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
