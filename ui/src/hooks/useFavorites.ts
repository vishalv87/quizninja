import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  addToFavorites,
  removeFromFavorites,
  getFavorites,
  checkIsFavorite,
} from "@/lib/api/favorites";
import { toast } from "sonner";

/**
 * Hook to get all favorite quiz IDs
 */
export function useFavorites() {
  return useQuery({
    queryKey: ["favorites"],
    queryFn: getFavorites,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

/**
 * Hook to check if a specific quiz is favorited
 */
export function useIsFavorite(quizId: string) {
  return useQuery({
    queryKey: ["is-favorite", quizId],
    queryFn: () => checkIsFavorite(quizId),
    enabled: !!quizId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

/**
 * Hook to toggle favorite status of a quiz
 */
export function useToggleFavorite() {
  const queryClient = useQueryClient();

  const addMutation = useMutation({
    mutationFn: (quizId: string) => addToFavorites(quizId),
    onSuccess: (_, quizId) => {
      queryClient.invalidateQueries({ queryKey: ["favorites"] });
      queryClient.invalidateQueries({ queryKey: ["is-favorite", quizId] });
      toast.success("Added to favorites");
    },
    onError: (error: any) => {
      toast.error("Failed to add to favorites", {
        description: error.message || "Could not add quiz to favorites.",
      });
    },
  });

  const removeMutation = useMutation({
    mutationFn: (quizId: string) => removeFromFavorites(quizId),
    onSuccess: (_, quizId) => {
      queryClient.invalidateQueries({ queryKey: ["favorites"] });
      queryClient.invalidateQueries({ queryKey: ["is-favorite", quizId] });
      toast.success("Removed from favorites");
    },
    onError: (error: any) => {
      toast.error("Failed to remove from favorites", {
        description: error.message || "Could not remove quiz from favorites.",
      });
    },
  });

  return {
    add: addMutation,
    remove: removeMutation,
    toggle: (quizId: string, isFavorite: boolean) => {
      if (isFavorite) {
        removeMutation.mutate(quizId);
      } else {
        addMutation.mutate(quizId);
      }
    },
  };
}
