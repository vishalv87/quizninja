import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi } from '@/lib/api/auth'
import { userApi } from '@/lib/api/user'
import type { Profile } from '@/types/auth'
import type { UserProfile } from '@/types/user'
import { toast } from 'sonner'

/**
 * Hook to fetch user profile
 */
export function useProfile() {
  return useQuery({
    queryKey: ['profile'],
    queryFn: async () => {
      const response = await authApi.getProfile()
      return response
    },
    retry: 1,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

/**
 * Hook to update user profile
 */
export function useUpdateProfile() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: Partial<Profile>) => {
      const response = await authApi.updateProfile(data)
      return response.data
    },
    onSuccess: (data) => {
      // Update the profile cache
      queryClient.setQueryData(['profile'], data)

      toast.success('Profile updated successfully', {
        description: 'Your changes have been saved.',
      })
    },
    onError: (error: Error) => {
      toast.error('Failed to update profile', {
        description: error.message || 'Please try again later.',
      })
    },
  })
}

/**
 * Hook to invalidate and refetch profile
 */
export function useRefreshProfile() {
  const queryClient = useQueryClient()

  return () => {
    queryClient.invalidateQueries({ queryKey: ['profile'] })
  }
}

/**
 * Hook to fetch another user's profile by user ID
 * Respects privacy settings and returns privacy-aware data
 */
export function useUserProfile(userId: string | undefined) {
  return useQuery({
    queryKey: ['user-profile', userId],
    queryFn: async () => {
      if (!userId) {
        throw new Error('User ID is required')
      }
      const response = await userApi.getUserProfile(userId)
      return response.data
    },
    enabled: !!userId, // Only fetch when userId is provided
    retry: 1,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}