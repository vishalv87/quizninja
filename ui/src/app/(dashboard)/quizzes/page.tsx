"use client";

import { useState, useEffect, useMemo } from "react";
import { QuizList } from "@/components/quiz/QuizList";
import { QuizFilters, QuizFilterValues } from "@/components/quiz/QuizFilters";
import { QuizSearch } from "@/components/quiz/QuizSearch";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useFavorites } from "@/hooks/useFavorites";
import { useDebounce } from "@/hooks/useDebounce";
import { useCompletedQuizMap } from "@/hooks/useCompletedQuizMap";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { Quiz } from "@/types/quiz";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

const ITEMS_PER_PAGE = 12;

export default function QuizzesPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [activeTab, setActiveTab] = useState("all");
  const [filters, setFilters] = useState<QuizFilterValues>({
    category: "",
    difficulty: "",
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

  const handleTabChange = (value: string) => {
    setActiveTab(value);
    setCurrentPage(1);
  };

  // Build filter object for API
  const apiFilters = {
    search: debouncedSearch || undefined,
    category: filters.category || undefined,
    difficulty: filters.difficulty || undefined,
    limit: ITEMS_PER_PAGE,
    offset: (currentPage - 1) * ITEMS_PER_PAGE,
  };

  // Fetch quizzes based on filters
  const {
    data: allQuizzes,
    isLoading: allQuizzesLoading,
    error: allQuizzesError,
  } = useQuizzes(apiFilters);

  // Fetch favorites for filtering
  const { data: favoritesData } = useFavorites();

  // Extract favorite quiz IDs from the favorites response (memoized to prevent infinite loops)
  const favoriteIds = useMemo(
    () => favoritesData?.favorites.map((fav) => fav.quiz_id) || [],
    [favoritesData]
  );

  // Extract all quiz IDs for fetching completion status
  const allQuizIds = useMemo(() => {
    const ids = new Set<string>();
    allQuizzes?.forEach((quiz) => ids.add(quiz.id));
    return Array.from(ids);
  }, [allQuizzes]);

  // Fetch completed quiz attempts for displaying completion status
  const { data: completedQuizMap } = useCompletedQuizMap(allQuizIds);

  // Unified filtering logic based on active tab
  useEffect(() => {
    if (!allQuizzes) {
      setFilteredQuizzes([]);
      return;
    }

    let filtered = [...allQuizzes];

    // Apply tab-based filtering
    switch (activeTab) {
      case "featured":
        filtered = filtered.filter((quiz) => quiz.is_featured);
        break;
      case "completed":
        filtered = filtered.filter((quiz) => completedQuizMap?.has(quiz.id));
        break;
      case "not-completed":
        filtered = filtered.filter((quiz) => !completedQuizMap?.has(quiz.id));
        break;
      case "favorites":
        filtered = filtered.filter((quiz) => favoriteIds.includes(quiz.id));
        break;
      case "all":
      default:
        // No additional filtering for "all" tab
        break;
    }

    setFilteredQuizzes(filtered);
  }, [activeTab, allQuizzes, completedQuizMap, favoriteIds]);

  // Pagination component (reusable for all tabs)
  const renderPagination = () => {
    if (!filteredQuizzes || filteredQuizzes.length === 0) return null;

    return (
      <div className="flex justify-center pt-4">
        <Pagination>
          <PaginationContent className="gap-1">
            <PaginationItem>
              <PaginationPrevious
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                className={`rounded-xl transition-all duration-300 hover:bg-violet-50 hover:text-violet-700 dark:hover:bg-violet-900/20 dark:hover:text-violet-400 ${
                  currentPage === 1 ? "pointer-events-none opacity-50" : "cursor-pointer"
                }`}
              />
            </PaginationItem>

            {/* Show page numbers */}
            {[...Array(Math.min(5, Math.ceil(100 / ITEMS_PER_PAGE)))].map((_, i) => {
              const pageNum = i + 1;
              const isActive = currentPage === pageNum;
              return (
                <PaginationItem key={pageNum}>
                  <PaginationLink
                    onClick={() => setCurrentPage(pageNum)}
                    isActive={isActive}
                    className={`cursor-pointer rounded-xl transition-all duration-300 ${
                      isActive
                        ? "bg-gradient-to-r from-violet-600 to-indigo-600 text-white border-0 shadow-md hover:from-violet-700 hover:to-indigo-700"
                        : "hover:bg-violet-50 hover:text-violet-700 dark:hover:bg-violet-900/20 dark:hover:text-violet-400"
                    }`}
                  >
                    {pageNum}
                  </PaginationLink>
                </PaginationItem>
              );
            })}

            <PaginationItem>
              <PaginationNext
                onClick={() => setCurrentPage((p) => p + 1)}
                className={`rounded-xl transition-all duration-300 hover:bg-violet-50 hover:text-violet-700 dark:hover:bg-violet-900/20 dark:hover:text-violet-400 ${
                  filteredQuizzes.length < ITEMS_PER_PAGE
                    ? "pointer-events-none opacity-50"
                    : "cursor-pointer"
                }`}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
    );
  };

  return (
    <div className="space-y-10 pb-10">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
          Explore Quizzes
        </h1>
        <p className="text-slate-500 dark:text-slate-400 mt-1">
          Challenge yourself with our diverse collection of quizzes
        </p>
      </div>

      {/* Search Bar */}
      <div className="bg-white/60 dark:bg-black/40 backdrop-blur-md p-3 rounded-2xl border border-white/20 dark:border-white/10 shadow-sm">
        <QuizSearch
          value={searchQuery}
          onChange={handleSearchChange}
          placeholder="Search quizzes by title or description..."
        />
      </div>

      <div className="container px-0 md:px-4">
        {/* Filters and Tabs */}
        <div className="flex flex-col gap-6 md:flex-row md:items-start md:justify-between mb-8">
          <Tabs value={activeTab} onValueChange={handleTabChange} className="w-full">
            <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between mb-6">
              <TabsList className="grid w-full max-w-4xl grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 bg-white/60 dark:bg-black/40 backdrop-blur-md border border-white/20 dark:border-white/10 p-1 rounded-xl shadow-sm">
                <TabsTrigger
                  value="all"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  All
                </TabsTrigger>
                <TabsTrigger
                  value="featured"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  Featured
                </TabsTrigger>
                <TabsTrigger
                  value="completed"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  Completed
                </TabsTrigger>
                <TabsTrigger
                  value="not-completed"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  Not Completed
                </TabsTrigger>
                <TabsTrigger
                  value="favorites"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  Favorites
                </TabsTrigger>
              </TabsList>

              <div className="flex-1 lg:max-w-2xl">
                <QuizFilters filters={filters} onFilterChange={handleFilterChange} />
              </div>
            </div>

            {/* All Tab */}
            <TabsContent value="all" className="space-y-8">
              <QuizList
                quizzes={filteredQuizzes}
                isLoading={allQuizzesLoading}
                error={allQuizzesError}
                completedQuizMap={completedQuizMap}
              />
              {renderPagination()}
            </TabsContent>

            {/* Featured Tab */}
            <TabsContent value="featured" className="space-y-8">
              <QuizList
                quizzes={filteredQuizzes}
                isLoading={allQuizzesLoading}
                error={allQuizzesError}
                completedQuizMap={completedQuizMap}
              />
              {renderPagination()}
            </TabsContent>

            {/* Completed Tab */}
            <TabsContent value="completed" className="space-y-8">
              <QuizList
                quizzes={filteredQuizzes}
                isLoading={allQuizzesLoading}
                error={allQuizzesError}
                completedQuizMap={completedQuizMap}
              />
              {renderPagination()}
            </TabsContent>

            {/* Not Completed Tab */}
            <TabsContent value="not-completed" className="space-y-8">
              <QuizList
                quizzes={filteredQuizzes}
                isLoading={allQuizzesLoading}
                error={allQuizzesError}
                completedQuizMap={completedQuizMap}
              />
              {renderPagination()}
            </TabsContent>

            {/* Favorites Tab */}
            <TabsContent value="favorites" className="space-y-8">
              <QuizList
                quizzes={filteredQuizzes}
                isLoading={allQuizzesLoading}
                error={allQuizzesError}
                completedQuizMap={completedQuizMap}
              />
              {renderPagination()}
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
