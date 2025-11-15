import { apiClient } from './client'
import type { LeaderboardEntry } from '@/types/api'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Leaderboard API Service
 * Handles leaderboard-related operations like fetching global rankings,
 * user rank, and leaderboard statistics
 */

/**
 * Leaderboard Response Type
 */
export interface LeaderboardResponse {
  leaderboard: LeaderboardEntry[]
  total: number
}

/**
 * User Rank Response Type
 */
export interface UserRankResponse {
  rank: number
  total_users: number
  user: {
    id: string
    full_name: string
    avatar_url?: string
  }
  total_points: number
  quizzes_completed: number
  achievements_unlocked: number
}

/**
 * Leaderboard Stats Response Type
 */
export interface LeaderboardStats {
  total_users: number
  total_points_distributed: number
  average_points: number
  top_user: {
    id: string
    full_name: string
    total_points: number
  } | null
}

/**
 * Get global leaderboard
 * Returns top users ranked by total points
 * @param limit - Maximum number of entries to return (default: 100)
 */
export async function getLeaderboard(limit: number = 100): Promise<LeaderboardResponse> {
  try {
    apiLogger.debug('[LEADERBOARD API] Fetching leaderboard', { limit })
    const response = await apiClient.get<LeaderboardResponse>(
      `${API_ENDPOINTS.LEADERBOARD.GET}?limit=${limit}`
    )
    apiLogger.info('[LEADERBOARD API] Leaderboard fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[LEADERBOARD API] Failed to fetch leaderboard', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch leaderboard')
  }
}

/**
 * Get leaderboard with achievement counts
 * Returns leaderboard with detailed achievement information
 * @param limit - Maximum number of entries to return (default: 100)
 */
export async function getLeaderboardWithAchievements(limit: number = 100): Promise<LeaderboardResponse> {
  try {
    apiLogger.debug('[LEADERBOARD API] Fetching leaderboard with achievements', { limit })
    const response = await apiClient.get<LeaderboardResponse>(
      `${API_ENDPOINTS.LEADERBOARD.WITH_ACHIEVEMENTS}?limit=${limit}`
    )
    apiLogger.info('[LEADERBOARD API] Leaderboard with achievements fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[LEADERBOARD API] Failed to fetch leaderboard with achievements', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch leaderboard with achievements')
  }
}

/**
 * Get current user's rank
 * Returns the user's position on the leaderboard
 */
export async function getUserRank(): Promise<APIResponse<UserRankResponse>> {
  try {
    apiLogger.debug('[LEADERBOARD API] Fetching user rank')
    const response = await apiClient.get<APIResponse<UserRankResponse>>(
      API_ENDPOINTS.LEADERBOARD.RANK
    )
    apiLogger.info('[LEADERBOARD API] User rank fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[LEADERBOARD API] Failed to fetch user rank', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user rank')
  }
}

/**
 * Get leaderboard statistics
 * Returns overall statistics about the leaderboard
 */
export async function getLeaderboardStats(): Promise<APIResponse<LeaderboardStats>> {
  try {
    apiLogger.debug('[LEADERBOARD API] Fetching leaderboard statistics')
    const response = await apiClient.get<APIResponse<LeaderboardStats>>(
      API_ENDPOINTS.LEADERBOARD.STATS
    )
    apiLogger.info('[LEADERBOARD API] Leaderboard stats fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[LEADERBOARD API] Failed to fetch leaderboard stats', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch leaderboard stats')
  }
}

/**
 * Export all leaderboard API functions
 */
export const leaderboardApi = {
  getLeaderboard,
  getLeaderboardWithAchievements,
  getUserRank,
  getLeaderboardStats,
}
