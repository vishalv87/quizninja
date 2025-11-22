"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { useIsFavorite, useToggleFavorite } from "@/hooks/useFavorites";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import { RelatedQuizzesCard } from "@/components/quiz/RelatedQuizzesCard";
import {
  BookOpen,
  Trophy,
  Star,
  ArrowLeft,
  Play,
  Target,
  Heart,
  Lightbulb,
  AlertTriangle,
  CheckCircle2,
} from "lucide-react";
import Link from "next/link";
import Image from "next/image";
import { toast } from "sonner";

// Quiz rules content
const QUIZ_RULES = [
  {
    title: "Answer All Questions",
    description: "Complete all questions to submit your quiz. Unanswered questions will be marked as incorrect."
  },
  {
    title: "Time Management",
    description: "Keep an eye on the timer. The quiz will auto-submit when time expires."
  },
  {
    title: "Review Your Answers",
    description: "You can review and change answers before final submission. Use the navigation to go back to previous questions."
  },
  {
    title: "Single Attempt Per Question",
    description: "You can only select one answer per question. Choose carefully!"
  },
  {
    title: "Fair Play",
    description: "Do not use external resources or assistance. This is a test of your knowledge."
  },
];

export default function QuizDetailPage() {
  const params = useParams();
  const router = useRouter();
  const quizId = Array.isArray(params.id) ? params.id[0] : params.id;

  const { data: quiz, isLoading, error } = useQuiz(quizId);
  const { data: isFavorite, isLoading: favoriteLoading } = useIsFavorite(quizId);
  const { toggle: toggleFavorite } = useToggleFavorite();

  // Loading state
  if (isLoading) {
    return (
      <div className="container py-8">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  // Error state
  if (error || !quiz) {
    return (
      <div className="container py-8">
        <EmptyState
          icon={BookOpen}
          title="Quiz Not Found"
          description="The quiz you're looking for doesn't exist or has been removed."
          action={
            <Button onClick={() => router.push("/quizzes")}>
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Quizzes
            </Button>
          }
        />
      </div>
    );
  }

  // Determine difficulty color with fallback
  const difficultyColor = quiz.difficulty ? {
    beginner: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
    intermediate: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
    advanced: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  }[quiz.difficulty] || "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300"
    : "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300";

  // Get category color for gradient
  const getCategoryGradient = (category: string | null | undefined) => {
    const gradients: Record<string, string> = {
      science: "from-blue-500 to-blue-700",
      mathematics: "from-purple-500 to-purple-700",
      history: "from-amber-500 to-amber-700",
      geography: "from-green-500 to-green-700",
      technology: "from-cyan-500 to-cyan-700",
      literature: "from-pink-500 to-pink-700",
    };
    return category ? (gradients[category.toLowerCase()] || "from-gray-500 to-gray-700") : "from-gray-500 to-gray-700";
  };

  return (
    <div className="container py-8 space-y-6 max-w-7xl">
      {/* Back Button */}
      <Button variant="ghost" onClick={() => router.push("/quizzes")}>
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to Quizzes
      </Button>

      {/* Quiz Header Card with Gradient Background */}
      <Card className={`bg-gradient-to-r ${getCategoryGradient(quiz.category)} text-white border-0 shadow-lg`}>
        <CardContent className="pt-6 space-y-4">
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 space-y-3">
              {/* Badges Row */}
              <div className="flex gap-2 flex-wrap">
                <Badge variant="secondary" className="bg-white/20 text-white border-white/30 backdrop-blur-sm">
                  {quiz.category?.toUpperCase() ?? "UNCATEGORIZED"}
                </Badge>
                {quiz.difficulty && (
                  <Badge variant="secondary" className="bg-white/20 text-white border-white/30 backdrop-blur-sm">
                    {quiz.difficulty.toUpperCase()}
                  </Badge>
                )}
                {/* Tags */}
                {quiz.tags && quiz.tags.length > 0 && (
                  <>
                    {quiz.tags.map((tag, index) => (
                      <Badge
                        key={index}
                        variant="secondary"
                        className="bg-white/10 text-white border-white/20 backdrop-blur-sm"
                      >
                        {tag}
                      </Badge>
                    ))}
                  </>
                )}
              </div>

              {/* Title and Description */}
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <h1 className="text-3xl md:text-4xl font-bold tracking-tight">
                    {quiz.title}
                  </h1>
                  {quiz.is_featured && (
                    <Star className="h-6 w-6 text-yellow-300 fill-yellow-300" />
                  )}
                </div>
                <p className="text-white/90 text-base md:text-lg">{quiz.description}</p>
              </div>

              {/* Stats Row */}
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6 pt-4">
                <div className="space-y-1">
                  <p className="text-sm text-white/70 font-medium">Questions</p>
                  <p className="text-2xl font-bold text-white">{quiz.question_count}</p>
                </div>
                <div className="space-y-1">
                  <p className="text-sm text-white/70 font-medium">Time Limit</p>
                  <p className="text-2xl font-bold text-white">
                    {quiz.time_limit ? `${quiz.time_limit} min` : 'No limit'}
                  </p>
                </div>
                <div className="space-y-1">
                  <p className="text-sm text-white/70 font-medium">Points</p>
                  <p className="text-2xl font-bold text-white">{quiz.points}</p>
                </div>
                <div className="col-span-2 md:col-span-1 flex items-end">
                  <Link href={`/quizzes/${quiz.id}/take`} className="w-full">
                    <Button size="lg" className="w-full bg-white text-gray-900 hover:bg-white/90">
                      <Play className="mr-2 h-5 w-5" />
                      Start Quiz
                    </Button>
                  </Link>
                </div>
              </div>
            </div>

            {/* Favorite Button */}
            <Button
              variant="outline"
              size="icon"
              onClick={() => {
                toggleFavorite(quizId, isFavorite || false);
                toast.success(
                  isFavorite ? "Removed from favorites" : "Added to favorites",
                  { duration: 2000 }
                );
              }}
              disabled={favoriteLoading}
              className="shrink-0 bg-white/20 border-white/30 hover:bg-white/30 text-white backdrop-blur-sm"
            >
              <Heart
                className={`h-5 w-5 ${
                  isFavorite ? "fill-red-300 text-red-300" : ""
                }`}
              />
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Desktop Grid Layout: Main Content + Sidebar */}
      <div className="grid lg:grid-cols-3 gap-6">
        {/* Main Content Area (Left Side - 2/3 on desktop) */}
        <div className="lg:col-span-2 space-y-6">

          {/* Quiz Rules */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CheckCircle2 className="h-5 w-5" />
                Quiz Rules
              </CardTitle>
              <CardDescription>Please read these rules carefully before starting the quiz</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-3">
                {QUIZ_RULES.map((rule, index) => (
                  <div key={index} className="flex gap-3">
                    <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold text-sm">
                      {index + 1}
                    </div>
                    <div className="flex-1 pt-1">
                      <h4 className="font-semibold text-sm">{rule.title}</h4>
                      <p className="text-sm text-muted-foreground mt-1">
                        {rule.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>

              {/* Important Notice */}
              <Alert className="bg-amber-50 dark:bg-amber-950 border-amber-200 dark:border-amber-800">
                <AlertTriangle className="h-4 w-4 text-amber-600 dark:text-amber-400" />
                <AlertDescription className="text-amber-900 dark:text-amber-100">
                  <strong>Important:</strong> Violating these rules may result in disqualification. Play fair and test your true knowledge!
                </AlertDescription>
              </Alert>
            </CardContent>
          </Card>

          <Separator />

          {/* Question Preview */}
          {quiz.questions && quiz.questions.length > 0 && (
            <>
              <Card className="border-2 border-dashed">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle className="flex items-center gap-2">
                      <Target className="h-5 w-5" />
                      Question Preview
                    </CardTitle>
                    <Badge variant="secondary">Sample Question</Badge>
                  </div>
                  <CardDescription>
                    Here's a sample question from this quiz (Question 1 of {quiz.question_count})
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {/* Question Image */}
                  {quiz.questions[0].image_url && (
                    <div className="relative w-full h-64 rounded-lg overflow-hidden bg-muted">
                      <Image
                        src={quiz.questions[0].image_url}
                        alt="Question visual"
                        fill
                        className="object-contain"
                        unoptimized
                      />
                    </div>
                  )}

                  {/* Question Text */}
                  <div className="p-4 bg-muted/50 rounded-lg">
                    <p className="font-medium text-lg">
                      {quiz.questions[0].question_text}
                    </p>
                  </div>

                  {/* Question Options */}
                  {quiz.questions[0].options && (
                    <div className="space-y-2">
                      {quiz.questions[0].options.map((optionText, index) => (
                        <div
                          key={`option-${index}`}
                          className="p-3 border rounded-lg bg-background cursor-not-allowed opacity-75"
                        >
                          <span className="font-medium mr-2">
                            {String.fromCharCode(65 + index)}.
                          </span>
                          {optionText}
                        </div>
                      ))}
                    </div>
                  )}

                  <div className="space-y-2">
                    <p className="text-xs text-muted-foreground italic">
                      * This is just a preview. Actual answers are hidden until you start the quiz.
                    </p>
                    {quiz.questions[0].explanation && (
                      <div className="flex items-center gap-2 text-xs text-blue-600 dark:text-blue-400">
                        <Lightbulb className="h-3 w-3" />
                        <span className="font-medium">Detailed explanations available after quiz completion</span>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>

              <Separator />
            </>
          )}

          {/* View Results Button for completed quizzes */}
          {quiz.user_has_attempted && (
            <Link href={`/quizzes/${quiz.id}/results`} className="w-full">
              <Button size="lg" variant="secondary" className="w-full">
                <Trophy className="mr-2 h-5 w-5" />
                View Results & History
              </Button>
            </Link>
          )}
        </div>

        {/* Sidebar (Right Side - 1/3 on desktop) */}
        <div className="lg:col-span-1">
          <div className="sticky top-6">
            <RelatedQuizzesCard
              currentQuizId={quiz.id}
              category={quiz.category}
              difficulty={quiz.difficulty}
            />
          </div>
        </div>
      </div>
    </div>
  );
}