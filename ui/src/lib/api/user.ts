import { apiClient } from './client'
import type { UserStats, UserProfile } from '@/types/user'
import type { QuizAttempt } from '@/types/quiz'
import type { APIResponse, PaginatedResponse, AttemptHistoryResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * User API Service
 * Handles user-related operations like stats and attempts
 */

/**
 * User Attempt Filters
 */
export interface UserAttemptFilters {
  limit?: number
  offset?: number
  category?: string
  difficulty?: string
  status?: 'completed' | 'in_progress' | 'abandoned'
  from_date?: string
  to_date?: string
}

/**
 * Get current user's statistics
 * Returns comprehensive statistics including points, quizzes taken, achievements, etc.
 */
export async function getUserStats(): Promise<APIResponse<UserStats>> {
  try {
    const response = await apiClient.get<APIResponse<UserStats>>(
      API_ENDPOINTS.USERS.STATS
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[USER API] Failed to fetch user stats', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user stats')
  }
}

/**
 * Get user's quiz attempt history with optional filters
 * Supports pagination and filtering by category, difficulty, status, and date range
 */
export async function getUserAttempts(
  filters?: UserAttemptFilters
): Promise<PaginatedResponse<QuizAttempt>> {
  try {
    const params = new URLSearchParams()
    if (filters?.limit) params.append('limit', filters.limit.toString())
    if (filters?.offset) params.append('offset', filters.offset.toString())
    if (filters?.category) params.append('category', filters.category)
    if (filters?.difficulty) params.append('difficulty', filters.difficulty)
    if (filters?.status) params.append('status', filters.status)
    if (filters?.from_date) params.append('from_date', filters.from_date)
    if (filters?.to_date) params.append('to_date', filters.to_date)

    const url = `${API_ENDPOINTS.USERS.ATTEMPTS}${params.toString() ? '?' + params.toString() : ''}`

    // The axios interceptor already unwraps response.data, so we get AttemptHistoryResponse directly
    // Cast through unknown because TypeScript types don't reflect the interceptor behavior
    const apiResponse = await apiClient.get<AttemptHistoryResponse<QuizAttempt>>(url)
    const response = apiResponse as unknown as AttemptHistoryResponse<QuizAttempt>

    // Transform backend response format to frontend PaginatedResponse format
    const paginatedData: PaginatedResponse<QuizAttempt> = {
      data: response.attempts || [],  // Map backend "attempts" field to frontend "data" field
      total: response.total || 0,
      limit: response.page_size || filters?.limit || 10,
      offset: ((response.page || 1) - 1) * (response.page_size || filters?.limit || 10),
      has_more: (response.page || 0) < (response.total_pages || 0)
    }

    return paginatedData
  } catch (error: any) {
    apiLogger.error('[USER API] Failed to fetch user attempts', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user attempts')
  }
}

/**
 * Get details of a specific attempt
 */
export async function getAttemptDetails(attemptId: string): Promise<APIResponse<QuizAttempt>> {
  try {
    const response = await apiClient.get<APIResponse<QuizAttempt>>(
      API_ENDPOINTS.USERS.ATTEMPT_DETAILS(attemptId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[USER API] Failed to fetch attempt details', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch attempt details')
  }
}

/**
 * Get another user's profile by user ID
 * Returns profile information respecting privacy settings
 */
export async function getUserProfile(userId: string): Promise<APIResponse<UserProfile>> {
  try {
    const response = await apiClient.get<APIResponse<UserProfile>>(
      API_ENDPOINTS.USERS.PROFILE(userId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[USER API] Failed to fetch user profile', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user profile')
  }
}

/**
 * Get another user's statistics by user ID
 * Returns stats only if user's privacy settings allow it
 */
export async function getUserStatsById(userId: string): Promise<APIResponse<UserStats>> {
  try {
    const response = await apiClient.get<APIResponse<UserStats>>(
      API_ENDPOINTS.USERS.USER_STATS(userId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[USER API] Failed to fetch user stats by ID', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user stats')
  }
}

/**
 * Export all user API functions
 */
export const userApi = {
  getUserStats,
  getUserAttempts,
  getAttemptDetails,
  getUserProfile,
  getUserStatsById,
}