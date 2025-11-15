"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { QuizList } from "@/components/quiz/QuizList";
import { QuizSearch } from "@/components/quiz/QuizSearch";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useDebounce } from "@/hooks/useDebounce";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { ArrowLeft } from "lucide-react";

const ITEMS_PER_PAGE = 12;

export default function CategoryQuizzesPage() {
  const params = useParams();
  const router = useRouter();
  const categoryId = params.categoryId as string;

  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);

  // Debounce search query
  const debouncedSearch = useDebounce(searchQuery, 300);

  // Reset to first page when search changes
  const handleSearchChange = (query: string) => {
    setSearchQuery(query);
    setCurrentPage(1);
  };

  // Build filter object for API
  const apiFilters = {
    category: categoryId,
    search: debouncedSearch || undefined,
    limit: ITEMS_PER_PAGE,
    offset: (currentPage - 1) * ITEMS_PER_PAGE,
  };

  // Fetch quizzes for this category
  const {
    data: quizzes,
    isLoading,
    error,
  } = useQuizzes(apiFilters);

  // Format category name for display (capitalize and replace underscores)
  const categoryName = categoryId
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");

  return (
    <div className="container py-8 space-y-6">
      {/* Back Button */}
      <Button variant="ghost" onClick={() => router.push("/quizzes")}>
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to All Quizzes
      </Button>

      {/* Page Header */}
      <div>
        <div className="flex items-center gap-3">
          <h1 className="text-3xl font-bold tracking-tight">{categoryName}</h1>
          <Badge variant="secondary" className="text-sm">
            {categoryId}
          </Badge>
        </div>
        <p className="text-muted-foreground mt-2">
          Browse all quizzes in the {categoryName} category
        </p>
      </div>

      <Separator />

      {/* Search Bar */}
      <QuizSearch
        value={searchQuery}
        onChange={handleSearchChange}
        placeholder={`Search ${categoryName} quizzes...`}
      />

      {/* Quiz List */}
      <div className="space-y-6">
        <QuizList quizzes={quizzes || []} isLoading={isLoading} error={error} />

        {/* Pagination */}
        {quizzes && quizzes.length > 0 && (
          <div className="flex justify-center">
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                    className={
                      currentPage === 1
                        ? "pointer-events-none opacity-50"
                        : "cursor-pointer"
                    }
                  />
                </PaginationItem>

                {/* Show page numbers */}
                {[...Array(Math.min(5, Math.ceil(100 / ITEMS_PER_PAGE)))].map(
                  (_, i) => {
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
                  }
                )}

                <PaginationItem>
                  <PaginationNext
                    onClick={() => setCurrentPage((p) => p + 1)}
                    className={
                      quizzes.length < ITEMS_PER_PAGE
                        ? "pointer-events-none opacity-50"
                        : "cursor-pointer"
                    }
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        )}
      </div>
    </div>
  );
}
