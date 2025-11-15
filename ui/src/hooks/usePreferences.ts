"use client";

import { useQuery, useMutation, useQueryClient, UseQueryResult, UseMutationResult } from "@tanstack/react-query";
import { getPreferences, updatePreferences } from "@/lib/api/preferences";
import type { UserPreferences } from "@/types/user";
import type { APIResponse } from "@/types/api";
import { toast } from "sonner";

/**
 * Hook to fetch current user's preferences
 * Returns user preferences including categories, difficulty, notifications, etc.
 *
 * @returns React Query result with user preferences data
 */
export function usePreferences(): UseQueryResult<APIResponse<UserPreferences>, Error> {
  return useQuery({
    queryKey: ["user", "preferences"],
    queryFn: getPreferences,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to update user preferences
 * Provides mutation function to update preferences and invalidates cache on success
 *
 * @returns React Query mutation result
 */
export function useUpdatePreferences(): UseMutationResult<
  APIResponse<UserPreferences>,
  Error,
  Partial<UserPreferences>
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: Partial<UserPreferences>) => updatePreferences(data),
    onSuccess: () => {
      // Invalidate preferences cache to trigger refetch
      queryClient.invalidateQueries({ queryKey: ["user", "preferences"] });
      toast.success("Preferences updated successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to update preferences");
    },
  });
}