-- Create review_reports table for user-submitted review reports
-- Includes unique reporter per review and moderation status

CREATE TABLE IF NOT EXISTS review_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL,
    reported_by UUID NOT NULL,
    reason_category VARCHAR(32) NOT NULL,
    reason_text VARCHAR(500),
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_review_reports_review FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    CONSTRAINT fk_review_reports_reported_by FOREIGN KEY (reported_by) REFERENCES users(public_id) ON DELETE CASCADE,
    CONSTRAINT chk_review_reports_category CHECK (reason_category IN ('offensive','spam','improper','other')),
    CONSTRAINT chk_review_reports_status CHECK (status IN ('pending','accepted','rejected'))
);

-- One report per user per review
CREATE UNIQUE INDEX IF NOT EXISTS idx_review_reports_unique_reporter ON review_reports(review_id, reported_by);

-- For admin triage
CREATE INDEX IF NOT EXISTS idx_review_reports_status_created ON review_reports(status, created_at DESC);

-- Track aggregated reports on reviews
ALTER TABLE reviews
    ADD COLUMN IF NOT EXISTS reports_count INTEGER NOT NULL DEFAULT 0;
