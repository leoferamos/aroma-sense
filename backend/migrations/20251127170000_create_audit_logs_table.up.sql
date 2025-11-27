-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
    user_id BIGINT,
    actor_id BIGINT,
    actor_type VARCHAR(50) NOT NULL DEFAULT 'user',
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    details TEXT NOT NULL DEFAULT '{}',
    old_values TEXT,
    new_values TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    compliance VARCHAR(100),
    severity VARCHAR(20) NOT NULL DEFAULT 'info',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Indexes for performance
    INDEX idx_audit_logs_user_id (user_id),
    INDEX idx_audit_logs_actor_id (actor_id),
    INDEX idx_audit_logs_resource (resource),
    INDEX idx_audit_logs_resource_id (resource_id),
    INDEX idx_audit_logs_timestamp (timestamp),
    INDEX idx_audit_logs_action (action),

    -- Foreign keys
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create index for efficient queries
CREATE INDEX idx_audit_logs_composite ON audit_logs (user_id, action, timestamp DESC);
CREATE INDEX idx_audit_logs_compliance ON audit_logs (compliance, timestamp DESC);

-- Add comments for documentation
COMMENT ON TABLE audit_logs IS 'Audit logs for LGPD compliance and system monitoring';
COMMENT ON COLUMN audit_logs.public_id IS 'Public UUID for external references';
COMMENT ON COLUMN audit_logs.user_id IS 'ID of the user being acted upon';
COMMENT ON COLUMN audit_logs.actor_id IS 'ID of the user performing the action';
COMMENT ON COLUMN audit_logs.actor_type IS 'Type of actor: user, admin, system';
COMMENT ON COLUMN audit_logs.action IS 'Action performed: login, update, delete, etc.';
COMMENT ON COLUMN audit_logs.resource IS 'Resource type: user, order, product, etc.';
COMMENT ON COLUMN audit_logs.resource_id IS 'ID of the specific resource';
COMMENT ON COLUMN audit_logs.details IS 'JSON details of the action';
COMMENT ON COLUMN audit_logs.compliance IS 'Compliance framework: LGPD, GDPR, etc.';
COMMENT ON COLUMN audit_logs.severity IS 'Severity level: info, warning, error, critical';