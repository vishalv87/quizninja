"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getDiscussions,
  getDiscussion,
  createDiscussion,
  updateDiscussion,
  deleteDiscussion,
  likeDiscussion,
  getDiscussionStats,
  getDiscussionReplies,
  createDiscussionReply,
  updateDiscussionReply,
  deleteDiscussionReply,
  likeDiscussionReply,
  type DiscussionFilters,
} from "@/lib/api/discussions";
import type {
  Discussion,
  DiscussionReply,
  CreateDiscussionRequest,
  CreateDiscussionReplyRequest,
} from "@/types/discussion";
import { QUERY_KEYS } from "@/lib/constants";
import { toast } from "sonner";

/**
 * Hook to fetch all discussions with optional filters
 *
 * @param filters - Optional filters for discussions (quiz_id, sort, limit, offset)
 * @returns React Query result with discussions data
 */
export function useDiscussions(filters?: DiscussionFilters) {
  return useQuery({
    queryKey: [QUERY_KEYS.DISCUSSIONS, filters],
    queryFn: () => getDiscussions(filters),
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch a specific discussion by ID
 *
 * @param discussionId - The ID of the discussion to fetch
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with discussion data
 */
export function useDiscussion(discussionId: string, enabled: boolean = true) {
  return useQuery({
    queryKey: [QUERY_KEYS.DISCUSSION, discussionId],
    queryFn: () => getDiscussion(discussionId),
    enabled: enabled && !!discussionId,
    staleTime: 1 * 60 * 1000, // 1 minute
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch discussion statistics
 *
 * @returns React Query result with discussion stats
 */
export function useDiscussionStats() {
  return useQuery({
    queryKey: [QUERY_KEYS.DISCUSSIONS, "stats"],
    queryFn: getDiscussionStats,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

/**
 * Hook to fetch replies for a specific discussion
 *
 * @param discussionId - The ID of the discussion
 * @param enabled - Whether to enable the query (default: true)
 * @returns React Query result with replies data
 */
export function useDiscussionReplies(
  discussionId: string,
  enabled: boolean = true
) {
  return useQuery({
    queryKey: [QUERY_KEYS.DISCUSSION, discussionId, "replies"],
    queryFn: () => getDiscussionReplies(discussionId),
    enabled: enabled && !!discussionId,
    staleTime: 1 * 60 * 1000, // 1 minute
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to create a new discussion
 *
 * @returns Mutation object with mutate function and status
 */
export function useCreateDiscussion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDiscussionRequest) => createDiscussion(data),
    onSuccess: (newDiscussion) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS, "stats"] });
      toast.success("Discussion created successfully");
      return newDiscussion;
    },
    onError: (error: any) => {
      toast.error("Failed to create discussion", {
        description: error.message || "Could not create the discussion.",
      });
    },
  });
}

/**
 * Hook to update a discussion
 *
 * @returns Mutation object with mutate function and status
 */
export function useUpdateDiscussion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: string;
      data: Partial<CreateDiscussionRequest>;
    }) => updateDiscussion(id, data),
    onSuccess: (updatedDiscussion) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS] });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, updatedDiscussion.id],
      });
      toast.success("Discussion updated successfully");
    },
    onError: (error: any) => {
      toast.error("Failed to update discussion", {
        description: error.message || "Could not update the discussion.",
      });
    },
  });
}

/**
 * Hook to delete a discussion
 *
 * @returns Mutation object with mutate function and status
 */
export function useDeleteDiscussion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (discussionId: string) => deleteDiscussion(discussionId),
    onSuccess: (_, discussionId) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS] });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS, "stats"] });
      queryClient.removeQueries({ queryKey: [QUERY_KEYS.DISCUSSION, discussionId] });
      toast.success("Discussion deleted successfully");
    },
    onError: (error: any) => {
      toast.error("Failed to delete discussion", {
        description: error.message || "Could not delete the discussion.",
      });
    },
  });
}

/**
 * Hook to like/unlike a discussion
 *
 * @returns Mutation object with mutate function and status
 */
export function useLikeDiscussion() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (discussionId: string) => likeDiscussion(discussionId),
    onSuccess: (_, discussionId) => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS] });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId],
      });
    },
    onError: (error: any) => {
      toast.error("Failed to update like status", {
        description: error.message || "Could not update the like status.",
      });
    },
  });
}

/**
 * Hook to create a reply to a discussion
 *
 * @returns Mutation object with mutate function and status
 */
export function useCreateReply() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      discussionId,
      data,
    }: {
      discussionId: string;
      data: CreateDiscussionReplyRequest;
    }) => createDiscussionReply(discussionId, data),
    onSuccess: (_, { discussionId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId, "replies"],
      });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS] });
      toast.success("Reply posted successfully");
    },
    onError: (error: any) => {
      toast.error("Failed to post reply", {
        description: error.message || "Could not post the reply.",
      });
    },
  });
}

/**
 * Hook to update a discussion reply
 *
 * @returns Mutation object with mutate function and status
 */
export function useUpdateReply() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      replyId,
      discussionId,
      data,
    }: {
      replyId: string;
      discussionId: string;
      data: CreateDiscussionReplyRequest;
    }) => updateDiscussionReply(replyId, data),
    onSuccess: (_, { discussionId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId, "replies"],
      });
      toast.success("Reply updated successfully");
    },
    onError: (error: any) => {
      toast.error("Failed to update reply", {
        description: error.message || "Could not update the reply.",
      });
    },
  });
}

/**
 * Hook to delete a discussion reply
 *
 * @returns Mutation object with mutate function and status
 */
export function useDeleteReply() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      replyId,
      discussionId,
    }: {
      replyId: string;
      discussionId: string;
    }) => deleteDiscussionReply(replyId),
    onSuccess: (_, { discussionId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId, "replies"],
      });
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.DISCUSSIONS] });
      toast.success("Reply deleted successfully");
    },
    onError: (error: any) => {
      toast.error("Failed to delete reply", {
        description: error.message || "Could not delete the reply.",
      });
    },
  });
}

/**
 * Hook to like/unlike a discussion reply
 *
 * @returns Mutation object with mutate function and status
 */
export function useLikeReply() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      replyId,
      discussionId,
    }: {
      replyId: string;
      discussionId: string;
    }) => likeDiscussionReply(replyId),
    onSuccess: (_, { discussionId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.DISCUSSION, discussionId, "replies"],
      });
    },
    onError: (error: any) => {
      toast.error("Failed to update like status", {
        description: error.message || "Could not update the like status.",
      });
    },
  });
}