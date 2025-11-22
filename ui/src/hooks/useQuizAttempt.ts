import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  startQuizAttempt,
  submitQuizAttempt,
  updateQuizAttempt,
  getAttemptDetails,
  getQuiz,
  getActiveAttemptForQuiz,
  abandonQuizSession,
} from "@/lib/api/quiz";
import { toast } from "sonner";
import type {
  QuizAnswer,
  QuizAttempt,
  QuizResults,
} from "@/types/quiz";

/**
 * Hook to start a quiz attempt
 */
export function useStartQuizAttempt() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (quizId: string) => startQuizAttempt(quizId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-attempts"] });
      queryClient.invalidateQueries({ queryKey: ["active-attempt"] });
    },
    onError: (error: any) => {
      toast.error("Failed to start quiz", {
        description: error.message || "Could not start the quiz. Please try again.",
      });
    },
  });
}

/**
 * Hook to submit a quiz attempt
 */
export function useSubmitQuizAttempt() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      quizId,
      attemptId,
      answers,
    }: {
      quizId: string;
      attemptId: string;
      answers: QuizAnswer[];
    }) => submitQuizAttempt(quizId, attemptId, answers),
    onSuccess: (data: QuizResults) => {
      queryClient.invalidateQueries({ queryKey: ["user-attempts"] });
      queryClient.invalidateQueries({ queryKey: ["active-attempt"] });
      queryClient.invalidateQueries({ queryKey: ["user-stats"] });

      toast.success("Quiz submitted successfully!", {
        description: `You scored ${data.percentage.toFixed(1)}%`,
      });
    },
    onError: (error: any) => {
      toast.error("Failed to submit quiz", {
        description: error.message || "Could not submit the quiz. Please try again.",
      });
    },
  });
}

/**
 * Hook to update quiz attempt
 */
export function useUpdateQuizAttempt() {
  return useMutation({
    mutationFn: ({
      quizId,
      attemptId,
      answers,
    }: {
      quizId: string;
      attemptId: string;
      answers: QuizAnswer[];
    }) => updateQuizAttempt(quizId, attemptId, answers),
    onError: (error: any) => {
      console.error("Failed to update quiz attempt:", error);
    },
  });
}

/**
 * Hook to fetch quiz attempt results
 * Fetches attempt details and quiz info, then constructs QuizResults object
 */
export function useQuizAttemptResults(attemptId: string) {
  return useQuery({
    queryKey: ["quiz-attempt-results", attemptId],
    queryFn: async (): Promise<QuizResults> => {
      // Fetch attempt details
      const attempt = await getAttemptDetails(attemptId);

      // Fetch quiz details
      const quiz = await getQuiz(attempt.quiz_id);

      // Calculate percentage and passed status
      const percentage = (attempt.correct_answers / attempt.total_questions) * 100;
      const passed = percentage >= 60; // Assuming 60% is passing

      return {
        attempt,
        quiz,
        percentage,
        passed,
      };
    },
    enabled: !!attemptId,
    staleTime: 0, // Always fetch fresh results
    refetchOnWindowFocus: false,
  });
}

/**
 * Hook to fetch active attempt for a specific quiz
 * Returns null if no active attempt exists
 */
export function useActiveAttempt(quizId: string) {
  return useQuery<QuizAttempt | null>({
    queryKey: ["active-attempt", quizId],
    queryFn: () => getActiveAttemptForQuiz(quizId),
    enabled: !!quizId,
    staleTime: 0, // Always check for fresh data
    refetchOnWindowFocus: true, // Refetch when user comes back to page
  });
}

/**
 * Hook to abandon a quiz session
 */
export function useAbandonQuizSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      quizId,
      attemptId,
    }: {
      quizId: string;
      attemptId: string;
    }) => abandonQuizSession(quizId, attemptId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["active-attempt"] });
      queryClient.invalidateQueries({ queryKey: ["user-attempts"] });
    },
    onError: (error: any) => {
      toast.error("Failed to abandon quiz", {
        description: error.message || "Could not abandon the quiz.",
      });
    },
  });
}
