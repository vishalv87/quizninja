"use client";

import { useQuery, useMutation, useQueryClient, UseQueryResult, UseMutationResult } from "@tanstack/react-query";
import {
  getUserAchievements,
  getUserAchievementsById,
  checkAchievements
} from "@/lib/api/achievements";
import type { UserAchievement } from "@/types/achievement";
import type { APIResponse } from "@/types/api";
import { showMultipleAchievementToasts } from "@/components/achievement/AchievementToast";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch current user's unlocked achievements
 * Returns only achievements that the user has unlocked
 *
 * @returns React Query result with user achievements
 */
export function useUserAchievements(): UseQueryResult<UserAchievement[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.ACHIEVEMENTS, "user"],
    queryFn: async () => {
      const response = await getUserAchievements();
      return Array.isArray(response.achievements) ? response.achievements : [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch a specific user's unlocked achievements
 *
 * @param userId - The ID of the user
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with user achievements
 */
export function useUserAchievementsById(
  userId: string,
  enabled: boolean = true
): UseQueryResult<UserAchievement[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.ACHIEVEMENTS, "user", userId],
    queryFn: async () => {
      const response = await getUserAchievementsById(userId);
      return Array.isArray(response.achievements) ? response.achievements : [];
    },
    enabled: enabled && !!userId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to check for newly unlocked achievements
 * Call this after completing quizzes or other activities
 * Shows toast notifications for newly unlocked achievements
 *
 * @returns Mutation function and state
 */
export function useCheckAchievements(): UseMutationResult<
  APIResponse<UserAchievement[]>,
  Error,
  void,
  unknown
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: checkAchievements,
    onSuccess: (data) => {
      // Invalidate achievement-related queries
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.ACHIEVEMENTS] });

      // Show custom toast notification for each newly unlocked achievement
      const newAchievements = data.data;
      if (newAchievements && newAchievements.length > 0) {
        showMultipleAchievementToasts(newAchievements);
      }
    },
    onError: (error: Error) => {
      // Silent fail - don't show error toast for achievement checks
      console.error("Failed to check achievements:", error);
    },
  });
}