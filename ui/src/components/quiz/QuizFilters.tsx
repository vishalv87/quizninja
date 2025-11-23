"use client";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Card } from "@/components/ui/card";
import { useCategories } from "@/hooks/useCategories";
import { Skeleton } from "@/components/ui/skeleton";

export interface QuizFilterValues {
  category: string;
  difficulty: string;
  isFeatured: boolean;
  showFavoritesOnly?: boolean;
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

  const handleFeaturedChange = (checked: boolean) => {
    onFilterChange({
      ...filters,
      isFeatured: checked,
    });
  };

  const handleFavoritesChange = (checked: boolean) => {
    onFilterChange({
      ...filters,
      showFavoritesOnly: checked,
    });
  };

  return (
    <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
      <div className="flex flex-1 gap-3 w-full sm:w-auto">
          {/* Category Filter */}
          <div className="w-full sm:w-[200px]">
            {categoriesLoading ? (
              <Skeleton className="h-10 w-full rounded-xl" />
            ) : (
              <Select
                value={filters.category || "all"}
                onValueChange={handleCategoryChange}
              >
                <SelectTrigger id="category-filter" className="bg-white/80 dark:bg-background/80 border-gray-200/60 dark:border-gray-700/60 rounded-xl transition-all duration-300 hover:border-primary/30 focus:ring-violet-500/20">
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
              <SelectTrigger id="difficulty-filter" className="bg-white/80 dark:bg-background/80 border-gray-200/60 dark:border-gray-700/60 rounded-xl transition-all duration-300 hover:border-primary/30 focus:ring-violet-500/20">
                <SelectValue placeholder="Difficulty" />
              </SelectTrigger>
              <SelectContent className="rounded-xl">
                <SelectItem value="all">All Difficulties</SelectItem>
                <SelectItem value="beginner">Beginner</SelectItem>
                <SelectItem value="intermediate">Intermediate</SelectItem>
                <SelectItem value="advanced">Advanced</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Toggles */}
        <div className="flex items-center gap-6 w-full sm:w-auto justify-end">
          <div className="flex items-center gap-2">
            <Switch
              id="featured-filter"
              checked={filters.isFeatured}
              onCheckedChange={handleFeaturedChange}
              className="data-[state=checked]:bg-violet-600"
            />
            <Label htmlFor="featured-filter" className="text-sm font-medium cursor-pointer text-muted-foreground hover:text-foreground transition-colors">
              Featured
            </Label>
          </div>

          <div className="flex items-center gap-2">
            <Switch
              id="favorites-filter"
              checked={filters.showFavoritesOnly || false}
              onCheckedChange={handleFavoritesChange}
              className="data-[state=checked]:bg-violet-600"
            />
            <Label htmlFor="favorites-filter" className="text-sm font-medium cursor-pointer text-muted-foreground hover:text-foreground transition-colors">
              Favorites
            </Label>
          </div>
        </div>
      </div>
  );
}