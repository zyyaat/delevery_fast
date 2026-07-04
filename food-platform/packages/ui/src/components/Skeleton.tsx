// Skeleton component (shimmer loading)

import { cn } from '../utils/cn'

export interface SkeletonProps {
  className?: string
  width?: string | number
  height?: string | number
  rounded?: 'sm' | 'md' | 'lg' | 'full'
}

const roundedClasses = {
  sm: 'rounded-sm',
  md: 'rounded-md',
  lg: 'rounded-lg',
  full: 'rounded-full',
}

export function Skeleton({ className, width, height, rounded = 'md' }: SkeletonProps) {
  return (
    <div
      className={cn('skeleton', roundedClasses[rounded], className)}
      style={{ width, height }}
      aria-hidden="true"
    />
  )
}
