import { apiClient } from './client'
import type { Notification, NotificationStats, NotificationFilter, NotificationListResponse } from '@/types/notification'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Notifications API Service
 * Handles notification-related operations like fetching notifications,
 * marking as read/unread, and managing notifications
 */

/**
 * Get all notifications with optional filters
 * @param filters - Optional filters for type, read status, pagination
 * @returns List of notifications with pagination info
 */
export async function getNotifications(filters?: NotificationFilter): Promise<NotificationListResponse> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Fetching notifications', { filters })
    const response = await apiClient.get<NotificationListResponse>(
      API_ENDPOINTS.NOTIFICATIONS.LIST,
      { params: filters }
    )
    apiLogger.info('[NOTIFICATIONS API] Notifications fetched successfully', {
      count: Array.isArray(response) ? response.length : 0
    })

    // Handle both array and object responses
    if (Array.isArray(response)) {
      return {
        notifications: response,
        total: response.length,
        page: 1,
        limit: filters?.limit || 20,
        hasMore: false
      }
    }

    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to fetch notifications', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch notifications')
  }
}

/**
 * Get notification statistics (total, unread, read counts)
 * @returns Notification statistics
 */
export async function getNotificationStats(): Promise<NotificationStats> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Fetching notification stats')
    const response = await apiClient.get<APIResponse<NotificationStats>>(
      API_ENDPOINTS.NOTIFICATIONS.STATS
    )
    apiLogger.info('[NOTIFICATIONS API] Stats fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to fetch stats', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch notification stats')
  }
}

/**
 * Get a single notification by ID
 * @param id - Notification ID
 * @returns Notification details
 */
export async function getNotification(id: string): Promise<Notification> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Fetching notification', { id })
    const response = await apiClient.get<APIResponse<Notification>>(
      API_ENDPOINTS.NOTIFICATIONS.GET(id)
    )
    apiLogger.info('[NOTIFICATIONS API] Notification fetched successfully', { id })
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to fetch notification', {
      id,
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch notification')
  }
}

/**
 * Mark a notification as read
 * @param id - Notification ID
 * @returns Updated notification
 */
export async function markNotificationAsRead(id: string): Promise<Notification> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Marking notification as read', { id })
    const response = await apiClient.put<APIResponse<Notification>>(
      API_ENDPOINTS.NOTIFICATIONS.MARK_READ(id)
    )
    apiLogger.info('[NOTIFICATIONS API] Notification marked as read', { id })
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to mark notification as read', {
      id,
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to mark notification as read')
  }
}

/**
 * Mark a notification as unread
 * @param id - Notification ID
 * @returns Updated notification
 */
export async function markNotificationAsUnread(id: string): Promise<Notification> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Marking notification as unread', { id })
    const response = await apiClient.put<APIResponse<Notification>>(
      API_ENDPOINTS.NOTIFICATIONS.MARK_UNREAD(id)
    )
    apiLogger.info('[NOTIFICATIONS API] Notification marked as unread', { id })
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to mark notification as unread', {
      id,
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to mark notification as unread')
  }
}

/**
 * Mark all notifications as read
 * @returns API response
 */
export async function markAllNotificationsAsRead(): Promise<APIResponse<void>> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Marking all notifications as read')
    const response = await apiClient.put<APIResponse<void>>(
      API_ENDPOINTS.NOTIFICATIONS.MARK_ALL_READ
    )
    apiLogger.info('[NOTIFICATIONS API] All notifications marked as read')
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to mark all notifications as read', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to mark all notifications as read')
  }
}

/**
 * Delete a notification
 * @param id - Notification ID
 * @returns API response
 */
export async function deleteNotification(id: string): Promise<APIResponse<void>> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Deleting notification', { id })
    const response = await apiClient.delete<APIResponse<void>>(
      API_ENDPOINTS.NOTIFICATIONS.DELETE(id)
    )
    apiLogger.info('[NOTIFICATIONS API] Notification deleted successfully', { id })
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to delete notification', {
      id,
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to delete notification')
  }
}

/**
 * Create a notification (admin/system use)
 * @param data - Notification data
 * @returns Created notification
 */
export async function createNotification(data: Partial<Notification>): Promise<Notification> {
  try {
    apiLogger.debug('[NOTIFICATIONS API] Creating notification', { data })
    const response = await apiClient.post<APIResponse<Notification>>(
      API_ENDPOINTS.NOTIFICATIONS.CREATE,
      data
    )
    apiLogger.info('[NOTIFICATIONS API] Notification created successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[NOTIFICATIONS API] Failed to create notification', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to create notification')
  }
}
