"use client";

import { useState } from "react";
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
import { MessageSquare, MessagesSquare, Users, TrendingUp, Search, Filter } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { DiscussionFilters } from "@/lib/api/discussions";
import type { DiscussionSort } from "@/constants";
import { GlassCard } from "@/components/common/GlassCard";
import { StatsCard } from "@/components/common/StatsCard";
import { StatsGrid } from "@/components/common/StatsGrid";

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

  const { data: stats, isLoading: statsLoading } = useDiscussionStats();
  const { data: quizzes = [] } = useQuizzes();

  const handleQuizFilter = (quizId: string) => {
    if (quizId === "all") {
      setFilters({ ...filters, quiz_id: undefined });
    } else {
      setFilters({ ...filters, quiz_id: quizId });
    }
  };

  const handleSortChange = (sort: DiscussionSort) => {
    setFilters({ ...filters, sort });
  };

  return (
    <div className="space-y-10 pb-10">
      {/* Header with action button */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
            Discussions
          </h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">
            Share insights, ask questions, and engage with the community
          </p>
        </div>
        <CreateDiscussionDialog />
      </div>

      {/* Stats Cards */}
      <StatsGrid columns={3}>
        <StatsCard
          title="Discussions"
          value={stats?.total_discussions ?? 0}
          description="Total discussions"
          icon={MessagesSquare}
          color="blue"
          loading={statsLoading}
        />
        <StatsCard
          title="Replies"
          value={stats?.total_replies ?? 0}
          description="Community responses"
          icon={MessageSquare}
          color="yellow"
          loading={statsLoading}
        />
        <StatsCard
          title="Your Posts"
          value={stats?.user_discussions ?? 0}
          description="Your contributions"
          icon={Users}
          color="green"
          loading={statsLoading}
        />
      </StatsGrid>

      {/* Filters */}
      <div className="container px-0 md:px-4">
        <GlassCard padding="none" rounded="2xl" className="mb-8">
          <div className="p-6 border-b border-white/10">
            <div className="flex items-center gap-2">
              <span className="bg-gradient-to-br from-violet-500 to-purple-600 text-white p-1.5 rounded-lg shadow-sm">
                <Filter className="h-4 w-4" />
              </span>
              <h2 className="text-lg font-bold tracking-tight text-slate-800 dark:text-slate-100">
                Filters
              </h2>
            </div>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
              Search and filter discussions by quiz or sort by popularity
            </p>
          </div>
          <div className="p-6">
            {/* Search Input */}
            <div className="mb-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                <Input
                  placeholder="Search discussions by title or content..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 bg-white/50 dark:bg-white/10 border-white/20 dark:border-white/10 backdrop-blur-sm"
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
                  <SelectTrigger className="bg-white/50 dark:bg-white/10 border-white/20 dark:border-white/10 backdrop-blur-sm">
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
                  className={`gap-2 ${
                    filters.sort === "recent"
                      ? "bg-gradient-to-r from-violet-600 to-indigo-600 text-white border-0 shadow-md"
                      : "bg-white/50 dark:bg-white/10 border-white/20 dark:border-white/10"
                  }`}
                >
                  <TrendingUp className="h-4 w-4" />
                  Recent
                </Button>
                <Button
                  variant={filters.sort === "popular" ? "default" : "outline"}
                  onClick={() => handleSortChange("popular")}
                  className={`gap-2 ${
                    filters.sort === "popular"
                      ? "bg-gradient-to-r from-violet-600 to-indigo-600 text-white border-0 shadow-md"
                      : "bg-white/50 dark:bg-white/10 border-white/20 dark:border-white/10"
                  }`}
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
                  <Badge variant="secondary" className="bg-white/50 dark:bg-white/10">
                    Search: &quot;{searchQuery}&quot;
                  </Badge>
                )}
                {filters.sort && (
                  <Badge variant="secondary" className="bg-white/50 dark:bg-white/10">
                    Sort: {filters.sort === "recent" ? "Recent" : "Popular"}
                  </Badge>
                )}
                {filters.quiz_id && (
                  <Badge variant="secondary" className="bg-white/50 dark:bg-white/10">
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
                  className="text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Clear All
                </Button>
              </div>
            )}
          </div>
        </GlassCard>

        {/* Discussions List */}
        <GlassCard padding="none" rounded="2xl">
          <div className="p-6 border-b border-white/10">
            <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
              <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
                <MessagesSquare className="h-4 w-4" />
              </span>
              {filters.quiz_id
                ? `Discussions about ${
                    quizzes.find((q) => q.id === filters.quiz_id)?.title
                  }`
                : "All Discussions"}
            </h2>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
              {isLoading
                ? "Loading discussions..."
                : `${discussions.length} discussion${
                    discussions.length !== 1 ? "s" : ""
                  } found`}
            </p>
          </div>
          <div className="p-6">
            <DiscussionList discussions={discussions} isLoading={isLoading} />
          </div>
        </GlassCard>
      </div>
    </div>
  );
}
