-- Roll back review_reports and reports_count

ALTER TABLE reviews
    DROP COLUMN IF EXISTS reports_count;

DROP INDEX IF EXISTS idx_review_reports_status_created;
DROP INDEX IF EXISTS idx_review_reports_unique_reporter;
DROP TABLE IF EXISTS review_reports;
