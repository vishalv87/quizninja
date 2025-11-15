"use client";

import { useState } from "react";
import Link from "next/link";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreVertical, Trash2, Eye, EyeOff } from "lucide-react";
import type { Notification } from "@/types/notification";
import {
  getNotificationIcon,
  getNotificationColor,
  formatNotificationTime,
  getNotificationActionUrl,
} from "@/lib/notification-utils";
import { cn } from "@/lib/utils";
import {
  useMarkNotificationAsRead,
  useMarkNotificationAsUnread,
  useDeleteNotification,
} from "@/hooks/useNotifications";

interface NotificationCardProps {
  notification: Notification;
  variant?: "default" | "compact";
}

export function NotificationCard({ notification, variant = "default" }: NotificationCardProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const markAsRead = useMarkNotificationAsRead();
  const markAsUnread = useMarkNotificationAsUnread();
  const deleteNotification = useDeleteNotification();

  const Icon = getNotificationIcon(notification.type);
  const colors = getNotificationColor(notification.type);
  const actionUrl = getNotificationActionUrl(notification.type, notification.data);
  const timeAgo = formatNotificationTime(notification.created_at);

  const handleMarkAsRead = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (!notification.is_read) {
      markAsRead.mutate(notification.id);
    }
  };

  const handleToggleRead = async (e?: React.MouseEvent) => {
    if (e) {
      e.preventDefault();
      e.stopPropagation();
    }
    if (notification.is_read) {
      markAsUnread.mutate(notification.id);
    } else {
      markAsRead.mutate(notification.id);
    }
  };

  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDeleting(true);
    deleteNotification.mutate(notification.id, {
      onSettled: () => setIsDeleting(false),
    });
  };

  const CardContent = (
    <div
      className={cn(
        "group relative flex gap-3 p-4 transition-colors hover:bg-accent/50",
        !notification.is_read && "bg-accent/30",
        isDeleting && "opacity-50 pointer-events-none"
      )}
      onClick={handleMarkAsRead}
    >
      {/* Unread indicator */}
      {!notification.is_read && (
        <div className="absolute left-0 top-0 bottom-0 w-1 bg-primary" />
      )}

      {/* Icon */}
      <div className={cn("flex-shrink-0 p-2 rounded-full", colors.bg)}>
        <Icon className={cn("h-5 w-5", colors.icon)} />
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0 space-y-1">
        <div className="flex items-start justify-between gap-2">
          <h4
            className={cn(
              "font-semibold text-sm leading-tight",
              !notification.is_read && "font-bold"
            )}
          >
            {notification.title}
          </h4>

          {/* Actions dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild onClick={(e) => e.stopPropagation()}>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0"
              >
                <MoreVertical className="h-4 w-4" />
                <span className="sr-only">More actions</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                onClick={handleToggleRead}
                disabled={markAsRead.isPending || markAsUnread.isPending}
              >
                {notification.is_read ? (
                  <>
                    <EyeOff className="mr-2 h-4 w-4" />
                    Mark as unread
                  </>
                ) : (
                  <>
                    <Eye className="mr-2 h-4 w-4" />
                    Mark as read
                  </>
                )}
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={handleDelete}
                disabled={deleteNotification.isPending}
                className="text-destructive focus:text-destructive"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <p className="text-sm text-muted-foreground line-clamp-2">
          {notification.message}
        </p>

        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <time dateTime={notification.created_at}>{timeAgo}</time>
        </div>
      </div>
    </div>
  );

  // If there's an action URL, wrap in a Link
  if (actionUrl) {
    return (
      <Card className="overflow-hidden hover:shadow-md transition-shadow">
        <Link href={actionUrl} className="block">
          {CardContent}
        </Link>
      </Card>
    );
  }

  // Otherwise, just return the card
  return (
    <Card className="overflow-hidden hover:shadow-md transition-shadow">
      {CardContent}
    </Card>
  );
}