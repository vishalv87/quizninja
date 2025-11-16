"use client";

import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type { Quiz } from "@/types/quiz";
import { BookOpen, Clock, Star, Trophy, Heart } from "lucide-react";
import Link from "next/link";
import { useIsFavorite, useToggleFavorite } from "@/hooks/useFavorites";
import { cn } from "@/lib/utils";

interface QuizCardProps {
  quiz: Quiz;
}

export function QuizCard({ quiz }: QuizCardProps) {
  // Check if quiz is favorited
  const { data: isFavorite, isLoading: isFavoriteLoading } = useIsFavorite(quiz.id);
  const { toggle } = useToggleFavorite();

  // Determine difficulty color
  const difficultyColor = {
    beginner: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
    intermediate:
      "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
    advanced: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  }[quiz.difficulty];

  // Handle favorite toggle
  const handleFavoriteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    toggle(quiz.id, isFavorite || false);
  };

  return (
    <Card className="hover:shadow-lg transition-shadow duration-300 flex flex-col h-full">
      <CardHeader className="space-y-2">
        <div className="flex items-start justify-between gap-2">
          <h3 className="text-xl font-bold line-clamp-2 flex-1">
            {quiz.title}
          </h3>
          <div className="flex items-center gap-1 flex-shrink-0">
            {quiz.is_featured && (
              <Star className="h-5 w-5 text-yellow-500 fill-yellow-500" />
            )}
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={handleFavoriteClick}
              disabled={isFavoriteLoading}
              aria-label={isFavorite ? "Remove from favorites" : "Add to favorites"}
            >
              <Heart
                className={cn(
                  "h-5 w-5 transition-colors",
                  isFavorite
                    ? "text-red-600 fill-red-600 dark:text-red-400 dark:fill-red-400"
                    : "text-muted-foreground hover:text-red-600 dark:hover:text-red-400"
                )}
              />
            </Button>
          </div>
        </div>
        <p className="text-sm text-muted-foreground line-clamp-2">
          {quiz.description}
        </p>
      </CardHeader>

      <CardContent className="flex-1">
        <div className="space-y-3">
          {/* Category and Difficulty */}
          <div className="flex gap-2 flex-wrap">
            <Badge variant="secondary">{quiz.category}</Badge>
            <Badge className={difficultyColor}>
              {quiz.difficulty.charAt(0).toUpperCase() +
                quiz.difficulty.slice(1)}
            </Badge>
          </div>

          {/* Quiz Stats */}
          <div className="grid grid-cols-2 gap-3 pt-2">
            <div className="flex items-center gap-2 text-sm">
              <BookOpen className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">
                {quiz.question_count} questions
              </span>
            </div>

            {quiz.time_limit_minutes && (
              <div className="flex items-center gap-2 text-sm">
                <Clock className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">
                  {quiz.time_limit_minutes} min
                </span>
              </div>
            )}

            <div className="flex items-center gap-2 text-sm">
              <Trophy className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">
                {quiz.total_points} points
              </span>
            </div>

            {quiz.average_score !== undefined && (
              <div className="flex items-center gap-2 text-sm">
                <Star className="h-4 w-4 text-muted-foreground" />
                <span className="text-muted-foreground">
                  {quiz.average_score.toFixed(0)}% avg
                </span>
              </div>
            )}
          </div>

          {/* Attempts count if available */}
          {quiz.attempts_count !== undefined && quiz.attempts_count > 0 && (
            <p className="text-xs text-muted-foreground pt-1">
              {quiz.attempts_count} attempt
              {quiz.attempts_count !== 1 ? "s" : ""}
            </p>
          )}
        </div>
      </CardContent>

      <CardFooter className="flex gap-2">
        <Link href={`/quizzes/${quiz.id}`} className="flex-1">
          <Button variant="outline" className="w-full">
            View Details
          </Button>
        </Link>
        <Link href={`/quizzes/${quiz.id}/take`} className="flex-1">
          <Button className="w-full">Start Quiz</Button>
        </Link>
      </CardFooter>
    </Card>
  );
}