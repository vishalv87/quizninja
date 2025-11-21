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
      console.log("[DEBUGGG] useQuizQuestions - Fetching questions for quizId:", quizId);
      const questions = await getQuizQuestions(quizId);
      console.log("[DEBUGGG] useQuizQuestions - Received questions:", questions);
      return questions;
    },
    enabled: !!quizId, // Only fetch if quiz ID is provided
    staleTime: 10 * 60 * 1000, // 10 minutes - questions don't change often
    refetchOnWindowFocus: false,
  });

  console.log("[DEBUGGG] useQuizQuestions - Hook result.data:", result.data);
  console.log("[DEBUGGG] useQuizQuestions - Hook result.isLoading:", result.isLoading);
  console.log("[DEBUGGG] useQuizQuestions - Hook result.error:", result.error);

  return result;
}
