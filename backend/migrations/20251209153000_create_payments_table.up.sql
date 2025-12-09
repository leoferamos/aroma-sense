-- Create payments table for gateway reconciliation
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    intent_id VARCHAR(255) NOT NULL UNIQUE,
    provider VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    order_public_id UUID NULL,
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    metadata JSONB,
    error_code VARCHAR(100),
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_payments_user_id FOREIGN KEY (user_id) REFERENCES users(public_id) ON DELETE CASCADE,
    CONSTRAINT fk_payments_order_public_id FOREIGN KEY (order_public_id) REFERENCES orders(public_id) ON DELETE SET NULL,
    CONSTRAINT check_payment_status CHECK (status IN ('pending','processing','succeeded','failed','canceled'))
);

CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_public_id ON payments(order_public_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
