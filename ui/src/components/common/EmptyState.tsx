import { LucideIcon, Inbox } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description?: string
  action?: {
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
    <Card className={className}>
      <div className="flex flex-col items-center justify-center p-8 text-center">
        <div className="rounded-full bg-muted p-4">
          <Icon className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="mt-4 text-lg font-semibold">{title}</h3>
        {description && (
          <p className="mt-2 text-sm text-muted-foreground">{description}</p>
        )}
        {action && (
          <Button onClick={action.onClick} className="mt-4">
            {action.label}
          </Button>
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
      <Icon className="h-12 w-12 text-muted-foreground" />
      <h4 className="mt-3 text-sm font-semibold">{title}</h4>
      {description && (
        <p className="mt-1 text-xs text-muted-foreground">{description}</p>
      )}
    </div>
  )
}