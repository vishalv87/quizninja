"use client";

import { useState } from "react";
import Link from "next/link";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Bell, CheckCheck, ArrowRight } from "lucide-react";
import { NotificationBell } from "./NotificationBell";
import {
  useNotifications,
  useMarkAllNotificationsAsRead,
} from "@/hooks/useNotifications";
import {
  getNotificationIcon,
  getNotificationColor,
  formatNotificationTime,
  getNotificationActionUrl,
} from "@/lib/notification-utils";
import { cn } from "@/lib/utils";
import type { Notification } from "@/types/notification";

export function NotificationDropdown() {
  const [open, setOpen] = useState(false);
  const { data: notifications = [], isLoading } = useNotifications({
    limit: 5,
  });
  const markAllAsRead = useMarkAllNotificationsAsRead();

  const handleMarkAllAsRead = () => {
    markAllAsRead.mutate();
  };

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger asChild>
        <div>
          <NotificationBell />
        </div>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-[380px] p-0">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b">
          <h3 className="font-semibold text-base">Notifications</h3>
          {notifications.length > 0 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={handleMarkAllAsRead}
              disabled={markAllAsRead.isPending}
              className="h-8 text-xs"
            >
              <CheckCheck className="h-3.5 w-3.5 mr-1" />
              Mark all read
            </Button>
          )}
        </div>

        {/* Notifications list */}
        <ScrollArea className="h-[400px]">
          {isLoading ? (
            <div className="p-4 space-y-4">
              {[...Array(3)].map((_, i) => (
                <NotificationItemSkeleton key={i} />
              ))}
            </div>
          ) : notifications.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
              <div className="rounded-full bg-muted p-4 mb-3">
                <Bell className="h-8 w-8 text-muted-foreground" />
              </div>
              <p className="text-sm text-muted-foreground">
                No new notifications
              </p>
            </div>
          ) : (
            <div className="divide-y">
              {notifications.map((notification) => (
                <NotificationItem
                  key={notification.id}
                  notification={notification}
                  onClick={() => setOpen(false)}
                />
              ))}
            </div>
          )}
        </ScrollArea>

        {/* Footer */}
        {notifications.length > 0 && (
          <>
            <DropdownMenuSeparator />
            <div className="p-2">
              <Link href="/notifications" onClick={() => setOpen(false)}>
                <Button variant="ghost" className="w-full justify-between h-9">
                  View all notifications
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function NotificationItem({
  notification,
  onClick,
}: {
  notification: Notification;
  onClick?: () => void;
}) {
  const Icon = getNotificationIcon(notification.type);
  const colors = getNotificationColor(notification.type);
  const actionUrl = getNotificationActionUrl(notification.type, notification.data);
  const timeAgo = formatNotificationTime(notification.created_at);

  const content = (
    <div
      className={cn(
        "flex gap-3 p-3 hover:bg-accent transition-colors cursor-pointer relative",
        !notification.is_read && "bg-accent/30"
      )}
      onClick={onClick}
    >
      {/* Unread indicator */}
      {!notification.is_read && (
        <div className="absolute left-0 top-0 bottom-0 w-1 bg-primary" />
      )}

      {/* Icon */}
      <div className={cn("flex-shrink-0 p-1.5 rounded-full h-fit", colors.bg)}>
        <Icon className={cn("h-4 w-4", colors.icon)} />
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0 space-y-1">
        <p
          className={cn(
            "text-sm leading-tight",
            !notification.is_read && "font-semibold"
          )}
        >
          {notification.title}
        </p>
        <p className="text-xs text-muted-foreground line-clamp-2">
          {notification.message}
        </p>
        <p className="text-xs text-muted-foreground">{timeAgo}</p>
      </div>
    </div>
  );

  if (actionUrl) {
    return <Link href={actionUrl}>{content}</Link>;
  }

  return content;
}

function NotificationItemSkeleton() {
  return (
    <div className="flex gap-3 p-3">
      <div className="h-8 w-8 rounded-full bg-muted flex-shrink-0" />
      <div className="flex-1 space-y-2">
        <div className="h-3 bg-muted rounded w-3/4" />
        <div className="h-2 bg-muted rounded w-full" />
        <div className="h-2 bg-muted rounded w-1/4" />
      </div>
    </div>
  );
}