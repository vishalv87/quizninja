"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { QuizCard } from "@/components/quiz/QuizCard";
import { useQuizzes } from "@/hooks/useQuizzes";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import { BookOpen, Sparkles } from "lucide-react";
import type { Quiz } from "@/types/quiz";

interface RelatedQuizzesCardProps {
  currentQuizId: string;
  category?: string | null;
  difficulty?: string | null;
}

export function RelatedQuizzesCard({
  currentQuizId,
  category,
  difficulty,
}: RelatedQuizzesCardProps) {
  // Fetch related quizzes based on category
  const { data: quizzes, isLoading } = useQuizzes({
    category: category || undefined,
    limit: 6,
  });

  // Filter out current quiz and limit to 3 results
  const relatedQuizzes = quizzes
    ?.filter((quiz) => quiz.id !== currentQuizId)
    .slice(0, 3);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Sparkles className="h-5 w-5" />
            Related Quizzes
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex justify-center py-8">
            <LoadingSpinner size="md" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!relatedQuizzes || relatedQuizzes.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Sparkles className="h-5 w-5" />
            Related Quizzes
          </CardTitle>
        </CardHeader>
        <CardContent>
          <EmptyState
            icon={BookOpen}
            title="No Related Quizzes"
            description="We couldn't find any related quizzes at the moment."
          />
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Sparkles className="h-5 w-5" />
          Related Quizzes
        </CardTitle>
        {category && (
          <p className="text-sm text-muted-foreground">
            More quizzes in {category.charAt(0).toUpperCase() + category.slice(1)}
          </p>
        )}
      </CardHeader>
      <CardContent className="space-y-4">
        {relatedQuizzes.map((quiz) => (
          <QuizCard key={quiz.id} quiz={quiz} />
        ))}
      </CardContent>
    </Card>
  );
}
