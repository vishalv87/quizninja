"use client";

import { useEffect, useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useNotificationStats } from "./useNotifications";
import {
  showNotificationToast,
  showMultipleNotificationsToast,
} from "@/components/notification/NotificationToast";
import type { Notification } from "@/types/notification";
import { QUERY_KEYS } from "@/lib/constants";

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
 * - Reads from existing query cache to avoid creating duplicate polling queries
 *
 * Usage:
 * Call this hook once in your app layout or provider to enable global notification toasts
 */
export function useNewNotificationToast() {
  const queryClient = useQueryClient();

  // Track previous unread count to detect increases
  const prevUnreadCountRef = useRef<number | undefined>();

  // Track if this is the initial mount to avoid showing toasts on first load
  const isInitialMountRef = useRef(true);

  // Get notification statistics (shared query key with NotificationBell — no extra request)
  const { data: stats } = useNotificationStats();

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

        // Read from the existing notification query cache (from NotificationDropdown)
        // instead of creating a separate polling query
        const cachedNotifications = queryClient.getQueryData<Notification[]>(
          [QUERY_KEYS.NOTIFICATIONS, { limit: 5 }]
        );
        const latestNotification = cachedNotifications?.[0];

        // Show appropriate toast based on count
        if (newCount === 1 && latestNotification) {
          showNotificationToast(latestNotification);
        } else if (newCount > 1) {
          showMultipleNotificationsToast(newCount, latestNotification);
        } else if (newCount === 1) {
          // No cached notification data available — show summary
          showMultipleNotificationsToast(newCount);
        }
      }
    }

    // Update previous count for next comparison
    prevUnreadCountRef.current = stats?.unread_notifications;
  }, [stats, queryClient]);

  // This hook doesn't return anything - it just shows toasts as a side effect
}
