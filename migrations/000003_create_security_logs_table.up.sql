CREATE TABLE security_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    description TEXT,
    ip_address VARCHAR(45),
    location TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_security_logs_user_id ON security_logs(user_id);
CREATE INDEX idx_security_logs_severity ON security_logs(severity);
CREATE INDEX idx_security_logs_created_at ON security_logs(created_at); 