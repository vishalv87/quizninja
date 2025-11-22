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
  const result = useQuery({
    queryKey: ["quiz-questions", quizId],
    queryFn: async () => {
      const questions = await getQuizQuestions(quizId);
      return questions;
    },
    enabled: !!quizId, // Only fetch if quiz ID is provided
    staleTime: 10 * 60 * 1000, // 10 minutes - questions don't change often
    refetchOnWindowFocus: false,
  });

  return result;
}
