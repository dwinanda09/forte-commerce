# ForteCommerce Frontend Setup

## Quick Start

### Installation

```bash
cd /Users/dwinanda/saas/fortecommerce/frontend
npm install
```

### Development

```bash
npm run dev
```

The frontend runs on `http://localhost:3000`.

Ensure the backend is running at `http://localhost:8080`.

### Build & Production

```bash
npm run build
npm start
```

### Environment Variables

Create `.env.local`:

```
BACKEND_URL=http://localhost:8080
```

## Project Structure

### `/app` — Next.js 15 App Router
- `layout.tsx` — Root layout with global styles
- `login/page.tsx` — Login page with demo credentials
- `(protected)/` — Protected routes with middleware
  - `page.tsx` — Products listing
  - `checkout/page.tsx` — Checkout submission
  - `checkout/[id]/page.tsx` — Checkout status polling
  - `orders/page.tsx` — Orders list
  - `orders/[id]/page.tsx` — Order detail with pay/cancel actions
  - `api/[...path]/route.ts` — Backend proxy with auth header injection

### `/components` — React Components
- `ProductCard.tsx` — Single product card with qty selector
- `CartDrawer.tsx` — Right-side drawer cart overlay
- `CheckoutStatus.tsx` — Checkout status display with timer
- `OrderStatusBadge.tsx` — Order status pill badge
- `CheckoutStatusBadge.tsx` — Checkout status pill badge
- `Navigation.tsx` — Top nav with logo, orders link, logout

### `/hooks` — Custom Hooks
- `useAuth.ts` — Login/logout with js-cookie
- `useCart.ts` — Zustand cart state management
- `useCheckoutPoller.ts` — SWR-based checkout polling with 2s interval

### `/lib` — Utilities
- `types.ts` — TypeScript interfaces for API responses and models
- `api.ts` — Fetch wrapper that hits `/api/[...path]` proxy routes

### Root Files
- `middleware.ts` — Auth check, redirect to /login if no token
- `tailwind.config.ts` — Mountain Mist color palette + custom tokens
- `globals.css` — Inter font, TailwindCSS directives, focus-visible ring
- `postcss.config.js` — Tailwind + Autoprefixer
- `next.config.ts` — Base Next.js config
- `tsconfig.json` — Strict mode enabled
- `package.json` — Next.js 15, React 19, SWR, Zustand, js-cookie
- `Dockerfile` — Multi-stage build for production

## Design System

### Colors (Mountain Mist, Light Luxury)
- **Teal** `#01796F` — Primary action (CTAs, links, accents)
- **Steel** `#6D8196` — Borders, muted text, secondary labels
- **Mist** `#B0C4DE` — Surface tints, dividers, light backgrounds
- **Graphite** `#5A5A5A` — Body text
- **Surface** `#F8FAFB` — Page background
- **White** `#FFFFFF` — Cards, modals

### Typography
- Font: **Inter** 400, 600, 700
- Max two weights: 400 (regular), 600/700 (semibold/bold)

### Spacing & Radius
- Radius: `sm: 4px`, `md: 8px`, `lg: 12px`
- Shadows: `card: 0 1px 3px rgba(0,0,0,.08)`, `modal: 0 10px 40px rgba(0,0,0,.12)`

### Interactive States
- **Focus-visible**: 2px solid teal outline + 2px offset
- **Hover buttons**: Teal hover shade `#015f58`
- **Disabled**: `opacity-40`, `cursor-not-allowed`

## API Routes

All requests go through `/api/[...path]/route.ts` proxy. Backend must run at `http://localhost:8080`.

### Implemented Endpoints
- `POST /api/v1/auth/login` → `{ token }`
- `GET /api/v1/products` → `Product[]`
- `POST /api/v1/checkout` → `{ checkout_id }`
- `GET /api/v1/checkout/:id` → `CheckoutSession`
- `POST /api/v1/checkout/:id/confirm` → `Order`
- `POST /api/v1/orders/:id/pay` → `Order`
- `POST /api/v1/orders/:id/cancel` → `Order`
- `GET /api/v1/orders/:id` → `Order`
- `GET /api/v1/orders` → `Order[]`

## Key Features

### 1. Authentication
- Login form with username/password
- Token stored in secure cookie (`auth_token`)
- Middleware protects routes, redirects to `/login`
- Logout clears cookie and redirects

### 2. Product Browsing
- SWR data fetching with loading/error states
- Responsive grid (1 mobile, 2 tablet, 3 desktop)
- SKU, name, price (teal), available qty
- Qty selector + "Add to Cart" button
- Out-of-stock state (disabled, muted)

### 3. Cart Management
- Zustand state store
- Add/remove/update qty
- Fixed FAB shows item count
- Right-side drawer with line items, subtotal
- "Proceed to Checkout" button

### 4. Checkout Flow
- Submit cart SKUs → receive `checkout_id`
- Poll checkout status every 2s while pending
- Show countdown timer (expires_at)
- On completion: display order summary, items, promotions, total
- "Confirm & Continue" button → creates order, redirects to `/orders/:id`

### 5. Order Management
- Orders list table with filtering
- Order detail page with items, promotions, totals
- Status-based actions:
  - **Pending**: Pay Now + Cancel Order buttons
  - **Paid**: Read-only display
  - **Cancelled**: Read-only display

### 6. Error Handling
- API errors caught and displayed as toast/modal
- Form validation on inputs
- Graceful loading states with spinner

## No Install Required

The task specifies **not** to run `npm install`. All files are ready. To start the project:

```bash
cd /Users/dwinanda/saas/fortecommerce/frontend
npm install  # Run this yourself
npm run dev
```

## Notes

- No `console.log` in production code (use proper logging if needed)
- No hardcoded secrets (env vars used for backend URL)
- No shadcn defaults (custom component styling with Mountain Mist palette)
- Strict TypeScript: no `any` types
- All interactive elements have focus-visible rings
- Responsive design: mobile-first Tailwind approach
- No secondary accent colors (teal is the only action color)
