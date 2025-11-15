import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getFeaturedQuizzes } from "@/lib/api/quiz";
import type { Quiz } from "@/types/quiz";

/**
 * Hook to fetch featured quizzes
 *
 * @returns React Query result with featured quizzes data
 */
export function useFeaturedQuizzes(): UseQueryResult<Quiz[], Error> {
  return useQuery({
    queryKey: ["quizzes", "featured"],
    queryFn: getFeaturedQuizzes,
    staleTime: 10 * 60 * 1000, // 10 minutes (featured quizzes change less frequently)
    refetchOnWindowFocus: false,
  });
}