import type { Config } from 'tailwindcss'

const config: Config = {
  content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}', './hooks/**/*.{ts,tsx}', './lib/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        steel: '#6D8196',
        mist: '#B0C4DE',
        teal: {
          DEFAULT: '#01796F',
          hover: '#015f58',
        },
        graphite: '#5A5A5A',
        surface: '#B0C4DE',
      },
      borderRadius: {
        sm: '4px',
        md: '8px',
        lg: '12px',
      },
      borderWidth: {
        '1.5': '1.5px',
      },
      boxShadow: {
        card: '0 1px 3px rgba(0,0,0,.08), 0 1px 2px rgba(0,0,0,.05)',
        modal: '0 10px 40px rgba(0,0,0,.12)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      keyframes: {
        'slide-up': {
          '0%': { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
      animation: {
        'slide-up': 'slide-up 0.2s ease-out',
      },
    },
  },
  plugins: [],
}

export default config
