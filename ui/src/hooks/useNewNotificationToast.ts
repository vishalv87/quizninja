"use client";

import { useEffect, useRef } from "react";
import { useNotificationStats, useNotifications } from "./useNotifications";
import {
  showNotificationToast,
  showMultipleNotificationsToast,
} from "@/components/notification/NotificationToast";

/**
 * Hook to detect and show toast notifications when new notifications arrive
 *
 * This hook monitors the unread notification count and shows a toast when it increases,
 * indicating that new notifications have arrived via polling.
 *
 * Features:
 * - Only shows toasts for NEW notifications (not on initial load)
 * - Shows individual toast for single notification
 * - Shows summary toast for multiple notifications
 * - Prevents toast spam by tracking what's already been shown
 *
 * Usage:
 * Call this hook once in your app layout or provider to enable global notification toasts
 */
export function useNewNotificationToast() {
  // Track previous unread count to detect increases
  const prevUnreadCountRef = useRef<number | undefined>();

  // Track if this is the initial mount to avoid showing toasts on first load
  const isInitialMountRef = useRef(true);

  // Get notification statistics (refetches every 60 seconds)
  const { data: stats } = useNotificationStats();

  // Get latest unread notifications to show in toast
  const { data: notifications = [] } = useNotifications({
    is_read: false,
    limit: 5, // Get latest 5 unread notifications
  });

  useEffect(() => {
    // Skip on initial mount to avoid showing toasts for existing notifications
    if (isInitialMountRef.current) {
      isInitialMountRef.current = false;
      prevUnreadCountRef.current = stats?.unread_notifications;
      return;
    }

    // If we have both previous and current counts
    if (prevUnreadCountRef.current !== undefined && stats) {
      const prevCount = prevUnreadCountRef.current;
      const currentCount = stats.unread_notifications;

      // Check if count increased (new notifications arrived)
      if (currentCount > prevCount) {
        const newCount = currentCount - prevCount;
        const latestNotification = notifications[0];

        // Show appropriate toast based on count
        if (newCount === 1 && latestNotification) {
          // Single new notification - show detailed toast
          showNotificationToast(latestNotification);
        } else if (newCount > 1) {
          // Multiple new notifications - show summary toast
          showMultipleNotificationsToast(newCount, latestNotification);
        }
      }
    }

    // Update previous count for next comparison
    prevUnreadCountRef.current = stats?.unread_notifications;
  }, [stats, notifications]);

  // This hook doesn't return anything - it just shows toasts as a side effect
}
