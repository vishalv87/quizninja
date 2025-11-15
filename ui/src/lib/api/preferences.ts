import { apiClient } from './client'
import type { UserPreferences } from '@/types/user'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Preferences API Service
 * Handles user preference-related operations
 */

/**
 * Get current user's preferences
 * Returns user preferences including categories, difficulty, notifications, etc.
 */
export async function getPreferences(): Promise<APIResponse<UserPreferences>> {
  try {
    apiLogger.debug('[PREFERENCES API] Fetching user preferences')
    const response = await apiClient.get<APIResponse<UserPreferences>>(
      API_ENDPOINTS.USERS.PREFERENCES.GET
    )
    apiLogger.info('[PREFERENCES API] User preferences fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[PREFERENCES API] Failed to fetch user preferences', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch preferences')
  }
}

/**
 * Update current user's preferences
 * Updates user preferences (partial update supported)
 */
export async function updatePreferences(
  data: Partial<UserPreferences>
): Promise<APIResponse<UserPreferences>> {
  try {
    apiLogger.debug('[PREFERENCES API] Updating user preferences', { data })
    const response = await apiClient.put<APIResponse<UserPreferences>>(
      API_ENDPOINTS.USERS.PREFERENCES.UPDATE,
      data
    )
    apiLogger.info('[PREFERENCES API] User preferences updated successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[PREFERENCES API] Failed to update user preferences', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to update preferences')
  }
}

/**
 * Export all preferences API functions
 */
export const preferencesApi = {
  getPreferences,
  updatePreferences,
}
