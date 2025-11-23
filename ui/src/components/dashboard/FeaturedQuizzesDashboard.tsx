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
      <Card className="border-none shadow-sm bg-white/50 backdrop-blur-sm h-full">
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
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <Skeleton key={i} className="h-24 rounded-xl" />
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
    <Card className="border-none shadow-sm bg-white/50 backdrop-blur-sm h-full flex flex-col">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
              Featured Quizzes
            </CardTitle>
            <CardDescription>Popular and recommended quizzes</CardDescription>
          </div>
          <Button variant="ghost" size="sm" className="text-muted-foreground hover:text-primary" asChild>
            <Link href="/quizzes?featured=true">
              View All
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent className="flex-1">
        <div className="space-y-4">
          {displayQuizzes.map((quiz: any) => (
            <Link key={quiz.id} href={`/quizzes/${quiz.id}`} className="block group">
              <div className="relative overflow-hidden rounded-xl border bg-card p-4 transition-all duration-300 hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5">
                <div className="flex gap-4">
                  <div className="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-lg bg-muted">
                    {/* Placeholder for quiz image if available, or a colored box */}
                    <div className="absolute inset-0 flex items-center justify-center bg-gradient-to-br from-indigo-100 to-purple-100 text-indigo-500 font-bold text-xl">
                      {quiz.title.charAt(0)}
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <h4 className="font-semibold leading-none tracking-tight mb-2 group-hover:text-primary transition-colors truncate">
                      {quiz.title}
                    </h4>
                    <p className="text-sm text-muted-foreground line-clamp-2 mb-2">
                      {quiz.description}
                    </p>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span className="bg-secondary px-2 py-0.5 rounded-full">
                        {quiz.category || 'General'}
                      </span>
                      <span>•</span>
                      <span>{quiz.questions_count || 0} questions</span>
                    </div>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
