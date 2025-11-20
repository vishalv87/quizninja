import { apiClient } from './client'
import type { UserStats, UserProfile } from '@/types/user'
import type { QuizAttempt } from '@/types/quiz'
import type { APIResponse, PaginatedResponse, AttemptHistoryResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * User API Service
 * Handles user-related operations like stats, attempts, and active sessions
 */

/**
 * Active Session Response Type (Frontend)
 */
export interface ActiveSession {
  id: string
  quiz_id: string
  quiz_title: string
  category: string
  difficulty: string
  started_at: string
  time_elapsed_seconds: number
  questions_answered: number
  total_questions: number
  status: 'in_progress' | 'paused'
}

/**
 * Backend Active Session Response Type
 * Note: Backend sends 'active' instead of 'in_progress' for session_state
 * Note: Backend sends 'time_spent_so_far' instead of 'time_elapsed_seconds'
 */
export interface BackendActiveSession {
  id: string
  quiz_id: string
  quiz_title: string
  quiz_category: string
  quiz_difficulty: string
  quiz_thumbnail?: string
  total_questions: number
  original_time_limit: number
  session_state: 'active' | 'paused'
  current_question_index: number
  created_at: string
  time_spent_so_far: number
}

/**
 * Backend Active Sessions Response Wrapper
 */
export interface ActiveSessionsResponse {
  sessions: BackendActiveSession[]
  total: number
  active_count: number
  paused_count: number
}

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
    apiLogger.debug('[USER API] Fetching user stats')
    const response = await apiClient.get<APIResponse<UserStats>>(
      API_ENDPOINTS.USERS.STATS
    )
    apiLogger.info('[USER API] User stats fetched successfully')
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
    apiLogger.debug('[USER API] Fetching user attempts', { filters })

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
    const response = await apiClient.get<AttemptHistoryResponse<QuizAttempt>>(url)

    // Transform backend response format to frontend PaginatedResponse format
    const paginatedData: PaginatedResponse<QuizAttempt> = {
      data: response.attempts || [],  // Map backend "attempts" field to frontend "data" field
      total: response.total || 0,
      limit: response.page_size || filters?.limit || 10,
      offset: ((response.page || 1) - 1) * (response.page_size || filters?.limit || 10),
      has_more: (response.page || 0) < (response.total_pages || 0)
    }

    apiLogger.info('[USER API] User attempts fetched successfully', {
      count: paginatedData.data.length,
      total: paginatedData.total,
    })

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
 * Get user's active quiz sessions
 * Returns all in-progress or paused quiz attempts
 */
export async function getActiveSessions(): Promise<ActiveSessionsResponse> {
  try {
    apiLogger.debug('[USER API] Fetching active sessions')
    const response = await apiClient.get<{ data: ActiveSessionsResponse }>(
      API_ENDPOINTS.USERS.ACTIVE_SESSIONS
    )

    // For this endpoint, the interceptor doesn't unwrap the data wrapper
    const activeSessions = (response as any).data || response

    apiLogger.info('[USER API] Active sessions fetched successfully', {
      count: activeSessions?.sessions?.length || 0,
    })
    return activeSessions
  } catch (error: any) {
    apiLogger.error('[USER API] Failed to fetch active sessions', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch active sessions')
  }
}

/**
 * Get details of a specific attempt
 */
export async function getAttemptDetails(attemptId: string): Promise<APIResponse<QuizAttempt>> {
  try {
    apiLogger.debug('[USER API] Fetching attempt details', { attemptId })
    const response = await apiClient.get<APIResponse<QuizAttempt>>(
      API_ENDPOINTS.USERS.ATTEMPT_DETAILS(attemptId)
    )
    apiLogger.info('[USER API] Attempt details fetched successfully')
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
    apiLogger.debug('[USER API] Fetching user profile', { userId })
    const response = await apiClient.get<APIResponse<UserProfile>>(
      API_ENDPOINTS.USERS.PROFILE(userId)
    )
    apiLogger.info('[USER API] User profile fetched successfully')
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
    apiLogger.debug('[USER API] Fetching user stats by ID', { userId })
    const response = await apiClient.get<APIResponse<UserStats>>(
      API_ENDPOINTS.USERS.USER_STATS(userId)
    )
    apiLogger.info('[USER API] User stats by ID fetched successfully')
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
  getActiveSessions,
  getAttemptDetails,
  getUserProfile,
  getUserStatsById,
}