"use client";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useCategories } from "@/hooks/useCategories";
import { Skeleton } from "@/components/ui/skeleton";
import { DIFFICULTY_OPTIONS } from "@/constants";

export interface QuizFilterValues {
  category: string;
  difficulty: string;
}

interface QuizFiltersProps {
  filters: QuizFilterValues;
  onFilterChange: (filters: QuizFilterValues) => void;
}

export function QuizFilters({ filters, onFilterChange }: QuizFiltersProps) {
  const { data: categories, isLoading: categoriesLoading } = useCategories();

  const handleCategoryChange = (value: string) => {
    onFilterChange({
      ...filters,
      category: value === "all" ? "" : value,
    });
  };

  const handleDifficultyChange = (value: string) => {
    onFilterChange({
      ...filters,
      difficulty: value === "all" ? "" : value,
    });
  };

  return (
    <div className="flex gap-3 items-center justify-center">
      {/* Category Filter */}
      <div className="w-full sm:w-[200px]">
        {categoriesLoading ? (
          <Skeleton className="h-10 w-full rounded-xl" />
        ) : (
          <Select
            value={filters.category || "all"}
            onValueChange={handleCategoryChange}
          >
            <SelectTrigger id="category-filter" className="bg-white/90 dark:bg-background/90 backdrop-blur-sm border-gray-200/50 dark:border-gray-700/50 rounded-xl transition-all duration-300 hover:border-violet-400/50 dark:hover:border-violet-600/50 focus:ring-violet-500/20 shadow-sm">
              <SelectValue placeholder="Category" />
            </SelectTrigger>
            <SelectContent className="rounded-xl">
              <SelectItem value="all">All Categories</SelectItem>
              {categories?.map((category) => (
                <SelectItem key={category.id} value={category.id}>
                  {category.display_name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
      </div>

      {/* Difficulty Filter */}
      <div className="w-full sm:w-[180px]">
        <Select
          value={filters.difficulty || "all"}
          onValueChange={handleDifficultyChange}
        >
          <SelectTrigger id="difficulty-filter" className="bg-white/90 dark:bg-background/90 backdrop-blur-sm border-gray-200/50 dark:border-gray-700/50 rounded-xl transition-all duration-300 hover:border-violet-400/50 dark:hover:border-violet-600/50 focus:ring-violet-500/20 shadow-sm">
            <SelectValue placeholder="Difficulty" />
          </SelectTrigger>
          <SelectContent className="rounded-xl">
            <SelectItem value="all">All Difficulties</SelectItem>
            {DIFFICULTY_OPTIONS.map(({ value, label }) => (
              <SelectItem key={value} value={value}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
