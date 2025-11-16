"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { useQuizActiveSession } from "@/hooks/useQuizAttempt";
import { useIsFavorite, useToggleFavorite } from "@/hooks/useFavorites";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Progress } from "@/components/ui/progress";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import {
  BookOpen,
  Trophy,
  Star,
  ArrowLeft,
  Play,
  Target,
  RotateCcw,
  Heart,
  PlayCircle,
  Lightbulb,
  AlertTriangle,
  FileText,
  CheckCircle2,
  HelpCircle,
  Clock,
  Users,
  TrendingUp,
  Award,
  BarChart3,
  Calendar,
  User,
  Tag,
  Image as ImageIcon,
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
  const { data: activeSession, isLoading: sessionLoading } = useQuizActiveSession(quizId);
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

  // Calculate progress for active session
  const sessionProgress = activeSession
    ? Math.round(((activeSession.current_question_index + 1) / quiz.question_count) * 100)
    : 0;

  // Format time remaining
  const formatTimeRemaining = (seconds?: number) => {
    if (!seconds) return "N/A";
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="container py-8 space-y-6 max-w-5xl">
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
                {activeSession && (
                  <Badge className={`${
                    activeSession.session_state === 'paused'
                      ? 'bg-orange-500 text-white'
                      : 'bg-green-500 text-white'
                  } border-0`}>
                    {activeSession.session_state.toUpperCase()}
                  </Badge>
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

              {/* Progress Bar for Active Session */}
              {activeSession && (
                <div className="space-y-2">
                  <Progress value={sessionProgress} className="h-2 bg-white/30" />
                  <div className="grid grid-cols-3 gap-4 text-sm">
                    <div>
                      <p className="text-white/70">Answered</p>
                      <p className="font-semibold">{activeSession.current_question_index + 1}</p>
                    </div>
                    <div>
                      <p className="text-white/70">Remaining</p>
                      <p className="font-semibold">{quiz.question_count - (activeSession.current_question_index + 1)}</p>
                    </div>
                    <div>
                      <p className="text-white/70">Time Left</p>
                      <p className="font-semibold">{formatTimeRemaining(activeSession.time_remaining)}</p>
                    </div>
                  </div>
                </div>
              )}

              {/* Stats Row for Normal Flow */}
              {!activeSession && (
                <div className="grid grid-cols-3 gap-6 pt-4">
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
                </div>
              )}
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

      {/* Previous Attempt Card */}
      {quiz.user_has_attempted && quiz.user_best_score !== undefined && (
        <Card className="bg-gradient-to-r from-amber-50 to-amber-100 dark:from-amber-950 dark:to-amber-900 border-amber-200 dark:border-amber-800">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-gradient-to-br from-yellow-400 to-amber-600 rounded-lg">
                <Trophy className="h-8 w-8 text-white" />
              </div>
              <div>
                <p className="text-3xl font-bold text-amber-900 dark:text-amber-100">
                  {quiz.user_best_score.toFixed(1)}%
                </p>
                <p className="text-sm text-amber-700 dark:text-amber-300">Your Best Score</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Quiz Metadata */}
      {(quiz.tags?.length || quiz.created_by || quiz.thumbnail_url || quiz.created_at) && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              Quiz Information
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid md:grid-cols-2 gap-6">
              {/* Left Column */}
              <div className="space-y-4">
                {/* Tags */}
                {quiz.tags && quiz.tags.length > 0 && (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Tag className="h-4 w-4" />
                      Tags
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {quiz.tags.map((tag, index) => (
                        <Badge key={index} variant="outline" className="text-xs">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}

                {/* Creator */}
                {quiz.created_by && (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <User className="h-4 w-4" />
                      Created By
                    </div>
                    <p className="text-sm font-medium">{quiz.created_by}</p>
                  </div>
                )}
              </div>

              {/* Right Column */}
              <div className="space-y-4">
                {/* Timestamps */}
                {quiz.created_at && (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Calendar className="h-4 w-4" />
                      Created
                    </div>
                    <p className="text-sm font-medium">
                      {new Date(quiz.created_at).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </p>
                  </div>
                )}

                {quiz.updated_at && quiz.updated_at !== quiz.created_at && (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Calendar className="h-4 w-4" />
                      Last Updated
                    </div>
                    <p className="text-sm font-medium">
                      {new Date(quiz.updated_at).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </p>
                  </div>
                )}
              </div>
            </div>

            {/* Thumbnail */}
            {quiz.thumbnail_url && (
              <div className="space-y-2 pt-2 border-t">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <ImageIcon className="h-4 w-4" />
                  Quiz Thumbnail
                </div>
                <div className="relative w-full h-48 rounded-lg overflow-hidden bg-muted">
                  <Image
                    src={quiz.thumbnail_url}
                    alt={quiz.title}
                    fill
                    className="object-cover"
                    unoptimized
                  />
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Community Statistics & Ratings */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Community Statistics */}
        {quiz.statistics && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="h-5 w-5" />
                Community Statistics
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-muted-foreground text-sm">
                    <Users className="h-4 w-4" />
                    Total Attempts
                  </div>
                  <p className="text-2xl font-bold">{quiz.statistics.total_attempts.toLocaleString()}</p>
                </div>
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-muted-foreground text-sm">
                    <TrendingUp className="h-4 w-4" />
                    Avg Score
                  </div>
                  <p className="text-2xl font-bold">{quiz.statistics.average_score.toFixed(1)}%</p>
                </div>
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-muted-foreground text-sm">
                    <Award className="h-4 w-4" />
                    Completion Rate
                  </div>
                  <p className="text-2xl font-bold">{quiz.statistics.completion_rate.toFixed(1)}%</p>
                </div>
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-muted-foreground text-sm">
                    <Clock className="h-4 w-4" />
                    Avg Time
                  </div>
                  <p className="text-2xl font-bold">{Math.floor(quiz.statistics.average_time / 60)}m</p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Rating Display */}
        {(quiz.average_rating !== undefined && quiz.total_ratings !== undefined) && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Star className="h-5 w-5" />
                Community Rating
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <p className="text-4xl font-bold">{quiz.average_rating.toFixed(1)}</p>
                    <div className="flex items-center">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <Star
                          key={star}
                          className={`h-5 w-5 ${
                            star <= Math.round(quiz.average_rating!)
                              ? "fill-yellow-400 text-yellow-400"
                              : "text-gray-300"
                          }`}
                        />
                      ))}
                    </div>
                  </div>
                  <p className="text-sm text-muted-foreground mt-1">
                    Based on {quiz.total_ratings} {quiz.total_ratings === 1 ? 'rating' : 'ratings'}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>

      {/* Article Summary & Rules Tabs */}
      <Tabs defaultValue="summary" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="summary" className="flex items-center gap-2">
            <FileText className="h-4 w-4" />
            Article Summary
          </TabsTrigger>
          <TabsTrigger value="rules" className="flex items-center gap-2">
            <HelpCircle className="h-4 w-4" />
            Rules
          </TabsTrigger>
        </TabsList>

        <TabsContent value="summary" className="space-y-4">
          {quiz.article_summary ? (
            <>
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <BookOpen className="h-5 w-5" />
                    Key Points to Remember
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <p className="text-muted-foreground leading-relaxed">
                    {quiz.article_summary}
                  </p>

                  {/* Pro Tip Callout */}
                  <Alert className="bg-blue-50 dark:bg-blue-950 border-blue-200 dark:border-blue-800">
                    <Lightbulb className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                    <AlertDescription className="text-blue-900 dark:text-blue-100">
                      <strong>Pro Tip:</strong> Review these key points before starting the quiz to refresh your knowledge and improve your score!
                    </AlertDescription>
                  </Alert>
                </CardContent>
              </Card>
            </>
          ) : (
            <Card>
              <CardContent className="pt-6">
                <EmptyState
                  icon={FileText}
                  title="No Article Summary Available"
                  description="This quiz doesn't have an article summary yet."
                />
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="rules" className="space-y-4">
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
        </TabsContent>
      </Tabs>

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
                  {quiz.questions[0].options.map((option, index) => (
                    <div
                      key={option.id}
                      className="p-3 border rounded-lg bg-background cursor-not-allowed opacity-75"
                    >
                      <span className="font-medium mr-2">
                        {String.fromCharCode(65 + index)}.
                      </span>
                      {option.option_text}
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

      {/* Action Buttons */}
      <div className="flex flex-col gap-4">
        <div className="flex gap-4">
          {activeSession ? (
            <>
              <Link href={`/quizzes/${quiz.id}/take?resume=true`} className="flex-1">
                <Button size="lg" className="w-full">
                  {activeSession.session_state === 'paused' ? (
                    <>
                      <RotateCcw className="mr-2 h-5 w-5" />
                      Resume Quiz
                    </>
                  ) : (
                    <>
                      <PlayCircle className="mr-2 h-5 w-5" />
                      Continue Quiz
                    </>
                  )}
                </Button>
              </Link>
              <Link href={`/quizzes/${quiz.id}/take`} className="flex-1">
                <Button size="lg" variant="outline" className="w-full">
                  <Play className="mr-2 h-5 w-5" />
                  Start New Attempt
                </Button>
              </Link>
            </>
          ) : (
            <Link href={`/quizzes/${quiz.id}/take`} className="flex-1">
              <Button size="lg" className="w-full" disabled={sessionLoading}>
                <Play className="mr-2 h-5 w-5" />
                {sessionLoading ? "Checking..." : "Start Quiz"}
              </Button>
            </Link>
          )}
        </div>

        {/* View Results Button for completed quizzes */}
        {quiz.user_has_attempted && !activeSession && (
          <Link href={`/quizzes/${quiz.id}/results`} className="w-full">
            <Button size="lg" variant="secondary" className="w-full">
              <Trophy className="mr-2 h-5 w-5" />
              View Results & History
            </Button>
          </Link>
        )}
      </div>
    </div>
  );
}