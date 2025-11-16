"use client";

import { QuizList } from "@/components/quiz/QuizList";
import { useFavorites } from "@/hooks/useFavorites";
import { Separator } from "@/components/ui/separator";
import { Heart } from "lucide-react";
import type { Quiz } from "@/types/quiz";

export default function FavoritesPage() {
  const { data: favoritesData, isLoading, error } = useFavorites();

  // Extract quiz data from favorites response
  // The API returns favorites with quiz objects already embedded
  const favoriteQuizzes: Quiz[] = favoritesData?.favorites.map((favorite) => ({
    id: favorite.quiz.id,
    title: favorite.quiz.title,
    description: favorite.quiz.description,
    category: favorite.quiz.category,
    difficulty: favorite.quiz.difficulty,
    time_limit: favorite.quiz.time_limit,
    question_count: favorite.quiz.question_count,
    points: favorite.quiz.points,
    is_featured: favorite.quiz.is_featured,
    tags: favorite.quiz.tags,
    thumbnail_url: favorite.quiz.thumbnail_url,
    created_at: favorite.quiz.created_at,
  })) || [];

  return (
    <div className="container py-8 space-y-6">
      {/* Page Header */}
      <div className="flex items-center gap-3">
        <div className="flex items-center justify-center w-12 h-12 bg-red-100 dark:bg-red-900/20 rounded-full">
          <Heart className="w-6 h-6 text-red-600 dark:text-red-400 fill-red-600 dark:fill-red-400" />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Favorite Quizzes</h1>
          <p className="text-muted-foreground mt-1">
            {isLoading
              ? "Loading your favorites..."
              : favoriteQuizzes.length === 0
              ? "You haven't favorited any quizzes yet"
              : `You have ${favoriteQuizzes.length} favorite ${favoriteQuizzes.length === 1 ? 'quiz' : 'quizzes'}`
            }
          </p>
        </div>
      </div>

      <Separator />

      {/* Quiz List */}
      <QuizList
        quizzes={favoriteQuizzes}
        isLoading={isLoading}
        error={error}
      />
    </div>
  );
}
