"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getChallengeStats } from "@/lib/api/challenges";
import type { ChallengeStats } from "@/types/challenge";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch challenge statistics for the current user
 * Returns stats like total challenges, wins, losses, and win rate
 *
 * @returns React Query result with challenge stats
 */
export function useChallengeStats(): UseQueryResult<ChallengeStats, Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.CHALLENGES, 'stats'],
    queryFn: async () => {
      const response = await getChallengeStats();
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true,
  });
}
