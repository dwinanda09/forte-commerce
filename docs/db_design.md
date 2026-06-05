# ForteCommerce Database Design

## Entity Relationship Diagram (ERD)

```
Users (1) ─── (∞) Orders
       │
       └─────── (∞) CartItems

Products (1) ─── (∞) OrderLines
       │
       └─────── (∞) CartItems
       │
       └─────── (∞) InventoryLevels

Orders (1) ─── (∞) OrderLines
      │
      └─────── (1) CheckoutSession

CheckoutSession (1) ─── (∞) CheckoutItems

InventoryLevels (1) ─── (∞) InventoryTransactions
```

### Core Tables

#### `users`
- `id` (UUID, PK)
- `email` (VARCHAR UNIQUE)
- `password_hash` (VARCHAR)
- `full_name` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `products`
- `id` (UUID, PK)
- `sku` (VARCHAR UNIQUE)
- `name` (VARCHAR)
- `description` (TEXT)
- `price` (DECIMAL)
- `active` (BOOLEAN)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `inventory_levels`
- `id` (UUID, PK)
- `product_id` (UUID, FK → products)
- `on_hand` (INTEGER) — physical inventory
- `reserved_qty` (INTEGER) — inventory held in active checkout sessions
- `available` (INTEGER) — computed as on_hand - reserved_qty
- `last_updated` (TIMESTAMP)
- Unique constraint: (product_id)

#### `checkout_sessions`
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `status` (ENUM: pending, reserved, payment_processing, completed, expired, abandoned)
- `total_amount` (DECIMAL)
- `expires_at` (TIMESTAMP)
- `payment_snapshot` (JSONB) — captures payment method & amount at checkout time
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `checkout_items`
- `id` (UUID, PK)
- `checkout_session_id` (UUID, FK → checkout_sessions)
- `product_id` (UUID, FK → products)
- `qty` (INTEGER)
- `price_snapshot` (DECIMAL) — price locked at checkout time
- `created_at` (TIMESTAMP)

#### `orders`
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `checkout_session_id` (UUID, FK → checkout_sessions)
- `status` (ENUM: pending, confirmed, shipped, delivered, cancelled)
- `total_amount` (DECIMAL)
- `order_snapshot` (JSONB) — user, shipping address, items at time of order creation
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

#### `order_lines`
- `id` (UUID, PK)
- `order_id` (UUID, FK → orders)
- `product_id` (UUID, FK → products)
- `qty` (INTEGER)
- `price_at_order` (DECIMAL) — price locked at order time
- `created_at` (TIMESTAMP)

#### `inventory_transactions`
- `id` (UUID, PK)
- `inventory_level_id` (UUID, FK → inventory_levels)
- `transaction_type` (ENUM: reserve, release, sell, restock, adjustment)
- `qty_delta` (INTEGER)
- `related_entity_id` (UUID) — checkout_session_id or order_id
- `notes` (TEXT)
- `created_at` (TIMESTAMP)

#### `cart_items`
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `product_id` (UUID, FK → products)
- `qty` (INTEGER)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

---

## Design Decisions

### 1. Reservation Pattern for Inventory Hold

**Problem:** When a user creates a checkout session with items, we need to guarantee that inventory is not oversold if multiple concurrent checkouts attempt to purchase the same product.

**Solution:** 
- `inventory_levels.reserved_qty` tracks quantities held by active checkout sessions
- `available = on_hand - reserved_qty` gives the true available quantity for sale
- On checkout creation, we atomically increment `reserved_qty` and check that `available >= requested`
- If checkout expires or is abandoned, we decrement `reserved_qty` (release)
- When an order is confirmed, we decrement both `reserved_qty` and `on_hand` (sell)

**Benefits:**
- Prevents overbooking under concurrent load
- Maintains visibility of what inventory is truly available
- Audit trail through `inventory_transactions`

### 2. Snapshot Pattern for Order Immutability

