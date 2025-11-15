import { apiClient } from './client'
import type { Achievement, UserAchievement, AchievementProgress, AchievementStats } from '@/types/achievement'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Achievements API Service
 * Handles achievement-related operations like fetching achievements,
 * checking progress, and unlocking achievements
 */

/**
 * Achievements List Response Type
 */
export interface AchievementsListResponse {
  achievements: Achievement[]
  total: number
}

/**
 * Achievement Progress Response Type
 */
export interface AchievementProgressResponse {
  progress: AchievementProgress[]
  total: number
}

/**
 * User Achievements Response Type
 */
export interface UserAchievementsResponse {
  achievements: UserAchievement[]
  total: number
}

/**
 * Get all available achievements
 * Returns all achievements in the system
 */
export async function getAchievements(): Promise<AchievementsListResponse> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Fetching all achievements')
    const response = await apiClient.get<AchievementsListResponse>(
      API_ENDPOINTS.ACHIEVEMENTS.LIST
    )
    apiLogger.info('[ACHIEVEMENTS API] Achievements fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to fetch achievements', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch achievements')
  }
}

/**
 * Get achievements by category
 * @param category - The category to filter by
 */
export async function getAchievementsByCategory(category: string): Promise<AchievementsListResponse> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Fetching achievements by category', { category })
    const response = await apiClient.get<AchievementsListResponse>(
      API_ENDPOINTS.ACHIEVEMENTS.BY_CATEGORY(category)
    )
    apiLogger.info('[ACHIEVEMENTS API] Achievements by category fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to fetch achievements by category', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch achievements by category')
  }
}

/**
 * Get current user's achievement progress
 * Returns progress on all achievements (both locked and unlocked)
 */
export async function getAchievementProgress(): Promise<AchievementProgressResponse> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Fetching achievement progress')
    const response = await apiClient.get<AchievementProgressResponse>(
      API_ENDPOINTS.ACHIEVEMENTS.PROGRESS
    )
    apiLogger.info('[ACHIEVEMENTS API] Achievement progress fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to fetch achievement progress', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch achievement progress')
  }
}

/**
 * Get current user's achievement statistics
 * Returns stats like total achievements, unlocked count, points earned
 */
export async function getAchievementStats(): Promise<APIResponse<AchievementStats>> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Fetching achievement statistics')
    const response = await apiClient.get<APIResponse<AchievementStats>>(
      API_ENDPOINTS.ACHIEVEMENTS.STATS
    )
    apiLogger.info('[ACHIEVEMENTS API] Achievement stats fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to fetch achievement stats', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch achievement stats')
  }
}

/**
 * Get current user's unlocked achievements
 * Returns only achievements that the user has unlocked
 */
export async function getUserAchievements(): Promise<UserAchievementsResponse> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Fetching user achievements')
    const response = await apiClient.get<UserAchievementsResponse>(
      API_ENDPOINTS.USERS.ACHIEVEMENTS
    )
    apiLogger.info('[ACHIEVEMENTS API] User achievements fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to fetch user achievements', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user achievements')
  }
}

/**
 * Get a specific user's achievements
 * @param userId - The ID of the user
 */
export async function getUserAchievementsById(userId: string): Promise<UserAchievementsResponse> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Fetching achievements for user', { userId })
    const response = await apiClient.get<UserAchievementsResponse>(
      API_ENDPOINTS.USERS.ACHIEVEMENTS_BY_USER(userId)
    )
    apiLogger.info('[ACHIEVEMENTS API] User achievements fetched successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to fetch user achievements', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch user achievements')
  }
}

/**
 * Check for newly unlocked achievements
 * Triggers backend to check if any achievements were unlocked based on recent activities
 * Returns any newly unlocked achievements
 */
export async function checkAchievements(): Promise<APIResponse<UserAchievement[]>> {
  try {
    apiLogger.debug('[ACHIEVEMENTS API] Checking for new achievements')
    const response = await apiClient.post<APIResponse<UserAchievement[]>>(
      API_ENDPOINTS.ACHIEVEMENTS.CHECK
    )
    apiLogger.info('[ACHIEVEMENTS API] Achievement check completed successfully')
    return response as any
  } catch (error: any) {
    apiLogger.error('[ACHIEVEMENTS API] Failed to check achievements', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to check achievements')
  }
}

/**
 * Export all achievements API functions
 */
export const achievementsApi = {
  getAchievements,
  getAchievementsByCategory,
  getAchievementProgress,
  getAchievementStats,
  getUserAchievements,
  getUserAchievementsById,
  checkAchievements,
}
