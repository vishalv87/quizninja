"use client";

import { useMutation, useQueryClient, UseMutationResult } from "@tanstack/react-query";
import {
  acceptChallenge,
  declineChallenge,
  linkAttempt,
  completeChallenge
} from "@/lib/api/challenges";
import type { Challenge } from "@/types/challenge";
import type { APIResponse } from "@/types/api";
import { toast } from "sonner";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to accept a challenge
 * Invalidates challenges lists on success
 *
 * @returns Mutation function and state
 */
export function useAcceptChallenge(): UseMutationResult<
  APIResponse<Challenge>,
  Error,
  string,
  unknown
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: acceptChallenge,
    onSuccess: () => {
      // Invalidate challenges to refresh the lists
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CHALLENGES] });
      toast.success("Challenge accepted! Good luck!");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to accept challenge");
    },
  });
}

/**
 * Hook to decline a challenge
 * Invalidates challenges lists on success
 *
 * @returns Mutation function and state
 */
export function useDeclineChallenge(): UseMutationResult<
  APIResponse<Challenge>,
  Error,
  string,
  unknown
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: declineChallenge,
    onSuccess: () => {
      // Invalidate challenges to refresh the lists
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CHALLENGES] });
      toast.success("Challenge declined");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to decline challenge");
    },
  });
}

/**
 * Interface for link attempt mutation
 */
interface LinkAttemptParams {
  challengeId: string;
  attemptId: string;
}

/**
 * Hook to link a quiz attempt to a challenge
 * This is called after completing a quiz for a challenge
 * Invalidates challenge data on success
 *
 * @returns Mutation function and state
 */
export function useLinkAttempt(): UseMutationResult<
  APIResponse<Challenge>,
  Error,
  LinkAttemptParams,
  unknown
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ challengeId, attemptId }: LinkAttemptParams) =>
      linkAttempt(challengeId, attemptId),
    onSuccess: () => {
      // Invalidate challenges to refresh the data
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CHALLENGES] });
      toast.success("Quiz attempt linked to challenge");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to link attempt to challenge");
    },
  });
}

/**
 * Hook to complete a challenge
 * Marks the challenge as finished
 * Invalidates challenges lists on success
 *
 * @returns Mutation function and state
 */
export function useCompleteChallenge(): UseMutationResult<
  APIResponse<Challenge>,
  Error,
  string,
  unknown
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: completeChallenge,
    onSuccess: () => {
      // Invalidate challenges to refresh the lists
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CHALLENGES] });
      toast.success("Challenge completed!");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to complete challenge");
    },
  });
}
