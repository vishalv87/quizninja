"use client";

import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { searchUsers, type UserSearchResult } from "@/lib/api/friends";
import type { APIResponse } from "@/types/api";

/**
 * Hook to search for users by name or email
 * Searches are only performed when query is at least 2 characters
 *
 * @param query - The search query string
 * @param enabled - Whether the query should run (defaults to query length >= 2)
 * @returns React Query result with search results
 */
export function useSearchUsers(
  query: string,
  enabled?: boolean
): UseQueryResult<UserSearchResult[], Error> {
  // Only search if query is at least 2 characters
  const shouldFetch = enabled !== undefined ? enabled : query.length >= 2;

  return useQuery({
    queryKey: ["users", "search", query],
    queryFn: async () => {
      const response = await searchUsers(query);
      // Extract the array from the APIResponse wrapper
      return Array.isArray(response.data) ? response.data : [];
    },
    enabled: shouldFetch,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false, // Don't refetch search results on window focus
  });
}
