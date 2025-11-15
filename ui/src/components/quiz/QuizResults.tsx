"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import type { QuizResults as QuizResultsType } from "@/types/quiz";
import {
  Trophy,
  Target,
  Clock,
  CheckCircle2,
  XCircle,
  Award,
  RotateCcw,
} from "lucide-react";
import Link from "next/link";
import { Progress } from "@/components/ui/progress";

interface QuizResultsProps {
  results: QuizResultsType;
  onRetry?: () => void;
}

export function QuizResults({ results, onRetry }: QuizResultsProps) {
  const { attempt, quiz, percentage, passed } = results;

  // Calculate stats
  const accuracy = (attempt.correct_answers / attempt.total_questions) * 100;
  const incorrectAnswers = attempt.total_questions - attempt.correct_answers;

  // Format time
  const timeSpentMinutes = attempt.time_spent_seconds
    ? Math.floor(attempt.time_spent_seconds / 60)
    : 0;
  const timeSpentSeconds = attempt.time_spent_seconds
    ? attempt.time_spent_seconds % 60
    : 0;

  return (
    <div className="space-y-6">
      {/* Header Card */}
      <Card className={passed ? "border-green-500" : "border-yellow-500"}>
        <CardHeader className="text-center space-y-4">
          {/* Icon */}
          <div className="flex justify-center">
            {passed ? (
              <div className="w-20 h-20 rounded-full bg-green-100 dark:bg-green-900 flex items-center justify-center">
                <Trophy className="w-10 h-10 text-green-600 dark:text-green-400" />
              </div>
            ) : (
              <div className="w-20 h-20 rounded-full bg-yellow-100 dark:bg-yellow-900 flex items-center justify-center">
                <Target className="w-10 h-10 text-yellow-600 dark:text-yellow-400" />
              </div>
            )}
          </div>

          {/* Title */}
          <div>
            <h1 className="text-3xl font-bold mb-2">
              {passed ? "Congratulations!" : "Good Effort!"}
            </h1>
            <p className="text-muted-foreground">
              {passed
                ? "You passed the quiz!"
                : "Keep practicing to improve your score"}
            </p>
          </div>

          {/* Score */}
          <div className="space-y-2">
            <div className="text-6xl font-bold">{percentage.toFixed(1)}%</div>
            <Badge
              variant={passed ? "default" : "secondary"}
              className="text-lg px-4 py-1"
            >
              {passed ? "Passed" : "Not Passed"}
            </Badge>
          </div>

          {/* Progress Bar */}
          <div className="max-w-md mx-auto w-full">
            <Progress
              value={percentage}
              className={`h-3 ${passed ? "bg-green-200 dark:bg-green-900" : ""}`}
            />
          </div>
        </CardHeader>
      </Card>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Score */}
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-lg">
                <Trophy className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-2xl font-bold">{attempt.score}</p>
                <p className="text-xs text-muted-foreground">
                  Points Earned
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Correct Answers */}
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
                <CheckCircle2 className="h-5 w-5 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <p className="text-2xl font-bold">{attempt.correct_answers}</p>
                <p className="text-xs text-muted-foreground">
                  Correct Answers
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Incorrect Answers */}
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-red-100 dark:bg-red-900 rounded-lg">
                <XCircle className="h-5 w-5 text-red-600 dark:text-red-400" />
              </div>
              <div>
                <p className="text-2xl font-bold">{incorrectAnswers}</p>
                <p className="text-xs text-muted-foreground">
                  Incorrect Answers
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Time Spent */}
        {attempt.time_spent_seconds !== undefined && (
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                  <Clock className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                </div>
                <div>
                  <p className="text-2xl font-bold">
                    {timeSpentMinutes}:{timeSpentSeconds.toString().padStart(2, "0")}
                  </p>
                  <p className="text-xs text-muted-foreground">Time Spent</p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>

      {/* Quiz Info */}
      <Card>
        <CardHeader>
          <CardTitle>Quiz Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Quiz Name</span>
            <span className="font-medium">{quiz.title}</span>
          </div>
          <Separator />
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Category</span>
            <Badge variant="secondary">{quiz.category}</Badge>
          </div>
          <Separator />
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Difficulty</span>
            <Badge>{quiz.difficulty}</Badge>
          </div>
          <Separator />
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Total Questions</span>
            <span className="font-medium">{attempt.total_questions}</span>
          </div>
          <Separator />
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground">Accuracy</span>
            <span className="font-medium">{accuracy.toFixed(1)}%</span>
          </div>
        </CardContent>
      </Card>

      {/* Actions */}
      <div className="flex flex-col sm:flex-row gap-4">
        <Link href="/quizzes" className="flex-1">
          <Button variant="outline" size="lg" className="w-full">
            Browse More Quizzes
          </Button>
        </Link>
        {onRetry && (
          <Button size="lg" onClick={onRetry} className="flex-1">
            <RotateCcw className="mr-2 h-4 w-4" />
            Retry Quiz
          </Button>
        )}
      </div>
    </div>
  );
}
