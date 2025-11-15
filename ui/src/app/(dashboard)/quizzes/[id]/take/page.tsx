"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { useQuizQuestions } from "@/hooks/useQuizQuestions";
import { useActiveAttempt } from "@/hooks/useQuizAttempt";
import { useQuizStore } from "@/store/quizStore";
import {
  useStartQuizAttempt,
  useSubmitQuizAttempt,
  useSaveQuizProgress,
  usePauseQuizAttempt,
  useResumeQuizAttempt,
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
  const saveProgressMutation = useSaveQuizProgress();
  const pauseAttemptMutation = usePauseQuizAttempt();
  const resumeAttemptMutation = useResumeQuizAttempt();

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

  // Initialize quiz
  useEffect(() => {
    if (quiz && questions && !currentQuiz) {
      if (isResuming && activeAttempt) {
        // Resume existing attempt
        toast.info("Resuming quiz from where you left off");
        startQuiz(quiz, activeAttempt);

        // Restore saved answers to the store
        if (activeAttempt.answers && activeAttempt.answers.length > 0) {
          activeAttempt.answers.forEach((answer) => {
            setAnswer(answer.question_id, answer);
          });
        }
      } else {
        // Start a new quiz attempt
        startAttemptMutation.mutate(quizId, {
          onSuccess: (attempt) => {
            startQuiz(quiz, attempt);
          },
        });
      }
    }
  }, [quiz, questions, currentQuiz, quizId, isResuming, activeAttempt]);

  // Setup quiz timer
  useQuizTimer(() => {
    // Auto-submit when time runs out
    handleSubmit();
  });

  // Auto-save progress every 30 seconds
  useEffect(() => {
    if (!currentAttempt) return;

    const interval = setInterval(() => {
      handleSaveProgress();
    }, 30000); // 30 seconds

    return () => clearInterval(interval);
  }, [currentAttempt, answers]);

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

    setAnswer(question.id, {
      question_id: question.id,
      selected_answer: answer,
      is_correct: false, // Will be determined on backend
      points_earned: 0, // Will be determined on backend
    });
  };

  const handleSaveProgress = () => {
    if (!currentAttempt) return;

    const answersList = Object.values(answers);
    saveProgressMutation.mutate({
      quizId,
      attemptId: currentAttempt.id,
      answers: answersList,
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
        onSuccess: (results) => {
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

    setPaused(true);
    pauseAttemptMutation.mutate({
      quizId,
      attemptId: currentAttempt.id,
    });
  };

  const handleResume = () => {
    if (!currentAttempt) return;

    setPaused(false);
    resumeAttemptMutation.mutate({
      quizId,
      attemptId: currentAttempt.id,
    });
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
    </div>
  );
}
