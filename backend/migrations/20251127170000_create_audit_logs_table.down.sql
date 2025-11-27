-- +migrate Down
DROP TABLE IF EXISTS audit_logs;
DROP EXTENSION IF EXISTS "uuid-ossp";