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
    <Card className="p-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Category Filter */}
        <div className="space-y-2">
          <Label htmlFor="category-filter">Category</Label>
          {categoriesLoading ? (
            <Skeleton className="h-10 w-full" />
          ) : (
            <Select
              value={filters.category || "all"}
              onValueChange={handleCategoryChange}
            >
              <SelectTrigger id="category-filter">
                <SelectValue placeholder="All Categories" />
              </SelectTrigger>
              <SelectContent>
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
        <div className="space-y-2">
          <Label htmlFor="difficulty-filter">Difficulty</Label>
          <Select
            value={filters.difficulty || "all"}
            onValueChange={handleDifficultyChange}
          >
            <SelectTrigger id="difficulty-filter">
              <SelectValue placeholder="All Difficulties" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Difficulties</SelectItem>
              <SelectItem value="beginner">Beginner</SelectItem>
              <SelectItem value="intermediate">Intermediate</SelectItem>
              <SelectItem value="advanced">Advanced</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Featured Filter */}
        <div className="space-y-2">
          <Label htmlFor="featured-filter">Featured Only</Label>
          <div className="flex items-center h-10 px-3 border rounded-md">
            <Switch
              id="featured-filter"
              checked={filters.isFeatured}
              onCheckedChange={handleFeaturedChange}
            />
            <Label htmlFor="featured-filter" className="ml-3 cursor-pointer">
              {filters.isFeatured ? "Yes" : "No"}
            </Label>
          </div>
        </div>

        {/* Favorites Filter */}
        <div className="space-y-2">
          <Label htmlFor="favorites-filter">Favorites Only</Label>
          <div className="flex items-center h-10 px-3 border rounded-md">
            <Switch
              id="favorites-filter"
              checked={filters.showFavoritesOnly || false}
              onCheckedChange={handleFavoritesChange}
            />
            <Label htmlFor="favorites-filter" className="ml-3 cursor-pointer">
              {filters.showFavoritesOnly ? "Yes" : "No"}
            </Label>
          </div>
        </div>
      </div>
    </Card>
  );
}