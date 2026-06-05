import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import useSWR from 'swr'
import { useCheckoutPoller } from './useCheckoutPoller'
import { api } from '@/lib/api'

vi.mock('swr')
vi.mock('@/lib/api', () => ({
  api: {
    getCheckout: vi.fn(),
  },
}))

describe('useCheckoutPoller', () => {
  const mockSession = {
    checkout_id: 'chk-123',
    status: 'pending',
    expires_at: new Date(Date.now() + 600000).toISOString(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('with valid checkout id', () => {
    it('fetches checkout data when id is provided', () => {
      const mockUseSWR = vi.fn().mockReturnValue({
        data: mockSession,
        isLoading: false,
        error: undefined,
      })
      vi.mocked(useSWR).mockImplementation(mockUseSWR)

      renderHook(() => useCheckoutPoller('chk-123'))

      expect(mockUseSWR).toHaveBeenCalledWith(
        'chk-123',
        expect.any(Function),
        expect.any(Object)
      )
    })

    it('returns session data when fetch succeeds', () => {
      vi.mocked(useSWR).mockReturnValue({
        data: mockSession,
        isLoading: false,
        error: undefined,
      } as any)

      const { result } = renderHook(() => useCheckoutPoller('chk-123'))

      expect(result.current.session).toEqual(mockSession)
      expect(result.current.isLoading).toBe(false)
      expect(result.current.error).toBeUndefined()
    })

    it('returns loading state during fetch', () => {
      vi.mocked(useSWR).mockReturnValue({
        data: undefined,
        isLoading: true,
        error: undefined,
      } as any)

      const { result } = renderHook(() => useCheckoutPoller('chk-123'))

      expect(result.current.isLoading).toBe(true)
      expect(result.current.session).toBeUndefined()
    })

    it('returns error when fetch fails', () => {
      const error = new Error('Fetch failed')
      vi.mocked(useSWR).mockReturnValue({
        data: undefined,
        isLoading: false,
        error,
      } as any)

      const { result } = renderHook(() => useCheckoutPoller('chk-123'))

      expect(result.current.error).toEqual(error)
      expect(result.current.session).toBeUndefined()
    })
  })

  describe('with null checkout id', () => {
    it('does not fetch when id is null', () => {
      const mockUseSWR = vi.fn().mockReturnValue({
        data: undefined,
        isLoading: false,
        error: undefined,
      })
      vi.mocked(useSWR).mockImplementation(mockUseSWR)

      renderHook(() => useCheckoutPoller(null))

      expect(mockUseSWR).toHaveBeenCalledWith(
        null,
        expect.any(Function),
        expect.any(Object)
      )
    })

    it('returns undefined session when id is null', () => {
      vi.mocked(useSWR).mockReturnValue({
        data: undefined,
        isLoading: false,
        error: undefined,
      } as any)

      const { result } = renderHook(() => useCheckoutPoller(null))

      expect(result.current.session).toBeUndefined()
      expect(result.current.isLoading).toBe(false)
    })
  })

  describe('refresh interval behavior', () => {
    it('polls every 2 seconds when data is undefined', () => {
      const refreshIntervalMock = vi.fn((data) => {
        if (!data) return 2000
        if (data.status === 'pending') return 2000
        return 0
      })

      const config = {
        refreshInterval: refreshIntervalMock,
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
      }

      vi.mocked(useSWR).mockReturnValue({
        data: undefined,
        isLoading: true,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      const refreshInterval = configArg.refreshInterval

      expect(refreshInterval(undefined)).toBe(2000)
    })

    it('polls every 2 seconds when status is pending', () => {
      const config = {
        refreshInterval: (data: any) => {
          if (!data) return 2000
          if (data.status === 'pending') return 2000
          return 0
        },
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
      }

      vi.mocked(useSWR).mockReturnValue({
        data: mockSession,
        isLoading: false,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      const refreshInterval = configArg.refreshInterval

      expect(refreshInterval(mockSession)).toBe(2000)
    })

    it('stops polling when status is completed', () => {
      const completedSession = {
        ...mockSession,
        status: 'completed',
      }

      const config = {
        refreshInterval: (data: any) => {
          if (!data) return 2000
          if (data.status === 'pending') return 2000
          return 0
        },
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
      }

      vi.mocked(useSWR).mockReturnValue({
        data: completedSession,
        isLoading: false,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      const refreshInterval = configArg.refreshInterval

      expect(refreshInterval(completedSession)).toBe(0)
    })

    it('stops polling when status is expired', () => {
      const expiredSession = {
        ...mockSession,
        status: 'expired',
      }

      const config = {
        refreshInterval: (data: any) => {
          if (!data) return 2000
          if (data.status === 'pending') return 2000
          return 0
        },
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
      }

      vi.mocked(useSWR).mockReturnValue({
        data: expiredSession,
        isLoading: false,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      const refreshInterval = configArg.refreshInterval

      expect(refreshInterval(expiredSession)).toBe(0)
    })

    it('stops polling when status is failed', () => {
      const failedSession = {
        ...mockSession,
        status: 'failed',
      }

      const config = {
        refreshInterval: (data: any) => {
          if (!data) return 2000
          if (data.status === 'pending') return 2000
          return 0
        },
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
      }

      vi.mocked(useSWR).mockReturnValue({
        data: failedSession,
        isLoading: false,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      const refreshInterval = configArg.refreshInterval

      expect(refreshInterval(failedSession)).toBe(0)
    })
  })

  describe('SWR configuration', () => {
    it('disables revalidateOnFocus', () => {
      vi.mocked(useSWR).mockReturnValue({
        data: undefined,
        isLoading: false,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      expect(configArg.revalidateOnFocus).toBe(false)
    })

    it('disables revalidateOnReconnect', () => {
      vi.mocked(useSWR).mockReturnValue({
        data: undefined,
        isLoading: false,
        error: undefined,
      } as any)

      renderHook(() => useCheckoutPoller('chk-123'))

      const configArg = vi.mocked(useSWR).mock.calls[0][2]
      expect(configArg.revalidateOnReconnect).toBe(false)
    })
  })

  describe('fetcher function', () => {
    it('calls api.getCheckout with correct id', () => {
      const mockApiResponse = { data: mockSession }
      vi.mocked(api.getCheckout).mockResolvedValue(mockApiResponse)

      const mockUseSWR = vi.fn((key, fetcher) => {
        if (key) {
          fetcher(key)
        }
        return { data: undefined, isLoading: false, error: undefined }
      })

      vi.mocked(useSWR).mockImplementation(mockUseSWR)

      renderHook(() => useCheckoutPoller('chk-123'))

      expect(vi.mocked(api.getCheckout)).toHaveBeenCalledWith('chk-123')
    })

    it('returns data property from api response', async () => {
      const mockApiResponse = { data: mockSession }
      vi.mocked(api.getCheckout).mockResolvedValue(mockApiResponse)

      let fetcherResult: any
      const mockUseSWR = vi.fn((key, fetcher) => {
        if (key) {
          fetcher(key).then((result: any) => {
            fetcherResult = result
          })
        }
        return { data: undefined, isLoading: false, error: undefined }
      })

      vi.mocked(useSWR).mockImplementation(mockUseSWR)

      renderHook(() => useCheckoutPoller('chk-123'))

      await waitFor(() => {
        expect(fetcherResult).toEqual(mockSession)
      }, { timeout: 100 })
    })
  })
})
