-- +migrate Down
ALTER TABLE users DROP COLUMN deactivation_reason;
ALTER TABLE users DROP COLUMN deactivation_notes;
ALTER TABLE users DROP COLUMN suspension_until;
ALTER TABLE users DROP COLUMN reactivation_requested;
ALTER TABLE users DROP COLUMN contestation_deadline;