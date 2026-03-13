import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getCategories } from "@/lib/api/quiz";
import type { Category } from "@/types/quiz";

/**
 * Hook to fetch all quiz categories
 *
 * @returns React Query result with categories data
 */
export function useCategories(): UseQueryResult<Category[], Error> {
  return useQuery({
    queryKey: ["categories"],
    queryFn: getCategories,
    staleTime: 30 * 60 * 1000, // 30 minutes (categories rarely change)
    refetchOnWindowFocus: false,
  });
}