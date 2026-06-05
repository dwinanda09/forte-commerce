import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { CheckoutStatusBadge } from './CheckoutStatusBadge'
import type { CheckoutSession } from '@/lib/types'

describe('CheckoutStatusBadge', () => {
  describe('pending status', () => {
    it('renders processing label for pending status', () => {
      const status: CheckoutSession['status'] = 'pending'
      render(<CheckoutStatusBadge status={status} />)

      expect(screen.getByText('Processing')).toBeInTheDocument()
    })

    it('applies pending background color', () => {
      const status: CheckoutSession['status'] = 'pending'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.bg-yellow-50')
      expect(badge).toBeInTheDocument()
    })

    it('applies pending text color', () => {
      const status: CheckoutSession['status'] = 'pending'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.text-yellow-900')
      expect(badge).toBeInTheDocument()
    })
  })

  describe('completed status', () => {
    it('renders completed label', () => {
      const status: CheckoutSession['status'] = 'completed'
      render(<CheckoutStatusBadge status={status} />)

      expect(screen.getByText('Completed')).toBeInTheDocument()
    })

    it('applies completed background color', () => {
      const status: CheckoutSession['status'] = 'completed'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.bg-teal-50')
      expect(badge).toBeInTheDocument()
    })

    it('applies completed text color', () => {
      const status: CheckoutSession['status'] = 'completed'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.text-teal')
      expect(badge).toBeInTheDocument()
    })
  })

  describe('expired status', () => {
    it('renders expired label', () => {
      const status: CheckoutSession['status'] = 'expired'
      render(<CheckoutStatusBadge status={status} />)

      expect(screen.getByText('Expired')).toBeInTheDocument()
    })

    it('applies expired background color', () => {
      const status: CheckoutSession['status'] = 'expired'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.bg-red-50')
      expect(badge).toBeInTheDocument()
    })

    it('applies expired text color', () => {
      const status: CheckoutSession['status'] = 'expired'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.text-red-900')
      expect(badge).toBeInTheDocument()
    })
  })

  describe('failed status', () => {
    it('renders failed label', () => {
      const status: CheckoutSession['status'] = 'failed'
      render(<CheckoutStatusBadge status={status} />)

      expect(screen.getByText('Failed')).toBeInTheDocument()
    })

    it('applies failed background color', () => {
      const status: CheckoutSession['status'] = 'failed'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.bg-red-50')
      expect(badge).toBeInTheDocument()
    })

    it('applies failed text color', () => {
      const status: CheckoutSession['status'] = 'failed'
      const { container } = render(<CheckoutStatusBadge status={status} />)

      const badge = container.querySelector('.text-red-900')
      expect(badge).toBeInTheDocument()
    })
  })

  describe('visual consistency', () => {
    it('all statuses have padding', () => {
      const statuses: CheckoutSession['status'][] = ['pending', 'completed', 'expired', 'failed']

      statuses.forEach((status) => {
        const { container } = render(<CheckoutStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('px-3')
        expect(badge).toHaveClass('py-1')
      })
    })

    it('all statuses are inline blocks', () => {
      const statuses: CheckoutSession['status'][] = ['pending', 'completed', 'expired', 'failed']

      statuses.forEach((status) => {
        const { container } = render(<CheckoutStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('inline-block')
      })
    })

    it('all statuses have text-sm font-medium', () => {
      const statuses: CheckoutSession['status'][] = ['pending', 'completed', 'expired', 'failed']

      statuses.forEach((status) => {
        const { container } = render(<CheckoutStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('text-sm')
        expect(badge).toHaveClass('font-medium')
      })
    })

    it('all statuses have rounded-full', () => {
      const statuses: CheckoutSession['status'][] = ['pending', 'completed', 'expired', 'failed']

      statuses.forEach((status) => {
        const { container } = render(<CheckoutStatusBadge status={status} />)
        const badge = container.querySelector('span')
        expect(badge).toHaveClass('rounded-full')
      })
    })
  })

  describe('different statuses render correctly', () => {
    it('renders each status independently', () => {
      const { rerender } = render(<CheckoutStatusBadge status="pending" />)
      expect(screen.getByText('Processing')).toBeInTheDocument()

      rerender(<CheckoutStatusBadge status="completed" />)
      expect(screen.getByText('Completed')).toBeInTheDocument()

      rerender(<CheckoutStatusBadge status="expired" />)
      expect(screen.getByText('Expired')).toBeInTheDocument()

      rerender(<CheckoutStatusBadge status="failed" />)
      expect(screen.getByText('Failed')).toBeInTheDocument()
    })
  })

  describe('error statuses consistency', () => {
    it('expired and failed statuses use same colors', () => {
      const { container: expiredContainer } = render(<CheckoutStatusBadge status="expired" />)
      const { container: failedContainer } = render(<CheckoutStatusBadge status="failed" />)

      const expiredBadge = expiredContainer.querySelector('span')
      const failedBadge = failedContainer.querySelector('span')

      expect(expiredBadge).toHaveClass('bg-red-50')
      expect(failedBadge).toHaveClass('bg-red-50')
      expect(expiredBadge).toHaveClass('text-red-900')
      expect(failedBadge).toHaveClass('text-red-900')
    })
  })

  describe('no extra elements', () => {
    it('renders only a single span element', () => {
      const { container } = render(<CheckoutStatusBadge status="pending" />)

      const spans = container.querySelectorAll('span')
      expect(spans).toHaveLength(1)
    })
  })

  describe('label distinctiveness', () => {
    it('pending uses "Processing" not "Pending"', () => {
      render(<CheckoutStatusBadge status="pending" />)

      expect(screen.getByText('Processing')).toBeInTheDocument()
      expect(screen.queryByText('Pending')).not.toBeInTheDocument()
    })

    it('completed is distinct from other labels', () => {
      const { rerender } = render(<CheckoutStatusBadge status="completed" />)
      expect(screen.getByText('Completed')).toBeInTheDocument()

      rerender(<CheckoutStatusBadge status="expired" />)
      expect(screen.queryByText('Completed')).not.toBeInTheDocument()
    })
  })
})
