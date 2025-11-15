"use client";

import { useMutation, useQueryClient, UseMutationResult } from "@tanstack/react-query";
import { createChallenge } from "@/lib/api/challenges";
import type { Challenge, CreateChallengeRequest } from "@/types/challenge";
import type { APIResponse } from "@/types/api";
import { toast } from "sonner";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to create a new challenge
 * Sends a challenge to another user for a specific quiz
 * Invalidates challenges list on success
 *
 * @returns Mutation function and state
 */
export function useCreateChallenge(): UseMutationResult<
  APIResponse<Challenge>,
  Error,
  CreateChallengeRequest,
  unknown
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createChallenge,
    onSuccess: () => {
      // Invalidate and refetch challenges
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.CHALLENGES] });
      toast.success("Challenge sent successfully!");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create challenge");
    },
  });
}
