"use client";

import { useEffect, useState } from "react";
import { QuizList } from "@/components/quiz/QuizList";
import { useFavorites } from "@/hooks/useFavorites";
import { useQuizzes } from "@/hooks/useQuizzes";
import { Separator } from "@/components/ui/separator";
import { Heart } from "lucide-react";
import type { Quiz } from "@/types/quiz";

export default function FavoritesPage() {
  const { data: favoriteIds, isLoading: favoritesLoading, error: favoritesError } = useFavorites();
  const { data: allQuizzes, isLoading: quizzesLoading, error: quizzesError } = useQuizzes({});
  const [favoriteQuizzes, setFavoriteQuizzes] = useState<Quiz[]>([]);

  // Filter quizzes that are in favorites
  useEffect(() => {
    if (favoriteIds && allQuizzes) {
      const filtered = allQuizzes.filter((quiz) =>
        favoriteIds.includes(quiz.id)
      );
      setFavoriteQuizzes(filtered);
    }
  }, [favoriteIds, allQuizzes]);

  const isLoading = favoritesLoading || quizzesLoading;
  const error = favoritesError || quizzesError;

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
