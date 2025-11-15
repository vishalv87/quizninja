"use client";

import { QuizCard } from "./QuizCard";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/common/EmptyState";
import type { Quiz } from "@/types/quiz";
import { BookOpen } from "lucide-react";

interface QuizListProps {
  quizzes: Quiz[];
  isLoading?: boolean;
  error?: Error | null;
}

export function QuizList({ quizzes, isLoading, error }: QuizListProps) {
  // Loading state
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {Array.from({ length: 6 }).map((_, index) => (
          <QuizCardSkeleton key={index} />
        ))}
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <EmptyState
        icon={BookOpen}
        title="Error Loading Quizzes"
        description={
          error.message || "Failed to load quizzes. Please try again."
        }
      />
    );
  }

  // Empty state
  if (!quizzes || quizzes.length === 0) {
    return (
      <EmptyState
        icon={BookOpen}
        title="No Quizzes Found"
        description="No quizzes match your criteria. Try adjusting your filters."
      />
    );
  }

  // Quiz list
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {quizzes.map((quiz) => (
        <QuizCard key={quiz.id} quiz={quiz} />
      ))}
    </div>
  );
}

// Skeleton loader for quiz cards
function QuizCardSkeleton() {
  return (
    <div className="border rounded-lg p-6 space-y-4">
      <div className="space-y-2">
        <Skeleton className="h-6 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </div>
      <div className="flex gap-2">
        <Skeleton className="h-6 w-20" />
        <Skeleton className="h-6 w-20" />
      </div>
      <div className="grid grid-cols-2 gap-3">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-full" />
      </div>
      <div className="flex gap-2">
        <Skeleton className="h-10 flex-1" />
        <Skeleton className="h-10 flex-1" />
      </div>
    </div>
  );
}