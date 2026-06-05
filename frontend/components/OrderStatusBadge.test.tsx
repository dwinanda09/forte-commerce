import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { OrderStatusBadge } from './OrderStatusBadge'
import type { Order } from '@/lib/types'

describe('OrderStatusBadge', () => {
  describe('pending status', () => {
    it('renders pending label', () => {
      const status: Order['status'] = 'pending'
      render(<OrderStatusBadge status={status} />)

      expect(screen.getByText('Pending')).toBeInTheDocument()
    })

    it('applies pending background color', () => {
      const status: Order['status'] = 'pending'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('.bg-yellow-50')
      expect(badge).toBeInTheDocument()
    })

    it('applies pending text color', () => {
      const status: Order['status'] = 'pending'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('.text-yellow-900')
      expect(badge).toBeInTheDocument()
    })

    it('has badge styling classes', () => {
      const status: Order['status'] = 'pending'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('span')
      expect(badge).toHaveClass('inline-block')
      expect(badge).toHaveClass('rounded-full')
      expect(badge).toHaveClass('text-sm')
      expect(badge).toHaveClass('font-medium')
    })
  })

  describe('paid status', () => {
    it('renders paid label', () => {
      const status: Order['status'] = 'paid'
      render(<OrderStatusBadge status={status} />)

      expect(screen.getByText('Paid')).toBeInTheDocument()
    })

    it('applies paid background color', () => {
      const status: Order['status'] = 'paid'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('.bg-teal-50')
      expect(badge).toBeInTheDocument()
    })

    it('applies paid text color', () => {
      const status: Order['status'] = 'paid'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('.text-teal')
      expect(badge).toBeInTheDocument()
    })
  })

  describe('cancelled status', () => {
    it('renders cancelled label', () => {
      const status: Order['status'] = 'cancelled'
      render(<OrderStatusBadge status={status} />)

      expect(screen.getByText('Cancelled')).toBeInTheDocument()
    })

    it('applies cancelled background color', () => {
      const status: Order['status'] = 'cancelled'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('.bg-gray-100')
      expect(badge).toBeInTheDocument()
    })

    it('applies cancelled text color', () => {
      const status: Order['status'] = 'cancelled'
      const { container } = render(<OrderStatusBadge status={status} />)

      const badge = container.querySelector('.text-gray-700')
      expect(badge).toBeInTheDocument()
    })
  })

  describe('visual consistency', () => {
    it('all statuses have padding', () => {
      const statuses: Order['status'][] = ['pending', 'paid', 'cancelled']

      statuses.forEach((status) => {
        const { container } = render(<OrderStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('px-3')
        expect(badge).toHaveClass('py-1')
      })
    })

    it('all statuses are inline blocks', () => {
      const statuses: Order['status'][] = ['pending', 'paid', 'cancelled']

      statuses.forEach((status) => {
        const { container } = render(<OrderStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('inline-block')
      })
    })

    it('all statuses have text-sm font-medium', () => {
      const statuses: Order['status'][] = ['pending', 'paid', 'cancelled']

      statuses.forEach((status) => {
        const { container } = render(<OrderStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('text-sm')
        expect(badge).toHaveClass('font-medium')
      })
    })
  })

  describe('different statuses render correctly', () => {
    it('renders each status independently', () => {
      const { rerender } = render(<OrderStatusBadge status="pending" />)
      expect(screen.getByText('Pending')).toBeInTheDocument()

      rerender(<OrderStatusBadge status="paid" />)
      expect(screen.getByText('Paid')).toBeInTheDocument()

      rerender(<OrderStatusBadge status="cancelled" />)
      expect(screen.getByText('Cancelled')).toBeInTheDocument()
    })
  })

  describe('edge cases', () => {
    it('handles uppercase status names', () => {
      const status: Order['status'] = 'paid'
      const { container } = render(<OrderStatusBadge status={status} />)

      const label = container.querySelector('span')
      expect(label?.textContent).toBe('Paid')
    })

    it('label text is always capitalized', () => {
      const statuses: Array<{ status: Order['status']; label: string }> = [
        { status: 'pending', label: 'Pending' },
        { status: 'paid', label: 'Paid' },
        { status: 'cancelled', label: 'Cancelled' },
      ]

      statuses.forEach(({ status, label }) => {
        const { container } = render(<OrderStatusBadge status={status} />)
        expect(container.querySelector('span')?.textContent).toBe(label)
      })
    })
  })

  describe('no extra elements', () => {
    it('renders only a single span element', () => {
      const { container } = render(<OrderStatusBadge status="pending" />)

      const spans = container.querySelectorAll('span')
      expect(spans).toHaveLength(1)
    })
  })
})
