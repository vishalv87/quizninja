import { apiClient } from "./client";
import { API_ENDPOINTS } from "./endpoints";
import { apiLogger } from "@/lib/logger";
import type { FavoritesApiResponse, FavoritesListResponse } from "@/types/favorites";

/**
 * Favorites API Service
 * Handles all favorite quiz-related API calls
 */

/**
 * Add a quiz to favorites
 */
export async function addToFavorites(quizId: string): Promise<void> {
  try {
    await apiClient.post(API_ENDPOINTS.FAVORITES.ADD, { quiz_id: quizId });
  } catch (error) {
    apiLogger.error("Error adding to favorites", { quizId, error });
    throw error;
  }
}

/**
 * Remove a quiz from favorites
 */
export async function removeFromFavorites(quizId: string): Promise<void> {
  try {
    await apiClient.delete(API_ENDPOINTS.FAVORITES.REMOVE(quizId));
  } catch (error) {
    apiLogger.error("Error removing from favorites", { quizId, error });
    throw error;
  }
}

/**
 * Get list of favorite quizzes
 * Returns paginated favorites with full quiz details embedded
 */
export async function getFavorites(): Promise<FavoritesListResponse> {
  try {
    // The apiClient interceptor unwraps axios response, returning backend payload directly
    // Backend wraps response in { data: ... }, so the actual response is FavoritesApiResponse
    // We cast through unknown because TypeScript types don't reflect the interceptor behavior
    const response = await apiClient.get<FavoritesApiResponse>(API_ENDPOINTS.FAVORITES.LIST);
    return (response as unknown as FavoritesApiResponse).data;
  } catch (error) {
    apiLogger.error("Error fetching favorites", error);
    throw error;
  }
}

/**
 * Check if a quiz is favorited
 */
export async function checkIsFavorite(quizId: string): Promise<boolean> {
  try {
    const response = await apiClient.get<{ data: { is_favorite: boolean } }>(
      API_ENDPOINTS.FAVORITES.CHECK(quizId)
    );
    // Cast through unknown because TypeScript types don't reflect the interceptor behavior
    return (response as unknown as { data: { is_favorite: boolean } }).data.is_favorite;
  } catch (error) {
    apiLogger.error("Error checking favorite status", { quizId, error });
    throw error;
  }
}
