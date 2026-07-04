// App layout with header + bottom navigation

import { type ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { cn } from '@food-platform/ui'

interface AppLayoutProps {
  children: ReactNode
  showHeader?: boolean
  showBottomNav?: boolean
}

export function AppLayout({
  children,
  showHeader = true,
  showBottomNav = true,
}: AppLayoutProps) {
  return (
    <div className="min-h-screen bg-bg-primary flex flex-col">
      {showHeader && <Header />}
      <main className="flex-1 pb-20 md:pb-0">{children}</main>
      {showBottomNav && <BottomNav />}
    </div>
  )
}

// ============ Header ============

function Header() {
  return (
    <header className="sticky top-0 z-sticky bg-surface border-b border-border h-14 flex items-center px-4 gap-3">
      {/* Address */}
      <Link to="/" className="flex items-center gap-2 flex-1 min-w-0">
        <span className="material-symbols-rounded text-primary text-xl">location_on</span>
        <div className="flex flex-col min-w-0">
          <span className="text-overline text-text-tertiary">التوصيل إلى</span>
          <span className="text-body-sm font-semibold text-text-primary truncate">
            الزمالك، القاهرة
          </span>
        </div>
      </Link>

      {/* Search */}
      <Link
        to="/search"
        className="flex items-center justify-center w-10 h-10 rounded-full hover:bg-bg-tertiary transition-colors"
        aria-label="بحث"
      >
        <span className="material-symbols-rounded text-text-secondary">search</span>
      </Link>

      {/* Profile */}
      <Link
        to="/profile"
        className="flex items-center justify-center w-10 h-10 rounded-full bg-primary/10 hover:bg-primary/20 transition-colors"
        aria-label="حسابي"
      >
        <span className="material-symbols-rounded text-primary">person</span>
      </Link>
    </header>
  )
}

// ============ Bottom Navigation ============

function BottomNav() {
  const location = useLocation()

  const navItems = [
    { path: '/', label: 'الرئيسية', icon: 'home' },
    { path: '/search', label: 'بحث', icon: 'search' },
    { path: '/orders', label: 'طلباتي', icon: 'receipt_long' },
    { path: '/profile', label: 'حسابي', icon: 'person' },
  ]

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-sticky bg-surface border-t border-border h-16 flex items-center justify-around px-2 md:hidden">
      {navItems.map((item) => {
        const isActive =
          item.path === '/'
            ? location.pathname === '/'
            : location.pathname.startsWith(item.path)

        return (
          <Link
            key={item.path}
            to={item.path}
            className={cn(
              'flex flex-col items-center justify-center gap-1 flex-1 h-full transition-colors',
              isActive ? 'text-primary' : 'text-text-tertiary',
            )}
          >
            <span className="material-symbols-rounded text-xl">{item.icon}</span>
            <span className="text-overline font-medium">{item.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}
