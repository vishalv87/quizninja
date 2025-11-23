import { useQuery } from "@tanstack/react-query";
import { getLatestCompletedAttempt } from "@/lib/api/quiz";
import type { QuizAttempt } from "@/types/quiz";

/**
 * Hook to fetch completed (non-abandoned) quiz attempts for a list of quiz IDs
 * and build a map of quizId -> latestAttempt for efficient lookup
 *
 * Uses the same endpoint as the quiz detail page which correctly filters out abandoned attempts
 */
export function useCompletedQuizMap(quizIds: string[]) {
  return useQuery<Map<string, QuizAttempt>>({
    queryKey: ["completed-quiz-map", quizIds],
    queryFn: async () => {
      const completedMap = new Map<string, QuizAttempt>();

      // Fetch completion status for each quiz in parallel
      const results = await Promise.all(
        quizIds.map(async (quizId) => {
          try {
            const attempt = await getLatestCompletedAttempt(quizId);
            return { quizId, attempt };
          } catch {
            // Quiz has no completed attempt, skip it
            return { quizId, attempt: null };
          }
        })
      );

      // Build the map from successful results
      for (const { quizId, attempt } of results) {
        if (attempt && attempt.status === "completed") {
          completedMap.set(quizId, attempt);
        }
      }

      return completedMap;
    },
    enabled: quizIds.length > 0,
    staleTime: 5 * 60 * 1000, // Cache for 5 minutes - completed attempts don't change often
    refetchOnWindowFocus: false,
  });
}
