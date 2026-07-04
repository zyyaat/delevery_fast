-- +migrate Up
-- Create orders table (partitioned by month for scalability)
CREATE TABLE IF NOT EXISTS orders (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number     VARCHAR(10) NOT NULL UNIQUE,
    customer_id      UUID NOT NULL,
    restaurant_id    UUID NOT NULL,
    driver_id        UUID,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending',
    subtotal         DECIMAL(10,2) NOT NULL,
    delivery_fee     DECIMAL(10,2) NOT NULL DEFAULT 0,
    service_fee      DECIMAL(10,2) NOT NULL DEFAULT 0,
    vat              DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount         DECIMAL(10,2) NOT NULL DEFAULT 0,
    total            DECIMAL(10,2) NOT NULL,
    payment_method   VARCHAR(20) NOT NULL,
    payment_status   VARCHAR(20) NOT NULL DEFAULT 'pending',
    delivery_address TEXT NOT NULL,
    latitude         DECIMAL(10,7) NOT NULL,
    longitude        DECIMAL(10,7) NOT NULL,
    eta_minutes      INT NOT NULL DEFAULT 35,
    scheduled_for    TIMESTAMPTZ,
    prep_started_at  TIMESTAMPTZ,
    picked_up_at     TIMESTAMPTZ,
    delivered_at     TIMESTAMPTZ,
    cancel_reason    TEXT,
    notes            TEXT,
    cashback_earned  DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'picked_up', 'delivered', 'cancelled', 'refunded')),
    CHECK (payment_method IN ('vodafone_cash', 'instapay', 'card', 'cod')),
    CHECK (payment_status IN ('pending', 'captured', 'failed', 'refunded', 'partial_refunded')),
    CHECK (total >= 0)
);

-- Indexes
CREATE INDEX idx_orders_customer ON orders (customer_id, created_at DESC);
CREATE INDEX idx_orders_restaurant ON orders (restaurant_id, status);
CREATE INDEX idx_orders_driver ON orders (driver_id, status);
CREATE INDEX idx_orders_status ON orders (status, created_at);
CREATE INDEX idx_orders_active_customer ON orders (customer_id) WHERE status IN ('pending', 'confirmed', 'preparing', 'ready', 'picked_up');
CREATE INDEX idx_orders_active_restaurant ON orders (restaurant_id) WHERE status IN ('pending', 'confirmed', 'preparing', 'ready');
CREATE INDEX idx_orders_created_at ON orders (created_at DESC);

COMMENT ON TABLE orders IS 'Customer food orders (partitioned by month in production)';
COMMENT ON COLUMN orders.order_number IS 'Short human-readable order number (e.g., A7X92F)';

-- Create order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id      UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id  UUID NOT NULL,
    name          VARCHAR(255) NOT NULL,
    quantity      INT NOT NULL CHECK (quantity > 0 AND quantity <= 50),
    unit_price    DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    modifiers     JSONB,
    notes         TEXT,
    line_total    DECIMAL(10,2) NOT NULL
);

CREATE INDEX idx_order_items_order ON order_items (order_id);
CREATE INDEX idx_order_items_menu_item ON order_items (menu_item_id);

COMMENT ON TABLE order_items IS 'Individual items within an order';

-- +migrate Down
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
