"use client";

import { useQuery, useMutation, useQueryClient, UseQueryResult, UseMutationResult } from "@tanstack/react-query";
import {
  getNotifications,
  getNotification,
  getNotificationStats,
  markNotificationAsRead,
  markNotificationAsUnread,
  markAllNotificationsAsRead,
  deleteNotification
} from "@/lib/api/notifications";
import type { Notification, NotificationFilter, NotificationStats } from "@/types/notification";
import type { NotificationType } from "@/constants";
import { QUERY_KEYS } from "@/lib/constants";
import { toast } from "sonner";

/**
 * Hook to fetch notifications with optional filters
 * Returns paginated list of notifications
 *
 * @param filters - Optional filters for type, read status, and pagination
 * @returns React Query result with notifications data
 */
export function useNotifications(filters?: NotificationFilter): UseQueryResult<Notification[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.NOTIFICATIONS, filters],
    queryFn: async () => {
      const response = await getNotifications(filters);
      return Array.isArray(response.notifications) ? response.notifications : [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes (notifications are time-sensitive)
    refetchOnWindowFocus: true, // Refetch when user returns to window
    refetchInterval: 60 * 1000, // Auto-refetch every 60 seconds for real-time feel
  });
}

/**
 * Hook to fetch a single notification by ID
 *
 * @param id - Notification ID
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with notification data
 */
export function useNotification(
  id: string,
  enabled: boolean = true
): UseQueryResult<Notification, Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.NOTIFICATIONS, id],
    queryFn: async () => {
      return await getNotification(id);
    },
    enabled: enabled && !!id,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook to fetch unread notifications only
 * Convenience wrapper around useNotifications
 *
 * @returns React Query result with unread notifications
 */
export function useUnreadNotifications(): UseQueryResult<Notification[], Error> {
  return useNotifications({ is_read: false });
}

/**
 * Hook to fetch notifications by type
 * Convenience wrapper around useNotifications
 *
 * @param type - Notification type to filter by
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with filtered notifications
 */
export function useNotificationsByType(
  type: NotificationType,
  enabled: boolean = true
): UseQueryResult<Notification[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.NOTIFICATIONS, "type", type],
    queryFn: async () => {
      const response = await getNotifications({ type });
      return Array.isArray(response.notifications) ? response.notifications : [];
    },
    enabled: enabled && !!type,
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch notification statistics
 *
 * @returns React Query result with notification stats
 */
export function useNotificationStats(): UseQueryResult<NotificationStats, Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.NOTIFICATIONS, "stats"],
    queryFn: getNotificationStats,
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
    refetchInterval: 60 * 1000, // Auto-refetch every 60 seconds
  });
}

/**
 * Hook to mark a notification as read
 *
 * @returns Mutation hook with optimistic update
 */
export function useMarkNotificationAsRead(): UseMutationResult<Notification, Error, string> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: markNotificationAsRead,
    onMutate: async (notificationId: string) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });

      // Snapshot previous value
      const previousNotifications = queryClient.getQueryData([QUERY_KEYS.NOTIFICATIONS]);

      // Optimistically update notification
      queryClient.setQueriesData(
        { queryKey: [QUERY_KEYS.NOTIFICATIONS] },
        (old: any) => {
          if (Array.isArray(old)) {
            return old.map((notification: Notification) =>
              notification.id === notificationId
                ? { ...notification, is_read: true }
                : notification
            );
          }
          return old;
        }
      );

      // Update stats optimistically
      queryClient.setQueryData(
        [QUERY_KEYS.NOTIFICATIONS, "stats"],
        (old: NotificationStats | undefined) => {
          if (old) {
            return {
              ...old,
              unread_notifications: Math.max(0, old.unread_notifications - 1),
              read_notifications: old.read_notifications + 1,
            };
          }
          return old;
        }
      );

      return { previousNotifications };
    },
    onError: (_err, _variables, context: any) => {
      // Rollback on error
      if (context?.previousNotifications) {
        queryClient.setQueryData([QUERY_KEYS.NOTIFICATIONS], context.previousNotifications);
      }
      toast.error("Failed to mark notification as read");
    },
    onSuccess: () => {
      // Invalidate all notification queries to refetch fresh data
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });
    },
  });
}

/**
 * Hook to mark a notification as unread
 *
 * @returns Mutation hook with optimistic update
 */
