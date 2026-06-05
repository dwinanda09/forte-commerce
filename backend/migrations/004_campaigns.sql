-- +goose Up
CREATE TABLE campaigns (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    priority    INT          NOT NULL DEFAULT 0,
    conditions  JSONB        NOT NULL DEFAULT '[]',
    actions     JSONB        NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_campaigns_active_priority ON campaigns(is_active, priority);

INSERT INTO campaigns (name, description, conditions, actions, priority) VALUES (
    'MacBook Pro Free Raspberry Pi',
    'Buy a MacBook Pro and get a Raspberry Pi B for free',
    '[{"type":"cart_has_sku","sku":"43N23P","min_qty":1},{"type":"cart_has_sku","sku":"234234","min_qty":1}]',
    '[{"type":"free_item","sku":"234234","trigger_sku":"43N23P"}]',
    10
);

INSERT INTO campaigns (name, description, conditions, actions, priority) VALUES (
    'Google Home Bundle (3 for 2)',
    'Buy 3 Google Home units, pay for only 2',
    '[{"type":"item_qty_gte","sku":"120P90","qty":3}]',
    '[{"type":"buy_n_get_m","sku":"120P90","buy_n":3,"pay_m":2}]',
    20
);

INSERT INTO campaigns (name, description, conditions, actions, priority) VALUES (
    'Alexa Speaker 10% Discount',
    'Buy 3 or more Alexa Speakers and get 10% off all Alexa items',
    '[{"type":"item_qty_gte","sku":"A304SD","qty":3}]',
    '[{"type":"pct_discount_on_sku","sku":"A304SD","pct":10}]',
    30
);

-- +goose Down
DROP TABLE campaigns;
