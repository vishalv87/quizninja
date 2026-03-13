"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { BookOpen, Heart, Shuffle, ArrowRight } from "lucide-react";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { Category } from "@/types/quiz";
import { cn } from "@/lib/utils";

interface CategoryCardProps {
  category: Category;
  isFavorite?: boolean;
  onToggleFavorite?: (categoryId: string) => void;
}

export function CategoryCard({
  category,
  isFavorite = false,
  onToggleFavorite,
}: CategoryCardProps) {
  const [imageError, setImageError] = useState(false);

  // Handle favorite toggle
  const handleFavoriteClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onToggleFavorite?.(category.id);
  };

  // Handle quick start click
  const handleQuickStartClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    // Navigate to quizzes filtered by this category
    window.location.href = `/quizzes?category=${category.name}`;
  };

  return (
    <div className="group relative h-full p-[2px] rounded-2xl bg-gradient-to-br from-violet-500 to-indigo-500 transition-all duration-300 hover:shadow-xl hover:-translate-y-1 shadow-md">
      <Card className="flex flex-col h-full overflow-hidden border-0 bg-white dark:bg-background rounded-[14px]">
        <CardHeader className="space-y-3 pb-3">
        <div className="flex items-start justify-between gap-2">
          <div className="flex items-center gap-3">
            {/* Category Icon */}
            {category.icon_url && !imageError ? (
              <div className="w-12 h-12 rounded-xl flex items-center justify-center overflow-hidden relative bg-violet-100 dark:bg-violet-900/30">
                <Image
                  src={category.icon_url}
                  alt={category.display_name}
                  fill
                  className="object-cover"
                  onError={() => setImageError(true)}
                  unoptimized
                />
              </div>
            ) : (
              <div className="w-12 h-12 rounded-xl flex items-center justify-center bg-violet-100 dark:bg-violet-900/30">
                <BookOpen className="h-6 w-6 text-violet-700 dark:text-violet-400" />
              </div>
            )}
            <div className="space-y-1">
              <h3 className="text-xl font-bold line-clamp-1 group-hover:text-violet-600 dark:group-hover:text-violet-400 transition-colors">
                {category.display_name}
              </h3>
              <Badge
                variant="secondary"
                className="font-medium bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-400"
              >
                {category.quiz_count} {category.quiz_count === 1 ? "Quiz" : "Quizzes"}
              </Badge>
            </div>
          </div>

          {/* Favorite Button */}
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 -mr-2 -mt-1 text-muted-foreground hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"
            onClick={handleFavoriteClick}
            aria-label={isFavorite ? "Remove from favorites" : "Add to favorites"}
          >
            <Heart
              className={cn(
                "h-5 w-5 transition-all duration-300",
                isFavorite
                  ? "fill-red-500 text-red-500 scale-110"
                  : "scale-100"
              )}
            />
          </Button>
        </div>

        <p className="text-sm text-muted-foreground line-clamp-2 leading-relaxed min-h-[2.5rem]">
          {category.description || "Explore quizzes in this category"}
        </p>
      </CardHeader>

      <CardContent className="flex-1 pb-0" />

      <CardFooter className="pt-0 pb-5 px-6 flex gap-3">
        {/* Browse Category Button */}
        <Button
          variant="outline"
          className="flex-1 rounded-xl border-violet-200 dark:border-violet-800 hover:bg-violet-50 dark:hover:bg-violet-900/20 hover:text-violet-700 dark:hover:text-violet-300 transition-all duration-300"
          asChild
        >
          <Link href={`/quizzes/category/${category.id}`}>
            Browse
            <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
          </Link>
        </Button>

        {/* Quick Start Button */}
        <Button
          className="flex-1 rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 shadow-md hover:shadow-lg transition-all duration-300 text-white"
          onClick={handleQuickStartClick}
          disabled={category.quiz_count === 0}
        >
          <Shuffle className="mr-2 h-4 w-4" />
          Quick Start
        </Button>
      </CardFooter>
      </Card>
    </div>
  );
}
