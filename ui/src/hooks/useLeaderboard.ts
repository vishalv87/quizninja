"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import {
  getLeaderboard,
  getLeaderboardWithAchievements,
  getUserRank,
  getLeaderboardStats
} from "@/lib/api/leaderboard";
import type { LeaderboardEntry } from "@/types/api";
import type { UserRankResponse, LeaderboardStats } from "@/lib/api/leaderboard";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuth } from "@/hooks/useAuth";

/**
 * Hook to fetch global leaderboard
 * Returns top users ranked by total points
 *
 * @param limit - Maximum number of entries to return (default: 50)
 * @returns React Query result with leaderboard data
 */
export function useLeaderboard(limit: number = 50): UseQueryResult<LeaderboardEntry[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.LEADERBOARD, limit],
    queryFn: async () => {
      const response = await getLeaderboard(limit);
      return Array.isArray(response.leaderboard) ? response.leaderboard : [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch leaderboard with achievement counts
 * Returns leaderboard with detailed achievement information
 *
 * @param limit - Maximum number of entries to return (default: 50)
 * @returns React Query result with leaderboard data
 */
export function useLeaderboardWithAchievements(limit: number = 50): UseQueryResult<LeaderboardEntry[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.LEADERBOARD, "achievements", limit],
    queryFn: async () => {
      const response = await getLeaderboardWithAchievements(limit);
      return Array.isArray(response.leaderboard) ? response.leaderboard : [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch current user's rank on the leaderboard
 * Returns the user's position and stats
 * Only fetches when user is authenticated
 *
 * @returns React Query result with user rank data
 */
export function useUserRank(): UseQueryResult<UserRankResponse | null, Error> {
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  return useQuery({
    queryKey: [QUERY_KEYS.LEADERBOARD, "rank"],
    queryFn: async () => {
      const response = await getUserRank();
      return response.data ?? null;
    },
    enabled: !authLoading && isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch leaderboard statistics
 * Returns overall stats about the leaderboard
 *
 * @returns React Query result with leaderboard stats
 */
export function useLeaderboardStats(): UseQueryResult<LeaderboardStats, Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.LEADERBOARD, "stats"],
    queryFn: async () => {
      const response = await getLeaderboardStats();
      return response.data ?? {
        total_users: 0,
        total_points_distributed: 0,
        average_points: 0,
        top_user: null
      };
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
    refetchOnWindowFocus: false,
  });
}