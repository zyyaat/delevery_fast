import type { Config } from 'tailwindcss'
import preset from '@food-platform/theme/tailwind.config'

const config: Config = {
  ...preset,
  content: ['./index.html', './src/**/*.{ts,tsx}'],
} as Config

export default config
