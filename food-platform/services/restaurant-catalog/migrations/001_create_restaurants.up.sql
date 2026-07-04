-- +migrate Up
-- Create restaurants table with PostGIS for geospatial queries
CREATE TABLE IF NOT EXISTS restaurants (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             VARCHAR(255) NOT NULL,
    slug             VARCHAR(255) NOT NULL UNIQUE,
    cuisine_types    VARCHAR(100)[] NOT NULL DEFAULT '{}',
    rating           DECIMAL(2,1) NOT NULL DEFAULT 0,
    rating_count     INT NOT NULL DEFAULT 0,
    logo_url         VARCHAR(500),
    cover_url        VARCHAR(500),
    location         GEOGRAPHY(POINT, 4326) NOT NULL,
    address          TEXT NOT NULL,
    city             VARCHAR(100) NOT NULL DEFAULT 'Cairo',
    is_open          BOOLEAN NOT NULL DEFAULT true,
    status           VARCHAR(30) NOT NULL DEFAULT 'pending_verification',
    eta_min_minutes  INT NOT NULL DEFAULT 20,
    eta_max_minutes  INT NOT NULL DEFAULT 40,
    delivery_fee     DECIMAL(10,2) NOT NULL DEFAULT 20.00,
    price_range      SMALLINT NOT NULL DEFAULT 2 CHECK (price_range >= 1 AND price_range <= 4),
    commission_rate  DECIMAL(4,3) NOT NULL DEFAULT 0.150,
    opens_at         VARCHAR(5) NOT NULL DEFAULT '10:00',
    closes_at        VARCHAR(5) NOT NULL DEFAULT '23:59',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Geospatial index for proximity queries
CREATE INDEX idx_restaurants_location ON restaurants USING GIST (location);

-- Other indexes
CREATE INDEX idx_restaurants_slug ON restaurants (slug);
CREATE INDEX idx_restaurants_status ON restaurants (status);
CREATE INDEX idx_restaurants_is_open ON restaurants (is_open) WHERE is_open = true;
CREATE INDEX idx_restaurants_cuisine ON restaurants USING GIN (cuisine_types);
CREATE INDEX idx_restaurants_rating ON restaurants (rating DESC, rating_count DESC);
CREATE INDEX idx_restaurants_name ON restaurants USING GIN (to_tsvector('arabic', name));

COMMENT ON TABLE restaurants IS 'Restaurant catalog with geospatial data';
COMMENT ON COLUMN restaurants.location IS 'Geographic coordinates (PostGIS GEOGRAPHY type)';
COMMENT ON COLUMN restaurants.cuisine_types IS 'Array of cuisine types (egyptian, italian, etc.)';

-- +migrate Down
DROP TABLE IF EXISTS restaurants;
