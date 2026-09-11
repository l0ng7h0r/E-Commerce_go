-- 006_create_orders_payments.sql
CREATE TABLE IF NOT EXISTS orders (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_amount   NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    status         VARCHAR(50) NOT NULL DEFAULT 'pending',
    logistic_branch VARCHAR(255),
    district       VARCHAR(50),
    logistic_company VARCHAR(50),
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id   UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    quantity   INT NOT NULL DEFAULT 1,
    price      NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id       UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount         NUMERIC(12, 2) NOT NULL,
    status         VARCHAR(50) NOT NULL DEFAULT 'completed',
    transaction_id VARCHAR(255),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
