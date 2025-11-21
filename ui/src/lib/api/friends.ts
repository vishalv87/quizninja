import { apiClient } from './client'
import type { Friend, FriendRequest } from '@/types/user'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Friends API Service
 * Handles friend-related operations like friend requests, friends list, and user search
 */

/**
 * User Search Result Type
 */
export interface UserSearchResult {
  id: string
  full_name: string
  email: string
  avatar_url?: string
  is_friend: boolean
  has_pending_request: boolean
  is_request_sent: boolean
}

/**
 * Friend Request Response Action
 */
export interface FriendRequestAction {
  action: 'accept' | 'decline'
}

/**
 * Friends List Response Type (matches backend FriendsListResponse)
 */
export interface FriendsListResponse {
  friends: Friend[]
  total: number
}

/**
 * Get current user's friends list
 * Returns all accepted friend relationships
 */
export async function getFriends(): Promise<FriendsListResponse> {
  try {
    const response = await apiClient.get<FriendsListResponse>(
      API_ENDPOINTS.FRIENDS.LIST
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to fetch friends list', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch friends')
  }
}

/**
 * Get friend requests (both sent and received)
 * Returns all pending friend requests
 */
export async function getFriendRequests(): Promise<APIResponse<FriendRequest[]>> {
  try {
    const response = await apiClient.get<APIResponse<FriendRequest[]>>(
      API_ENDPOINTS.FRIENDS.REQUESTS.LIST
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to fetch friend requests', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch friend requests')
  }
}

/**
 * Send a friend request to another user
 * @param userId - The ID of the user to send the request to
 */
export async function sendFriendRequest(userId: string): Promise<APIResponse<FriendRequest>> {
  try {
    const response = await apiClient.post<APIResponse<FriendRequest>>(
      API_ENDPOINTS.FRIENDS.REQUESTS.SEND,
      { requested_user_id: userId }
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to send friend request', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to send friend request')
  }
}

/**
 * Accept a friend request
 * @param requestId - The ID of the friend request to accept
 */
export async function acceptFriendRequest(requestId: string): Promise<APIResponse<Friend>> {
  try {
    const response = await apiClient.put<APIResponse<Friend>>(
      API_ENDPOINTS.FRIENDS.REQUESTS.RESPOND(requestId),
      { status: 'accepted' }
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to accept friend request', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to accept friend request')
  }
}

/**
 * Decline a friend request
 * @param requestId - The ID of the friend request to decline
 */
export async function declineFriendRequest(requestId: string): Promise<APIResponse<null>> {
  try {
    const response = await apiClient.put<APIResponse<null>>(
      API_ENDPOINTS.FRIENDS.REQUESTS.RESPOND(requestId),
      { status: 'rejected' }
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to decline friend request', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to decline friend request')
  }
}

/**
 * Cancel a sent friend request
 * @param requestId - The ID of the friend request to cancel
 */
export async function cancelFriendRequest(requestId: string): Promise<APIResponse<null>> {
  try {
    const response = await apiClient.delete<APIResponse<null>>(
      API_ENDPOINTS.FRIENDS.REQUESTS.CANCEL(requestId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to cancel friend request', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to cancel friend request')
  }
}

/**
 * Remove a friend
 * @param friendId - The ID of the friend to remove
 */
export async function removeFriend(friendId: string): Promise<APIResponse<null>> {
  try {
    const response = await apiClient.delete<APIResponse<null>>(
      API_ENDPOINTS.FRIENDS.REMOVE(friendId)
    )
    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to remove friend', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to remove friend')
  }
}

/**
 * Search for users by name or email
 * @param query - The search query
 */
export async function searchUsers(query: string): Promise<APIResponse<UserSearchResult[]>> {
  try {
    const params = new URLSearchParams()
    params.append('q', query)

    const url = `${API_ENDPOINTS.FRIENDS.SEARCH}?${params.toString()}`

    const response = await apiClient.get<APIResponse<UserSearchResult[]>>(url)

    return response as any
  } catch (error: any) {
    apiLogger.error('[FRIENDS API] Failed to search users', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Failed to search users')
  }
}

/**
 * Export all friends API functions
 */
export const friendsApi = {
  getFriends,
  getFriendRequests,
  sendFriendRequest,
  acceptFriendRequest,
  declineFriendRequest,
  cancelFriendRequest,
  removeFriend,
  searchUsers,
}
