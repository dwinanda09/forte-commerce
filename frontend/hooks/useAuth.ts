'use client'

import { useState, useEffect } from 'react'
import Cookies from 'js-cookie'

function parseRole(token: string): string {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role ?? 'buyer'
  } catch {
    return 'buyer'
  }
}

export function useAuth() {
  const [role, setRole] = useState<string>('buyer')

  useEffect(() => {
    const token = Cookies.get('auth_token')
    if (token) {
      setRole(parseRole(token))
    }
  }, [])

  function login(token: string) {
    Cookies.set('auth_token', token, { expires: 1, sameSite: 'strict' })
    setRole(parseRole(token))
  }

  function logout() {
    Cookies.remove('auth_token')
    setRole('buyer')
    window.location.href = '/login'
  }

  function getToken() {
    return Cookies.get('auth_token')
  }

  return { login, logout, getToken, role }
}
