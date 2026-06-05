# ForteCommerce 🛒

[![Go 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![Next.js 15](https://img.shields.io/badge/Next.js-15-black?logo=next.js&logoColor=white)](https://nextjs.org)
[![PostgreSQL 16](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![RabbitMQ 3](https://img.shields.io/badge/RabbitMQ-3-FF6600?logo=rabbitmq&logoColor=white)](https://www.rabbitmq.com)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://www.docker.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A full-stack, event-driven e-commerce platform with async checkout processing, intelligent promotion engine, and role-based multi-seller support. Built with Go, Next.js, PostgreSQL, and RabbitMQ.

---

## 📑 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture Overview](#-architecture-overview)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Environment Variables](#-environment-variables)
- [Database Migrations](#-database-migrations)
- [Available Make Commands](#-available-make-commands)
- [API Overview](#-api-overview)
- [Architecture Diagrams](#-architecture-diagrams)
- [Campaign System](#-campaign-system)
- [Testing](#-testing)
- [API Documentation](#-api-documentation)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🚀 Features

### Buyer Features
- User registration and JWT-based authentication
- Browse products with real-time inventory visibility
- Add items to cart with SKU-based selection
- Async checkout with 15-minute session expiry
- Automatic promotion evaluation and discount application
- Order confirmation and payment tracking
- View order history and details
- Cancel orders and restore inventory
- Real-time checkout status polling

### Seller Features
- Product management (create, read, update, delete)
- Inventory tracking with reserved quantity management
- Campaign creation and management
- Flexible promotion conditions and actions
- Campaign priority-based execution
- Campaign activation/deactivation toggle
- Real-time inventory updates

### System Features
- Event-driven architecture with RabbitMQ message queue
- Asynchronous checkout processing with transaction isolation
- Automatic checkout session expiry (60-second background worker)
- Row-level locking for inventory consistency
- JWT-based authorization with role-based access control (RBAC)
- Bcrypt password hashing with per-user salt
- Structured logging and error tracking
- Comprehensive test coverage (≥90%)
- Containerized deployment with Docker Compose

---

## 🏗️ Tech Stack

| Layer | Technology | Version | Purpose |
|-------|-----------|---------|---------|
| **Backend** | Go | 1.23+ | Core API server |
| | Echo Framework | v4 | HTTP routing & middleware |
| | sqlx | Latest | Database abstraction |
| | JWT | Custom | Token-based auth |
| | bcrypt | Golang std | Password hashing |
| | RabbitMQ AMQP | Go driver | Message queue integration |
| **Frontend** | Next.js | 15 | React framework |
| | TypeScript | Latest | Type safety |
| | Tailwind CSS | Latest | Utility-first styling |
| | Vitest | Latest | Unit testing |
| | Playwright | Latest | E2E testing |
| **Database** | PostgreSQL | 16 | Primary data store |
| **Message Queue** | RabbitMQ | 3 (Management) | Async task processing |
| **Infrastructure** | Docker Compose | Latest | Container orchestration |

---

## 🏗️ Architecture Overview

```mermaid
graph TB
    subgraph Client["🖥️ Client Layer"]
        FE["Frontend<br/>Next.js :3000"]
    end

    subgraph API["🔌 API Gateway"]
        BE["Backend API<br/>Echo :8080"]
    end

    subgraph Queue["📨 Message Queue"]
        MQ["RabbitMQ :14046<br/>Management :14047"]
    end

    subgraph Workers["⚙️ Workers"]
        CW["Checkout Worker<br/>Async Processing"]
        EW["Expiry Worker<br/>60s Tick"]
    end

    subgraph Storage["💾 Data Layer"]
        DB["PostgreSQL :14045<br/>forte_commerce"]
    end

    FE -->|HTTP/HTTPS| BE
    BE -->|Query/Update| DB
    BE -->|Publish Jobs| MQ
    MQ -->|Consume Jobs| CW
    CW -->|Query/Update| DB
    EW -->|Query/Update| DB
    BE -.->|Check Status| MQ

    style Client fill:#e1f5ff
    style API fill:#fff3e0
    style Queue fill:#f3e5f5
    style Workers fill:#e8f5e9
    style Storage fill:#fce4ec
```

**Data Flow:**
1. **Frontend** → REST API calls to Backend
2. **Backend** → Validates requests, stores data in PostgreSQL, publishes async jobs to RabbitMQ
3. **Checkout Worker** → Consumes checkout jobs, evaluates promotions, reserves inventory, updates PostgreSQL
4. **Expiry Worker** → Periodically releases expired checkout sessions and inventory reservations
5. **RabbitMQ Management UI** → Monitor queues at `http://localhost:15672` (default: guest/guest)

---

## 📁 Project Structure

### Backend (`/backend`)

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go                 # Entry point, server initialization
│   └── migrate/
│       └── main.go                 # Migration runner
├── internal/
│   ├── config/
│   │   └── config.go               # Environment & app configuration
│   ├── domain/
│   │   ├── models.go               # Core domain entities
│   │   ├── errors.go               # Domain-specific errors
│   │   └── constants.go            # Enums & constants
│   ├── handler/
│   │   ├── auth.go                 # Auth endpoints
│   │   ├── product.go              # Product endpoints
│   │   ├── checkout.go             # Checkout endpoints
│   │   ├── order.go                # Order endpoints
│   │   └── campaign.go             # Campaign endpoints
│   ├── middleware/
│   │   ├── auth.go                 # JWT validation
│   │   ├── request_id.go           # Request ID tracking
│   │   └── error_handler.go        # Global error handling
│   ├── promotion/
│   │   ├── engine.go               # Promotion evaluation logic
│   │   ├── conditions.go           # Condition evaluators
│   │   └── actions.go              # Action executors
│   ├── queue/
│   │   ├── checkout.go             # Checkout job handler
│   │   └── expiry.go               # Expiry worker
│   ├── resource/
│   │   ├── repository.go           # Data access layer
│   │   └── postgres.go             # PostgreSQL connection
│   ├── router/
│   │   └── router.go               # Route definitions
│   └── usecase/
│       ├── auth.go                 # Authentication logic
│       ├── product.go              # Product business logic
│       ├── checkout.go             # Checkout orchestration
│       ├── order.go                # Order business logic
│       └── campaign.go             # Campaign management
├── util/
│   ├── logger.go                   # Structured logging
│   ├── response.go                 # API response formatting
│   └── error_wrapper.go            # Error utilities
├── migrations/
│   └── *.sql                       # Goose database migrations
├── Dockerfile                      # Backend container image
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
└── main.go                         # Bootstrap (via cmd/server)
```

### Frontend (`/frontend`)

```
frontend/
├── app/
│   ├── layout.tsx                  # Root layout
│   ├── page.tsx                    # Home page
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (buyer)/
│   │   ├── products/page.tsx
│   │   ├── checkout/page.tsx
│   │   ├── checkout/[id]/page.tsx
│   │   └── orders/page.tsx
│   └── (seller)/
│       ├── products/page.tsx
│       ├── products/new/page.tsx
│       ├── campaigns/page.tsx
│       └── campaigns/new/page.tsx
├── components/
│   ├── ProductCard.tsx
│   ├── CartDrawer.tsx
│   ├── CheckoutStatus.tsx
│   ├── CampaignForm.tsx
│   ├── Navigation.tsx
│   └── [other UI components]
├── hooks/
│   ├── useAuth.ts
│   ├── useCart.ts
│   └── [custom hooks]
├── lib/
│   ├── api.ts                      # API client
│   ├── auth.ts                     # Auth utilities
│   └── utils.ts
├── styles/
│   └── globals.css
├── e2e/
│   └── [Playwright tests]
├── Dockerfile                      # Frontend container image
├── next.config.ts                  # Next.js configuration
├── tailwind.config.ts              # Tailwind CSS config
├── package.json
└── tsconfig.json
```

---

## 🚀 Getting Started

### Prerequisites

- **Docker** or **OrbStack** (macOS recommended for better performance)
- **Go 1.23+** (for local backend development)
- **Node.js 18+** (for local frontend development)
- **PostgreSQL client tools** (optional, for direct DB access)

### Quick Start: Docker (All Services)

Clone the repository and start all services:

```bash
git clone https://github.com/yourusername/fortecommerce.git
cd fortecommerce

# Start PostgreSQL, RabbitMQ, Backend, and Frontend
make dev
```

Services will be available at:
- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **RabbitMQ Management:** http://localhost:15672 (guest/guest)
- **PostgreSQL:** localhost:14045

### Quick Start: Local Development (Hot Reload)

For faster iteration with file watching:

```bash
# Copy environment template
cp .env.example .env

# Start infrastructure (Docker) + run BE & FE locally with hot-reload
make local
```

**What happens:**
1. PostgreSQL and RabbitMQ start in Docker
2. Backend runs locally with file watching (will restart on code changes)
3. Frontend runs locally with Next.js hot reload (will refresh on code changes)

Press `Ctrl+C` to stop. Run `make stop` to clean up Docker containers.

### Manual Setup (Advanced)

If you prefer more control:

```bash
# Start infrastructure
make docker-up

# In separate terminals:
make dev-backend
make dev-frontend
```

---

## 📋 Environment Variables

Create a `.env` file in the project root with the following variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL hostname |
| `DB_PORT` | `14045` | PostgreSQL port |
| `DB_USER` | `forte` | PostgreSQL username |
| `DB_PASSWORD` | `forte123` | PostgreSQL password |
| `DB_NAME` | `forte_commerce` | PostgreSQL database name |
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:14046/` | RabbitMQ connection string |
| `JWT_SECRET` | `change-me-to-a-secure-random-string` | JWT signing secret (⚠️ CHANGE IN PRODUCTION) |
| `PORT` | `8080` | Backend server port |
| `BACKEND_URL` | `http://localhost:8080` | Backend URL for frontend |

### Frontend Environment (`.env.local`)

Create `frontend/.env.local`:

```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

---

## 🗄️ Database Migrations

ForteCommerce uses **Goose** for database migrations.

### Install Goose

```bash
make install-goose
```

### Run Migrations

```bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status
```

Migrations are stored in `/backend/migrations/*.sql`. See the schema definition for table structure.

---

## 📌 Available Make Commands

```bash
# Development & Deployment
make dev                   # Stop all, start infra → BE → FE in Docker
make local                 # Start infra in Docker, run BE+FE locally with hot-reload
make stop                  # Stop all containers and local processes

# Infrastructure
make docker-up             # Start PostgreSQL + RabbitMQ only
make docker-down           # Stop all containers
make compose-up            # Start all services via docker-compose
make compose-down          # Stop containers and remove volumes

# Database Migrations
make install-goose         # Install goose migration tool
make migrate-up            # Apply all pending migrations
make migrate-down          # Rollback last migration
make migrate-status        # Show migration status

# Local Development
make dev-backend           # Run backend with hot-reload (cmd: go run ./cmd/server)
make dev-frontend          # Run frontend with hot-reload (cmd: npm run dev)

# Testing & Quality
make test                  # Run backend unit tests with race detector
make test-coverage         # Run tests with coverage (requires ≥90%)
make build                 # Build backend binary to ./bin/server
make clean                 # Remove build artifacts
make lint                  # Run go vet on backend

# Help
make help                  # Display this help message
```

---

## 💡 API Overview

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All endpoints (except `/auth/login` and `/auth/register`) require JWT in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

### Response Envelope

All API responses follow this format:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Example"
  },
  "meta": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2024-06-05T15:30:45Z"
  }
}
```

On error:

```json
{
  "success": false,
  "data": null,
  "meta": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2024-06-05T15:30:45Z",
    "error": "Insufficient inventory"
  }
}
```

### Endpoints

#### Authentication

| Method | Endpoint | Auth | Description | Request |
|--------|----------|------|-------------|---------|
| `POST` | `/auth/login` | None | Login with credentials, returns JWT | `{ "username": "john", "password": "secret" }` |
| `POST` | `/auth/register` | None | Register new user account | `{ "username": "john", "password": "secret", "email": "john@example.com" }` |

#### Products (Public)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/products` | JWT | List all products with inventory |
| `GET` | `/products/:id` | JWT | Get product details |

#### Products (Seller Only)

| Method | Endpoint | Auth | Description | Request |
|--------|----------|------|-------------|---------|
| `POST` | `/seller/products` | JWT+Seller | Create new product | `{ "sku": "SKU-001", "name": "Widget", "price": 29.99, "inventory_qty": 100 }` |
| `PUT` | `/seller/products/:id` | JWT+Seller | Update product | `{ "name": "Widget Pro", "price": 39.99, "inventory_qty": 150 }` |
| `DELETE` | `/seller/products/:id` | JWT+Seller | Delete product | |

#### Checkout (Async)

| Method | Endpoint | Auth | Description | Request |
|--------|----------|------|-------------|---------|
| `POST` | `/checkout` | JWT | Submit checkout (async, queued) | `{ "items": [{ "sku": "SKU-001", "qty": 2 }, { "sku": "SKU-002", "qty": 1 }] }` |
| `GET` | `/checkout/:id` | JWT | Get checkout session status | |
| `POST` | `/checkout/:id/confirm` | JWT | Confirm checkout → create order | |

#### Orders

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/orders` | JWT | List user's orders |
| `GET` | `/orders/:id` | JWT | Get order details |
| `POST` | `/orders/:id/pay` | JWT | Mark order as paid |
| `POST` | `/orders/:id/cancel` | JWT | Cancel order & restore inventory |

#### Campaigns (Public List)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/campaigns` | JWT | List all active campaigns |

#### Campaigns (Seller Only)

| Method | Endpoint | Auth | Description | Request |
|--------|----------|------|-------------|---------|
| `GET` | `/seller/campaigns` | JWT+Seller | List seller's campaigns | |
| `GET` | `/seller/campaigns/:id` | JWT+Seller | Get campaign details | |
| `POST` | `/seller/campaigns` | JWT+Seller | Create campaign | See [Campaign System](#-campaign-system) |
| `PUT` | `/seller/campaigns/:id` | JWT+Seller | Update campaign | |
| `DELETE` | `/seller/campaigns/:id` | JWT+Seller | Delete campaign | |
| `PATCH` | `/seller/campaigns/:id/toggle` | JWT+Seller | Toggle active status | `{ "is_active": true }` |

---

## 🏛️ Architecture Diagrams

### 11a. Authentication Flow

```mermaid
sequenceDiagram
    actor User
    participant FE as Frontend
    participant BE as Backend API
    participant DB as PostgreSQL

    User->>FE: Enter credentials
    FE->>BE: POST /api/v1/auth/login
    BE->>DB: SELECT user WHERE username=?
    DB-->>BE: User record + password_hash + salt
    BE->>BE: bcrypt verify (sha256(password+salt))
    alt Valid credentials
        BE->>BE: Sign JWT (user_id, username, role, exp)
        BE-->>FE: 200 { token: "eyJ..." }
        FE->>FE: Store token in localStorage
    else Invalid credentials
        BE-->>FE: 401 Invalid credentials
    end
```

### 11b. Async Checkout Flow (Core)

```mermaid
sequenceDiagram
    actor Buyer
    participant FE as Frontend
    participant BE as Backend API
    participant MQ as RabbitMQ
    participant W as Checkout Worker
    participant DB as PostgreSQL

    Buyer->>FE: Add items to cart & checkout
    FE->>BE: POST /api/v1/checkout { items: [SKU1, SKU2] }
    BE->>DB: INSERT checkout_session (status=pending, expires_in=15min)
    BE->>MQ: Publish CheckoutJob { checkout_id, items }
    BE-->>FE: 202 { checkout_id: "uuid" }

    loop Poll until status != pending
        FE->>BE: GET /api/v1/checkout/{id}
        BE->>DB: SELECT checkout_session WHERE id=?
        DB-->>BE: Session record
        BE-->>FE: { status: "pending" }
    end

    Note over W,MQ: Async Processing (Background)
    MQ->>W: Consume CheckoutJob
    W->>DB: BEGIN TRANSACTION
    W->>DB: SELECT products WHERE sku IN (...) FOR UPDATE
    DB-->>W: Product rows (row-level locked)
    W->>W: Validate inventory availability
    W->>DB: SELECT active campaigns ORDER BY priority DESC
    W->>W: Evaluate campaign conditions (ALL must pass)
    W->>W: Apply matching campaign actions → compute discounts
    W->>DB: UPDATE products SET reserved_qty += qty
    W->>DB: UPDATE checkout_session SET status=completed, result=JSON
    W->>DB: COMMIT

    FE->>BE: GET /api/v1/checkout/{id}
    BE-->>FE: { status: "completed", result: { items, total, promotions } }

    Buyer->>FE: Review & confirm order
    FE->>BE: POST /api/v1/checkout/{id}/confirm
    BE->>DB: BEGIN TRANSACTION
    BE->>DB: INSERT order (from checkout result)
    BE->>DB: UPDATE products SET inventory_qty -= qty, reserved_qty -= qty
    BE->>DB: COMMIT
    BE-->>FE: 200 { order: { id, status, items, total } }
```

### 11c. Order Payment & Cancellation Flow

```mermaid
sequenceDiagram
    actor Buyer
    participant BE as Backend API
    participant DB as PostgreSQL

    alt Pay Order
        Buyer->>BE: POST /api/v1/orders/{id}/pay
        BE->>DB: UPDATE orders SET status=paid
        BE-->>Buyer: 200 { status: "paid" }
    else Cancel Order
        Buyer->>BE: POST /api/v1/orders/{id}/cancel
        BE->>DB: BEGIN TRANSACTION
        BE->>DB: UPDATE products SET inventory_qty += qty (restore stock)
        BE->>DB: UPDATE orders SET status=cancelled
        BE->>DB: COMMIT
        BE-->>Buyer: 200 { status: "cancelled" }
    end
```

### 11d. Checkout Expiry Worker Flow

```mermaid
sequenceDiagram
    participant EW as Expiry Worker
    participant DB as PostgreSQL

    loop Every 60 seconds
        EW->>DB: SELECT checkout_sessions WHERE status=pending AND expires_at < NOW()
        DB-->>EW: Expired sessions list
        loop For each expired session
            EW->>DB: BEGIN TRANSACTION
            EW->>DB: SELECT products WHERE reserved (from session) FOR UPDATE
            EW->>DB: UPDATE products SET reserved_qty -= qty (release reservation)
            EW->>DB: UPDATE checkout_session SET status=expired
            EW->>DB: COMMIT
        end
    end
```

### 12a. Checkout Processing Logic

```mermaid
flowchart TD
    A[CheckoutJob Received] --> B[Parse item SKUs]
    B --> C[Lock product rows FOR UPDATE]
    C --> D{All products exist?}
    D -- No --> E[Set status=failed: product not found]
    D -- Yes --> F{Sufficient inventory?}
    F -- No --> G[Set status=failed: insufficient stock]
    F -- Yes --> H[Fetch active campaigns by priority DESC]
    H --> I[For each campaign, evaluate ALL conditions]
    I --> J{All conditions pass?}
    J -- No --> K[Skip campaign]
    K --> I
    J -- Yes --> L[Apply campaign actions → compute discounts]
    L --> M[Update reserved_qty in inventory]
    M --> N[Save CheckoutResult as JSON]
    N --> O[Set status=completed]
    O --> P[Release DB lock & COMMIT]
```

### 12b. Promotion Engine Logic

```mermaid
flowchart TD
    A["Cart + Promotions List"] --> B["Sort campaigns by priority DESC"]
    B --> C["For each campaign (in order)"]
    C --> D{"is_active = true?"}
    D -- No --> E["Skip campaign"]
    D -- Yes --> F["Evaluate ALL conditions"]
    F --> G{"ALL conditions pass?"}
    G -- No --> E
    G -- Yes --> H["Apply each action"]
    H --> I["Add to promotions_applied list"]
    I --> J{"More campaigns?"}
    J -- Yes --> C
    J -- No --> K["Return final cart + total discounts"]
```

### 12c. Auth & Role Check Flow

```mermaid
flowchart LR
    A["Incoming Request"] --> B{"Authorization header present?"}
    B -- No --> C["401 Unauthorized"]
    B -- Yes --> D["Parse Bearer token"]
    D --> E{"JWT valid + not expired?"}
    E -- No --> C
    E -- Yes --> F["Extract claims: user_id, username, role"]
    F --> G{"Is seller-only route?"}
    G -- No --> H["Proceed to handler"]
    G -- Yes --> I{"role == seller?"}
    I -- No --> J["403 Forbidden"]
    I -- Yes --> H
```

---

## 💰 Campaign System

The promotion engine evaluates campaigns based on flexible conditions and actions. Campaigns are stored as JSON documents and processed during checkout.

### Condition Types

Conditions are **AND-ed together**. All must pass for the campaign to apply.

| Type | Parameters | Description | Example |
|------|-----------|-------------|---------|
| `cart_has_sku` | `sku: string`, `min_qty: int` | Cart must contain SKU with at least min_qty | `{ "type": "cart_has_sku", "sku": "SKU-001", "min_qty": 2 }` |
| `item_qty_gte` | `sku: string`, `qty: int` | Specific SKU quantity >= qty | `{ "type": "item_qty_gte", "sku": "SKU-002", "qty": 3 }` |
| `cart_total_gte` | `amount: decimal` | Cart subtotal >= amount (before discounts) | `{ "type": "cart_total_gte", "amount": 100.00 }` |
| `cart_item_count_gte` | `count: int` | Number of unique SKUs >= count | `{ "type": "cart_item_count_gte", "count": 3 }` |

### Action Types

Actions modify the cart total or apply item-level discounts.

| Type | Parameters | Description | Example |
|------|-----------|-------------|---------|
| `free_item` | `sku: string`, `trigger_sku: string` | Get `sku` free when `trigger_sku` is in cart | `{ "type": "free_item", "sku": "SKU-BONUS", "trigger_sku": "SKU-001" }` |
| `buy_n_get_m` | `sku: string`, `buy_n: int`, `pay_m: int` | Buy `buy_n` of `sku`, pay for only `pay_m` | `{ "type": "buy_n_get_m", "sku": "SKU-001", "buy_n": 3, "pay_m": 2 }` |
| `pct_discount_on_sku` | `sku: string`, `pct: int` | Apply `pct`% discount on specific SKU | `{ "type": "pct_discount_on_sku", "sku": "SKU-001", "pct": 20 }` |
| `pct_discount_on_cart` | `pct: int` | Apply `pct`% discount on entire cart total | `{ "type": "pct_discount_on_cart", "pct": 10 }` |
| `fixed_discount` | `amount: decimal` | Fixed amount off cart total | `{ "type": "fixed_discount", "amount": 15.00 }` |

### Example Campaign

**"Buy 3 Pay 2" Promotion:**

```json
{
  "id": "camp-001",
  "seller_id": "seller-123",
  "name": "Buy 3 Pay 2",
  "description": "Buy 3 of SKU-001, pay for only 2",
  "is_active": true,
  "priority": 10,
  "conditions": [
    {
      "type": "item_qty_gte",
      "sku": "SKU-001",
      "qty": 3
    }
  ],
  "actions": [
    {
      "type": "buy_n_get_m",
      "sku": "SKU-001",
      "buy_n": 3,
      "pay_m": 2
    }
  ],
  "created_at": "2024-06-05T10:00:00Z",
  "updated_at": "2024-06-05T10:00:00Z"
}
```

**"Free Shipping Item" Promotion:**

```json
{
  "id": "camp-002",
  "seller_id": "seller-123",
  "name": "Free Shipping with Premium",
  "description": "Buy Premium SKU, get free shipping item",
  "is_active": true,
  "priority": 5,
  "conditions": [
    {
      "type": "cart_has_sku",
      "sku": "SKU-PREMIUM",
      "min_qty": 1
    },
    {
      "type": "cart_total_gte",
      "amount": 50.00
    }
  ],
  "actions": [
    {
      "type": "free_item",
      "sku": "SKU-SHIPPING",
      "trigger_sku": "SKU-PREMIUM"
    }
  ]
}
```

**"Summer Flash Sale" Promotion:**

```json
{
  "id": "camp-003",
  "seller_id": "seller-456",
  "name": "Summer Flash Sale",
  "description": "25% off all summer items",
  "is_active": true,
  "priority": 3,
  "conditions": [],
  "actions": [
    {
      "type": "pct_discount_on_cart",
      "pct": 25
    }
  ]
}
```

---

## ✅ Testing

### Backend Tests

Run backend unit tests with race detection:

```bash
make test
```

### Test Coverage (Required: ≥90%)

Generate coverage report:

```bash
make test-coverage
```

This will:
1. Run all tests with coverage tracking
2. Enforce ≥90% code coverage (fails if below)
3. Generate `backend/coverage.html` for visual inspection

View coverage report:

```bash
open backend/coverage.html
```

### Frontend Tests

Run frontend unit tests and E2E tests:

```bash
cd frontend

# Unit tests (Vitest)
npm run test

# E2E tests (Playwright)
npm run test:e2e

# Test coverage
npm run test:coverage
```

---

## 📚 API Documentation

### Postman Collection

Import the Postman collection to test all endpoints:

- **File:** `docs/fortecommerce.postman_collection.json`
- **Import:** Open Postman → File → Import → Select JSON file
- **Environments:** Update `{{base_url}}` and `{{token}}` variables

### Swagger UI

Generated API documentation with Swagger/OpenAPI:

- **File:** `docs/swagger/index.html`
- **Open:** Double-click to open in browser, or serve with `python -m http.server 8000` in `/docs`

### cURL Examples

**Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "secret"
  }'
```

**List Products:**

```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Submit Checkout:**

```bash
curl -X POST http://localhost:8080/api/v1/checkout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "sku": "SKU-001", "qty": 2 },
      { "sku": "SKU-002", "qty": 1 }
    ]
  }'
```

---

## 🚀 Deployment

This project follows **Trunk-Based Development (TBD)** as its branching strategy.

### Trunk-Based Development

All engineers commit directly to `main` (the trunk) or via short-lived feature branches (max 1–2 days). Long-running branches are avoided.

```
main (trunk)
 ├── feature/add-discount-engine   ← short-lived, merged within 1-2 days
 ├── fix/checkout-expiry-bug       ← short-lived, merged within 1-2 days
 └── ...
```

### Key Principles

- **`main` is always deployable.** Every commit merged to `main` must pass all tests and be production-ready.
- **Small, frequent commits.** Prefer many small PRs over large batches of changes.
- **Feature flags over long branches.** Incomplete features are hidden behind flags, not isolated in separate branches.
- **No long-lived branches.** Feature branches must be merged or deleted within 1–2 days.

### Branch Naming

| Pattern | Purpose |
|---------|---------|
| `main` | Production trunk — always green |
| `feature/<short-description>` | New features (max 2 days) |
| `fix/<short-description>` | Bug fixes |
| `chore/<short-description>` | Tooling, deps, config changes |

### CI/CD Flow

```
Developer opens PR (from short-lived branch)
       │
       ▼
  Peer Review (required approval)
       │
       ▼
  CI checks pass (lint, tests, build)
       │
       ▼
  Merge PR → main
       │
       └── Grab merged commit hash (e.g. a3f92c1)
                    │
                    ▼
       Trigger GitHub Actions workflow manually
       (workflow_dispatch) with commit hash as input
                    │
                    ▼
       GitHub Actions deploys main @ <commit_hash>
```

**Deploy trigger (GitHub Actions `workflow_dispatch`):**

```yaml
on:
  workflow_dispatch:
    inputs:
      commit_sha:
        description: 'Commit hash to deploy (from merged PR)'
        required: true
        type: string

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{ github.event.inputs.commit_sha }}
      # ... build & deploy steps
```

### Running in Production

```bash
make compose-up    # start all services (postgres, rabbitmq, backend, frontend)
make compose-down  # stop and remove volumes
```

---

## 🤝 Contributing

1. **Fork** the repository
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Write tests** for your changes (TDD approach)
4. **Commit** with clear messages (`git commit -m 'feat: add amazing feature'`)
5. **Push** to your fork (`git push origin feature/amazing-feature`)
6. **Create a Pull Request** with a detailed description

### Code Standards

- **Backend:** Go, follow `gofmt` and `go vet`
- **Frontend:** TypeScript, follow ESLint rules
- **Tests:** Minimum 80% coverage required
- **Commits:** Follow Conventional Commits format

See [contributing guidelines](CONTRIBUTING.md) for more details.

---

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

## 🙏 Support

For questions, issues, or suggestions:

- Open an [Issue](../../issues)
- Start a [Discussion](../../discussions)
- Email: support@fortecommerce.local

---

**Last Updated:** June 5, 2024  
**Maintainer:** [Your Name](https://github.com/yourusername)
