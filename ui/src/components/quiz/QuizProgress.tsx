"use client";

import { Progress } from "@/components/ui/progress";
import { Card } from "@/components/ui/card";

interface QuizProgressProps {
  currentQuestion: number; // 0-indexed
  totalQuestions: number;
  answeredCount: number;
}

export function QuizProgress({
  currentQuestion,
  totalQuestions,
  answeredCount,
}: QuizProgressProps) {
  const progressPercentage = ((currentQuestion + 1) / totalQuestions) * 100;
  const answeredPercentage = (answeredCount / totalQuestions) * 100;

  return (
    <Card className="p-4">
      <div className="space-y-3">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium">Progress</p>
            <p className="text-xs text-muted-foreground">
              Question {currentQuestion + 1} of {totalQuestions}
            </p>
          </div>
          <div className="text-right">
            <p className="text-sm font-medium">
              {answeredCount}/{totalQuestions}
            </p>
            <p className="text-xs text-muted-foreground">Answered</p>
          </div>
        </div>

        {/* Progress Bar */}
        <Progress value={progressPercentage} className="h-2" />

        {/* Stats */}
        <div className="grid grid-cols-2 gap-2 text-xs text-muted-foreground">
          <div>
            {Math.round(progressPercentage)}% complete
          </div>
          <div className="text-right">
            {totalQuestions - answeredCount} remaining
          </div>
        </div>
      </div>
    </Card>
  );
}
