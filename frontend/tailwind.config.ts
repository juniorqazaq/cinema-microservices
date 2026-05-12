import type { Config } from 'tailwindcss'

/** MOVIE HOUSE–style: deep navy, card #101827, blue #3b82f6 (no gradients / blur). */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        page: '#050a18',
        card: '#101827',
        card2: '#1a2332',
        accent: '#3b82f6',
        accentHover: '#2563eb',
        accentDim: '#172554',
        vip: '#38bdf8',
        vipDim: '#0c1a24',
        danger: '#f87171',
        dangerDim: '#241010',
        border: '#1e293b',
        border2: '#334155',
        muted: '#9ca3af',
      },
      fontFamily: {
        sans: ['"DM Sans"', 'system-ui', 'sans-serif'],
      },
      fontSize: {
        body: ['13px', { lineHeight: '1.5' }],
        'card-title': ['15px', { lineHeight: '1.4' }],
        section: ['20px', { lineHeight: '1.3' }],
        hero: ['32px', { lineHeight: '1.2' }],
      },
    },
  },
  plugins: [],
} satisfies Config
