import { apiClient } from "./client";
import { API_ENDPOINTS } from "./endpoints";
import type {
  Discussion,
  DiscussionReply,
  CreateDiscussionRequest,
  CreateDiscussionReplyRequest,
} from "@/types/discussion";
import { apiLogger } from "@/lib/logger";

/**
 * Discussions API Service
 * Handles all discussion-related API calls
 */

export interface DiscussionFilters {
  quiz_id?: string;
  sort?: "recent" | "popular";
  limit?: number;
  offset?: number;
}

export interface DiscussionStats {
  total_discussions: number;
  total_replies: number;
  user_discussions: number;
}

// ============ DISCUSSIONS ============

/**
 * Get all discussions with optional filters
 */
export async function getDiscussions(
  filters?: DiscussionFilters
): Promise<Discussion[]> {
  try {
    apiLogger.debug("Fetching discussions with filters", filters);
    const response = (await apiClient.get<{ discussions: Discussion[] }>(
      API_ENDPOINTS.DISCUSSIONS.LIST,
      {
        params: filters,
      }
    )) as unknown as { discussions: Discussion[] };
    apiLogger.debug("Discussions fetched successfully", {
      count: response.discussions.length,
    });
    return response.discussions;
  } catch (error) {
    apiLogger.error("Error fetching discussions", error);
    throw error;
  }
}

/**
 * Get a single discussion by ID with replies
 */
export async function getDiscussion(id: string): Promise<Discussion> {
  try {
    apiLogger.debug("Fetching discussion", { discussionId: id });
    const response = (await apiClient.get<{ discussion: Discussion }>(
      API_ENDPOINTS.DISCUSSIONS.GET(id)
    )) as unknown as { discussion: Discussion };
    apiLogger.debug("Discussion fetched successfully", { discussionId: id });
    return response.discussion;
  } catch (error) {
    apiLogger.error("Error fetching discussion", { discussionId: id, error });
    throw error;
  }
}

/**
 * Create a new discussion
 */
export async function createDiscussion(
  data: CreateDiscussionRequest
): Promise<Discussion> {
  try {
    apiLogger.debug("Creating discussion", { quizId: data.quiz_id });
    const response = (await apiClient.post<{ discussion: Discussion }>(
      API_ENDPOINTS.DISCUSSIONS.CREATE,
      data
    )) as unknown as { discussion: Discussion };
    apiLogger.debug("Discussion created successfully", {
      discussionId: response.discussion.id,
    });
    return response.discussion;
  } catch (error) {
    apiLogger.error("Error creating discussion", error);
    throw error;
  }
}

/**
 * Update a discussion
 */
export async function updateDiscussion(
  id: string,
  data: Partial<CreateDiscussionRequest>
): Promise<Discussion> {
  try {
    apiLogger.debug("Updating discussion", { discussionId: id });
    const response = (await apiClient.put<{ discussion: Discussion }>(
      API_ENDPOINTS.DISCUSSIONS.UPDATE(id),
      data
    )) as unknown as { discussion: Discussion };
    apiLogger.debug("Discussion updated successfully", { discussionId: id });
    return response.discussion;
  } catch (error) {
    apiLogger.error("Error updating discussion", { discussionId: id, error });
    throw error;
  }
}

/**
 * Delete a discussion
 */
export async function deleteDiscussion(id: string): Promise<void> {
  try {
    apiLogger.debug("Deleting discussion", { discussionId: id });
    await apiClient.delete(API_ENDPOINTS.DISCUSSIONS.DELETE(id));
    apiLogger.debug("Discussion deleted successfully", { discussionId: id });
  } catch (error) {
    apiLogger.error("Error deleting discussion", { discussionId: id, error });
    throw error;
  }
}

/**
 * Like or unlike a discussion
 */
