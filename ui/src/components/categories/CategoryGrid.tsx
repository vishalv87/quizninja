"use client";

import { CategoryCard } from "./CategoryCard";
import { useCategoryFavorites } from "@/hooks/useCategoryFavorites";
import type { Category } from "@/types/quiz";
import { FolderX } from "lucide-react";

interface CategoryGridProps {
  categories: Category[];
  searchQuery?: string;
}

export function CategoryGrid({ categories, searchQuery }: CategoryGridProps) {
  const { isFavorite, toggleFavorite } = useCategoryFavorites();

  if (categories.length === 0) {
    return (
      <div className="text-center py-16 px-4">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-violet-100 dark:bg-violet-900/20 mb-4">
          <FolderX className="h-8 w-8 text-violet-600 dark:text-violet-400" />
        </div>
        <h3 className="text-lg font-semibold mb-2">No categories found</h3>
        <p className="text-muted-foreground max-w-md mx-auto">
          {searchQuery
            ? `No categories match "${searchQuery}". Try a different search term.`
            : "No categories available at the moment. Check back later!"}
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {categories.map((category) => (
        <CategoryCard
          key={category.id}
          category={category}
          isFavorite={isFavorite(category.id)}
          onToggleFavorite={toggleFavorite}
        />
      ))}
    </div>
  );
}
