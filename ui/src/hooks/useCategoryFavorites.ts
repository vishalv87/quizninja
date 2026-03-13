import { useState, useEffect, useCallback } from "react";

const STORAGE_KEY = "category-favorites";

/**
 * Hook to manage category favorites using localStorage
 * Provides persistent favorites that survive page refreshes
 */
export function useCategoryFavorites() {
  const [favorites, setFavorites] = useState<Set<string>>(new Set());
  const [isLoaded, setIsLoaded] = useState(false);

  // Load favorites from localStorage on mount
  useEffect(() => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored);
        setFavorites(new Set(Array.isArray(parsed) ? parsed : []));
      }
    } catch (error) {
      console.error("Error loading category favorites:", error);
    }
    setIsLoaded(true);
  }, []);

  // Persist favorites to localStorage whenever they change
  useEffect(() => {
    if (isLoaded) {
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(Array.from(favorites)));
      } catch (error) {
        console.error("Error saving category favorites:", error);
      }
    }
  }, [favorites, isLoaded]);

  const isFavorite = useCallback(
    (categoryId: string) => favorites.has(categoryId),
    [favorites]
  );

  const toggleFavorite = useCallback((categoryId: string) => {
    setFavorites((prev) => {
      const next = new Set(prev);
      if (next.has(categoryId)) {
        next.delete(categoryId);
      } else {
        next.add(categoryId);
      }
      return next;
    });
  }, []);

  const addFavorite = useCallback((categoryId: string) => {
    setFavorites((prev) => new Set([...Array.from(prev), categoryId]));
  }, []);

  const removeFavorite = useCallback((categoryId: string) => {
    setFavorites((prev) => {
      const next = new Set(prev);
      next.delete(categoryId);
      return next;
    });
  }, []);

  return {
    favorites,
    isFavorite,
    toggleFavorite,
    addFavorite,
    removeFavorite,
    isLoaded,
  };
}
