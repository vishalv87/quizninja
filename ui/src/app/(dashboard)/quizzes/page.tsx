"use client";

import { useState } from "react";
import { QuizList } from "@/components/quiz/QuizList";
import { QuizFilters, QuizFilterValues } from "@/components/quiz/QuizFilters";
import { QuizSearch } from "@/components/quiz/QuizSearch";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useFeaturedQuizzes } from "@/hooks/useFeaturedQuizzes";
import { useDebounce } from "@/hooks/useDebounce";
import { Separator } from "@/components/ui/separator";
import { Card } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
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
  });
  const [currentPage, setCurrentPage] = useState(1);

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

  return (
    <div className="container py-8 space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Quizzes</h1>
        <p className="text-muted-foreground mt-2">
          Explore quizzes and test your knowledge
        </p>
      </div>

      <Separator />

      {/* Search Bar */}
      <QuizSearch
        value={searchQuery}
        onChange={handleSearchChange}
        placeholder="Search quizzes by title or description..."
      />

      {/* Filters */}
      <QuizFilters filters={filters} onFilterChange={handleFilterChange} />

      {/* Tabs for All Quizzes and Featured */}
      <Tabs defaultValue="all" className="w-full">
        <TabsList className="grid w-full max-w-md grid-cols-2">
          <TabsTrigger value="all">All Quizzes</TabsTrigger>
          <TabsTrigger value="featured">Featured</TabsTrigger>
        </TabsList>

        {/* All Quizzes Tab */}
        <TabsContent value="all" className="mt-6 space-y-6">
          <QuizList
            quizzes={allQuizzes || []}
            isLoading={allQuizzesLoading}
            error={allQuizzesError}
          />

          {/* Pagination for All Quizzes */}
          {allQuizzes && allQuizzes.length > 0 && (
            <div className="flex justify-center">
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
                        allQuizzes.length < ITEMS_PER_PAGE
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
        <TabsContent value="featured" className="mt-6">
          <QuizList
            quizzes={featuredQuizzes || []}
            isLoading={featuredLoading}
            error={featuredError}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}