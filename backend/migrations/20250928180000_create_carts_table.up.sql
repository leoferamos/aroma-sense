-- Migration to create the carts table
CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT fk_carts_user_id FOREIGN KEY (user_id) REFERENCES users(public_id) ON DELETE CASCADE,
    CONSTRAINT unique_user_cart UNIQUE(user_id)
);

-- Index for faster lookups by user_id
CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);
