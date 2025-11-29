"use client";

import { QuizList } from "@/components/quiz/QuizList";
import { useFavorites } from "@/hooks/useFavorites";
import { Heart } from "lucide-react";
import type { Quiz } from "@/types/quiz";
import { PageHero } from "@/components/common/PageHero";
import { GlassCard } from "@/components/common/GlassCard";

export default function FavoritesPage() {
  const { data: favoritesData, isLoading, error } = useFavorites();

  // Extract quiz data from favorites response
  // The API returns favorites with quiz objects already embedded
  const favoriteQuizzes: Quiz[] = favoritesData?.favorites.map((favorite) => ({
    id: favorite.quiz.id,
    title: favorite.quiz.title,
    description: favorite.quiz.description,
    category: favorite.quiz.category,
    difficulty: favorite.quiz.difficulty as Quiz["difficulty"],
    time_limit: favorite.quiz.time_limit,
    question_count: favorite.quiz.question_count,
    points: favorite.quiz.points,
    is_featured: favorite.quiz.is_featured,
    tags: favorite.quiz.tags,
    thumbnail_url: favorite.quiz.thumbnail_url,
    created_at: favorite.quiz.created_at,
    updated_at: favorite.quiz.created_at,
  })) || [];

  const getDescription = () => {
    if (isLoading) return "Loading your favorite quizzes...";
    if (favoriteQuizzes.length === 0) return "Save your favorite quizzes here for quick access. Start exploring and add some!";
    return `You have ${favoriteQuizzes.length} favorite ${favoriteQuizzes.length === 1 ? 'quiz' : 'quizzes'}. Ready to challenge yourself?`;
  };

  return (
    <div className="space-y-10 pb-10">
      {/* Hero Section */}
      <PageHero
        title="Favorites"
        icon="❤️"
        description={getDescription()}
      />

      <div className="container px-0 md:px-4">
        <GlassCard padding="none" rounded="2xl">
          <div className="p-6 border-b border-white/10">
            <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
              <span className="bg-gradient-to-br from-red-500 to-pink-500 text-white p-1.5 rounded-lg shadow-sm">
                <Heart className="h-4 w-4" />
              </span>
              Your Favorite Quizzes
            </h2>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
              {isLoading
                ? "Loading..."
                : favoriteQuizzes.length === 0
                ? "No favorites yet"
                : `${favoriteQuizzes.length} ${favoriteQuizzes.length === 1 ? 'quiz' : 'quizzes'} saved`}
            </p>
          </div>
          <div className="p-6">
            <QuizList
              quizzes={favoriteQuizzes}
              isLoading={isLoading}
              error={error}
            />
          </div>
        </GlassCard>
      </div>
    </div>
  );
}
