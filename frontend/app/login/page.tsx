'use client'

import Image from 'next/image'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { api } from '@/lib/api'
import { useAuth } from '@/hooks/useAuth'

export default function LoginPage() {
  const router = useRouter()
  const { login } = useAuth()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setIsLoading(true)

    try {
      const res = await api.login(username, password)
      login(res.data.token)
      router.replace('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-surface flex items-center justify-center px-6 py-12">
      <div className="bg-white rounded-lg shadow-modal max-w-sm w-full p-8">
        {/* Logo */}
        <div className="flex items-center justify-center gap-3 mb-8">
          <Image src="/fc-logo.png" alt="ForteCommerce" width={40} height={40} priority />
          <span className="font-bold text-xl text-graphite">ForteCommerce</span>
        </div>

        <h1 className="text-2xl font-bold text-graphite text-center mb-8">
          Sign In
        </h1>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
            <p className="text-red-700 text-sm">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          {/* Username */}
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-graphite mb-2"
            >
              Username
            </label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className="w-full px-4 py-2 border border-mist rounded-md text-graphite placeholder-steel focus-visible:outline-2 focus-visible:outline-teal bg-white"
              placeholder="Enter your username"
            />
          </div>

          {/* Password */}
          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-graphite mb-2"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-4 py-2 border border-mist rounded-md text-graphite placeholder-steel focus-visible:outline-2 focus-visible:outline-teal bg-white"
              placeholder="Enter your password"
            />
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-teal text-white py-2.5 rounded-md font-semibold hover:bg-teal-hover transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal mt-6"
          >
            {isLoading ? 'Signing In...' : 'Sign In'}
          </button>
        </form>

        {/* Demo Credentials */}
        <div className="mt-8 pt-8 border-t border-mist">
          <p className="text-xs text-steel text-center mb-3">Demo Credentials:</p>
          <div className="bg-surface rounded-md p-3 space-y-1 text-xs font-mono text-graphite">
            <p>Buyer — username: <span className="text-teal font-semibold">demo1</span> / password: <span className="text-teal font-semibold">demo1</span></p>
            <p>Seller — username: <span className="text-teal font-semibold">seller1</span> / password: <span className="text-teal font-semibold">seller1</span></p>
          </div>
        </div>
      </div>
    </div>
  )
}
