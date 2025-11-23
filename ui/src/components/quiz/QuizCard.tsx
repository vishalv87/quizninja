"use client";

import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type { Quiz, QuizAttempt } from "@/types/quiz";
import { BookOpen, Clock, Star, Trophy, Heart, CheckCircle, XCircle, ArrowRight } from "lucide-react";
import Link from "next/link";
import { useIsFavorite, useToggleFavorite } from "@/hooks/useFavorites";
import { cn } from "@/lib/utils";

interface QuizCardProps {
  quiz: Quiz;
  completedAttempt?: QuizAttempt;
}

export function QuizCard({ quiz, completedAttempt }: QuizCardProps) {
  // Check if quiz is favorited
  const { data: isFavorite, isLoading: isFavoriteLoading } = useIsFavorite(quiz.id);
  const { toggle } = useToggleFavorite();

  // Derive completion status - only consider truly completed (non-abandoned) attempts
  const isCompleted = completedAttempt && completedAttempt.status === "completed";
  const percentage = completedAttempt?.percentage_score ?? 0;
  const passed = percentage >= 60;

  // Determine difficulty color
  const difficultyColor = {
    beginner: "bg-green-100 text-green-700 border-green-200",
    intermediate: "bg-yellow-100 text-yellow-700 border-yellow-200",
    advanced: "bg-red-100 text-red-700 border-red-200",
  }[quiz.difficulty];

  // Handle favorite toggle
  const handleFavoriteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    toggle(quiz.id, isFavorite || false);
  };

  return (
    <Card className="group relative flex flex-col h-full overflow-hidden border-0 shadow-md transition-all duration-300 hover:shadow-xl hover:-translate-y-1 bg-white dark:bg-background rounded-2xl">
      {/* Top Decoration Bar */}
      <div className={cn("h-1.5 w-full bg-gradient-to-r", 
        quiz.difficulty === 'beginner' ? "from-green-400 to-emerald-500" :
        quiz.difficulty === 'intermediate' ? "from-yellow-400 to-orange-500" :
        "from-red-500 to-rose-600"
      )} />
      
      <CardHeader className="space-y-3 pb-3">
        <div className="flex items-start justify-between gap-2">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <Badge variant="outline" className={cn("capitalize font-medium border", difficultyColor)}>
                {quiz.difficulty}
              </Badge>
              {quiz.is_featured && (
                <Badge variant="secondary" className="bg-yellow-100 text-yellow-800 border-yellow-200 gap-1">
                  <Star className="h-3 w-3 fill-yellow-500 text-yellow-500" />
                  Featured
                </Badge>
              )}
            </div>
            <h3 className="text-xl font-bold line-clamp-2 group-hover:text-violet-600 dark:group-hover:text-violet-400 transition-colors">
              {quiz.title}
            </h3>
          </div>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 -mr-2 text-muted-foreground hover:text-red-500 hover:bg-red-50"
            onClick={handleFavoriteClick}
            disabled={isFavoriteLoading}
            aria-label={isFavorite ? "Remove from favorites" : "Add to favorites"}
          >
            <Heart
              className={cn(
                "h-5 w-5 transition-all duration-300",
                isFavorite
                  ? "fill-red-500 text-red-500 scale-110"
                  : "scale-100"
              )}
            />
          </Button>
        </div>
        <p className="text-sm text-muted-foreground line-clamp-2 leading-relaxed">
          {quiz.description}
        </p>
      </CardHeader>

      <CardContent className="flex-1 pb-4">
        <div className="grid grid-cols-2 gap-3 py-3 border-t border-b border-gray-100/50 dark:border-gray-800/50">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <BookOpen className="h-4 w-4 text-violet-500" />
            <span>{quiz.question_count} questions</span>
          </div>

          {quiz.time_limit && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4 text-violet-500" />
              <span>{quiz.time_limit} min</span>
            </div>
          )}

          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Trophy className="h-4 w-4 text-violet-500" />
            <span>{quiz.points} points</span>
          </div>

          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Badge variant="secondary" className="px-1.5 py-0 h-5 text-xs font-normal bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300">
              {quiz.category}
            </Badge>
          </div>
        </div>
      </CardContent>

      <CardFooter className="pt-0 pb-5 px-6 flex flex-col gap-3">
        {/* Completion Status Badge */}
        {isCompleted && (
          <div className={cn(
            "w-full flex items-center justify-between p-2.5 rounded-xl text-sm font-medium",
            passed ? "bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400" : "bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-400"
          )}>
            <div className="flex items-center gap-1.5">
              {passed ? <CheckCircle className="h-4 w-4" /> : <XCircle className="h-4 w-4" />}
              <span>{passed ? "Passed" : "Failed"}</span>
            </div>
            <span className="font-bold">
              {completedAttempt?.score}/{completedAttempt?.total_points}
            </span>
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex gap-3 w-full">
          {isCompleted ? (
            <>
              <Button variant="outline" className="flex-1 rounded-xl border-violet-200 dark:border-violet-800 hover:bg-violet-50 dark:hover:bg-violet-900/20 hover:text-violet-700 dark:hover:text-violet-300 transition-all duration-300" asChild>
                <Link href={`/quizzes/${quiz.id}`}>Details</Link>
              </Button>
              <Button className="flex-1 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 shadow-md hover:shadow-lg transition-all duration-300 text-white" asChild>
                <Link href={`/quizzes/${quiz.id}/results/${completedAttempt?.id}`}>
                  <Trophy className="mr-2 h-4 w-4" />
                  Results
                </Link>
              </Button>
            </>
          ) : (
            <>
              <Button variant="outline" className="flex-1 rounded-xl border-violet-200 dark:border-violet-800 hover:bg-violet-50 dark:hover:bg-violet-900/20 hover:text-violet-700 dark:hover:text-violet-300 transition-all duration-300" asChild>
                <Link href={`/quizzes/${quiz.id}`}>Details</Link>
              </Button>
              <Button className="flex-1 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 shadow-md hover:shadow-lg transition-all duration-300 text-white" asChild>
                <Link href={`/quizzes/${quiz.id}/take`}>
                  Start Quiz
                  <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                </Link>
              </Button>
            </>
          )}
        </div>
      </CardFooter>
    </Card>
  );
}