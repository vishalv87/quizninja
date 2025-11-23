"use client";

import { useState, useEffect, useMemo } from "react";
import { QuizList } from "@/components/quiz/QuizList";
import { QuizFilters, QuizFilterValues } from "@/components/quiz/QuizFilters";
import { QuizSearch } from "@/components/quiz/QuizSearch";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useFeaturedQuizzes } from "@/hooks/useFeaturedQuizzes";
import { useFavorites } from "@/hooks/useFavorites";
import { useDebounce } from "@/hooks/useDebounce";
import { useCompletedQuizMap } from "@/hooks/useCompletedQuizMap";
import { Separator } from "@/components/ui/separator";
import { Card } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { Quiz } from "@/types/quiz";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

const ITEMS_PER_PAGE = 12;

export default function QuizzesPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [filters, setFilters] = useState<QuizFilterValues>({
    category: "",
    difficulty: "",
    isFeatured: false,
    showFavoritesOnly: false,
  });
  const [currentPage, setCurrentPage] = useState(1);
  const [filteredQuizzes, setFilteredQuizzes] = useState<Quiz[]>([]);

  // Debounce search query to avoid excessive API calls
  const debouncedSearch = useDebounce(searchQuery, 300);

  // Reset to first page when filters change
  const handleFilterChange = (newFilters: QuizFilterValues) => {
    setFilters(newFilters);
    setCurrentPage(1);
  };

  const handleSearchChange = (query: string) => {
    setSearchQuery(query);
    setCurrentPage(1);
  };

  // Build filter object for API
  const apiFilters = {
    search: debouncedSearch || undefined,
    category: filters.category || undefined,
    difficulty: filters.difficulty || undefined,
    is_featured: filters.isFeatured || undefined,
    limit: ITEMS_PER_PAGE,
    offset: (currentPage - 1) * ITEMS_PER_PAGE,
  };

  // Fetch quizzes based on filters
  const {
    data: allQuizzes,
    isLoading: allQuizzesLoading,
    error: allQuizzesError,
  } = useQuizzes(apiFilters);

  // Fetch featured quizzes separately
  const {
    data: featuredQuizzes,
    isLoading: featuredLoading,
    error: featuredError,
  } = useFeaturedQuizzes();

  // Fetch favorites for filtering
  const { data: favoritesData } = useFavorites();

  // Extract favorite quiz IDs from the favorites response (memoized to prevent infinite loops)
  const favoriteIds = useMemo(
    () => favoritesData?.favorites.map((fav) => fav.quiz_id) || [],
    [favoritesData]
  );

  // Extract all quiz IDs for fetching completion status (combine all and featured)
  const allQuizIds = useMemo(() => {
    const ids = new Set<string>();
    allQuizzes?.forEach((quiz) => ids.add(quiz.id));
    featuredQuizzes?.forEach((quiz) => ids.add(quiz.id));
    return Array.from(ids);
  }, [allQuizzes, featuredQuizzes]);

  // Fetch completed quiz attempts for displaying completion status
  const { data: completedQuizMap } = useCompletedQuizMap(allQuizIds);

  // Filter quizzes by favorites when showFavoritesOnly is true
  useEffect(() => {
    if (filters.showFavoritesOnly && favoriteIds.length > 0 && allQuizzes) {
      const filtered = allQuizzes.filter((quiz) =>
        favoriteIds.includes(quiz.id)
      );
      setFilteredQuizzes(filtered);
    } else {
      setFilteredQuizzes(allQuizzes || []);
    }
  }, [filters.showFavoritesOnly, favoriteIds, allQuizzes]);

  return (
    <div className="space-y-8 pb-10">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-r from-blue-600 to-cyan-600 p-8 text-white shadow-xl lg:p-12">
        <div className="relative z-10 max-w-3xl">
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Explore Quizzes
          </h1>
          <p className="text-lg text-blue-100 mb-8 max-w-2xl">
            Challenge yourself with our diverse collection of quizzes. Filter by category, difficulty, or search for specific topics to test your knowledge.
          </p>
          
          {/* Search Bar Embedded in Hero */}
          <div className="bg-white/10 backdrop-blur-md p-2 rounded-2xl border border-white/20 shadow-lg">
            <QuizSearch
              value={searchQuery}
              onChange={handleSearchChange}
              placeholder="Search quizzes by title or description..."
              className="border-none bg-transparent text-white placeholder:text-blue-200 focus-visible:ring-0 focus-visible:ring-offset-0"
            />
          </div>
        </div>
        
        {/* Decorative background elements */}
        <div className="absolute right-0 top-0 -mt-10 -mr-10 h-64 w-64 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute bottom-0 right-20 -mb-10 h-40 w-40 rounded-full bg-cyan-400/20 blur-2xl" />
      </div>

      <div className="container px-0 md:px-4">
        {/* Filters and Tabs */}
        <div className="flex flex-col gap-6 md:flex-row md:items-start md:justify-between mb-8">
          <Tabs defaultValue="all" className="w-full">
            <div className="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between mb-6">
              <TabsList className="grid w-full max-w-md grid-cols-2">
                <TabsTrigger value="all">All Quizzes</TabsTrigger>
                <TabsTrigger value="featured">Featured</TabsTrigger>
              </TabsList>
              
              <div className="flex-1 lg:max-w-2xl">
                <QuizFilters filters={filters} onFilterChange={handleFilterChange} />
              </div>
            </div>

            {/* All Quizzes Tab */}
            <TabsContent value="all" className="space-y-8">
              <QuizList
                quizzes={filteredQuizzes}
                isLoading={allQuizzesLoading}
                error={allQuizzesError}
                completedQuizMap={completedQuizMap}
              />

              {/* Pagination for All Quizzes */}
              {filteredQuizzes && filteredQuizzes.length > 0 && (
                <div className="flex justify-center pt-4">
                  <Pagination>
                    <PaginationContent>
                      <PaginationItem>
                        <PaginationPrevious
                          onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                          className={
                            currentPage === 1 ? "pointer-events-none opacity-50" : "cursor-pointer"
                          }
                        />
                      </PaginationItem>

                      {/* Show page numbers */}
                      {[...Array(Math.min(5, Math.ceil(100 / ITEMS_PER_PAGE)))].map((_, i) => {
                        const pageNum = i + 1;
                        return (
                          <PaginationItem key={pageNum}>
                            <PaginationLink
                              onClick={() => setCurrentPage(pageNum)}
                              isActive={currentPage === pageNum}
                              className="cursor-pointer"
                            >
                              {pageNum}
                            </PaginationLink>
                          </PaginationItem>
                        );
                      })}

                      <PaginationItem>
                        <PaginationNext
                          onClick={() => setCurrentPage((p) => p + 1)}
                          className={
                            filteredQuizzes.length < ITEMS_PER_PAGE
                              ? "pointer-events-none opacity-50"
                              : "cursor-pointer"
                          }
                        />
                      </PaginationItem>
                    </PaginationContent>
                  </Pagination>
                </div>
              )}
            </TabsContent>

            {/* Featured Quizzes Tab */}
            <TabsContent value="featured">
              <QuizList
                quizzes={featuredQuizzes || []}
                isLoading={featuredLoading}
                error={featuredError}
                completedQuizMap={completedQuizMap}
              />
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}