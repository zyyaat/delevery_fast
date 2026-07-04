#!/bin/bash
# Script to create remaining 6 web app skeletons

set -e

APPS_ROOT="/home/z/my-project/food-platform/apps"
PORT=5174

declare -A APP_NAMES=(
  ["driver-web"]="تطبيق المندوب"
  ["restaurant-web"]="تطبيق المطعم"
  ["support-web"]="تطبيق الدعم"
  ["command-center"]="مركز القيادة"
  ["employee-portal"]="بوابة الموظفين"
  ["field-supervisor-web"]="تطبيق المشرف الميداني"
)

for app in driver-web restaurant-web support-web command-center employee-portal field-supervisor-web; do
  NAME="${APP_NAMES[$app]}"
  echo "=== Creating $app (port $PORT) ==="

  APP_DIR="$APPS_ROOT/$app"

  # package.json
  cat > "$APP_DIR/package.json" << EOF
{
  "name": "@food-platform/$app",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite --port $PORT",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "test": "vitest run",
    "lint": "eslint src --ext .ts,.tsx",
    "type-check": "tsc --noEmit",
    "clean": "rm -rf dist node_modules/.vite"
  },
  "dependencies": {
    "@food-platform/api-client": "workspace:*",
    "@food-platform/auth": "workspace:*",
    "@food-platform/hooks": "workspace:*",
    "@food-platform/types": "workspace:*",
    "@food-platform/ui": "workspace:*",
    "@food-platform/utils": "workspace:*",
    "@food-platform/theme": "workspace:*",
    "@food-platform/eslint-config": "workspace:*",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "react-router-dom": "^6.24.0",
    "@tanstack/react-query": "^5.51.0",
    "zustand": "^4.5.0",
    "axios": "^1.7.0",
    "zod": "^3.23.0",
    "react-hook-form": "^7.52.0",
    "@hookform/resolvers": "^3.9.0",
    "tailwindcss": "^3.4.0",
    "tailwindcss-rtl": "^0.9.0",
    "clsx": "^2.1.0",
    "tailwind-merge": "^2.4.0",
    "date-fns": "^3.6.0",
    "@sentry/react": "^8.20.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.0",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.0",
    "vite": "^5.3.0",
    "vitest": "^2.0.0",
    "@testing-library/react": "^16.0.0",
    "@testing-library/jest-dom": "^6.4.0",
    "@testing-library/user-event": "^14.5.0",
    "@playwright/test": "^1.45.0",
    "eslint": "^8.57.0",
    "typescript": "^5.5.0",
    "@typescript-eslint/eslint-plugin": "^7.0.0",
    "@typescript-eslint/parser": "^7.0.0",
    "eslint-plugin-react": "^7.34.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-import": "^2.29.0",
    "eslint-config-prettier": "^9.1.0",
    "eslint-import-resolver-typescript": "^3.6.0"
  }
}
EOF

  # tsconfig.json
  cat > "$APP_DIR/tsconfig.json" << 'EOF'
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "outDir": "./dist",
    "rootDir": "./src",
    "jsx": "react-jsx",
    "types": ["vite/client", "vitest/globals"]
  },
  "include": ["src", "vite.config.ts"],
  "exclude": ["dist", "node_modules"]
}
EOF

  # vite.config.ts
  cat > "$APP_DIR/vite.config.ts" << EOF
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: $PORT,
    host: true,
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
})
EOF

  # .eslintrc.cjs
  cat > "$APP_DIR/.eslintrc.cjs" << 'EOF'
module.exports = {
  root: true,
  extends: ['@food-platform/eslint-config'],
  parserOptions: {
    project: ['./tsconfig.json'],
  },
}
EOF

  # .env.example
  cat > "$APP_DIR/.env.example" << EOF
# $NAME
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8081
VITE_APP_NAME=$app
EOF

  # index.html
  cat > "$APP_DIR/index.html" << EOF
<!DOCTYPE html>
<html lang="ar" dir="rtl">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
    <meta name="theme-color" content="#FF5722" />
    <title>$NAME — Food Platform</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
EOF

  # src/main.tsx
  cat > "$APP_DIR/src/main.tsx" << 'EOF'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import { initApiClient } from '@food-platform/api-client'
import { zustandTokenStorage } from '@food-platform/auth'
import '@food-platform/theme/globals.css'

const apiURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

initApiClient({
  baseURL: apiURL,
  tokenStorage: zustandTokenStorage,
  refreshTokenUrl: '/api/v1/auth/refresh',
})

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
)
EOF

  # src/App.tsx (placeholder)
  cat > "$APP_DIR/src/App.tsx" << EOF
export default function App() {
  return (
    <div className="min-h-screen bg-bg-primary flex items-center justify-center p-4">
      <div className="text-center">
        <h1 className="text-h1 font-bold text-text-primary mb-4">$NAME</h1>
        <p className="text-body text-text-secondary">Skeleton جاهز — هتتعمل في Phase حسب الـ roadmap</p>
        <p className="text-caption text-text-tertiary mt-4">راجع docs/ui-ux/$app/UI-SPEC.md للتفاصيل</p>
      </div>
    </div>
  )
}
EOF

  # tailwind.config.ts
  cat > "$APP_DIR/tailwind.config.ts" << 'EOF'
import type { Config } from 'tailwindcss'
import preset from '@food-platform/theme/tailwind.config'

const config: Config = {
  ...preset,
  content: ['./index.html', './src/**/*.{ts,tsx}'],
} as Config

export default config
EOF

  # README.md
  cat > "$APP_DIR/README.md" << EOF
# $NAME

Skeleton ready. See [UI Spec](../../docs/ui-ux/$app/UI-SPEC.md) for detailed specs.

## Quick Start

\`\`\`bash
cp .env.example .env
pnpm dev
\`\`\`

App runs at http://localhost:$PORT
EOF

  echo "  ✓ $app done"
  PORT=$((PORT + 1))
done

echo ""
echo "✅ All 6 web app skeletons created!"
