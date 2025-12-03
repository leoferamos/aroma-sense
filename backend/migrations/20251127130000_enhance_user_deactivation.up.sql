-- +migrate Up
ALTER TABLE users ADD COLUMN deactivation_reason VARCHAR(50);
ALTER TABLE users ADD COLUMN deactivation_notes TEXT;
ALTER TABLE users ADD COLUMN suspension_until TIMESTAMP;
ALTER TABLE users ADD COLUMN reactivation_requested BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN contestation_deadline TIMESTAMP;