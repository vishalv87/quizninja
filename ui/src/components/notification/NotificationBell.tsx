"use client";

import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useNotificationStats } from "@/hooks/useNotifications";
import { cn } from "@/lib/utils";

interface NotificationBellProps {
  onClick?: () => void;
  className?: string;
  showBadge?: boolean;
}

export function NotificationBell({
  onClick,
  className,
  showBadge = true,
}: NotificationBellProps) {
  const { data: stats, isLoading } = useNotificationStats();

  const unreadCount = stats?.unread_notifications || 0;
  const hasUnread = unreadCount > 0;

  return (
    <Button
      variant="ghost"
      size="icon"
      className={cn("relative rounded-xl hover:bg-primary/10 transition-all duration-300", className)}
      onClick={onClick}
    >
      <Bell className={cn("h-5 w-5", hasUnread && "text-primary")} />
      <span className="sr-only">
        Notifications {hasUnread ? `(${unreadCount} unread)` : ""}
      </span>

      {/* Unread badge */}
      {showBadge && hasUnread && !isLoading && (
        <Badge
          variant="destructive"
          className="absolute -top-1 -right-1 h-5 min-w-[1.25rem] px-1 flex items-center justify-center text-xs font-semibold rounded-full"
        >
          {unreadCount > 99 ? "99+" : unreadCount}
        </Badge>
      )}

      {/* Animated pulse for new notifications */}
      {showBadge && hasUnread && (
        <span className="absolute -top-1 -right-1 flex h-3 w-3">
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-destructive opacity-75"></span>
        </span>
      )}
    </Button>
  );
}