// Tailwind config — shared across all apps
// Import this in each app's tailwind.config.ts: import preset from '@food-platform/theme/tailwind.config'

import type { Config } from 'tailwindcss'
import rtl from 'tailwindcss-rtl'

const config: Config = {
  darkMode: 'class',
  content: ['./src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#FF5722',
          dark: '#E64A19',
          light: '#FFAB91',
        },
        secondary: {
          DEFAULT: '#00897B',
          dark: '#00695C',
          light: '#B2DFDB',
        },
        success: {
          DEFAULT: '#2E7D32',
          light: '#C8E6C9',
        },
        warning: {
          DEFAULT: '#F9A825',
          light: '#FFF9C4',
        },
        error: {
          DEFAULT: '#D32F2F',
          light: '#FFCDD2',
        },
        info: {
          DEFAULT: '#0288D1',
          light: '#B3E5FC',
        },
        bg: {
          primary: '#FAFAFA',
          secondary: '#F5F5F5',
          tertiary: '#EEEEEE',
        },
        surface: '#FFFFFF',
        border: {
          DEFAULT: '#E0E0E0',
          strong: '#BDBDBD',
        },
        text: {
          primary: '#212121',
          secondary: '#616161',
          tertiary: '#9E9E9E',
          disabled: '#BDBDBD',
        },
        // Dark mode accents (for command center, employee portal)
        accent: {
          cyan: '#00D4FF',
          purple: '#B14EFF',
        },
      },
      fontFamily: {
        sans: ['Cairo', 'Inter', 'sans-serif'],
        en: ['Inter', 'sans-serif'],
        ar: ['Cairo', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      fontSize: {
        'display-lg': ['48px', { lineHeight: '56px', fontWeight: '800' }],
        'display-md': ['40px', { lineHeight: '48px', fontWeight: '800' }],
        'display-sm': ['32px', { lineHeight: '40px', fontWeight: '700' }],
        h1: ['28px', { lineHeight: '36px', fontWeight: '700' }],
        h2: ['24px', { lineHeight: '32px', fontWeight: '700' }],
        h3: ['20px', { lineHeight: '28px', fontWeight: '600' }],
        h4: ['18px', { lineHeight: '24px', fontWeight: '600' }],
        'body-lg': ['16px', { lineHeight: '24px' }],
        body: ['14px', { lineHeight: '20px' }],
        'body-sm': ['13px', { lineHeight: '18px' }],
        caption: ['12px', { lineHeight: '16px', fontWeight: '500' }],
        overline: ['11px', { lineHeight: '16px', fontWeight: '600' }],
        'mono-lg': ['16px', { lineHeight: '24px', fontWeight: '600' }],
        mono: ['14px', { lineHeight: '20px', fontWeight: '500' }],
        'mono-sm': ['12px', { lineHeight: '16px', fontWeight: '500' }],
      },
      spacing: {
        1: '4px',
        2: '8px',
        3: '12px',
        4: '16px',
        5: '20px',
        6: '24px',
        8: '32px',
        10: '40px',
        12: '48px',
        16: '64px',
        20: '80px',
        24: '96px',
      },
      borderRadius: {
        sm: '4px',
        md: '8px',
        lg: '12px',
        xl: '16px',
        '2xl': '24px',
        full: '9999px',
      },
      boxShadow: {
        sm: '0 1px 2px rgba(0,0,0,0.05)',
        md: '0 2px 8px rgba(0,0,0,0.08)',
        lg: '0 8px 24px rgba(0,0,0,0.12)',
        xl: '0 16px 48px rgba(0,0,0,0.16)',
        '2xl': '0 24px 64px rgba(0,0,0,0.20)',
      },
      transitionDuration: {
        fast: '100ms',
        normal: '200ms',
        slow: '300ms',
        slower: '450ms',
        slowest: '600ms',
      },
      transitionTimingFunction: {
        standard: 'cubic-bezier(0.4, 0, 0.2, 1)',
        decelerate: 'cubic-bezier(0, 0, 0.2, 1)',
        accelerate: 'cubic-bezier(0.4, 0, 1, 1)',
        spring: 'cubic-bezier(0.5, 1.5, 0.5, 1)',
      },
      zIndex: {
        dropdown: '1000',
        sticky: '1100',
        fixed: '1200',
        'modal-backdrop': '1300',
        modal: '1400',
        popover: '1500',
        toast: '1600',
        tooltip: '1700',
      },
      maxWidth: {
        container: {
          sm: '640px',
          md: '768px',
          lg: '1024px',
          xl: '1280px',
          full: '1440px',
        },
      },
      screens: {
        sm: '640px',
        md: '768px',
        lg: '1024px',
        xl: '1280px',
        '2xl': '1536px',
      },
      keyframes: {
        'fade-in': {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'fade-out': {
          '0%': { opacity: '1' },
          '100%': { opacity: '0' },
        },
        'slide-up': {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        'slide-down': {
          '0%': { transform: 'translateY(-20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        'scale-in': {
          '0%': { transform: 'scale(0.95)', opacity: '0' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-1000px 0' },
          '100%': { backgroundPosition: '1000px 0' },
        },
        pulse: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.5' },
        },
        shake: {
          '0%, 100%': { transform: 'translateX(0)' },
          '25%': { transform: 'translateX(-5px)' },
          '75%': { transform: 'translateX(5px)' },
        },
      },
      animation: {
        'fade-in': 'fade-in 200ms cubic-bezier(0, 0, 0.2, 1)',
        'fade-out': 'fade-out 150ms cubic-bezier(0.4, 0, 1, 1)',
        'slide-up': 'slide-up 300ms cubic-bezier(0, 0, 0.2, 1)',
        'slide-down': 'slide-down 300ms cubic-bezier(0, 0, 0.2, 1)',
        'scale-in': 'scale-in 200ms cubic-bezier(0, 0, 0.2, 1)',
        shimmer: 'shimmer 1500ms linear infinite',
        pulse: 'pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shake: 'shake 300ms cubic-bezier(0.4, 0, 0.2, 1)',
      },
    },
  },
  plugins: [rtl],
}

export default config
