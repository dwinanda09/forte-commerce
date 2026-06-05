import type { NextConfig } from 'next'
import withPWAInit from '@ducanh2912/next-pwa'

const withPWA = withPWAInit({
  dest: 'public',
  disable: true,
})

const config: NextConfig = {
  images: {
    unoptimized: true,
  },
}

export default withPWA(config)
