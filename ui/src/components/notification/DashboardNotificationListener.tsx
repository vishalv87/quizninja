"use client";

import { useNewNotificationToast } from "@/hooks/useNewNotificationToast";

/**
 * Dashboard Notification Listener
 *
 * This component listens for new notifications and shows toast alerts when they arrive.
 * It should be included once in the dashboard layout to enable global notification toasts.
 *
 * Features:
 * - Monitors notification count changes
 * - Shows toast notifications for new arrivals
 * - Handles single and multiple notifications appropriately
 */
export function DashboardNotificationListener() {
  // Initialize the notification toast hook
  useNewNotificationToast();

  // This component doesn't render anything - it just runs the hook
  return null;
}