**Problem:** Orders must be immutable to ensure accuracy in fulfillment and dispute resolution. However, product prices and user details change over time.

**Solution:**
- `checkout_sessions.payment_snapshot` → JSON capturing payment method, amount, and timestamp
- `orders.order_snapshot` → JSON capturing user details, shipping address, and order metadata at order creation
- `order_lines.price_at_order` → price locked at order creation time
- `checkout_items.price_snapshot` → price locked at checkout session time

**Benefits:**
- Orders remain accurate even if products are deleted or prices change
- Complete audit trail of what customer agreed to pay
- No need for historical tables; history is captured inline
- Faster reads (no joins needed for historical context)

### 3. Checkout Session State Machine

States: `pending` → `reserved` → `payment_processing` → `completed` OR `abandoned`

- **pending**: Session created, items selected but not yet reserved
- **reserved**: Inventory reserved; user is completing payment
- **payment_processing**: Payment is being processed
- **completed**: Order confirmed, payment successful
- **abandoned**: Session expired or user cancelled (releases inventory reservation)
- **expired**: Session timed out; inventory released

Transitions are enforced in the application layer. No invalid state transitions allowed.

### 4. Order Status Machine

States: `pending` → `confirmed` → `shipped` → `delivered` OR `cancelled`

- **pending**: Order created, awaiting fulfillment
- **confirmed**: Order confirmed, ready to ship
- **shipped**: Package is in transit
- **delivered**: Order arrived at customer
- **cancelled**: Order cancelled (refund issued, inventory returned)

### 5. Index Strategy

```sql
-- Products
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_active ON products(active);

-- Inventory
CREATE INDEX idx_inventory_product ON inventory_levels(product_id);

-- Orders
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_created ON orders(created_at DESC);
CREATE INDEX idx_orders_status ON orders(status);

-- Checkout Sessions
CREATE INDEX idx_checkout_user ON checkout_sessions(user_id);
CREATE INDEX idx_checkout_status ON checkout_sessions(status);
CREATE INDEX idx_checkout_expires ON checkout_sessions(expires_at) WHERE status IN ('pending', 'reserved');

-- Order Lines
CREATE INDEX idx_orderlines_order ON order_lines(order_id);
CREATE INDEX idx_orderlines_product ON order_lines(product_id);

-- Cart Items
CREATE INDEX idx_cart_user ON cart_items(user_id);

-- Inventory Transactions
CREATE INDEX idx_invtx_level ON inventory_transactions(inventory_level_id);
CREATE INDEX idx_invtx_type ON inventory_transactions(transaction_type);
CREATE INDEX idx_invtx_created ON inventory_transactions(created_at DESC);
```

**Rationale:**
- User-based queries (orders, cart, checkout) benefit from (user_id) indexes
- Status-based filtering (pending orders, active checkouts) needs status indexes
- Time-based queries (recent orders, expired sessions) need (created_at / expires_at) indexes
- Partial index on `expires_at` reduces bloat for completed checkouts
- Foreign key indexes improve join performance

---

## Concurrency Considerations

1. **Inventory Reservation:**
   - Use database-level locking or transactions with `FOR UPDATE` on `inventory_levels` to ensure atomicity
   - Check `available >= requested_qty` within the same transaction

2. **Checkout Session:**
   - Session creation and inventory reservation happen in one transaction
   - Expiry is handled by a background job or TTL-based cleanup

3. **Order Confirmation:**
   - Atomic transition: reserve → confirmed + inventory deduction in single transaction
   - Prevents race conditions where multiple confirmations attempt to process the same session

---

## Migration and Constraints

- All foreign keys use `ON DELETE CASCADE` for cleanup (or `ON DELETE RESTRICT` for safety, application-handled)
- Created_at/updated_at are set by the application or database defaults
- Unique constraints on SKU, user email to prevent duplicates
- Check constraints on integer quantities (>= 0) to prevent negative inventory
