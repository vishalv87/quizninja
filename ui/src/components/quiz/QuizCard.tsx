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
    <Card className="group relative flex flex-col h-full overflow-hidden border-gray-200/60 transition-all duration-300 hover:shadow-xl hover:-translate-y-1 bg-white">
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
            <h3 className="text-xl font-bold line-clamp-2 group-hover:text-primary transition-colors">
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
        <div className="grid grid-cols-2 gap-3 py-3 border-t border-b border-gray-50">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <BookOpen className="h-4 w-4 text-primary/70" />
            <span>{quiz.question_count} questions</span>
          </div>

          {quiz.time_limit && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4 text-primary/70" />
              <span>{quiz.time_limit} min</span>
            </div>
          )}

          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Trophy className="h-4 w-4 text-primary/70" />
            <span>{quiz.points} points</span>
          </div>

          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Badge variant="secondary" className="px-1.5 py-0 h-5 text-xs font-normal">
              {quiz.category}
            </Badge>
          </div>
        </div>
      </CardContent>

      <CardFooter className="pt-0 pb-5 px-6 flex flex-col gap-3">
        {/* Completion Status Badge */}
        {isCompleted && (
          <div className={cn(
            "w-full flex items-center justify-between p-2 rounded-lg text-sm font-medium",
            passed ? "bg-green-50 text-green-700" : "bg-red-50 text-red-700"
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
              <Button variant="outline" className="flex-1 border-primary/20 hover:bg-primary/5 hover:text-primary" asChild>
                <Link href={`/quizzes/${quiz.id}`}>Details</Link>
              </Button>
              <Button className="flex-1 bg-primary hover:bg-primary/90 shadow-sm" asChild>
                <Link href={`/quizzes/${quiz.id}/results/${completedAttempt?.id}`}>
                  <Trophy className="mr-2 h-4 w-4" />
                  Results
                </Link>
              </Button>
            </>
          ) : (
            <>
              <Button variant="outline" className="flex-1 border-primary/20 hover:bg-primary/5 hover:text-primary" asChild>
                <Link href={`/quizzes/${quiz.id}`}>Details</Link>
              </Button>
              <Button className="flex-1 bg-primary hover:bg-primary/90 shadow-sm group-hover:shadow-md transition-all" asChild>
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