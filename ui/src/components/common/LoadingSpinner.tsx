import { Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  className?: string
  text?: string
}

const sizeClasses = {
  sm: 'h-4 w-4',
  md: 'h-8 w-8',
  lg: 'h-12 w-12',
  xl: 'h-16 w-16',
}

export function LoadingSpinner({ size = 'md', className, text }: LoadingSpinnerProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-2">
      <Loader2
        className={cn('animate-spin text-primary', sizeClasses[size], className)}
      />
      {text && (
        <p className="text-sm text-muted-foreground">{text}</p>
      )}
    </div>
  )
}

/**
 * Full page loading spinner
 */
export function LoadingPage({ text = 'Loading...' }: { text?: string }) {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <LoadingSpinner size="lg" text={text} />
    </div>
  )
}

/**
 * Loading overlay for sections
 */
export function LoadingOverlay({ text }: { text?: string }) {
  return (
    <div className="absolute inset-0 z-50 flex items-center justify-center bg-background/80 backdrop-blur-sm">
      <LoadingSpinner size="lg" text={text} />
    </div>
  )
}