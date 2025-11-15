import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getQuizQuestions } from "@/lib/api/quiz";
import type { Question } from "@/types/quiz";

/**
 * Hook to fetch questions for a quiz
 *
 * @param quizId - Quiz ID
 * @returns React Query result with questions data
 */
export function useQuizQuestions(
  quizId: string
): UseQueryResult<Question[], Error> {
  return useQuery({
    queryKey: ["quiz-questions", quizId],
    queryFn: () => getQuizQuestions(quizId),
    enabled: !!quizId, // Only fetch if quiz ID is provided
    staleTime: 10 * 60 * 1000, // 10 minutes - questions don't change often
    refetchOnWindowFocus: false,
  });
}
