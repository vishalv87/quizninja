"use client";

import { useState, useMemo } from "react";
import { Loader2, Search, X } from "lucide-react";
import { useCategories } from "@/hooks/useCategories";
import { CategoryGrid } from "@/components/categories/CategoryGrid";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useDebounce } from "@/hooks/useDebounce";

export default function CategoriesPage() {
  const { data: categories = [], isLoading, error } = useCategories();
  const [searchQuery, setSearchQuery] = useState("");

  // Debounce search query
  const debouncedSearch = useDebounce(searchQuery, 300);

  // Filter categories based on search
  const filteredCategories = useMemo(() => {
    if (!debouncedSearch) return categories;

    const query = debouncedSearch.toLowerCase();
    return categories.filter(
      (category) =>
        category.display_name.toLowerCase().includes(query) ||
        category.description?.toLowerCase().includes(query) ||
        category.name.toLowerCase().includes(query)
    );
  }, [categories, debouncedSearch]);

  const handleClearSearch = () => {
    setSearchQuery("");
  };

  return (
    <div className="space-y-10 pb-10">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
          Browse Categories
        </h1>
        <p className="text-slate-500 dark:text-slate-400 mt-1">
          Explore quizzes across different topics and subjects
        </p>
      </div>

      {/* Search Bar */}
      <div className="bg-white/60 dark:bg-black/40 backdrop-blur-md p-3 rounded-2xl border border-white/20 dark:border-white/10 shadow-sm">
        <div className="relative w-full">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <Input
            type="text"
            placeholder="Search categories by name or topic..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 pr-10 border-none bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
          />
          {searchQuery && (
            <Button
              variant="ghost"
              size="icon"
              className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-600"
              onClick={handleClearSearch}
            >
              <X className="h-4 w-4" />
              <span className="sr-only">Clear search</span>
            </Button>
          )}
        </div>
      </div>

      <div className="container px-0 md:px-4">
        {/* Results Summary */}
        {!isLoading && !error && (
          <div className="mb-6 flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              {debouncedSearch ? (
                <>
                  {filteredCategories.length} {filteredCategories.length === 1 ? "category" : "categories"} found
                  {filteredCategories.length !== categories.length && (
                    <span className="text-muted-foreground/60"> out of {categories.length}</span>
                  )}
                </>
              ) : (
                <>
                  {categories.length} {categories.length === 1 ? "category" : "categories"} available
                </>
              )}
            </p>
          </div>
        )}

        {/* Loading State */}
        {isLoading && (
          <Card className="border-0 shadow-md rounded-2xl">
            <CardContent className="p-12 flex items-center justify-center gap-3">
              <Loader2 className="h-8 w-8 animate-spin text-violet-600" />
              <p className="text-muted-foreground">Loading categories...</p>
            </CardContent>
          </Card>
        )}

        {/* Error State */}
        {error && (
          <Card className="border-destructive border-0 shadow-md rounded-2xl">
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
          <CategoryGrid
            categories={filteredCategories}
            searchQuery={debouncedSearch}
          />
        )}
      </div>
    </div>
  );
}
