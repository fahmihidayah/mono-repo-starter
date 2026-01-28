-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    picture VARCHAR(500),
    issuer_name VARCHAR(100),
    scope TEXT,
    sub VARCHAR(255),
    access_token TEXT,
    passkey_credential_id VARCHAR(255),
    passkey_public_key TEXT,
    passkey_counter INTEGER DEFAULT 0,
    passkey_transports VARCHAR(255),
    passkey_device_type VARCHAR(50),
    passkey_backed_up BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_issuer_name ON accounts(issuer_name);
