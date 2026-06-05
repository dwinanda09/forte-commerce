import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { CartDrawer } from './CartDrawer'
import type { CartItem } from '@/hooks/useCart'

describe('CartDrawer', () => {
  const mockOnClose = vi.fn()
  const mockOnRemove = vi.fn()
  const mockOnUpdateQty = vi.fn()
  const mockOnClear = vi.fn()

  const mockItems: CartItem[] = [
    { sku: 'SKU-001', name: 'Product 1', price: 100, qty: 2 },
    { sku: 'SKU-002', name: 'Product 2', price: 50, qty: 1 },
  ]

  describe('closed state', () => {
    it('does not render overlay when closed', () => {
      const { container } = render(
        <CartDrawer
          items={[]}
          isOpen={false}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const overlay = container.querySelector('.bg-black\\/20')
      expect(overlay).not.toBeInTheDocument()
    })

    it('drawer is hidden when not open', () => {
      const { container } = render(
        <CartDrawer
          items={[]}
          isOpen={false}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const drawer = container.querySelector('.translate-x-full')
      expect(drawer).toBeInTheDocument()
    })
  })

  describe('open state', () => {
    it('renders overlay when open', () => {
      const { container } = render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const overlay = container.querySelector('.bg-black\\/20')
      expect(overlay).toBeInTheDocument()
    })

    it('drawer is visible when open', () => {
      const { container } = render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const drawer = container.querySelector('.translate-x-0')
      expect(drawer).toBeInTheDocument()
    })

    it('calls onClose when overlay is clicked', async () => {
      const { container } = render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const overlay = container.querySelector('.bg-black\\/20')
      if (overlay) {
        fireEvent.click(overlay)
        expect(mockOnClose).toHaveBeenCalled()
      }
    })

    it('calls onClose when close button is clicked', async () => {
      const user = userEvent.setup()
      render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const closeButton = screen.getByRole('button', { name: '✕' })
      await user.click(closeButton)

      expect(mockOnClose).toHaveBeenCalled()
    })
  })

  describe('empty cart', () => {
    it('shows empty state message', () => {
      render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByText('Your cart is empty')).toBeInTheDocument()
    })

    it('does not show checkout and clear buttons when empty', () => {
      render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.queryByRole('button', { name: /Proceed to Checkout/i })).not.toBeInTheDocument()
      expect(screen.queryByRole('button', { name: /Clear Cart/i })).not.toBeInTheDocument()
    })
  })

  describe('cart with items', () => {
    it('displays cart title with item count', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByText('Cart')).toBeInTheDocument()
      expect(screen.getByText('(3)')).toBeInTheDocument()
    })

    it('displays all cart items', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByText('Product 1')).toBeInTheDocument()
      expect(screen.getByText('Product 2')).toBeInTheDocument()
      expect(screen.getByText('SKU-001')).toBeInTheDocument()
      expect(screen.getByText('SKU-002')).toBeInTheDocument()
    })

    it('displays item quantities', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const qtySpans = screen.getAllByText(/^[0-9]$/)
      expect(qtySpans.length).toBeGreaterThan(0)
    })

    it('displays subtotal correctly', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByText('$250.00')).toBeInTheDocument()
    })

    it('shows checkout and clear buttons', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByRole('button', { name: /Proceed to Checkout/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Clear Cart/i })).toBeInTheDocument()
    })

    it('calls onClear when clear cart button is clicked', async () => {
      const user = userEvent.setup()
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const clearButton = screen.getByRole('button', { name: /Clear Cart/i })
      await user.click(clearButton)

      expect(mockOnClear).toHaveBeenCalled()
    })

    it('calls onRemove when remove button is clicked', async () => {
      const user = userEvent.setup()
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const removeButtons = screen.getAllByRole('button', { name: /Remove/i })
      await user.click(removeButtons[0])

      expect(mockOnRemove).toHaveBeenCalledWith('SKU-001')
    })

    it('calls onUpdateQty when increase qty button is clicked', async () => {
      const user = userEvent.setup()
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const plusButtons = screen.getAllByRole('button', { name: '+' })
      await user.click(plusButtons[0])

      expect(mockOnUpdateQty).toHaveBeenCalledWith('SKU-001', 3)
    })

    it('calls onUpdateQty with decreased qty when decrease button is clicked', async () => {
      const user = userEvent.setup()
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      const minusButtons = screen.getAllByRole('button', { name: '−' })
      await user.click(minusButtons[0])

      expect(mockOnUpdateQty).toHaveBeenCalledWith('SKU-001', 1)
    })
  })

  describe('item total calculation', () => {
    it('calculates item totals correctly', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByText('$200.00')).toBeInTheDocument()
      expect(screen.getByText('$50.00')).toBeInTheDocument()
    })
  })

  describe('header', () => {
    it('displays cart header title', () => {
      render(
        <CartDrawer
          items={mockItems}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.getByText('Cart')).toBeInTheDocument()
    })

    it('does not show item count in header when cart is empty', () => {
      render(
        <CartDrawer
          items={[]}
          isOpen={true}
          onClose={mockOnClose}
          onRemove={mockOnRemove}
          onUpdateQty={mockOnUpdateQty}
          onClear={mockOnClear}
        />
      )

      expect(screen.queryByText(/\(\d+\)/)).not.toBeInTheDocument()
    })
  })
})
