CREATE TABLE IF NOT EXISTS coupons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coupon_code TEXT UNIQUE NOT NULL,  -- Unique code for the coupon
    discount_percentage DECIMAL(5,2),  -- Discount in percentage (e.g., 15% off)
    expiration_date TIMESTAMP,         -- Expiration timestamp
    is_active BOOLEAN DEFAULT TRUE,    -- Whether the coupon is active
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
