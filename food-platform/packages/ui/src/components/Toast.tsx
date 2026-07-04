import { useEffect, useState } from 'react'
import { cn } from '../utils/cn'

export type ToastVariant = 'success' | 'error' | 'warning' | 'info'

export interface ToastProps {
  message: string
  variant?: ToastVariant
  duration?: number
  onClose?: () => void
}

const variantClasses: Record<ToastVariant, string> = {
  success: 'bg-success-light text-success border-success/30',
  error: 'bg-error-light text-error border-error/30',
  warning: 'bg-warning-light text-warning border-warning/30',
  info: 'bg-info-light text-info border-info/30',
}

const variantIcons: Record<ToastVariant, string> = {
  success: '✅',
  error: '❌',
  warning: '⚠️',
  info: 'ℹ️',
}

export function Toast({ message, variant = 'info', duration = 3000, onClose }: ToastProps) {
  const [visible, setVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false)
      onClose?.()
    }, duration)
    return () => clearTimeout(timer)
  }, [duration, onClose])

  if (!visible) return null

  return (
    <div
      className={cn(
        'fixed bottom-4 left-1/2 -translate-x-1/2 z-toast',
        'flex items-center gap-2 px-4 py-3 rounded-lg border shadow-lg',
        'animate-slide-up',
        variantClasses[variant],
      )}
      role="alert"
    >
      <span>{variantIcons[variant]}</span>
      <span className="text-body font-medium">{message}</span>
    </div>
  )
}
