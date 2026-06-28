-- Users table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    role_id INT NOT NULL DEFAULT 3,
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    verification_token_expires TIMESTAMP,
    phone_otp VARCHAR(6),
    phone_otp_expires TIMESTAMP,
    otp_token VARCHAR(255),
    otp_expires TIMESTAMP,
    otp_resend_count INT DEFAULT 0,
    otp_resend_last_time TIMESTAMP,
    account_locked BOOLEAN DEFAULT FALSE,
    lock_reason VARCHAR(255),
    locked_until TIMESTAMP,
    failed_login_attempts INT DEFAULT 0,
    last_failed_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- Detailed audit logs with device information
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT,
    action VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    device_info VARCHAR(255),
    browser_info VARCHAR(255),
    location VARCHAR(255),
    status VARCHAR(50),
    details TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Login history table
CREATE TABLE login_history (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    logout_time TIMESTAMP,
    ip_address VARCHAR(45),
    device_info VARCHAR(255),
    browser_info VARCHAR(255),
    location VARCHAR(255),
    status VARCHAR(50) DEFAULT 'success',
    failure_reason VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Password change logs
CREATE TABLE password_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    device_info VARCHAR(255),
    status VARCHAR(50) DEFAULT 'success',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Email/Phone change logs
CREATE TABLE contact_change_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    change_type VARCHAR(50),
    old_value VARCHAR(255),
    new_value VARCHAR(255),
    ip_address VARCHAR(45),
    status VARCHAR(50) DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Profile update logs
CREATE TABLE profile_update_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    field_name VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    ip_address VARCHAR(45),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Insert default roles
INSERT INTO roles (id, name) VALUES 
(1, 'Admin'),
(2, 'Moderator'),
(3, 'User');

-- Insert default admin account (password: admin123)
-- Using bcrypt hash of "admin123"
INSERT INTO users (email, password, role_id, email_verified, phone_verified, account_locked, created_at) VALUES 
('cliffe026@gmail.com', '$2a$12$eiSeagkNAtETvXpjmoTW.eaVryMGqvOlJ3XfY1Sh.RbnYcentADR.', 1, true, true, false, CURRENT_TIMESTAMP);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_locked ON users(account_locked);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_login_history_user ON login_history(user_id);
CREATE INDEX idx_login_history_time ON login_history(login_time);
CREATE INDEX idx_password_logs_user ON password_logs(user_id);
CREATE INDEX idx_contact_change_user ON contact_change_logs(user_id);
CREATE INDEX idx_profile_update_user ON profile_update_logs(user_id);
