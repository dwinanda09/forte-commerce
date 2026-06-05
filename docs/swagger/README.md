# ForteCommerce API Documentation

This directory contains the OpenAPI 3.0.3 specification and interactive Swagger UI for the ForteCommerce backend API.

## Files

- **openapi.yaml** - Complete OpenAPI 3.0.3 specification for the ForteCommerce API
- **index.html** - Self-contained Swagger UI interface for exploring and testing the API

## Quick Start

### Option 1: View in Browser

Simply open `index.html` in a web browser:

```bash
open index.html  # macOS
xdg-open index.html  # Linux
start index.html  # Windows
```

The HTML file loads the Swagger UI from CDN and references the OpenAPI spec from `openapi.yaml`.

### Option 2: Use a Local Server

To avoid CORS issues, serve the files via HTTP:

```bash
# Python 3
python3 -m http.server 8000

# Python 2
python -m SimpleHTTPServer 8000

# Node.js (with http-server package)
npx http-server
```

Then navigate to `http://localhost:8000/` or `http://localhost:8000/index.html`

### Option 3: Validate the Spec

Use online validators or tools:

- [Swagger Editor](https://editor.swagger.io) - Paste contents of `openapi.yaml`
- [ReDoc](https://redoc.ly/docs) - View spec with ReDoc (higher quality rendering)
- [OpenAPI Validator](https://www.openapis.org/)

## API Overview

### Authentication

All endpoints except `/auth/login` require a JWT bearer token:

```
Authorization: Bearer <your_jwt_token>
```

Obtain a token via `/auth/login` with username and password.

### Base URL

- Development: `http://localhost:8080/api/v1`
- Production: `https://api.fortecommerce.com/api/v1`

### Main Features

#### Products (CRUD)
- `GET /products` - List all products
- `POST /seller/products` - Create product (seller only)
- `PUT /seller/products/{id}` - Update product (seller only)
- `DELETE /seller/products/{id}` - Delete product (seller only)

#### Checkout (Async)
- `POST /checkout` - Submit checkout (returns 202, processing async)
- `GET /checkout/{id}` - Get checkout status and results
- `POST /checkout/{id}/confirm` - Confirm checkout and create order

#### Orders
- `GET /orders` - List user's orders
- `GET /orders/{id}` - Get order details
- `POST /orders/{id}/pay` - Mark order as paid
- `POST /orders/{id}/cancel` - Cancel order

#### Campaigns (Promotions)
- `GET /campaigns` - List active campaigns
- `GET /seller/campaigns` - List seller's campaigns (seller only)
- `POST /seller/campaigns` - Create campaign (seller only)
- `PUT /seller/campaigns/{id}` - Update campaign (seller only)
- `DELETE /seller/campaigns/{id}` - Delete campaign (seller only)
- `PATCH /seller/campaigns/{id}/toggle` - Toggle campaign active status (seller only)

## Async Checkout Flow

ForteCommerce uses asynchronous checkout to handle complex promotion calculations:

1. **Submit** - `POST /checkout` with items → returns `checkout_id` (202 Accepted)
2. **Wait** - Backend processes promotions asynchronously
3. **Poll** - `GET /checkout/{id}` until `status` = `completed`
4. **Confirm** - `POST /checkout/{id}/confirm` → creates order

The checkout session includes:
- Calculated subtotal, discounts, and total
- Applied promotions with details
- Item list with prices

## Promotions & Campaigns

Sellers can create dynamic promotion campaigns with:

**Condition Types:**
- `cart_has_sku` - Cart must contain a specific SKU
- `item_qty_gte` - Item quantity ≥ threshold
- `cart_total_gte` - Cart total ≥ amount
- `cart_item_count_gte` - Number of items ≥ count

**Action Types:**
- `free_item` - Give a free item
- `buy_n_get_m` - Buy N, get M free
- `pct_discount_on_sku` - Percentage discount on SKU
- `pct_discount_on_cart` - Percentage discount on entire cart
- `fixed_discount` - Fixed amount discount

## Response Format

All responses follow a consistent envelope:

```json
{
  "success": true,
  "data": { /* actual response data */ },
  "meta": {
    "request_id": "uuid-string",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

Error responses:
```json
{
  "success": false,
  "message": "Error description",
  "meta": {
    "request_id": "uuid-string",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

## Testing

In Swagger UI, you can:

1. Click "Authorize" to enter your JWT token
2. Expand any endpoint to see its details
3. Click "Try it out" to test (if enabled)
4. View example requests and responses

## Documentation Structure

The OpenAPI spec includes:

- **Detailed descriptions** for every endpoint
- **Complete request/response schemas** with examples
- **Error codes and status codes** for each endpoint
- **Authentication requirements** clearly marked
- **Role-based access control** (buyer/seller) documented
- **Inline code examples** for complex flows

## Updating the Specification

To update the OpenAPI spec:

1. Edit `openapi.yaml`
2. Follow OpenAPI 3.0.3 format
3. Keep examples and descriptions in sync with code
4. Validate the YAML syntax

To update the Swagger UI styling:

1. Edit the `<style>` section in `index.html`
2. Modify CDN versions in `<script>` tags if needed
3. Add custom JavaScript in the `<script>` block at the bottom

## Swagger UI CDN

This documentation uses Swagger UI from unpkg CDN:

- `swagger-ui-dist` v5 - Latest stable version
- CSS: `https://unpkg.com/swagger-ui-dist@5/swagger-ui.css`
- JS Bundle: `https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js`
- Standalone Preset: `https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js`

## Resources

- [OpenAPI 3.0.3 Specification](https://spec.openapis.org/oas/v3.0.3)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [API Design Best Practices](https://swagger.io/resources/articles/best-practices-in-api-design/)
