-- +migrate Up
CREATE TABLE IF NOT EXISTS drivers (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID NOT NULL UNIQUE,
    name             VARCHAR(255) NOT NULL,
    phone            VARCHAR(20) NOT NULL,
    vehicle_type     VARCHAR(20) NOT NULL DEFAULT 'motorcycle',
    vehicle_plate    VARCHAR(20),
    license_number   VARCHAR(50),
    kyc_status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    status           VARCHAR(20) NOT NULL DEFAULT 'offline',
    tier             VARCHAR(20) NOT NULL DEFAULT 'standard',
    rating           DECIMAL(2,1) NOT NULL DEFAULT 0,
    rating_count     INT NOT NULL DEFAULT 0,
    acceptance_rate  DECIMAL(4,3) NOT NULL DEFAULT 0,
    completion_rate  DECIMAL(4,3) NOT NULL DEFAULT 0,
    trust_score      INT NOT NULL DEFAULT 50 CHECK (trust_score >= 0 AND trust_score <= 100),
    total_earnings   DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_deliveries INT NOT NULL DEFAULT 0,
    latitude         DECIMAL(10,7),
    longitude        DECIMAL(10,7),
    heading          DECIMAL(6,2) DEFAULT 0,
    speed            DECIMAL(6,2) DEFAULT 0,
    last_online_at   TIMESTAMPTZ,
    photo_url        VARCHAR(500),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CHECK (vehicle_type IN ('motorcycle', 'car', 'bicycle')),
    CHECK (kyc_status IN ('pending', 'verified', 'rejected')),
    CHECK (status IN ('offline', 'online', 'on_break', 'on_delivery', 'suspended')),
    CHECK (tier IN ('platinum', 'gold', 'silver', 'standard'))
);

CREATE INDEX idx_drivers_user_id ON drivers (user_id);
CREATE INDEX idx_drivers_status ON drivers (status);
CREATE INDEX idx_drivers_kyc ON drivers (kyc_status);
CREATE INDEX idx_drivers_tier ON drivers (tier);
CREATE INDEX idx_drivers_location ON drivers USING GIST (ST_SetSRID(ST_MakePoint(COALESCE(longitude, 0), COALESCE(latitude, 0)), 4326)) WHERE status = 'online' AND kyc_status = 'verified';

COMMENT ON TABLE drivers IS 'Delivery drivers with KYC, tier system, and real-time location';

-- +migrate Down
DROP TABLE IF EXISTS drivers;
