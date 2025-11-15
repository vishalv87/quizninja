"use client";

import { useQuery, useMutation, useQueryClient, UseQueryResult, UseMutationResult } from "@tanstack/react-query";
import {
  getFriendRequests,
  sendFriendRequest,
  acceptFriendRequest,
  declineFriendRequest,
  cancelFriendRequest
} from "@/lib/api/friends";
import type { FriendRequest, Friend } from "@/types/user";
import type { APIResponse } from "@/types/api";
import { toast } from "sonner";
import { QUERY_KEYS } from "@/lib/constants";

/**
 * Hook to fetch friend requests (both sent and received)
 * Returns all pending friend requests
 *
 * @returns React Query result with friend requests data
 */
export function useFriendRequests(): UseQueryResult<FriendRequest[], Error> {
  return useQuery({
    queryKey: [QUERY_KEYS.FRIEND_REQUESTS],
    queryFn: async () => {
      const response = await getFriendRequests();
      // Extract the array from the APIResponse wrapper
      return Array.isArray(response.data) ? response.data : [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: true,
  });
}

/**
 * Hook to send a friend request
 * Invalidates friend requests list on success
 *
 * @returns Mutation function and state
 */
export function useSendFriendRequest(): UseMutationResult<APIResponse<FriendRequest>, Error, string, unknown> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: sendFriendRequest,
    onSuccess: () => {
      // Invalidate and refetch friend requests
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.FRIEND_REQUESTS] });
      toast.success("Friend request sent successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to send friend request");
    },
  });
}

/**
 * Hook to accept a friend request
 * Invalidates friend requests and friends list on success
 *
 * @returns Mutation function and state
 */
export function useAcceptFriendRequest(): UseMutationResult<APIResponse<Friend>, Error, string, unknown> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: acceptFriendRequest,
    onSuccess: () => {
      // Invalidate and refetch both friend requests and friends list
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.FRIEND_REQUESTS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.FRIENDS] });
      toast.success("Friend request accepted");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to accept friend request");
    },
  });
}

/**
 * Hook to decline a friend request
 * Invalidates friend requests list on success
 *
 * @returns Mutation function and state
 */
export function useDeclineFriendRequest(): UseMutationResult<APIResponse<null>, Error, string, unknown> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: declineFriendRequest,
    onSuccess: () => {
      // Invalidate and refetch friend requests
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.FRIEND_REQUESTS] });
      toast.success("Friend request declined");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to decline friend request");
    },
  });
}

/**
 * Hook to cancel a sent friend request
 * Invalidates friend requests list on success
 *
 * @returns Mutation function and state
 */
export function useCancelFriendRequest(): UseMutationResult<APIResponse<null>, Error, string, unknown> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: cancelFriendRequest,
    onSuccess: () => {
      // Invalidate and refetch friend requests
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.FRIEND_REQUESTS] });
      toast.success("Friend request canceled");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to cancel friend request");
    },
  });
}
