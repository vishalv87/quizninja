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
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-violet-600 via-indigo-600 to-purple-700 p-8 text-white shadow-2xl shadow-indigo-500/30 border border-white/10 lg:p-12">
        <div className="relative z-10 max-w-3xl">
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Browse Categories
          </h1>
          <p className="text-lg text-indigo-100 mb-8 max-w-2xl">
            Explore quizzes across different topics and subjects. Find your area of interest and challenge yourself with our diverse collection.
          </p>

          {/* Search Bar Embedded in Hero */}
          <div className="bg-white/10 backdrop-blur-md p-2 rounded-2xl border border-white/20 shadow-lg">
            <div className="relative w-full">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-indigo-200" />
              <Input
                type="text"
                placeholder="Search categories by name or topic..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 pr-10 border-none bg-transparent text-white placeholder:text-indigo-200 focus-visible:ring-0 focus-visible:ring-offset-0"
              />
              {searchQuery && (
                <Button
                  variant="ghost"
                  size="icon"
                  className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 text-indigo-200 hover:bg-white/20 hover:text-white"
                  onClick={handleClearSearch}
                >
                  <X className="h-4 w-4" />
                  <span className="sr-only">Clear search</span>
                </Button>
              )}
            </div>
          </div>
        </div>

        {/* Decorative background elements */}
        <div className="absolute right-0 top-0 -mt-20 -mr-20 h-96 w-96 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute bottom-0 right-20 -mb-20 h-64 w-64 rounded-full bg-indigo-400/20 blur-3xl" />
        <div className="absolute left-10 bottom-10 h-32 w-32 rounded-full bg-purple-400/20 blur-2xl" />
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
