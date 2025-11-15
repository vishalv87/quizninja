"use client";

import { NotificationCard } from "./NotificationCard";
import { Skeleton } from "@/components/ui/skeleton";
import { Bell } from "lucide-react";
import type { Notification } from "@/types/notification";
import { groupNotificationsByDate } from "@/lib/notification-utils";

interface NotificationListProps {
  notifications: Notification[];
  isLoading?: boolean;
  emptyMessage?: string;
  groupByDate?: boolean;
}

export function NotificationList({
  notifications,
  isLoading = false,
  emptyMessage = "No notifications yet",
  groupByDate = false,
}: NotificationListProps) {
  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(5)].map((_, i) => (
          <NotificationSkeleton key={i} />
        ))}
      </div>
    );
  }

  // Empty state
  if (!notifications || notifications.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="rounded-full bg-muted p-6 mb-4">
          <Bell className="h-12 w-12 text-muted-foreground" />
        </div>
        <h3 className="text-lg font-semibold mb-2">No Notifications</h3>
        <p className="text-sm text-muted-foreground max-w-sm">
          {emptyMessage}
        </p>
      </div>
    );
  }

  // Group by date if requested
  if (groupByDate) {
    const grouped = groupNotificationsByDate(notifications);

    return (
      <div className="space-y-6">
        {grouped.today.length > 0 && (
          <div>
            <h3 className="text-sm font-semibold text-muted-foreground mb-3 px-2">
              Today
            </h3>
            <div className="space-y-3">
              {grouped.today.map((notification) => (
                <NotificationCard
                  key={notification.id}
                  notification={notification}
                />
              ))}
            </div>
          </div>
        )}

        {grouped.yesterday.length > 0 && (
          <div>
            <h3 className="text-sm font-semibold text-muted-foreground mb-3 px-2">
              Yesterday
            </h3>
            <div className="space-y-3">
              {grouped.yesterday.map((notification) => (
                <NotificationCard
                  key={notification.id}
                  notification={notification}
                />
              ))}
            </div>
          </div>
        )}

        {grouped.thisWeek.length > 0 && (
          <div>
            <h3 className="text-sm font-semibold text-muted-foreground mb-3 px-2">
              This Week
            </h3>
            <div className="space-y-3">
              {grouped.thisWeek.map((notification) => (
                <NotificationCard
                  key={notification.id}
                  notification={notification}
                />
              ))}
            </div>
          </div>
        )}

        {grouped.earlier.length > 0 && (
          <div>
            <h3 className="text-sm font-semibold text-muted-foreground mb-3 px-2">
              Earlier
            </h3>
            <div className="space-y-3">
              {grouped.earlier.map((notification) => (
                <NotificationCard
                  key={notification.id}
                  notification={notification}
                />
              ))}
            </div>
          </div>
        )}
      </div>
    );
  }

  // Regular list (no grouping)
  return (
    <div className="space-y-3">
      {notifications.map((notification) => (
        <NotificationCard key={notification.id} notification={notification} />
      ))}
    </div>
  );
}

function NotificationSkeleton() {
  return (
    <div className="flex gap-3 p-4 border rounded-lg">
      <Skeleton className="h-12 w-12 rounded-full flex-shrink-0" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-1/4" />
      </div>
    </div>
  );
}