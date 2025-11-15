import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getQuiz } from "@/lib/api/quiz";
import type { Quiz } from "@/types/quiz";

/**
 * Hook to fetch a single quiz by ID
 *
 * @param id - Quiz ID
 * @returns React Query result with quiz data
 */
export function useQuiz(id: string): UseQueryResult<Quiz, Error> {
  return useQuery({
    queryKey: ["quiz", id],
    queryFn: () => getQuiz(id),
    enabled: !!id, // Only fetch if ID is provided
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}