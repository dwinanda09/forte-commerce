import { type NextRequest, NextResponse } from 'next/server'
import { cookies } from 'next/headers'

const BACKEND_URL = process.env.BACKEND_URL ?? 'http://localhost:8080'

async function handler(req: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const { path } = await params
  const backendPath = path.join('/')
  const url = `${BACKEND_URL}/api/${backendPath}${req.nextUrl.search}`

  const cookieStore = await cookies()
  const token = cookieStore.get('auth_token')?.value

  const headers = new Headers()
  headers.set('Content-Type', 'application/json')
  headers.set('X-Request-ID', crypto.randomUUID())
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const body = req.method !== 'GET' && req.method !== 'HEAD' ? await req.text() : undefined

  const res = await fetch(url, {
    method: req.method,
    headers,
    body: body || undefined,
  })

  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}

export { handler as GET, handler as POST, handler as PUT, handler as DELETE, handler as PATCH }
