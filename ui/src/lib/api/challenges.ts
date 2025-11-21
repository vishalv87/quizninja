import { apiClient } from './client'
import type { Challenge, CreateChallengeRequest, ChallengeStats } from '@/types/challenge'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Challenges API Service
 * Handles challenge-related operations like creating challenges, accepting/declining,
 * and managing challenge lifecycle
 */

/**
 * Link Attempt Request Type
 */
export interface LinkAttemptRequest {
  attempt_id: string
}

/**
 * Challenges List Response Type
 */
export interface ChallengesListResponse {
  challenges: Challenge[]
  total: number
}

/**
 * Get all challenges for the current user
 * Returns all challenges (as challenger or opponent)
 */
export async function getChallenges(): Promise<ChallengesListResponse> {
  try {
    const response = await apiClient.get<ChallengesListResponse>(
      API_ENDPOINTS.CHALLENGES.LIST
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to fetch challenges', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch challenges')
  }
}

/**
 * Get challenge statistics for the current user
 * Returns stats like total challenges, wins, losses, win rate
 */
export async function getChallengeStats(): Promise<APIResponse<ChallengeStats>> {
  try {
    const response = await apiClient.get<APIResponse<ChallengeStats>>(
      API_ENDPOINTS.CHALLENGES.STATS
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to fetch challenge stats', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch challenge stats')
  }
}

/**
 * Get pending challenges (challenges waiting for acceptance)
 * Returns challenges with status 'pending'
 */
export async function getPendingChallenges(): Promise<ChallengesListResponse> {
  try {
    const response = await apiClient.get<ChallengesListResponse>(
      API_ENDPOINTS.CHALLENGES.PENDING
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to fetch pending challenges', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch pending challenges')
  }
}

/**
 * Get active challenges (accepted challenges in progress)
 * Returns challenges with status 'accepted'
 */
export async function getActiveChallenges(): Promise<ChallengesListResponse> {
  try {
    const response = await apiClient.get<ChallengesListResponse>(
      API_ENDPOINTS.CHALLENGES.ACTIVE
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to fetch active challenges', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch active challenges')
  }
}

/**
 * Get completed challenges
 * Returns challenges with status 'completed'
 */
export async function getCompletedChallenges(): Promise<ChallengesListResponse> {
  try {
    const response = await apiClient.get<ChallengesListResponse>(
      API_ENDPOINTS.CHALLENGES.COMPLETED
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to fetch completed challenges', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch completed challenges')
  }
}

/**
 * Get a specific challenge by ID
 * @param challengeId - The ID of the challenge to retrieve
 */
export async function getChallenge(challengeId: string): Promise<APIResponse<Challenge>> {
  try {
    const response = await apiClient.get<APIResponse<Challenge>>(
      API_ENDPOINTS.CHALLENGES.GET(challengeId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to fetch challenge details', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch challenge details')
  }
}

/**
 * Create a new challenge
 * @param data - The challenge creation data (opponent_id and quiz_id)
 */
export async function createChallenge(data: CreateChallengeRequest): Promise<APIResponse<Challenge>> {
  try {
    const response = await apiClient.post<APIResponse<Challenge>>(
      API_ENDPOINTS.CHALLENGES.CREATE,
      data
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to create challenge', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to create challenge')
  }
}

/**
 * Accept a challenge
 * @param challengeId - The ID of the challenge to accept
 */
export async function acceptChallenge(challengeId: string): Promise<APIResponse<Challenge>> {
  try {
    const response = await apiClient.put<APIResponse<Challenge>>(
      API_ENDPOINTS.CHALLENGES.ACCEPT(challengeId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to accept challenge', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to accept challenge')
  }
}

/**
 * Decline a challenge
 * @param challengeId - The ID of the challenge to decline
 */
export async function declineChallenge(challengeId: string): Promise<APIResponse<Challenge>> {
  try {
    const response = await apiClient.put<APIResponse<Challenge>>(
      API_ENDPOINTS.CHALLENGES.DECLINE(challengeId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to decline challenge', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to decline challenge')
  }
}

/**
 * Link a quiz attempt to a challenge
 * @param challengeId - The ID of the challenge
 * @param attemptId - The ID of the quiz attempt to link
 */
export async function linkAttempt(challengeId: string, attemptId: string): Promise<APIResponse<Challenge>> {
  try {
    const response = await apiClient.post<APIResponse<Challenge>>(
      API_ENDPOINTS.CHALLENGES.LINK_ATTEMPT(challengeId),
      { attempt_id: attemptId }
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to link attempt', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to link attempt')
  }
}

/**
 * Complete a challenge (mark as finished)
 * @param challengeId - The ID of the challenge to complete
 */
export async function completeChallenge(challengeId: string): Promise<APIResponse<Challenge>> {
  try {
    const response = await apiClient.put<APIResponse<Challenge>>(
      API_ENDPOINTS.CHALLENGES.COMPLETE(challengeId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[CHALLENGES API] Failed to complete challenge', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to complete challenge')
  }
}

/**
 * Export all challenges API functions
 */
export const challengesApi = {
  getChallenges,
  getChallengeStats,
  getPendingChallenges,
  getActiveChallenges,
  getCompletedChallenges,
  getChallenge,
  createChallenge,
  acceptChallenge,
  declineChallenge,
  linkAttempt,
  completeChallenge,
}
