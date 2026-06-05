import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderHook, act } from '@testing-library/react'

vi.mock('js-cookie')

describe('useAuth', () => {
  let mockCookies: any

  beforeEach(async () => {
    const cookiesModule = await import('js-cookie')
    mockCookies = cookiesModule.default as any

    mockCookies.set = vi.fn()
    mockCookies.get = vi.fn()
    mockCookies.remove = vi.fn()

    delete (window as any).location
    window.location = { href: '' } as any
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('login', () => {
    it('sets auth token cookie with correct options', async () => {
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      act(() => {
        result.current.login('test-token-123')
      })

      expect(mockCookies.set).toHaveBeenCalledWith('auth_token', 'test-token-123', {
        expires: 1,
        sameSite: 'strict',
      })
    })

    it('accepts various token formats', async () => {
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      const tokens = ['bearer-token', 'jwt.token.here', 'simple']

      for (const token of tokens) {
        mockCookies.set.mockClear()
        act(() => {
          result.current.login(token)
        })

        expect(mockCookies.set).toHaveBeenCalledWith('auth_token', token, {
          expires: 1,
          sameSite: 'strict',
        })
      }
    })
  })

  describe('logout', () => {
    it('removes auth token cookie and redirects to login', async () => {
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      act(() => {
        result.current.logout()
      })

      expect(mockCookies.remove).toHaveBeenCalledWith('auth_token')
      expect(window.location.href).toBe('/login')
    })
  })

  describe('getToken', () => {
    it('retrieves token from cookie', async () => {
      mockCookies.get.mockReturnValueOnce('stored-token')
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      const token = result.current.getToken()

      expect(token).toBe('stored-token')
      expect(mockCookies.get).toHaveBeenCalledWith('auth_token')
    })

    it('returns undefined when token not set', async () => {
      mockCookies.get.mockReturnValueOnce(undefined)
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      const token = result.current.getToken()

      expect(token).toBeUndefined()
    })

    it('handles empty string cookie value', async () => {
      mockCookies.get.mockReturnValueOnce('')
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      const token = result.current.getToken()

      expect(token).toBe('')
    })
  })

  describe('hook lifecycle', () => {
    it('hook returns login, logout, and getToken methods', async () => {
      const { useAuth } = await import('./useAuth')
      const { result } = renderHook(() => useAuth())

      expect(typeof result.current.login).toBe('function')
      expect(typeof result.current.logout).toBe('function')
      expect(typeof result.current.getToken).toBe('function')
    })
  })

  it('handles multiple login calls in sequence', async () => {
    const { useAuth } = await import('./useAuth')
    const { result } = renderHook(() => useAuth())

    act(() => {
      result.current.login('token-1')
    })

    act(() => {
      result.current.login('token-2')
    })

    act(() => {
      result.current.login('token-3')
    })

    expect(mockCookies.set).toHaveBeenCalledTimes(3)
    expect(mockCookies.set.mock.calls[2]).toEqual(['auth_token', 'token-3', { expires: 1, sameSite: 'strict' }])
  })
})
