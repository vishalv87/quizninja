"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import {
  getAchievements,
  getAchievementsByCategory
} from "@/lib/api/achievements";
import type { Achievement } from "@/types/achievement";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch all available achievements
 * Returns all achievements in the system
 *
 * @returns React Query result with achievements data
 */
export function useAchievements(): UseQueryResult<Achievement[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.ACHIEVEMENTS],
    queryFn: async () => {
      const response = await getAchievements();
      return Array.isArray(response.achievements) ? response.achievements : [];
    },
    staleTime: 10 * 60 * 1000, // 10 minutes (achievements don't change often)
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch achievements filtered by category
 *
 * @param category - The category to filter by
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with filtered achievements
 */
export function useAchievementsByCategory(
  category: string,
  enabled: boolean = true
): UseQueryResult<Achievement[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.ACHIEVEMENTS, "category", category],
    queryFn: async () => {
      const response = await getAchievementsByCategory(category);
      return Array.isArray(response.achievements) ? response.achievements : [];
    },
    enabled: enabled && !!category,
    staleTime: 10 * 60 * 1000, // 10 minutes
    refetchOnWindowFocus: false,
  });
}
