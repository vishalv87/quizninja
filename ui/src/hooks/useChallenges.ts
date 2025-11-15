"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import {
  getChallenges,
  getPendingChallenges,
  getActiveChallenges,
  getCompletedChallenges,
  getChallenge
} from "@/lib/api/challenges";
import type { Challenge } from "@/types/challenge";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch all challenges for the current user
 * Returns challenges where user is either challenger or opponent
 *
 * @returns React Query result with challenges data
 */
export function useChallenges(): UseQueryResult<Challenge[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.CHALLENGES],
    queryFn: async () => {
      const response = await getChallenges();
      return Array.isArray(response.challenges) ? response.challenges : [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch pending challenges (waiting for acceptance)
 * Returns challenges with status 'pending'
 *
 * @returns React Query result with pending challenges
 */
export function usePendingChallenges(): UseQueryResult<Challenge[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.CHALLENGES, 'pending'],
    queryFn: async () => {
      const response = await getPendingChallenges();
      return Array.isArray(response.challenges) ? response.challenges : [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch active challenges (accepted and in progress)
 * Returns challenges with status 'accepted'
 *
 * @returns React Query result with active challenges
 */
export function useActiveChallenges(): UseQueryResult<Challenge[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.CHALLENGES, 'active'],
    queryFn: async () => {
      const response = await getActiveChallenges();
      return Array.isArray(response.challenges) ? response.challenges : [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch completed challenges
 * Returns challenges with status 'completed'
 *
 * @returns React Query result with completed challenges
 */
export function useCompletedChallenges(): UseQueryResult<Challenge[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.CHALLENGES, 'completed'],
    queryFn: async () => {
      const response = await getCompletedChallenges();
      return Array.isArray(response.challenges) ? response.challenges : [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutes (completed challenges change less frequently)
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to fetch a specific challenge by ID
 *
 * @param challengeId - The ID of the challenge to fetch
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with challenge data
 */
export function useChallenge(
  challengeId: string,
  enabled: boolean = true
): UseQueryResult<Challenge, Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.CHALLENGE, challengeId],
    queryFn: async () => {
      const response = await getChallenge(challengeId);
      return response.data;
    },
    enabled: enabled && !!challengeId,
    staleTime: 1 * 60 * 1000, // 1 minute
    refetchOnWindowFocus: true,
  });
}
