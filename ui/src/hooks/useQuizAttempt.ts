import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  startQuizAttempt,
  submitQuizAttempt,
  saveQuizProgress,
  updateQuizAttempt,
  getAttemptDetails,
  getQuiz,
  getActiveAttemptForQuiz,
  pauseQuizAttempt,
  resumeQuizAttempt,
} from "@/lib/api/quiz";
import { toast } from "sonner";
import type { QuizAnswer, QuizAttempt, QuizResults } from "@/types/quiz";

/**
 * Hook to start a quiz attempt
 */
export function useStartQuizAttempt() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (quizId: string) => startQuizAttempt(quizId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-attempts"] });
      queryClient.invalidateQueries({ queryKey: ["active-sessions"] });
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
      queryClient.invalidateQueries({ queryKey: ["active-sessions"] });
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
 * Hook to save quiz progress
 */
export function useSaveQuizProgress() {
  return useMutation({
    mutationFn: ({
      quizId,
      attemptId,
      answers,
    }: {
      quizId: string;
      attemptId: string;
      answers: QuizAnswer[];
    }) => saveQuizProgress(quizId, attemptId, answers),
    onSuccess: () => {
      toast.success("Progress saved");
    },
    onError: (error: any) => {
      toast.error("Failed to save progress", {
        description: error.message || "Could not save your progress.",
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
 * Hook to pause a quiz attempt
 */
export function usePauseQuizAttempt() {
  return useMutation({
    mutationFn: ({
      quizId,
      attemptId,
    }: {
      quizId: string;
      attemptId: string;
    }) => pauseQuizAttempt(quizId, attemptId),
    onSuccess: () => {
      toast.success("Quiz paused");
    },
    onError: (error: any) => {
      toast.error("Failed to pause quiz", {
        description: error.message || "Could not pause the quiz.",
      });
    },
  });
}

/**
 * Hook to resume a quiz attempt
 */
export function useResumeQuizAttempt() {
  return useMutation({
    mutationFn: ({
      quizId,
      attemptId,
    }: {
      quizId: string;
      attemptId: string;
    }) => resumeQuizAttempt(quizId, attemptId),
    onSuccess: () => {
      toast.success("Quiz resumed");
    },
    onError: (error: any) => {
      toast.error("Failed to resume quiz", {
        description: error.message || "Could not resume the quiz.",
      });
    },
  });
}
