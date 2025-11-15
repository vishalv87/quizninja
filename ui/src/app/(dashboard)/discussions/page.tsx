"use client";

import { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { DiscussionList } from "@/components/discussion/DiscussionList";
import { CreateDiscussionDialog } from "@/components/discussion/CreateDiscussionDialog";
import { useDiscussions, useDiscussionStats } from "@/hooks/useDiscussions";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useDebounce } from "@/hooks/useDebounce";
import { MessageSquare, MessagesSquare, Users, TrendingUp, Search } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { DiscussionFilters } from "@/lib/api/discussions";

export default function DiscussionsPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [filters, setFilters] = useState<DiscussionFilters>({
    sort: "recent",
  });

  // Debounce search query
  const debouncedSearch = useDebounce(searchQuery, 300);

  // Combine filters with debounced search
  const apiFilters = {
    ...filters,
    search: debouncedSearch || undefined,
  };

  const {
    data: discussions = [],
    isLoading,
    error,
  } = useDiscussions(apiFilters);

  const { data: stats } = useDiscussionStats();
  const { data: quizzes = [] } = useQuizzes();

  const handleQuizFilter = (quizId: string) => {
    if (quizId === "all") {
      setFilters({ ...filters, quiz_id: undefined });
    } else {
      setFilters({ ...filters, quiz_id: quizId });
    }
  };

  const handleSortChange = (sort: "recent" | "popular") => {
    setFilters({ ...filters, sort });
  };

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h1 className="text-4xl font-bold mb-2">Discussions</h1>
            <p className="text-muted-foreground">
              Share insights, ask questions, and engage with the community
            </p>
          </div>
          <CreateDiscussionDialog />
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid gap-4 md:grid-cols-3 mt-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Total Discussions
                </CardTitle>
                <MessagesSquare className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {stats.total_discussions}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Total Replies
                </CardTitle>
                <MessageSquare className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_replies}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Your Discussions
                </CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {stats.user_discussions}
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>

      {/* Filters */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Filters</CardTitle>
          <CardDescription>
            Search and filter discussions by quiz or sort by popularity
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Search Input */}
          <div className="mb-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search discussions by title or content..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          <div className="flex flex-col sm:flex-row gap-4">
            {/* Quiz Filter */}
            <div className="flex-1">
              <Select
                value={filters.quiz_id || "all"}
                onValueChange={handleQuizFilter}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Filter by quiz" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Quizzes</SelectItem>
                  {quizzes.map((quiz) => (
                    <SelectItem key={quiz.id} value={quiz.id}>
                      {quiz.title}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Sort Options */}
            <div className="flex gap-2">
              <Button
                variant={filters.sort === "recent" ? "default" : "outline"}
                onClick={() => handleSortChange("recent")}
                className="gap-2"
              >
                <TrendingUp className="h-4 w-4" />
                Recent
              </Button>
              <Button
                variant={filters.sort === "popular" ? "default" : "outline"}
                onClick={() => handleSortChange("popular")}
                className="gap-2"
              >
                <MessageSquare className="h-4 w-4" />
                Popular
              </Button>
            </div>
          </div>

          {/* Active Filters Display */}
          {(filters.quiz_id || filters.sort || searchQuery) && (
            <div className="flex gap-2 mt-4 flex-wrap">
              {searchQuery && (
                <Badge variant="secondary">
                  Search: &quot;{searchQuery}&quot;
                </Badge>
              )}
              {filters.sort && (
                <Badge variant="secondary">
                  Sort: {filters.sort === "recent" ? "Recent" : "Popular"}
                </Badge>
              )}
              {filters.quiz_id && (
                <Badge variant="secondary">
                  Quiz:{" "}
                  {quizzes.find((q) => q.id === filters.quiz_id)?.title ||
                    "Unknown"}
                </Badge>
              )}
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setFilters({ sort: "recent" });
                  setSearchQuery("");
                }}
              >
                Clear All
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Discussions List */}
      <Card>
        <CardHeader>
          <CardTitle>
            {filters.quiz_id
              ? `Discussions about ${
                  quizzes.find((q) => q.id === filters.quiz_id)?.title
                }`
              : "All Discussions"}
          </CardTitle>
          <CardDescription>
            {isLoading
              ? "Loading discussions..."
              : `${discussions.length} discussion${
                  discussions.length !== 1 ? "s" : ""
                } found`}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DiscussionList discussions={discussions} isLoading={isLoading} />
        </CardContent>
      </Card>
    </div>
  );
}
