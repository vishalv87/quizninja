"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getUserStats, getUserAttempts, type UserAttemptFilters } from "@/lib/api/user";
import type { UserStats } from "@/types/user";
import type { QuizAttempt } from "@/types/quiz";
import type { APIResponse } from "@/types/api";

/**
 * Hook to fetch current user's statistics
 * Returns comprehensive statistics including points, quizzes taken, achievements, etc.
 *
 * @returns React Query result with user stats data
 */
export function useUserStats(): UseQueryResult<APIResponse<UserStats>, Error> {
  return useQuery({
    queryKey: ["user", "stats"],
    queryFn: getUserStats,
    staleTime: 2 * 60 * 1000, // 2 minutes - stats change frequently
    refetchOnWindowFocus: true, // Refetch when user comes back to the app
  });
}

/**
 * Hook to fetch user's quiz attempt history with optional filters
 * Supports pagination and filtering by category, difficulty, status, and date range
 * Note: Each user has only ONE attempt per quiz (can be in_progress, paused, completed, or abandoned)
 *
 * @param filters - Optional filters for attempts (pagination, category, difficulty, etc.)
 * @returns React Query result with attempts data (one attempt per quiz)
 */
export function useUserAttempts(
  filters?: UserAttemptFilters
): UseQueryResult<QuizAttempt[], Error> {
  return useQuery({
    queryKey: ["user", "attempts", filters],
    queryFn: async () => {
      const response = await getUserAttempts(filters);
      // Extract the data array from the PaginatedResponse wrapper
      // Add safety fallback to ensure we never return undefined
      return response?.data || [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}
