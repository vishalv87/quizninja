"use client";

import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface QuizNavigationProps {
  currentQuestion: number; // 0-indexed
  totalQuestions: number;
  onPrevious: () => void;
  onNext: () => void;
  onSubmit?: () => void;
  isLoading?: boolean;
  canPrevious?: boolean;
  canNext?: boolean;
  showSubmit?: boolean;
}

export function QuizNavigation({
  currentQuestion,
  totalQuestions,
  onPrevious,
  onNext,
  onSubmit,
  isLoading = false,
  canPrevious = true,
  canNext = true,
  showSubmit = false,
}: QuizNavigationProps) {
  const isFirst = currentQuestion === 0;
  const isLast = currentQuestion === totalQuestions - 1;

  return (
    <div className="flex items-center justify-between gap-4">
      {/* Previous Button */}
      <Button
        variant="outline"
        size="lg"
        onClick={onPrevious}
        disabled={isFirst || !canPrevious || isLoading}
      >
        <ChevronLeft className="h-4 w-4 mr-2" />
        Previous
      </Button>

      {/* Question Indicators */}
      <div className="flex items-center gap-1 overflow-x-auto max-w-md">
        {Array.from({ length: totalQuestions }).map((_, index) => (
          <div
            key={index}
            className={`
              w-2 h-2 rounded-full transition-all
              ${
                index === currentQuestion
                  ? "w-8 bg-primary"
                  : index < currentQuestion
                    ? "bg-primary/50"
                    : "bg-muted"
              }
            `}
          />
        ))}
      </div>

      {/* Next or Submit Button */}
      {isLast && showSubmit ? (
        <Button
          size="lg"
          onClick={onSubmit}
          disabled={isLoading}
        >
          {isLoading ? "Submitting..." : "Submit Quiz"}
        </Button>
      ) : (
        <Button
          variant="outline"
          size="lg"
          onClick={onNext}
          disabled={isLast || !canNext || isLoading}
        >
          Next
          <ChevronRight className="h-4 w-4 ml-2" />
        </Button>
      )}
    </div>
  );
}
