"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import {
  getAchievementProgress,
  getAchievementStats
} from "@/lib/api/achievements";
import type { AchievementProgress, AchievementStats } from "@/types/achievement";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch user's achievement progress
 * Returns progress on all achievements (both locked and unlocked)
 *
 * @returns React Query result with achievement progress data
 */
export function useAchievementProgress(): UseQueryResult<AchievementProgress[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.ACHIEVEMENTS, "progress"],
    queryFn: async () => {
      const response = await getAchievementProgress();
      return Array.isArray(response.progress) ? response.progress : [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch user's achievement statistics
 * Returns stats like total achievements, unlocked count, points earned
 *
 * @returns React Query result with achievement stats
 */
export function useAchievementStats(): UseQueryResult<AchievementStats, Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.ACHIEVEMENTS, "stats"],
    queryFn: async () => {
      const response = await getAchievementStats();
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true,
  });
}