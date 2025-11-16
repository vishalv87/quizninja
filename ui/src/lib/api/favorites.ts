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
    apiLogger.debug("Adding quiz to favorites", { quizId });
    await apiClient.post(API_ENDPOINTS.FAVORITES.ADD, { quiz_id: quizId });
    apiLogger.debug("Quiz added to favorites", { quizId });
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
    apiLogger.debug("Removing quiz from favorites", { quizId });
    await apiClient.delete(API_ENDPOINTS.FAVORITES.REMOVE(quizId));
    apiLogger.debug("Quiz removed from favorites", { quizId });
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
    apiLogger.debug("Fetching favorites");
    const response = await apiClient.get<FavoritesApiResponse>(API_ENDPOINTS.FAVORITES.LIST);
    apiLogger.debug("Favorites fetched", { count: response.data.favorites.length });
    return response.data;
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
    apiLogger.debug("Checking if quiz is favorite", { quizId });
    const response = await apiClient.get<{ data: { is_favorite: boolean } }>(
      API_ENDPOINTS.FAVORITES.CHECK(quizId)
    );
    apiLogger.debug("Favorite status checked", {
      quizId,
      isFavorite: response.data.is_favorite,
    });
    return response.data.is_favorite;
  } catch (error) {
    apiLogger.error("Error checking favorite status", { quizId, error });
    throw error;
  }
}
