-- =============================================================================
-- reset_test_data.sql
-- Resets the database to the exact state described in the PDF spec so the
-- three basic test scenarios can be re-run cleanly.
--
-- Usage:
--   make seed-reset
--   # or directly:
--   psql "$DB_URL" -f backend/scripts/reset_test_data.sql
-- =============================================================================

-- ── 1. Clear transactional data (orders → sessions to honour FK order) ───────
DELETE FROM orders;
DELETE FROM checkout_sessions;

-- ── 2. Reset product inventory to PDF-spec values ────────────────────────────
--
--   SKU     | Name            | Price     | Inventory Qty
--   --------|-----------------|-----------|---------------
--   120P90  | Google Home     | $49.99    | 10
--   43N23P  | MacBook Pro     | $5,399.99 | 5
--   A304SD  | Alexa Speaker   | $109.50   | 10
--   234234  | Raspberry Pi B  | $30.00    | 2
--
UPDATE products SET inventory_qty = 10, reserved_qty = 0 WHERE sku = '120P90';
UPDATE products SET inventory_qty = 5,  reserved_qty = 0 WHERE sku = '43N23P';
UPDATE products SET inventory_qty = 10, reserved_qty = 0 WHERE sku = 'A304SD';
UPDATE products SET inventory_qty = 2,  reserved_qty = 0 WHERE sku = '234234';

-- ── 3. Restore campaigns to PDF-spec values ──────────────────────────────────
--
-- UPSERT by name so re-running is idempotent even if UUIDs differ.
-- Deactivate any unknown campaigns first so stale rows don't interfere.
--
UPDATE campaigns SET is_active = false
WHERE name NOT IN (
    'MacBook Pro Free Raspberry Pi',
    'Google Home Bundle (3 for 2)',
    'Alexa Speaker 10% Discount'
);

INSERT INTO campaigns (name, description, conditions, actions, priority, is_active)
VALUES (
    'MacBook Pro Free Raspberry Pi',
    'Buy a MacBook Pro and get a Raspberry Pi B for free',
    '[{"type":"cart_has_sku","sku":"43N23P","min_qty":1},{"type":"cart_has_sku","sku":"234234","min_qty":1}]',
    '[{"type":"free_item","sku":"234234","trigger_sku":"43N23P"}]',
    10,
    true
)
ON CONFLICT (name) DO UPDATE SET
    description = EXCLUDED.description,
    conditions  = EXCLUDED.conditions,
    actions     = EXCLUDED.actions,
    priority    = EXCLUDED.priority,
    is_active   = EXCLUDED.is_active,
    updated_at  = NOW();

INSERT INTO campaigns (name, description, conditions, actions, priority, is_active)
VALUES (
    'Google Home Bundle (3 for 2)',
    'Buy 3 Google Home units, pay for only 2',
    '[{"type":"item_qty_gte","sku":"120P90","qty":3}]',
    '[{"type":"buy_n_get_m","sku":"120P90","buy_n":3,"pay_m":2}]',
    20,
    true
)
ON CONFLICT (name) DO UPDATE SET
    description = EXCLUDED.description,
    conditions  = EXCLUDED.conditions,
    actions     = EXCLUDED.actions,
    priority    = EXCLUDED.priority,
    is_active   = EXCLUDED.is_active,
    updated_at  = NOW();

INSERT INTO campaigns (name, description, conditions, actions, priority, is_active)
VALUES (
    'Alexa Speaker 10% Discount',
    'Buy 3 or more Alexa Speakers and get 10% off all Alexa items',
    '[{"type":"item_qty_gte","sku":"A304SD","qty":3}]',
    '[{"type":"pct_discount_on_sku","sku":"A304SD","pct":10}]',
    30,
    true
)
ON CONFLICT (name) DO UPDATE SET
    description = EXCLUDED.description,
    conditions  = EXCLUDED.conditions,
    actions     = EXCLUDED.actions,
    priority    = EXCLUDED.priority,
    is_active   = EXCLUDED.is_active,
    updated_at  = NOW();

-- ── 4. Test credentials ───────────────────────────────────────────────────────
--
--   Role    | Username  | Password
--   --------|-----------|----------
--   buyer   | demo1     | demo1
--   buyer   | demo2     | demo2
--   seller  | seller1   | seller1
--
-- ── 4. Expected test-case outputs ────────────────────────────────────────────
--
--   Scenario 1 — MacBook Pro + Raspberry Pi B
--     POST /api/v1/checkout  { "items": ["43N23P", "234234"] }
--     subtotal          : $5,429.99
--     total_discount    : $30.00   (Raspberry Pi B is free)
--     total             : $5,399.99
--     promotion         : MacBook Pro Free Raspberry Pi
--
--   Scenario 2 — 3× Google Home (3-for-2 bundle)
--     POST /api/v1/checkout  { "items": ["120P90","120P90","120P90"] }
--     subtotal          : $149.97
--     total_discount    : $49.99   (1 unit free)
--     total             : $99.98
--     promotion         : Google Home Bundle (3 for 2)
--
--   Scenario 3 — 3× Alexa Speaker (10% quantity discount)
--     POST /api/v1/checkout  { "items": ["A304SD","A304SD","A304SD"] }
--     subtotal          : $328.50
--     total_discount    : $32.85   (10% of $328.50)
--     total             : $295.65
--     promotion         : Alexa Speaker 10% Discount
-- =============================================================================
