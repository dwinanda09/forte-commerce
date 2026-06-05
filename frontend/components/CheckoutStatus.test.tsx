import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { CheckoutStatus } from './CheckoutStatus'
import type { CheckoutSession } from '@/lib/types'

describe('CheckoutStatus', () => {
  const mockOnConfirm = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('pending status', () => {
    it('shows processing message for pending checkout', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'pending',
        expires_at: new Date(Date.now() + 600000).toISOString(),
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Processing your checkout...')).toBeInTheDocument()
    })

    it('displays spinning loader for pending status', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'pending',
        expires_at: new Date(Date.now() + 600000).toISOString(),
      }

      const { container } = render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      const spinner = container.querySelector('.animate-spin')
      expect(spinner).toBeInTheDocument()
    })

    it('shows expiration timer for pending checkout', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'pending',
        expires_at: new Date(Date.now() + 120000).toISOString(),
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Expires in')).toBeInTheDocument()
    })

    it('renders timer with minutes and seconds', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'pending',
        expires_at: new Date(Date.now() + 300000).toISOString(),
      }

      const { container } = render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      const timerSpan = container.querySelector('.font-mono')
      expect(timerSpan).toBeInTheDocument()
    })
  })

  describe('completed status', () => {
    const completedSession: CheckoutSession = {
      checkout_id: 'chk-123',
      status: 'completed',
      expires_at: new Date().toISOString(),
      result: {
        items: [
          { sku: 'SKU-001', name: 'Product 1', qty: 2, price: 100, total: 200 },
          { sku: 'SKU-002', name: 'Product 2', qty: 1, price: 50, total: 50 },
        ],
        promotions_applied: [
          { name: 'Summer Sale', description: '10% off', discount: 25 },
        ],
        subtotal: 250,
        total_discount: 25,
        total: 225,
      },
    }

    it('displays order summary heading', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Order Summary')).toBeInTheDocument()
    })

    it('shows items section with all items', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Product 1')).toBeInTheDocument()
      expect(screen.getByText('Product 2')).toBeInTheDocument()
    })

    it('displays item quantities', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText(/Qty: 2/)).toBeInTheDocument()
      expect(screen.getByText(/Qty: 1/)).toBeInTheDocument()
    })

    it('shows item prices and totals', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('$100.00 × 2')).toBeInTheDocument()
      expect(screen.getByText('$200.00')).toBeInTheDocument()
    })

    it('displays promotions section when applied', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Promotions Applied')).toBeInTheDocument()
      expect(screen.getByText('Summer Sale')).toBeInTheDocument()
      expect(screen.getByText('10% off')).toBeInTheDocument()
    })

    it('shows discount amount in promotion', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getAllByText('-$25.00')).toHaveLength(2)
    })

    it('displays subtotal, discount, and total', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Subtotal')).toBeInTheDocument()
      expect(screen.getByText('Total Discount')).toBeInTheDocument()
      expect(screen.getByText('Total')).toBeInTheDocument()
    })

    it('shows confirm button', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByRole('button', { name: /Confirm & Continue to Payment/i })).toBeInTheDocument()
    })

    it('calls onConfirm when confirm button is clicked', async () => {
      const user = userEvent.setup()
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      const confirmButton = screen.getByRole('button', { name: /Confirm & Continue to Payment/i })
      await user.click(confirmButton)

      expect(mockOnConfirm).toHaveBeenCalled()
    })

    it('disables confirm button when confirming', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={true} />)

      const confirmButton = screen.getByRole('button', { name: /Confirming/i })
      expect(confirmButton).toBeDisabled()
    })

    it('shows confirming state text', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={true} />)

      expect(screen.getByText('Confirming...')).toBeInTheDocument()
    })

    it('hides promotions section when no promotions applied', () => {
      const sessionNoPromos: CheckoutSession = {
        ...completedSession,
        result: {
          ...completedSession.result!,
          promotions_applied: [],
        },
      }

      render(<CheckoutStatus session={sessionNoPromos} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.queryByText('Promotions Applied')).not.toBeInTheDocument()
    })

    it('does not show discount row when discount is 0', () => {
      const sessionNoDiscount: CheckoutSession = {
        ...completedSession,
        result: {
          ...completedSession.result!,
          total_discount: 0,
        },
      }

      render(<CheckoutStatus session={sessionNoDiscount} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.queryByText('Total Discount')).not.toBeInTheDocument()
    })
  })

  describe('expired status', () => {
    it('displays expired message', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'expired',
        expires_at: new Date().toISOString(),
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Checkout Expired')).toBeInTheDocument()
    })

    it('shows return to products link', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'expired',
        expires_at: new Date().toISOString(),
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByRole('button', { name: /Return to Products/i })).toBeInTheDocument()
    })

    it('displays error message when provided', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'expired',
        expires_at: new Date().toISOString(),
        error_message: 'Session timeout exceeded',
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Session timeout exceeded')).toBeInTheDocument()
    })
  })

  describe('failed status', () => {
    it('displays failed message', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'failed',
        expires_at: new Date().toISOString(),
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Checkout Failed')).toBeInTheDocument()
    })

    it('shows return to products link for failed', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'failed',
        expires_at: new Date().toISOString(),
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByRole('button', { name: /Return to Products/i })).toBeInTheDocument()
    })

    it('displays error message when provided for failed', () => {
      const session: CheckoutSession = {
        checkout_id: 'chk-123',
        status: 'failed',
        expires_at: new Date().toISOString(),
        error_message: 'Payment processing failed',
      }

      render(<CheckoutStatus session={session} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Payment processing failed')).toBeInTheDocument()
    })
  })

  describe('completed badge', () => {
    const completedSession: CheckoutSession = {
      checkout_id: 'chk-123',
      status: 'completed',
      expires_at: new Date().toISOString(),
      result: {
        items: [],
        promotions_applied: [],
        subtotal: 0,
        total_discount: 0,
        total: 0,
      },
    }

    it('shows completed badge in order summary', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getByText('Completed')).toBeInTheDocument()
    })
  })

  describe('formatting', () => {
    const completedSession: CheckoutSession = {
      checkout_id: 'chk-123',
      status: 'completed',
      expires_at: new Date().toISOString(),
      result: {
        items: [{ sku: 'SKU-001', name: 'Product', qty: 1, price: 99.99, total: 99.99 }],
        promotions_applied: [],
        subtotal: 99.99,
        total_discount: 0,
        total: 99.99,
      },
    }

    it('formats prices with 2 decimal places', () => {
      render(<CheckoutStatus session={completedSession} onConfirm={mockOnConfirm} isConfirming={false} />)

      expect(screen.getAllByText('$99.99').length).toBeGreaterThan(0)
    })
  })
})
