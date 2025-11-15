"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuiz } from "@/hooks/useQuiz";
import { useActiveAttempt } from "@/hooks/useQuizAttempt";
import { useIsFavorite, useToggleFavorite } from "@/hooks/useFavorites";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import {
  BookOpen,
  Clock,
  Trophy,
  Star,
  ArrowLeft,
  Play,
  TrendingUp,
  Target,
  Users,
  RotateCcw,
  Heart,
} from "lucide-react";
import Link from "next/link";

export default function QuizDetailPage() {
  const params = useParams();
  const router = useRouter();
  const quizId = params.id as string;

  const { data: quiz, isLoading, error } = useQuiz(quizId);
  const { data: activeAttempt, isLoading: attemptLoading } = useActiveAttempt(quizId);
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

  // Determine difficulty color
  const difficultyColor = {
    easy: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
    medium:
      "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
    hard: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  }[quiz.difficulty];

  return (
    <div className="container py-8 space-y-6 max-w-5xl">
      {/* Back Button */}
      <Button variant="ghost" onClick={() => router.push("/quizzes")}>
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to Quizzes
      </Button>

      {/* Quiz Header */}
      <div className="space-y-4">
        <div className="flex items-start justify-between gap-4">
          <div className="space-y-2 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="text-4xl font-bold tracking-tight">
                {quiz.title}
              </h1>
              {quiz.is_featured && (
                <Star className="h-6 w-6 text-yellow-500 fill-yellow-500" />
              )}
            </div>
            <p className="text-lg text-muted-foreground">{quiz.description}</p>
          </div>

          {/* Favorite Button */}
          <Button
            variant="outline"
            size="icon"
            onClick={() => toggleFavorite(quizId, isFavorite || false)}
            disabled={favoriteLoading}
            className="shrink-0"
          >
            <Heart
              className={`h-5 w-5 ${
                isFavorite ? "fill-red-500 text-red-500" : ""
              }`}
            />
          </Button>
        </div>

        {/* Badges */}
        <div className="flex gap-2 flex-wrap">
          <Badge variant="secondary" className="text-sm">
            {quiz.category}
          </Badge>
          <Badge className={`${difficultyColor} text-sm`}>
            {quiz.difficulty.charAt(0).toUpperCase() + quiz.difficulty.slice(1)}
          </Badge>
        </div>
      </div>

      <Separator />

      {/* Quiz Info Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <BookOpen className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-2xl font-bold">{quiz.question_count}</p>
                <p className="text-xs text-muted-foreground">Questions</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {quiz.time_limit_minutes && (
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                  <Clock className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{quiz.time_limit_minutes}</p>
                  <p className="text-xs text-muted-foreground">Minutes</p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                <Trophy className="h-5 w-5 text-yellow-600 dark:text-yellow-400" />
              </div>
              <div>
                <p className="text-2xl font-bold">{quiz.total_points}</p>
                <p className="text-xs text-muted-foreground">Points</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
                <Target className="h-5 w-5 text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <p className="text-2xl font-bold">
                  {quiz.points_per_question}
                </p>
                <p className="text-xs text-muted-foreground">Per Question</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Statistics */}
      {(quiz.attempts_count !== undefined || quiz.average_score !== undefined) && (
        <>
          <Separator />
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <TrendingUp className="h-5 w-5" />
                Quiz Statistics
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {quiz.attempts_count !== undefined && (
                  <div className="flex items-center gap-3">
                    <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
                      <Users className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                    </div>
                    <div>
                      <p className="text-3xl font-bold">{quiz.attempts_count}</p>
                      <p className="text-sm text-muted-foreground">
                        Total Attempts
                      </p>
                    </div>
                  </div>
                )}

                {quiz.average_score !== undefined && (
                  <div className="flex items-center gap-3">
                    <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
                      <Star className="h-6 w-6 text-green-600 dark:text-green-400" />
                    </div>
                    <div>
                      <p className="text-3xl font-bold">
                        {quiz.average_score.toFixed(1)}%
                      </p>
                      <p className="text-sm text-muted-foreground">
                        Average Score
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </>
      )}

      <Separator />

      {/* Action Buttons */}
      <div className="flex gap-4">
        {activeAttempt ? (
          <>
            <Link href={`/quizzes/${quiz.id}/take?resume=true`} className="flex-1">
              <Button size="lg" className="w-full">
                <RotateCcw className="mr-2 h-5 w-5" />
                Resume Quiz
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
            <Button size="lg" className="w-full" disabled={attemptLoading}>
              <Play className="mr-2 h-5 w-5" />
              {attemptLoading ? "Checking..." : "Start Quiz"}
            </Button>
          </Link>
        )}
      </div>
    </div>
  );
}