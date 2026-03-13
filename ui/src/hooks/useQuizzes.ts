import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getQuizzes } from "@/lib/api/quiz";
import type { Quiz, QuizFilters } from "@/types/quiz";

/**
 * Hook to fetch quizzes with optional filters
 *
 * @param filters - Optional filters for quizzes
 * @returns React Query result with quizzes data
 */
export function useQuizzes(
  filters?: QuizFilters
): UseQueryResult<Quiz[], Error> {
  return useQuery({
    queryKey: ["quizzes", filters],
    queryFn: () => getQuizzes(filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
  });
}