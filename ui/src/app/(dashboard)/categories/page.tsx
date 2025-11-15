"use client";

import { Loader2, FolderOpen } from "lucide-react";
import { useCategories } from "@/hooks/useCategories";
import { CategoryGrid } from "@/components/categories/CategoryGrid";
import { Card, CardContent } from "@/components/ui/card";

export default function CategoriesPage() {
  const { data: categories = [], isLoading, error } = useCategories();

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
          <FolderOpen className="h-8 w-8" />
          Browse Categories
        </h1>
        <p className="text-muted-foreground mt-2">
          Explore quizzes across different topics and subjects
        </p>
      </div>

      {/* Loading State */}
      {isLoading && (
        <Card>
          <CardContent className="p-12 flex items-center justify-center gap-3">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="text-muted-foreground">Loading categories...</p>
          </CardContent>
        </Card>
      )}

      {/* Error State */}
      {error && (
        <Card className="border-destructive">
          <CardContent className="p-12 text-center">
            <p className="text-destructive font-semibold">Error loading categories</p>
            <p className="text-sm text-muted-foreground mt-2">
              {error instanceof Error ? error.message : "Failed to load categories. Please try again."}
            </p>
          </CardContent>
        </Card>
      )}

      {/* Categories Grid */}
      {!isLoading && !error && (
        <>
          {categories.length > 0 && (
            <div className="mb-4">
              <p className="text-sm text-muted-foreground">
                {categories.length} {categories.length === 1 ? "category" : "categories"} available
              </p>
            </div>
          )}
          <CategoryGrid categories={categories} />
        </>
      )}
    </div>
  );
}
