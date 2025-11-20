"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getUserStats, getUserAttempts, getActiveSessions, type UserAttemptFilters, type ActiveSession, type BackendActiveSession } from "@/lib/api/user";
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

/**
 * Hook to fetch user's active quiz sessions
 * Returns all in-progress or paused quiz attempts
 * Note: Each user has only ONE attempt per quiz (can be in_progress, paused, completed, or abandoned)
 *
 * @returns React Query result with active sessions data
 */
export function useActiveSessions(): UseQueryResult<ActiveSession[], Error> {
  return useQuery({
    queryKey: ["user", "active-sessions"],
    queryFn: async () => {
      const response = await getActiveSessions();
      const backendSessions = response?.sessions || [];

      // Transform backend session format to frontend ActiveSession format
      const transformedSessions: ActiveSession[] = backendSessions.map((session: BackendActiveSession): ActiveSession => ({
        id: session.id,
        quiz_id: session.quiz_id,
        quiz_title: session.quiz_title,
        category: session.quiz_category,
        difficulty: session.quiz_difficulty,
        started_at: session.created_at,
        time_elapsed_seconds: session.time_spent_so_far,
        questions_answered: session.current_question_index,
        total_questions: session.total_questions,
        status: session.session_state === 'active' ? 'in_progress' : session.session_state,
      }));

      return transformedSessions;
    },
    staleTime: 1 * 60 * 1000, // 1 minute - active sessions change frequently
    refetchOnWindowFocus: true, // Refetch when user comes back
    refetchInterval: 60 * 1000, // Auto-refetch every minute
  });
}