export function useMarkNotificationAsUnread(): UseMutationResult<Notification, Error, string> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: markNotificationAsUnread,
    onMutate: async (notificationId: string) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });

      // Snapshot previous value
      const previousNotifications = queryClient.getQueryData([QUERY_KEYS.NOTIFICATIONS]);

      // Optimistically update notification
      queryClient.setQueriesData(
        { queryKey: [QUERY_KEYS.NOTIFICATIONS] },
        (old: any) => {
          if (Array.isArray(old)) {
            return old.map((notification: Notification) =>
              notification.id === notificationId
                ? { ...notification, is_read: false }
                : notification
            );
          }
          return old;
        }
      );

      // Update stats optimistically
      queryClient.setQueryData(
        [QUERY_KEYS.NOTIFICATIONS, "stats"],
        (old: NotificationStats | undefined) => {
          if (old) {
            return {
              ...old,
              unread_notifications: old.unread_notifications + 1,
              read_notifications: Math.max(0, old.read_notifications - 1),
            };
          }
          return old;
        }
      );

      return { previousNotifications };
    },
    onError: (_err, _variables, context: any) => {
      // Rollback on error
      if (context?.previousNotifications) {
        queryClient.setQueryData([QUERY_KEYS.NOTIFICATIONS], context.previousNotifications);
      }
      toast.error("Failed to mark notification as unread");
    },
    onSuccess: () => {
      // Invalidate all notification queries to refetch fresh data
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });
    },
  });
}

/**
 * Hook to mark all notifications as read
 *
 * @returns Mutation hook with cache invalidation
 */
export function useMarkAllNotificationsAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: markAllNotificationsAsRead,
    onMutate: async () => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });

      // Snapshot previous value
      const previousNotifications = queryClient.getQueryData([QUERY_KEYS.NOTIFICATIONS]);

      // Optimistically update all notifications
      queryClient.setQueriesData(
        { queryKey: [QUERY_KEYS.NOTIFICATIONS] },
        (old: any) => {
          if (Array.isArray(old)) {
            return old.map((notification: Notification) => ({
              ...notification,
              is_read: true,
            }));
          }
          return old;
        }
      );

      // Update stats optimistically
      queryClient.setQueryData(
        [QUERY_KEYS.NOTIFICATIONS, "stats"],
        (old: NotificationStats | undefined) => {
          if (old) {
            return {
              ...old,
              unread_notifications: 0,
              read_notifications: old.total_notifications,
            };
          }
          return old;
        }
      );

      return { previousNotifications };
    },
    onError: (_err, _variables, context: any) => {
      // Rollback on error
      if (context?.previousNotifications) {
        queryClient.setQueryData([QUERY_KEYS.NOTIFICATIONS], context.previousNotifications);
      }
      toast.error("Failed to mark all notifications as read");
    },
    onSuccess: () => {
      // Invalidate all notification queries to refetch fresh data
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });
      toast.success("All notifications marked as read");
    },
  });
}

/**
 * Hook to delete a notification
 *
 * @returns Mutation hook with optimistic update
 */
export function useDeleteNotification() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteNotification,
    onMutate: async (notificationId: string) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });

      // Snapshot previous value
      const previousNotifications = queryClient.getQueryData([QUERY_KEYS.NOTIFICATIONS]);

      // Get the notification before deleting to update stats
      let deletedNotification: Notification | undefined;
      queryClient.setQueriesData(
        { queryKey: [QUERY_KEYS.NOTIFICATIONS] },
        (old: any) => {
          if (Array.isArray(old)) {
            deletedNotification = old.find((n: Notification) => n.id === notificationId);
            return old.filter((notification: Notification) => notification.id !== notificationId);
          }
          return old;
        }
      );

      // Update stats optimistically
      queryClient.setQueryData(
        [QUERY_KEYS.NOTIFICATIONS, "stats"],
        (old: NotificationStats | undefined) => {
          if (old && deletedNotification) {
            return {
              ...old,
              total_notifications: Math.max(0, old.total_notifications - 1),
              unread_notifications: deletedNotification.is_read ? old.unread_notifications : Math.max(0, old.unread_notifications - 1),
              read_notifications: deletedNotification.is_read ? Math.max(0, old.read_notifications - 1) : old.read_notifications,
            };
          }
          return old;
        }
      );

      return { previousNotifications };
    },
    onError: (_err, _variables, context: any) => {
      // Rollback on error
      if (context?.previousNotifications) {
        queryClient.setQueryData([QUERY_KEYS.NOTIFICATIONS], context.previousNotifications);
      }
      toast.error("Failed to delete notification");
    },
    onSuccess: () => {
      // Invalidate all notification queries to refetch fresh data
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.NOTIFICATIONS] });
      toast.success("Notification deleted");
    },
  });
}