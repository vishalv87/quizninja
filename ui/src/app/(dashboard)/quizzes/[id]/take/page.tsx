"use client";

import { useCallback, useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { useQuizQuestions } from "@/hooks/useQuizQuestions";
import { useActiveAttempt } from "@/hooks/useQuizAttempt";
import { useQuizStore } from "@/store/quizStore";
import {
  useStartQuizAttempt,
  useSubmitQuizAttempt,
  useSaveSessionProgress,
  usePauseQuizSession,
  useResumeQuizSession,
  useAbandonQuizSession,
} from "@/hooks/useQuizAttempt";
import type { PauseSessionRequest, SaveProgressRequest, AttemptAnswer } from "@/types/quiz";
import { useQuizTimer } from "@/hooks/useQuizTimer";
import { QuestionCard } from "@/components/quiz/QuestionCard";
import { QuizTimer } from "@/components/quiz/QuizTimer";
import { QuizProgress } from "@/components/quiz/QuizProgress";
import { QuizNavigation } from "@/components/quiz/QuizNavigation";
import { Button } from "@/components/ui/button";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { BookOpen, Save, X, Pause, Play } from "lucide-react";
import { toast } from "sonner";

export default function QuizTakingPage() {
  const params = useParams();
  const router = useRouter();
  const searchParams = useSearchParams();
  const quizId = params.id as string;
  const isResuming = searchParams.get("resume") === "true";

  const { data: quiz, isLoading: quizLoading } = useQuiz(quizId);
  const {
    data: questions,
    isLoading: questionsLoading,
    error: questionsError,
  } = useQuizQuestions(quizId);
  const { data: activeAttempt, isLoading: activeAttemptLoading } =
    useActiveAttempt(quizId);
  const startAttemptMutation = useStartQuizAttempt();
  const submitAttemptMutation = useSubmitQuizAttempt();
  const saveProgressMutation = useSaveSessionProgress();
  const pauseSessionMutation = usePauseQuizSession();
  const resumeSessionMutation = useResumeQuizSession();
  const abandonSessionMutation = useAbandonQuizSession();

  const {
    currentQuiz,
    currentAttempt,
    currentQuestionIndex,
    answers,
    timeRemaining,
    isPaused,
    nextQuestion,
    previousQuestion,
    setAnswer,
    getAnswer,
    startQuiz,
    resetQuiz,
    setPaused,
  } = useQuizStore();

  const [showExitDialog, setShowExitDialog] = useState(false);
  const [showSubmitDialog, setShowSubmitDialog] = useState(false);
  const [show409Dialog, setShow409Dialog] = useState(false);
  const [existingAttemptId, setExistingAttemptId] = useState<string | null>(null);

  // Initialize quiz - Following Flutter app's pattern: Check session FIRST
  useEffect(() => {
    if (quiz && questions && !currentQuiz) {
      // Case 1: User is resuming - Load the existing attempt
      if (isResuming && activeAttempt) {
        toast.info("Resuming quiz from where you left off");
        startQuiz(quiz, activeAttempt);

        // Restore saved answers to the store
        if (activeAttempt.answers && activeAttempt.answers.length > 0) {
          activeAttempt.answers.forEach((answer) => {
            setAnswer(answer.question_id, answer);
          });
        }
        return;
      }

      // Case 2: Active attempt exists BUT user is not resuming (clicked "Start New Attempt")
      // Show dialog immediately - don't try to start (following Flutter pattern)
      if (activeAttempt && !isResuming) {
        setExistingAttemptId(activeAttempt.id);
        setShow409Dialog(true);
        return;
      }

      // Case 3: No active attempt - Safe to start new quiz
      if (!activeAttempt) {
        startAttemptMutation.mutate(quizId, {
          onSuccess: (attempt) => {
            startQuiz(quiz, attempt);
          },
          onError: (error: any) => {
            // 409 error should be rare now (race condition fallback)
            if (error.response?.status === 409 || error.status === 409) {
              const errorData = error.response?.data || error.data;
              const attemptId = errorData?.error?.details?.existing_attempt_id;
              if (attemptId) {
                setExistingAttemptId(attemptId);
                setShow409Dialog(true);
              } else {
                toast.error("Active attempt exists", {
                  description: "You have an active attempt for this quiz. Please resume it.",
                });
                router.push(`/quizzes/${quizId}`);
              }
            }
          },
        });
      }
    }
  }, [quiz, questions, currentQuiz, quizId, isResuming, activeAttempt, router, setAnswer, startAttemptMutation, startQuiz]);

  // Save progress handler - defined before useEffect that uses it
  const handleSaveProgress = useCallback(() => {
    // Guard against race condition: ensure attemptId exists before saving
    if (!currentAttempt || !currentAttempt.id) return;

    const answersList: AttemptAnswer[] = Object.values(answers).map((answer: any) => ({
      question_id: answer.question_id,
      selected_answer: answer.selected_answer,
      is_correct: answer.is_correct,
      points_earned: answer.points_earned,
    }));

    const progressRequest: SaveProgressRequest = {
      current_question_index: currentQuestionIndex,
      current_answers: answersList,
      time_spent_so_far: quiz?.time_limit ?
        (quiz.time_limit * 60) - (timeRemaining || 0) : 0,
      time_remaining: timeRemaining || undefined,
    };

    saveProgressMutation.mutate({
      quizId,
      attemptId: currentAttempt.id,
      progressRequest,
    });
  }, [currentAttempt, answers, currentQuestionIndex, quiz?.time_limit, timeRemaining, saveProgressMutation, quizId]);

  // Setup quiz timer
  useQuizTimer(() => {
    // Auto-submit when time runs out
    handleSubmit();
  });

  // Auto-save progress every 10 seconds (like Flutter app)
  useEffect(() => {
    if (!currentAttempt || isPaused) return;

    const interval = setInterval(() => {
      handleSaveProgress();
    }, 10000); // 10 seconds

    return () => clearInterval(interval);
  }, [currentAttempt, isPaused, handleSaveProgress]);

  // Prevent accidental page close
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (currentAttempt) {
        e.preventDefault();
        e.returnValue = "";
      }
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => window.removeEventListener("beforeunload", handleBeforeUnload);
  }, [currentAttempt]);

  const handleAnswerChange = (answer: string) => {
    if (!questions || questions.length === 0) return;
    const question = questions[currentQuestionIndex];
    if (!question) return;

    // Calculate the selected option index for backend submission
    let selectedOptionIndex: number | undefined;
    if (question.options && question.question_type === "multiple_choice") {
      selectedOptionIndex = question.options.indexOf(answer);
      if (selectedOptionIndex === -1) selectedOptionIndex = undefined;
    }

    setAnswer(question.id, {
      question_id: question.id,
      selected_answer: answer,
      selected_option_index: selectedOptionIndex,
      is_correct: false, // Will be determined on backend
      points_earned: 0, // Will be determined on backend
    });
  };

  const handleSubmit = () => {
    setShowSubmitDialog(true);
  };

  const confirmSubmit = () => {
    if (!currentAttempt) return;

    const answersList = Object.values(answers);
    submitAttemptMutation.mutate(
      {
        quizId,
        attemptId: currentAttempt.id,
        answers: answersList,
      },
      {
        onSuccess: () => {
          resetQuiz();
          router.push(`/quizzes/${quizId}/results/${currentAttempt.id}`);
        },
      }
    );
  };

  const handleExit = () => {
    setShowExitDialog(true);
  };

  const confirmExit = () => {
    handleSaveProgress();
    resetQuiz();
    router.push(`/quizzes/${quizId}`);
  };

  const handlePause = () => {
    if (!currentAttempt) return;

    const answersList: AttemptAnswer[] = Object.values(answers).map((answer: any) => ({
      question_id: answer.question_id,
      selected_answer: answer.selected_answer,
      is_correct: answer.is_correct,
      points_earned: answer.points_earned,
    }));

    const pauseRequest: PauseSessionRequest = {
      current_question_index: currentQuestionIndex,
      current_answers: answersList,
      time_spent_so_far: quiz?.time_limit ?
        (quiz.time_limit * 60) - (timeRemaining || 0) : 0,
      time_remaining: timeRemaining || undefined,
    };

    pauseSessionMutation.mutate(
      {
        quizId,
        attemptId: currentAttempt.id,
        pauseRequest,
      },
      {
        onSuccess: () => {
          setPaused(true);
          router.push(`/quizzes/${quizId}`);
        },
      }
    );
  };

  const handleResume = () => {
    if (!currentAttempt) return;

    setPaused(false);
    resumeSessionMutation.mutate({
      quizId,
      attemptId: currentAttempt.id,
    });
  };

  // Handle resuming from 409 conflict dialog
  const handleResumeExisting = () => {
    if (!existingAttemptId) return;
    setShow409Dialog(false);
    router.push(`/quizzes/${quizId}/take?resume=true`);
  };

  // Handle abandoning from 409 conflict dialog
  const handleAbandonExisting = () => {
    if (!existingAttemptId) return;

    abandonSessionMutation.mutate(
      {
        quizId,
        attemptId: existingAttemptId,
      },
      {
        onSuccess: () => {
          setShow409Dialog(false);
          setExistingAttemptId(null);
          // Retry starting the quiz
          startAttemptMutation.mutate(quizId, {
            onSuccess: (attempt) => {
              if (quiz) {
                startQuiz(quiz, attempt);
              }
            },
          });
        },
      }
    );
  };

  // Loading state
  if (
    quizLoading ||
    questionsLoading ||
    (isResuming && activeAttemptLoading) ||
    startAttemptMutation.isPending
  ) {
    return (
      <div className="container py-8">
        <LoadingSpinner size="lg" />
        <p className="text-center mt-4 text-muted-foreground">
          {isResuming ? "Loading your saved progress..." : "Loading quiz..."}
        </p>
      </div>
    );
  }

  // Error state
  if (questionsError) {
    return (
      <div className="container py-8">
        <EmptyState
          icon={BookOpen}
          title="Error Loading Questions"
          description={
            questionsError instanceof Error
              ? questionsError.message
              : "Failed to load quiz questions. Please try again."
          }
          action={
            <Button onClick={() => router.push("/quizzes")}>
              Back to Quizzes
            </Button>
          }
        />
      </div>
    );
  }

  if (!quiz || !questions || !currentQuiz) {
    return (
      <div className="container py-8">
        <EmptyState
          icon={BookOpen}
          title="Quiz Not Found"
          description="Unable to load the quiz. Please try again."
          action={
            <Button onClick={() => router.push("/quizzes")}>
              Back to Quizzes
            </Button>
          }
        />
      </div>
    );
  }

  // DEBUG: Log the questions data
  console.log("[DEBUGGG] Take Page - questions from hook:", questions);
  console.log("[DEBUGGG] Take Page - questions type:", typeof questions);
  console.log("[DEBUGGG] Take Page - questions is array:", Array.isArray(questions));
  console.log("[DEBUGGG] Take Page - questions length:", questions?.length);
  console.log("[DEBUGGG] Take Page - currentQuestionIndex:", currentQuestionIndex);

  const currentQuestion = questions[currentQuestionIndex];
  console.log("[DEBUGGG] Take Page - currentQuestion:", currentQuestion);
  console.log("[DEBUGGG] Take Page - currentQuestion?.question_text:", currentQuestion?.question_text);
  console.log("[DEBUGGG] Take Page - currentQuestion?.options:", currentQuestion?.options);

  const currentAnswer = currentQuestion ? getAnswer(currentQuestion.id) : undefined;
  const answeredCount = Object.keys(answers).length;

  return (
    <div className="container py-8 max-w-4xl space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{currentQuiz.title}</h1>
          <p className="text-sm text-muted-foreground">{currentQuiz.category}</p>
        </div>
        <div className="flex gap-2">
          {isPaused ? (
            <Button variant="default" size="sm" onClick={handleResume}>
              <Play className="h-4 w-4 mr-2" />
              Resume
            </Button>
          ) : (
            <Button variant="outline" size="sm" onClick={handlePause}>
              <Pause className="h-4 w-4 mr-2" />
              Pause
            </Button>
          )}
          <Button variant="outline" size="sm" onClick={handleSaveProgress}>
            <Save className="h-4 w-4 mr-2" />
            Save Progress
          </Button>
          <Button variant="ghost" size="sm" onClick={handleExit}>
            <X className="h-4 w-4 mr-2" />
            Exit
          </Button>
        </div>
      </div>

      {/* Timer and Progress */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <QuizTimer timeRemaining={timeRemaining} isPaused={isPaused} />
        <QuizProgress
          currentQuestion={currentQuestionIndex}
          totalQuestions={currentQuiz.question_count}
          answeredCount={answeredCount}
        />
      </div>

      {/* Question */}
      {currentQuestion && (
        <QuestionCard
          question={currentQuestion}
          questionNumber={currentQuestionIndex + 1}
          selectedAnswer={currentAnswer?.selected_answer}
          onAnswerChange={handleAnswerChange}
        />
      )}

      {/* Navigation */}
      <QuizNavigation
        currentQuestion={currentQuestionIndex}
        totalQuestions={currentQuiz.question_count}
        onPrevious={previousQuestion}
        onNext={nextQuestion}
        onSubmit={handleSubmit}
        isLoading={submitAttemptMutation.isPending}
        showSubmit={currentQuestionIndex === currentQuiz.question_count - 1}
      />

      {/* Exit Confirmation Dialog */}
      <AlertDialog open={showExitDialog} onOpenChange={setShowExitDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Exit Quiz?</AlertDialogTitle>
            <AlertDialogDescription>
              Your progress will be saved. You can resume this quiz later from where
              you left off.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmExit}>
              Save and Exit
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Submit Confirmation Dialog */}
      <AlertDialog open={showSubmitDialog} onOpenChange={setShowSubmitDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Submit Quiz?</AlertDialogTitle>
            <AlertDialogDescription>
              You have answered {answeredCount} out of {currentQuiz.question_count}{" "}
              questions. Are you sure you want to submit?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmSubmit}>
              Submit Quiz
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 409 Conflict Dialog - Active Attempt Exists */}
      <AlertDialog open={show409Dialog} onOpenChange={setShow409Dialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Active Quiz Attempt Found</AlertDialogTitle>
            <AlertDialogDescription>
              You have an active attempt for this quiz. Would you like to resume your
              existing attempt or abandon it to start fresh?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="flex-col sm:flex-row gap-2">
            <AlertDialogCancel onClick={() => router.push(`/quizzes/${quizId}`)}>
              Go Back
            </AlertDialogCancel>
            <Button
              variant="outline"
              onClick={handleAbandonExisting}
              disabled={abandonSessionMutation.isPending}
            >
              {abandonSessionMutation.isPending ? "Abandoning..." : "Abandon & Start Fresh"}
            </Button>
            <AlertDialogAction onClick={handleResumeExisting}>
              Resume Existing
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
