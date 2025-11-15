"use client";

import Link from "next/link";
import { ArrowRight, Star } from "lucide-react";
import { useFeaturedQuizzes } from "@/hooks/useQuiz";
import { QuizCard } from "@/components/quiz/QuizCard";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";

export function FeaturedQuizzesDashboard() {
  const { data: quizzes, isLoading } = useFeaturedQuizzes();

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                Featured Quizzes
              </CardTitle>
              <CardDescription>Popular and recommended quizzes</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {[...Array(3)].map((_, i) => (
              <Skeleton key={i} className="h-48" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!quizzes || quizzes.length === 0) {
    return null;
  }

  // Show only first 3 featured quizzes on dashboard
  const displayQuizzes = quizzes.slice(0, 3);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
              Featured Quizzes
            </CardTitle>
            <CardDescription>Popular and recommended quizzes</CardDescription>
          </div>
          <Button variant="ghost" size="sm" asChild>
            <Link href="/quizzes?featured=true">
              View All
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {displayQuizzes.map((quiz: any) => (
            <QuizCard key={quiz.id} quiz={quiz} />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
