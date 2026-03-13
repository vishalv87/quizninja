"use client";

import { Card } from "@/components/ui/card";
import { Clock, AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";

interface QuizTimerProps {
  timeRemaining: number | null; // in seconds
}

export function QuizTimer({ timeRemaining }: QuizTimerProps) {
  // No time limit
  if (timeRemaining === null) {
    return (
      <Card className="p-4">
        <div className="flex items-center gap-3">
          <Clock className="h-5 w-5 text-muted-foreground" />
          <div>
            <p className="text-sm font-medium">No Time Limit</p>
            <p className="text-xs text-muted-foreground">
              Take your time to answer
            </p>
          </div>
        </div>
      </Card>
    );
  }

  // Format time as MM:SS
  const minutes = Math.floor(timeRemaining / 60);
  const seconds = timeRemaining % 60;
  const formattedTime = `${minutes.toString().padStart(2, "0")}:${seconds
    .toString()
    .padStart(2, "0")}`;

  // Determine urgency level
  const isUrgent = timeRemaining <= 60; // Last minute
  const isCritical = timeRemaining <= 10; // Last 10 seconds

  return (
    <Card
      className={cn(
        "p-4 transition-colors",
        isCritical && "border-red-500 bg-red-50 dark:bg-red-950/20",
        isUrgent &&
          !isCritical &&
          "border-yellow-500 bg-yellow-50 dark:bg-yellow-950/20"
      )}
    >
      <div className="flex items-center gap-3">
        {isCritical ? (
          <AlertCircle className="h-5 w-5 text-red-500 animate-pulse" />
        ) : (
          <Clock
            className={cn(
              "h-5 w-5",
              isUrgent
                ? "text-yellow-600 dark:text-yellow-400"
                : "text-muted-foreground"
            )}
          />
        )}
        <div className="flex-1">
          <p
            className={cn(
              "text-2xl font-bold font-mono",
              isCritical && "text-red-600 dark:text-red-400",
              isUrgent && !isCritical && "text-yellow-600 dark:text-yellow-400"
            )}
          >
            {formattedTime}
          </p>
          <p className="text-xs text-muted-foreground">Time Remaining</p>
        </div>
      </div>
    </Card>
  );
}
