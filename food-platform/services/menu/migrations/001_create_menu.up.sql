-- +migrate Up
-- Create menu_categories table
CREATE TABLE IF NOT EXISTS menu_categories (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    restaurant_id UUID NOT NULL,
    name          VARCHAR(255) NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_menu_categories_restaurant ON menu_categories (restaurant_id);

-- Create menu_items table
CREATE TABLE IF NOT EXISTS menu_items (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    restaurant_id     UUID NOT NULL,
    category_id       UUID NOT NULL REFERENCES menu_categories(id) ON DELETE CASCADE,
    name              VARCHAR(255) NOT NULL,
    description       TEXT,
    price             DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    image_url         VARCHAR(500),
    is_available      BOOLEAN NOT NULL DEFAULT true,
    prep_time_minutes INT NOT NULL DEFAULT 10,
    rating            DECIMAL(2,1) NOT NULL DEFAULT 0,
    rating_count      INT NOT NULL DEFAULT 0,
    is_most_ordered   BOOLEAN NOT NULL DEFAULT false,
    display_order     INT NOT NULL DEFAULT 0,
    modifiers         JSONB, -- Array of modifier objects
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_menu_items_restaurant ON menu_items (restaurant_id);
CREATE INDEX idx_menu_items_category ON menu_items (category_id);
CREATE INDEX idx_menu_items_available ON menu_items (restaurant_id) WHERE is_available = true;

COMMENT ON TABLE menu_items IS 'Food items in restaurant menus';
COMMENT ON COLUMN menu_items.modifiers IS 'JSON array of customization options (size, extras, etc.)';

-- +migrate Down
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS menu_categories;
