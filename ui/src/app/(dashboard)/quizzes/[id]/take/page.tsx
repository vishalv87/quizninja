"use client";

import { useEffect, useState, useRef } from "react";
import { useParams, useRouter } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { useQuizQuestions } from "@/hooks/useQuizQuestions";
import { useActiveAttempt } from "@/hooks/useQuizAttempt";
import { useQuizStore } from "@/store/quizStore";
import {
  useStartQuizAttempt,
  useSubmitQuizAttempt,
  useAbandonQuizSession,
} from "@/hooks/useQuizAttempt";
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
import { BookOpen, X } from "lucide-react";

export default function QuizTakingPage() {
  const params = useParams();
  const router = useRouter();
  const quizId = params.id as string;

  const { data: quiz, isLoading: quizLoading } = useQuiz(quizId);
  const {
    data: questions,
    isLoading: questionsLoading,
    error: questionsError,
  } = useQuizQuestions(quizId);
  const { data: activeAttempt, isLoading: activeAttemptLoading } = useActiveAttempt(quizId);
  const startAttemptMutation = useStartQuizAttempt();
  const submitAttemptMutation = useSubmitQuizAttempt();
  const abandonSessionMutation = useAbandonQuizSession();

  const {
    currentQuiz,
    currentAttempt,
    currentQuestionIndex,
    answers,
    timeRemaining,
    nextQuestion,
    previousQuestion,
    setAnswer,
    getAnswer,
    startQuiz,
    resetQuiz,
  } = useQuizStore();

  const [showExitDialog, setShowExitDialog] = useState(false);
  const [showSubmitDialog, setShowSubmitDialog] = useState(false);
  const [isInitializing, setIsInitializing] = useState(true);

  // Ref to prevent multiple initialization attempts
  const hasInitialized = useRef(false);

  // Extract stable mutate functions to avoid dependency issues
  const { mutate: startAttempt } = startAttemptMutation;
  const { mutate: abandonSession } = abandonSessionMutation;

  // Initialize quiz - Auto-abandon any existing attempt and start fresh
  useEffect(() => {
    // Guard: Only initialize once per component mount
    if (hasInitialized.current) {
      return;
    }

    // Wait for data to be ready
    if (!quiz || !questions || activeAttemptLoading) {
      return;
    }

    // Don't re-initialize if quiz is already set up
    if (currentQuiz) {
      setIsInitializing(false);
      return;
    }

    // Mark as initialized to prevent re-runs
    hasInitialized.current = true;
    setIsInitializing(true);

    // Helper function to start a new quiz attempt
    const startNewAttempt = () => {
      startAttempt(quizId, {
        onSuccess: (attempt) => {
          startQuiz(quiz, attempt);
          setIsInitializing(false);
        },
        onError: () => {
          setIsInitializing(false);
          hasInitialized.current = false; // Allow retry on error
          router.push(`/quizzes/${quizId}`);
        },
      });
    };

    // If there's an active attempt with a valid ID, auto-abandon it first
    if (activeAttempt?.id) {
      abandonSession(
        { quizId, attemptId: activeAttempt.id },
        {
          onSuccess: () => {
            // After abandoning, start fresh
            startNewAttempt();
          },
          onError: () => {
            // Even if abandon fails, try to start fresh
            startNewAttempt();
          },
        }
      );
    } else {
      // No active attempt - Start fresh
      startNewAttempt();
    }
  }, [
    quiz,
    questions,
    currentQuiz,
    quizId,
    activeAttempt?.id,
    activeAttemptLoading,
    router,
    startAttempt,
    startQuiz,
    abandonSession,
  ]);

  // Setup quiz timer
  useQuizTimer(() => {
    // Auto-submit when time runs out
    handleSubmit();
  });

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
    // Abandon the attempt when exiting
    if (currentAttempt) {
      abandonSessionMutation.mutate(
        { quizId, attemptId: currentAttempt.id },
        {
          onSettled: () => {
            resetQuiz();
            router.push(`/quizzes/${quizId}`);
          },
        }
      );
    } else {
      resetQuiz();
      router.push(`/quizzes/${quizId}`);
    }
  };

  // Loading state
  if (
    quizLoading ||
    questionsLoading ||
    activeAttemptLoading ||
    isInitializing ||
    startAttemptMutation.isPending ||
    abandonSessionMutation.isPending
  ) {
    return (
      <div className="container py-8">
        <LoadingSpinner size="lg" />
        <p className="text-center mt-4 text-muted-foreground">Loading quiz...</p>
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
          action={<Button onClick={() => router.push("/quizzes")}>Back to Quizzes</Button>}
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
          action={<Button onClick={() => router.push("/quizzes")}>Back to Quizzes</Button>}
        />
      </div>
    );
  }

  const currentQuestion = questions[currentQuestionIndex];

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
        <Button variant="ghost" size="sm" onClick={handleExit}>
          <X className="h-4 w-4 mr-2" />
          Exit
        </Button>
      </div>

      {/* Timer and Progress */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <QuizTimer timeRemaining={timeRemaining} />
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
              Are you sure you want to exit? Your attempt will be abandoned and you will need to
              start fresh next time.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmExit}>Exit Quiz</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Submit Confirmation Dialog */}
      <AlertDialog open={showSubmitDialog} onOpenChange={setShowSubmitDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Submit Quiz?</AlertDialogTitle>
            <AlertDialogDescription>
              You have answered {answeredCount} out of {currentQuiz.question_count} questions. Are
              you sure you want to submit?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmSubmit}>Submit Quiz</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
