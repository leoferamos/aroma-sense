-- Create reviews table
-- Using UUID id and referencing users(public_id) and products(id)

CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id INTEGER NOT NULL,
    user_id UUID NOT NULL,
    rating INTEGER NOT NULL,
    comment TEXT,
    status VARCHAR(16) NOT NULL DEFAULT 'published',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_reviews_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_reviews_user_id FOREIGN KEY (user_id) REFERENCES users(public_id) ON DELETE CASCADE,
    CONSTRAINT chk_reviews_rating CHECK (rating >= 1 AND rating <= 5),
    CONSTRAINT chk_reviews_status CHECK (status IN ('published','hidden','flagged'))
);

-- Users can review a product only once
CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_unique_product_user ON reviews(product_id, user_id);

-- For fetching published reviews by product
CREATE INDEX IF NOT EXISTS idx_reviews_product_status ON reviews(product_id, status);
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC);
