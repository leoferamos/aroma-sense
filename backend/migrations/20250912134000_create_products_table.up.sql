-- Migration to create the products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    brand VARCHAR(64) NOT NULL,
    weight FLOAT NOT NULL,
    description TEXT,
    price FLOAT NOT NULL,
    image_url VARCHAR(256),
    category VARCHAR(64) NOT NULL,
    notes TEXT,
    stock_quantity INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
