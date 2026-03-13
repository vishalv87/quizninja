import { LucideIcon, Inbox } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { ReactNode } from 'react'

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description?: string
  action?: ReactNode | {
    label: string
    onClick: () => void
  }
  className?: string
}

export function EmptyState({
  icon: Icon = Inbox,
  title,
  description,
  action,
  className,
}: EmptyStateProps) {
  return (
    <Card className={`border-none shadow-sm bg-white/50 dark:bg-background/50 backdrop-blur-sm rounded-xl ${className || ''}`}>
      <div className="flex flex-col items-center justify-center p-8 text-center">
        <div className="rounded-2xl bg-gradient-to-br from-violet-500/10 to-indigo-500/10 p-4 shadow-sm">
          <Icon className="h-8 w-8 text-violet-600 dark:text-violet-400" />
        </div>
        <h3 className="mt-4 text-lg font-semibold">{title}</h3>
        {description && (
          <p className="mt-2 text-sm text-muted-foreground">{description}</p>
        )}
        {action && (
          <div className="mt-4">
            {typeof action === 'object' && 'label' in action ? (
              <Button onClick={action.onClick}>
                {action.label}
              </Button>
            ) : (
              action
            )}
          </div>
        )}
      </div>
    </Card>
  )
}

/**
 * Compact empty state for smaller sections
 */
export function EmptyStateCompact({
  icon: Icon = Inbox,
  title,
  description,
}: Omit<EmptyStateProps, 'action' | 'className'>) {
  return (
    <div className="flex flex-col items-center justify-center p-6 text-center">
      <div className="rounded-xl bg-gradient-to-br from-violet-500/10 to-indigo-500/10 p-3">
        <Icon className="h-8 w-8 text-violet-600 dark:text-violet-400" />
      </div>
      <h4 className="mt-3 text-sm font-semibold">{title}</h4>
      {description && (
        <p className="mt-1 text-xs text-muted-foreground">{description}</p>
      )}
    </div>
  )
}