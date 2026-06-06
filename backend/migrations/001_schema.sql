-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(50)  UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    password_salt VARCHAR(64)  NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'buyer',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    sku           VARCHAR(50)   UNIQUE NOT NULL,
    name          VARCHAR(255)  NOT NULL,
    price         NUMERIC(10,2) NOT NULL,
    inventory_qty INT           NOT NULL DEFAULT 0,
    reserved_qty  INT           NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS checkout_sessions (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    items         JSONB       NOT NULL,
    result        JSONB,
    error_message TEXT,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id                  UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    checkout_session_id UUID          NOT NULL REFERENCES checkout_sessions(id),
    status              VARCHAR(20)   NOT NULL DEFAULT 'pending',
    items               JSONB         NOT NULL,
    promotions_applied  JSONB         NOT NULL DEFAULT '[]',
    subtotal            NUMERIC(10,2) NOT NULL,
    total_discount      NUMERIC(10,2) NOT NULL DEFAULT 0,
    total               NUMERIC(10,2) NOT NULL,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS campaigns (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL DEFAULT '',
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    priority    INT          NOT NULL DEFAULT 0,
    conditions  JSONB        NOT NULL DEFAULT '[]',
    actions     JSONB        NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_products_sku                    ON products(sku);
CREATE INDEX        IF NOT EXISTS idx_checkout_sessions_status_expires ON checkout_sessions(status, expires_at) WHERE status = 'pending';
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_checkout_session_id      ON orders(checkout_session_id);
CREATE INDEX        IF NOT EXISTS idx_campaigns_active_priority        ON campaigns(is_active, priority);

-- +goose Down
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS checkout_sessions;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
