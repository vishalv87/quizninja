"use client";

import { useQuery, useMutation, useQueryClient, UseQueryResult, UseMutationResult } from "@tanstack/react-query";
import { getFriends, removeFriend } from "@/lib/api/friends";
import type { Friend } from "@/types/user";
import type { APIResponse } from "@/types/api";
import { toast } from "sonner";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch current user's friends list
 * Returns all accepted friend relationships
 *
 * @returns React Query result with friends data
 */
export function useFriends(): UseQueryResult<Friend[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.FRIENDS],
    queryFn: async () => {
      const response = await getFriends();
      // Extract the array from the FriendsListResponse (friends field, not data)
      return Array.isArray(response.friends) ? response.friends : [];
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to remove a friend
 * Invalidates friends list on success
 *
 * @returns Mutation function and state
 */
export function useRemoveFriend(): UseMutationResult<APIResponse<null>, Error, string, unknown> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: removeFriend,
    onSuccess: () => {
      // Invalidate and refetch friends list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.FRIENDS] });
      toast.success("Friend removed successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to remove friend");
    },
  });
}