export async function likeDiscussion(id: string): Promise<void> {
  try {
    apiLogger.debug("Toggling discussion like", { discussionId: id });
    await apiClient.put(API_ENDPOINTS.DISCUSSIONS.LIKE(id));
    apiLogger.debug("Discussion like toggled successfully", {
      discussionId: id,
    });
  } catch (error) {
    apiLogger.error("Error toggling discussion like", {
      discussionId: id,
      error,
    });
    throw error;
  }
}

/**
 * Get discussion statistics
 */
export async function getDiscussionStats(): Promise<DiscussionStats> {
  try {
    apiLogger.debug("Fetching discussion stats");
    const response = (await apiClient.get<DiscussionStats>(
      API_ENDPOINTS.DISCUSSIONS.STATS
    )) as unknown as DiscussionStats;
    apiLogger.debug("Discussion stats fetched successfully");
    return response;
  } catch (error) {
    apiLogger.error("Error fetching discussion stats", error);
    throw error;
  }
}

// ============ REPLIES ============

/**
 * Get replies for a discussion
 */
export async function getDiscussionReplies(
  discussionId: string
): Promise<DiscussionReply[]> {
  try {
    apiLogger.debug("Fetching discussion replies", { discussionId });
    const response = (await apiClient.get<{ replies: DiscussionReply[] }>(
      API_ENDPOINTS.DISCUSSIONS.REPLIES.LIST(discussionId)
    )) as unknown as { replies: DiscussionReply[] };
    apiLogger.debug("Discussion replies fetched successfully", {
      discussionId,
      count: response.replies.length,
    });
    return response.replies;
  } catch (error) {
    apiLogger.error("Error fetching discussion replies", {
      discussionId,
      error,
    });
    throw error;
  }
}

/**
 * Create a reply to a discussion
 */
export async function createDiscussionReply(
  discussionId: string,
  data: CreateDiscussionReplyRequest
): Promise<DiscussionReply> {
  try {
    apiLogger.debug("Creating discussion reply", { discussionId });
    const response = (await apiClient.post<DiscussionReply>(
      API_ENDPOINTS.DISCUSSIONS.REPLIES.CREATE(discussionId),
      data
    )) as unknown as DiscussionReply;
    apiLogger.debug("Discussion reply created successfully", {
      discussionId,
      replyId: response.id,
    });
    return response;
  } catch (error) {
    apiLogger.error("Error creating discussion reply", {
      discussionId,
      error,
    });
    throw error;
  }
}

/**
 * Update a reply
 */
export async function updateDiscussionReply(
  replyId: string,
  data: CreateDiscussionReplyRequest
): Promise<DiscussionReply> {
  try {
    apiLogger.debug("Updating discussion reply", { replyId });
    const response = (await apiClient.put<DiscussionReply>(
      API_ENDPOINTS.DISCUSSIONS.REPLIES.UPDATE(replyId),
      data
    )) as unknown as DiscussionReply;
    apiLogger.debug("Discussion reply updated successfully", { replyId });
    return response;
  } catch (error) {
    apiLogger.error("Error updating discussion reply", { replyId, error });
    throw error;
  }
}

/**
 * Delete a reply
 */
export async function deleteDiscussionReply(replyId: string): Promise<void> {
  try {
    apiLogger.debug("Deleting discussion reply", { replyId });
    await apiClient.delete(API_ENDPOINTS.DISCUSSIONS.REPLIES.DELETE(replyId));
    apiLogger.debug("Discussion reply deleted successfully", { replyId });
  } catch (error) {
    apiLogger.error("Error deleting discussion reply", { replyId, error });
    throw error;
  }
}

/**
 * Like or unlike a reply
 */
export async function likeDiscussionReply(replyId: string): Promise<void> {
  try {
    apiLogger.debug("Toggling discussion reply like", { replyId });
    await apiClient.put(API_ENDPOINTS.DISCUSSIONS.REPLIES.LIKE(replyId));
    apiLogger.debug("Discussion reply like toggled successfully", { replyId });
  } catch (error) {
    apiLogger.error("Error toggling discussion reply like", {
      replyId,
      error,
    });
    throw error;
  }
}