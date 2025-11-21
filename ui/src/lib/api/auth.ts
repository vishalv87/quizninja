import { apiClient } from './client'
import { supabase, signIn as supabaseSignIn, signUp as supabaseSignUp, signOut as supabaseSignOut } from '@/lib/supabase/client'
import type { User, Profile, LoginCredentials, RegisterData } from '@/types/auth'
import type { APIResponse } from '@/types/api'
import { API_ENDPOINTS } from './endpoints'
import { apiLogger } from '@/lib/logger'

/**
 * Authentication API Service
 * Handles authentication operations with both Supabase and the backend API
 */

/**
 * Login user with email and password
 * 1. Authenticate with Supabase
 * 2. Send user data to backend API to sync/create user profile
 *
 * Backend returns: { user: User, message?: string }
 */
export async function login(credentials: LoginCredentials): Promise<{ user: User; message?: string }> {
  try {
    // Step 1: Authenticate with Supabase
    const authResult = await supabaseSignIn(
      credentials.email,
      credentials.password
    )

    if (!authResult || !authResult.user) {
      apiLogger.error('[AUTH API] Supabase authentication failed - no user returned')
      throw new Error('Authentication failed')
    }

    // Extract user metadata
    const supabaseUser = authResult.user
    const userName = supabaseUser.user_metadata?.full_name ||
                     supabaseUser.user_metadata?.name ||
                     supabaseUser.email?.split('@')[0] ||
                     'User'

    // Step 2: Call backend API login endpoint
    // Send Supabase user ID and name to sync with backend
    // Note: apiClient response interceptor already unwraps response.data
    const response = await apiClient.post<APIResponse<{ user: User; profile: Profile }>>(
      API_ENDPOINTS.AUTH.LOGIN,
      {
        supabase_user_id: supabaseUser.id,
        name: userName,
        email: supabaseUser.email,
      }
    )

    // Response is already unwrapped by the interceptor
    return response as any
  } catch (error: any) {
    apiLogger.error('[AUTH API] Login failed', {
      message: error.message,
      responseData: error.response?.data,
    })
    throw new Error(error.response?.data?.message || error.message || 'Login failed')
  }
}

/**
 * Register a new user
 * 1. Create account in Supabase
 * 2. Send registration data to backend API to create user profile
 *
 * Backend returns: { user: User, message?: string }
 */
export async function register(data: RegisterData): Promise<{ user: User; message?: string }> {
  try {
    // Step 1: Create account in Supabase
    const authResult = await supabaseSignUp(
      data.email,
      data.password,
      data.name
    )

    if (!authResult || !authResult.user) {
      throw new Error('Registration failed')
    }

    // Step 2: Call backend API register endpoint
    // Send Supabase user ID and user data to create backend profile
    // Note: apiClient response interceptor already unwraps response.data
    const response = await apiClient.post<APIResponse<{ user: User; profile: Profile }>>(
      API_ENDPOINTS.AUTH.REGISTER,
      {
        supabase_user_id: authResult.user.id,
        name: data.name,
        email: authResult.user.email,
      }
    )

    // Response is already unwrapped by the interceptor
    return response as any
  } catch (error: any) {
    throw new Error(error.response?.data?.message || error.message || 'Registration failed')
  }
}

/**
 * Logout current user
 * 1. Call backend API logout endpoint
 * 2. Sign out from Supabase
 */
export async function logout(): Promise<void> {
  try {
    // Call backend API logout endpoint (optional, for logging/cleanup)
    try {
      await apiClient.post(API_ENDPOINTS.AUTH.LOGOUT)
    } catch {
      // Continue with logout even if backend call fails
    }

    // Sign out from Supabase
    await supabaseSignOut()
  } catch (error: any) {
    throw new Error(error.message || 'Logout failed')
  }
}

/**
 * Get current user's profile
 */
export async function getProfile(): Promise<APIResponse<Profile>> {
  try {
    // Note: apiClient response interceptor already unwraps response.data
    const response = await apiClient.get<APIResponse<Profile>>(
      API_ENDPOINTS.PROFILE.GET
    )
    return response as any
  } catch (error: any) {
    throw new Error(error.response?.data?.message || error.message || 'Failed to fetch profile')
  }
}

/**
 * Update current user's profile
 */
export async function updateProfile(data: Partial<Profile>): Promise<APIResponse<Profile>> {
  try {
    // Note: apiClient response interceptor already unwraps response.data
    const response = await apiClient.put<APIResponse<Profile>>(
      API_ENDPOINTS.PROFILE.UPDATE,
      data
    )
    return response as any
  } catch (error: any) {
    throw new Error(error.response?.data?.message || error.message || 'Failed to update profile')
  }
}

/**
 * Check if user session is valid
 */
export async function checkSession(): Promise<boolean> {
  try {
    const { data: { session }, error } = await supabase.auth.getSession()

    if (error || !session) {
      return false
    }

    return true
  } catch (error) {
    return false
  }
}

/**
 * Get current session
 */
export async function getSession() {
  try {
    const { data: { session }, error } = await supabase.auth.getSession()

    if (error) {
      throw new Error(error.message)
    }

    return session
  } catch (error: any) {
    throw new Error(error.message || 'Failed to get session')
  }
}

/**
 * Refresh the current session
 */
export async function refreshSession() {
  try {
    const { data: { session }, error } = await supabase.auth.refreshSession()

    if (error) {
      throw new Error(error.message)
    }

    return session
  } catch (error: any) {
    throw new Error(error.message || 'Failed to refresh session')
  }
}

/**
 * Send password reset email
 */
export async function resetPassword(email: string): Promise<void> {
  try {
    const { error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: `${window.location.origin}/auth/reset-password`,
    })

    if (error) {
      throw new Error(error.message)
    }
  } catch (error: any) {
    throw new Error(error.message || 'Failed to send password reset email')
  }
}

/**
 * Update password
 */
export async function updatePassword(newPassword: string): Promise<void> {
  try {
    const { error } = await supabase.auth.updateUser({
      password: newPassword,
    })

    if (error) {
      throw new Error(error.message)
    }
  } catch (error: any) {
    throw new Error(error.message || 'Failed to update password')
  }
}

export const authApi = {
  login,
  register,
  logout,
  getProfile,
  updateProfile,
  checkSession,
  getSession,
  refreshSession,
  resetPassword,
  updatePassword,
}