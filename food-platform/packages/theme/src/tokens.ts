// Design tokens — single source of truth for all styling

export const tokens = {
  // ============ Colors (Light Mode) ============
  colors: {
    // Primary brand
    primary: '#FF5722',
    primaryDark: '#E64A19',
    primaryLight: '#FFAB91',

    // Secondary
    secondary: '#00897B',
    secondaryDark: '#00695C',
    secondaryLight: '#B2DFDB',

    // Semantic
    success: '#2E7D32',
    successLight: '#C8E6C9',
    warning: '#F9A825',
    warningLight: '#FFF9C4',
    error: '#D32F2F',
    errorLight: '#FFCDD2',
    info: '#0288D1',
    infoLight: '#B3E5FC',

    // Neutrals (warm Egyptian aesthetic)
    bgPrimary: '#FAFAFA',
    bgSecondary: '#F5F5F5',
    bgTertiary: '#EEEEEE',
    surface: '#FFFFFF',

    // Borders
    border: '#E0E0E0',
    borderStrong: '#BDBDBD',

    // Text
    textPrimary: '#212121',
    textSecondary: '#616161',
    textTertiary: '#9E9E9E',
    textDisabled: '#BDBDBD',
  },

  // ============ Dark Mode Colors ============
  darkColors: {
    bgPrimary: '#0A0E1A',
    bgSecondary: '#131A2E',
    bgTertiary: '#1A2240',
    surface: '#1A2240',
    border: '#2A3656',
    borderStrong: '#3A4A6E',
    textPrimary: '#EAF0FF',
    textSecondary: '#9BA8C7',
    textTertiary: '#5C6B8E',
    accent: '#00D4FF',
    accent2: '#B14EFF',
  },

  // ============ Typography ============
  typography: {
    fontFamily: {
      sans: 'Cairo, Inter, sans-serif',
      en: 'Inter, sans-serif',
      ar: 'Cairo, sans-serif',
      mono: 'JetBrains Mono, monospace',
    },
    fontSize: {
      displayLg: '48px',
      displayMd: '40px',
      displaySm: '32px',
      h1: '28px',
      h2: '24px',
      h3: '20px',
      h4: '18px',
      bodyLg: '16px',
      body: '14px',
      bodySm: '13px',
      caption: '12px',
      overline: '11px',
      monoLg: '16px',
      mono: '14px',
      monoSm: '12px',
    },
    lineHeight: {
      displayLg: '56px',
      displayMd: '48px',
      displaySm: '40px',
      h1: '36px',
      h2: '32px',
      h3: '28px',
      h4: '24px',
      bodyLg: '24px',
      body: '20px',
      bodySm: '18px',
      caption: '16px',
      overline: '16px',
    },
    fontWeight: {
      regular: 400,
      medium: 500,
      semibold: 600,
      bold: 700,
      extrabold: 800,
    },
  },

  // ============ Spacing ============
  spacing: {
    0: '0',
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

  // ============ Border Radius ============
  borderRadius: {
    none: '0',
    sm: '4px',
    md: '8px',
    lg: '12px',
    xl: '16px',
    '2xl': '24px',
    full: '9999px',
  },

  // ============ Shadows ============
  shadows: {
    none: 'none',
    sm: '0 1px 2px rgba(0,0,0,0.05)',
    md: '0 2px 8px rgba(0,0,0,0.08)',
    lg: '0 8px 24px rgba(0,0,0,0.12)',
    xl: '0 16px 48px rgba(0,0,0,0.16)',
    '2xl': '0 24px 64px rgba(0,0,0,0.20)',
  },

  // ============ Animation ============
  animation: {
    durationFast: '100ms',
    durationNormal: '200ms',
    durationSlow: '300ms',
    durationSlower: '450ms',
    durationSlowest: '600ms',
    easeStandard: 'cubic-bezier(0.4, 0, 0.2, 1)',
    easeDecelerate: 'cubic-bezier(0, 0, 0.2, 1)',
    easeAccelerate: 'cubic-bezier(0.4, 0, 1, 1)',
    easeSpring: 'cubic-bezier(0.5, 1.5, 0.5, 1)',
  },

  // ============ Breakpoints ============
  breakpoints: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px',
  },

  // ============ Z-Index ============
  zIndex: {
    base: 0,
    dropdown: 1000,
    sticky: 1100,
    fixed: 1200,
    modalBackdrop: 1300,
    modal: 1400,
    popover: 1500,
    toast: 1600,
    tooltip: 1700,
  },

  // ============ Container ============
  container: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    full: '1440px',
  },
} as const

export type Tokens = typeof tokens
export type ColorTokens = typeof tokens.colors
export type TypographyTokens = typeof tokens.typography
export type SpacingTokens = typeof tokens.spacing
export type AnimationTokens = typeof tokens.animation
