-- +goose Up

-- ── Products ──────────────────────────────────────────────────────────────────
INSERT INTO products (sku, name, price, inventory_qty) VALUES
    ('120P90', 'Google Home',    49.99,   10),
    ('43N23P', 'MacBook Pro',    5399.99,  5),
    ('A304SD', 'Alexa Speaker',  109.50,  10),
    ('234234', 'Raspberry Pi B',  30.00,   2)
ON CONFLICT (sku) DO NOTHING;

-- ── Users ─────────────────────────────────────────────────────────────────────
-- Passwords match username (bcrypt cost 12, zero salt)
INSERT INTO users (username, password_hash, password_salt) VALUES
    ('admin',  '$2a$12$6h1.RXCXtvmP.cjBdQfBZOJvKBxe2dZPwxIKJHvFdnKWJGT6cfn2W', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo1',  '$2a$12$AZm.sJU2rSvSfxjs.QiS0uWpQvmd1A.03sRYp/h0J9OWB13HzteuS', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo2',  '$2a$12$Fb8ipre94YraFBlne51e7OpJxedAwGhmHJ7q4MFwH9d7izAeVH1fO', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo3',  '$2a$12$VIglvK2NcFJyRKyDAOTC2OcYboLZVCDHRqq/HreHSHGoHb0.iHglC', '0000000000000000000000000000000000000000000000000000000000000000'),
    ('demo4',  '$2a$12$idMlFvRkCNiT1ZIOPQjMO.g6Vt7.foBl3nLR9gHRjqUuHkkG380zy', '0000000000000000000000000000000000000000000000000000000000000000')
ON CONFLICT (username) DO NOTHING;

INSERT INTO users (username, password_hash, password_salt, role) VALUES (
    'seller1',
    '$2a$12$pryu/2vk1sJPT1bepJfbreG8uQyHhORFHUCsgCOYXAZzJfoH6WJyy',
    '0000000000000000000000000000000000000000000000000000000000000000',
    'seller'
) ON CONFLICT (username) DO NOTHING;

-- ── Campaigns ─────────────────────────────────────────────────────────────────
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
DELETE FROM campaigns;
DELETE FROM products WHERE sku IN ('120P90', '43N23P', 'A304SD', '234234');
DELETE FROM users WHERE username IN ('admin', 'demo1', 'demo2', 'demo3', 'demo4', 'seller1');
