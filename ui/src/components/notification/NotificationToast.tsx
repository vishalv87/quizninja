"use client";

import Link from "next/link";
import { toast } from "sonner";
import { X } from "lucide-react";
import { Button } from "@/components/ui/button";
import type { Notification } from "@/types/notification";
import {
  getNotificationIcon,
  getNotificationColor,
  formatNotificationTime,
  getNotificationActionUrl,
} from "@/lib/notification-utils";
import { cn } from "@/lib/utils";

interface NotificationToastContentProps {
  notification: Notification;
  onDismiss?: () => void;
}

/**
 * Custom toast content for new notifications
 * Displays rich notification information with type-specific styling
 */
function NotificationToastContent({ notification, onDismiss }: NotificationToastContentProps) {
  const Icon = getNotificationIcon(notification.type);
  const colors = getNotificationColor(notification.type);
  const actionUrl = getNotificationActionUrl(notification.type, notification.data);
  const timeAgo = formatNotificationTime(notification.created_at);

  const content = (
    <div className="flex items-start gap-3 w-full">
      {/* Icon with type-specific colors */}
      <div className={cn("flex-shrink-0 p-2 rounded-full", colors.bg)}>
        <Icon className={cn("h-5 w-5", colors.icon)} />
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2 mb-1">
          <p className="font-semibold text-sm">New Notification</p>
          {onDismiss && (
            <Button
              variant="ghost"
              size="icon"
              className="h-5 w-5 -mr-1 -mt-1"
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onDismiss();
              }}
            >
              <X className="h-3 w-3" />
              <span className="sr-only">Dismiss</span>
            </Button>
          )}
        </div>
        <p className="font-medium text-sm leading-tight">{notification.title}</p>
        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
          {notification.message}
        </p>
        <p className="text-xs text-muted-foreground mt-1">{timeAgo}</p>
      </div>
    </div>
  );

  // If there's an action URL, wrap in a Link
  if (actionUrl) {
    return (
      <Link href={actionUrl} className="block hover:opacity-90 transition-opacity">
        {content}
      </Link>
    );
  }

  return content;
}

/**
 * Shows a custom toast notification for new notifications
 * Includes type-specific icon and colors, notification details, and optional action link
 *
 * @param notification - The new notification to display
 */
export function showNotificationToast(notification: Notification) {
  toast.custom(
    (t) => (
      <div className="bg-background border border-border rounded-lg shadow-lg p-4 max-w-md w-full">
        <NotificationToastContent
          notification={notification}
          onDismiss={() => toast.dismiss(t)}
        />
      </div>
    ),
    {
      duration: 5000,
      position: "top-right",
    }
  );
}

/**
 * Shows a summary toast when multiple notifications arrive at once
 * Instead of showing individual toasts, shows a single toast with the count
 *
 * @param count - Number of new notifications
 * @param latestNotification - The most recent notification (optional)
 */
export function showMultipleNotificationsToast(count: number, latestNotification?: Notification) {
  if (count === 1 && latestNotification) {
    // If only one notification, show the full toast
    showNotificationToast(latestNotification);
    return;
  }

  // For multiple notifications, show a summary toast
  toast.custom(
    (t) => (
      <div className="bg-background border border-border rounded-lg shadow-lg p-4 max-w-md w-full">
        <div className="flex items-start gap-3 w-full">
          {latestNotification ? (
            <>
              {/* Show latest notification icon */}
              <div
                className={cn(
                  "flex-shrink-0 p-2 rounded-full",
                  getNotificationColor(latestNotification.type).bg
                )}
              >
                {(() => {
                  const Icon = getNotificationIcon(latestNotification.type);
                  return (
                    <Icon
                      className={cn(
                        "h-5 w-5",
                        getNotificationColor(latestNotification.type).icon
                      )}
                    />
                  );
                })()}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-2 mb-1">
                  <p className="font-semibold text-sm">{count} New Notifications</p>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-5 w-5 -mr-1 -mt-1"
                    onClick={(e) => {
                      e.preventDefault();
                      toast.dismiss(t);
                    }}
                  >
                    <X className="h-3 w-3" />
                    <span className="sr-only">Dismiss</span>
                  </Button>
                </div>
                <p className="font-medium text-sm leading-tight">
                  {latestNotification.title}
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  and {count - 1} more
                </p>
              </div>
            </>
          ) : (
            <>
              {/* Generic icon when no notification provided */}
              <div className="flex-shrink-0 p-2 rounded-full bg-primary/10">
                <div className="h-5 w-5 rounded-full bg-primary" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-2 mb-1">
                  <p className="font-semibold text-sm">{count} New Notifications</p>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-5 w-5 -mr-1 -mt-1"
                    onClick={(e) => {
                      e.preventDefault();
                      toast.dismiss(t);
                    }}
                  >
                    <X className="h-3 w-3" />
                    <span className="sr-only">Dismiss</span>
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  You have {count} new notifications
                </p>
              </div>
            </>
          )}
        </div>
      </div>
    ),
    {
      duration: 5000,
      position: "top-right",
    }
  );
}
