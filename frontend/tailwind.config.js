/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        navy: {
          'primary': '#050505',
          'secondary': '#0C0C0C',
          'tertiary': '#121212',
          'hover': '#1C1C1C',
          'border': '#202020',
        },
        teal: {
          'accent': '#00C9A7',
          'hover': '#00E6BF',
          'dark': '#009A82',
        },
        accent: {
          blue: '#3B82F6',
          amber: '#F59E0B',
          red: '#EF4444',
          green: '#10B981',
        },
        text: {
          primary: '#E2E8F0',
          secondary: '#94A3B8',
          muted: '#4B5563',
        },
      },
      fontFamily: {
        'ui': ['Inter', 'system-ui', 'sans-serif'],
        'mono': ['DM Mono', 'Menlo', 'monospace'],
      },
      fontSize: {
        'xs': ['0.75rem', { lineHeight: '1rem' }],
        'sm': ['0.875rem', { lineHeight: '1.25rem' }],
        'base': ['1rem', { lineHeight: '1.5rem' }],
        'lg': ['1.125rem', { lineHeight: '1.75rem' }],
        'xl': ['1.25rem', { lineHeight: '1.75rem' }],
      },
      animation: {
        'fade-in': 'fadeIn 0.2s ease-in-out',
        'slide-in': 'slideIn 0.2s ease-out',
        'pulse-subtle': 'pulseSubtle 2s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideIn: {
          '0%': { transform: 'translateX(-10px)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        },
        pulseSubtle: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.7' },
        },
      },
    },
  },
  plugins: [],
}
