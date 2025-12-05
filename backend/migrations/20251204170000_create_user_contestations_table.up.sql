-- Create user_contestations table for tracking account deactivation contestations
CREATE TABLE user_contestations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    requested_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    reviewed_by INTEGER REFERENCES users(id), -- admin id
    review_notes TEXT
);

CREATE INDEX idx_user_contestations_user_id ON user_contestations(user_id);
CREATE INDEX idx_user_contestations_status ON user_contestations(status);