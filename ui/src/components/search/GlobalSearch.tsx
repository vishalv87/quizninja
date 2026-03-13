"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Search,
  Loader2,
  BookOpen,
  Users,
  MessageSquare,
  ExternalLink,
} from "lucide-react";
import { useDebounce } from "@/hooks/useDebounce";
import { useSearchUsers } from "@/hooks/useSearchUsers";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useDiscussions } from "@/hooks/useDiscussions";
import { QuizCard } from "@/components/quiz/QuizCard";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import Link from "next/link";
import type { Quiz } from "@/types/quiz";
import type { Discussion } from "@/types/discussion";

interface GlobalSearchProps {
  onClose?: () => void;
}

export function GlobalSearch({ onClose }: GlobalSearchProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const debouncedQuery = useDebounce(searchQuery, 400);

  // Search across different entities
  const { data: quizResults = [], isLoading: quizzesLoading } = useQuizzes({
    search: debouncedQuery || undefined,
    limit: 6,
  });

  const { data: userResults = [], isLoading: usersLoading } =
    useSearchUsers(debouncedQuery);

  const { data: discussionResults = [], isLoading: discussionsLoading } = useDiscussions({
    search: debouncedQuery || undefined,
    limit: 6,
  });

  const showResults = debouncedQuery.length >= 2;
  const isLoading = quizzesLoading || usersLoading || discussionsLoading;

  // Generate initials for avatar
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  return (
    <div className="w-full max-w-4xl mx-auto overflow-hidden rounded-2xl shadow-2xl">
      {/* Gradient Header with Search Input */}
      <div className="relative overflow-hidden bg-gradient-to-br from-violet-600 via-indigo-600 to-purple-700 p-6 text-white">
        {/* Decorative blur elements */}
        <div className="absolute right-0 top-0 -mt-10 -mr-10 h-40 w-40 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute bottom-0 left-10 h-20 w-20 rounded-full bg-purple-400/20 blur-2xl" />
        <div className="absolute top-1/2 right-20 h-16 w-16 rounded-full bg-indigo-400/20 blur-2xl" />

        <div className="relative z-10">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Search className="h-5 w-5" />
            Search QuizNinja
          </h2>
          <div className="relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-indigo-200" />
            <Input
              placeholder="Search for quizzes, users, or discussions..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-12 h-14 text-lg bg-white/10 border-white/20 text-white placeholder:text-indigo-200 rounded-xl backdrop-blur-sm focus:bg-white/20 focus:border-white/40 focus-visible:ring-0 focus-visible:ring-offset-0"
              autoFocus
            />
            {isLoading && (
              <Loader2 className="absolute right-4 top-1/2 -translate-y-1/2 h-5 w-5 text-indigo-200 animate-spin" />
            )}
          </div>
        </div>
      </div>

      {/* Glassmorphism Body */}
      <div className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-xl border-t-0 border border-white/20 dark:border-white/10 p-6">
        {showResults ? (
          <Tabs defaultValue="quizzes" className="w-full">
            <TabsList className="grid w-full grid-cols-3 bg-slate-100/80 dark:bg-black/40 backdrop-blur-md border border-slate-200/50 dark:border-white/10 p-1.5 rounded-xl mb-6">
              <TabsTrigger
                value="quizzes"
                className="gap-2 rounded-lg data-[state=active]:bg-white dark:data-[state=active]:bg-slate-800 data-[state=active]:shadow-md transition-all duration-300"
              >
                <div className="p-1 rounded-md bg-blue-500/10">
                  <BookOpen className="h-4 w-4 text-blue-500" />
                </div>
                Quizzes
                {quizResults.length > 0 && (
                  <Badge variant="secondary" className="ml-1 bg-blue-100 text-blue-700 dark:bg-blue-900/50 dark:text-blue-300">
                    {quizResults.length}
                  </Badge>
                )}
              </TabsTrigger>
              <TabsTrigger
                value="users"
                className="gap-2 rounded-lg data-[state=active]:bg-white dark:data-[state=active]:bg-slate-800 data-[state=active]:shadow-md transition-all duration-300"
              >
                <div className="p-1 rounded-md bg-purple-500/10">
                  <Users className="h-4 w-4 text-purple-500" />
                </div>
                Users
                {userResults.length > 0 && (
                  <Badge variant="secondary" className="ml-1 bg-purple-100 text-purple-700 dark:bg-purple-900/50 dark:text-purple-300">
                    {userResults.length}
                  </Badge>
                )}
              </TabsTrigger>
              <TabsTrigger
                value="discussions"
                className="gap-2 rounded-lg data-[state=active]:bg-white dark:data-[state=active]:bg-slate-800 data-[state=active]:shadow-md transition-all duration-300"
              >
                <div className="p-1 rounded-md bg-amber-500/10">
                  <MessageSquare className="h-4 w-4 text-amber-500" />
                </div>
                Discussions
                {discussionResults.length > 0 && (
                  <Badge variant="secondary" className="ml-1 bg-amber-100 text-amber-700 dark:bg-amber-900/50 dark:text-amber-300">
                    {discussionResults.length}
                  </Badge>
                )}
              </TabsTrigger>
            </TabsList>

            {/* Quizzes Tab */}
            <TabsContent value="quizzes" className="mt-0">
              {quizzesLoading ? (
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-8 flex items-center justify-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin text-violet-600" />
                  <p className="text-slate-500 dark:text-slate-400">Searching quizzes...</p>
                </div>
              ) : quizResults.length > 0 ? (
                <div className="space-y-4">
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Found {quizResults.length} {quizResults.length === 1 ? "quiz" : "quizzes"}
                  </p>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {quizResults.map((quiz: Quiz) => (
                      <QuizCard key={quiz.id} quiz={quiz} />
                    ))}
                  </div>
                  {quizResults.length >= 6 && (
                    <div className="text-center pt-4">
                      <Link href={`/quizzes?search=${encodeURIComponent(debouncedQuery)}`}>
                        <Button
                          variant="outline"
                          onClick={onClose}
                          className="bg-gradient-to-r from-violet-600 to-indigo-600 text-white border-0 hover:from-violet-700 hover:to-indigo-700"
                        >
                          View All Results
                          <ExternalLink className="ml-2 h-4 w-4" />
                        </Button>
                      </Link>
                    </div>
                  )}
                </div>
              ) : (
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-8 text-center">
                  <div className="p-3 rounded-xl bg-blue-500/10 w-fit mx-auto mb-3">
                    <BookOpen className="h-8 w-8 text-blue-500" />
                  </div>
                  <p className="text-slate-500 dark:text-slate-400">
                    No quizzes found matching &quot;{debouncedQuery}&quot;
                  </p>
                </div>
              )}
            </TabsContent>

            {/* Users Tab */}
            <TabsContent value="users" className="mt-0">
              {usersLoading ? (
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-8 flex items-center justify-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin text-violet-600" />
                  <p className="text-slate-500 dark:text-slate-400">Searching users...</p>
                </div>
              ) : userResults.length > 0 ? (
                <div className="space-y-4">
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Found {userResults.length} {userResults.length === 1 ? "user" : "users"}
                  </p>
                  <div className="space-y-2">
                    {userResults.map((user: any) => (
                      <div
                        key={user.id}
                        className="group bg-slate-50 dark:bg-slate-800/50 border border-slate-200/50 dark:border-slate-700/50 rounded-xl p-4 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300"
                      >
                        <div className="flex items-center gap-3">
                          <Avatar className="h-12 w-12 ring-2 ring-purple-500/20">
                            <AvatarImage src={user.avatar_url} alt={user.full_name} />
                            <AvatarFallback className="bg-gradient-to-br from-purple-500 to-pink-500 text-white font-semibold">
                              {getInitials(user.full_name)}
                            </AvatarFallback>
                          </Avatar>
                          <div className="flex-1 min-w-0">
                            <h4 className="font-semibold truncate text-slate-800 dark:text-slate-100">{user.full_name}</h4>
                            <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                              {user.email}
                            </p>
                          </div>
                          <Link href="/friends">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={onClose}
                              className="bg-purple-50 dark:bg-purple-900/30 border-purple-200 dark:border-purple-700 text-purple-700 dark:text-purple-300 hover:bg-purple-100 dark:hover:bg-purple-900/50"
                            >
                              View Profile
                            </Button>
                          </Link>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : (
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-8 text-center">
                  <div className="p-3 rounded-xl bg-purple-500/10 w-fit mx-auto mb-3">
                    <Users className="h-8 w-8 text-purple-500" />
                  </div>
                  <p className="text-slate-500 dark:text-slate-400">
                    No users found matching &quot;{debouncedQuery}&quot;
                  </p>
                </div>
              )}
            </TabsContent>

            {/* Discussions Tab */}
            <TabsContent value="discussions" className="mt-0">
              {discussionsLoading ? (
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-8 flex items-center justify-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin text-violet-600" />
                  <p className="text-slate-500 dark:text-slate-400">Searching discussions...</p>
                </div>
              ) : discussionResults.length > 0 ? (
                <div className="space-y-4">
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Found {discussionResults.length} {discussionResults.length === 1 ? "discussion" : "discussions"}
                  </p>
                  <div className="space-y-2">
                    {discussionResults.map((discussion: Discussion) => (
                      <div
                        key={discussion.id}
                        className="group bg-slate-50 dark:bg-slate-800/50 border border-slate-200/50 dark:border-slate-700/50 rounded-xl p-4 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300"
                      >
                        <Link href={`/discussions/${discussion.id}`} onClick={onClose}>
                          <div className="space-y-2">
                            <div className="flex items-start justify-between gap-2">
                              <h4 className="font-semibold line-clamp-2 flex-1 text-slate-800 dark:text-slate-100">
                                {discussion.title}
                              </h4>
                              <Badge variant="outline" className="shrink-0 bg-amber-50 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 border-amber-200 dark:border-amber-700">
                                <MessageSquare className="h-3 w-3 mr-1" />
                                {discussion.replies_count || 0}
                              </Badge>
                            </div>
                            {discussion.content && (
                              <p className="text-sm text-slate-500 dark:text-slate-400 line-clamp-2">
                                {discussion.content}
                              </p>
                            )}
                            <div className="flex items-center gap-3 text-xs text-slate-400 dark:text-slate-500">
                              <span>By {discussion.user?.name || "Unknown"}</span>
                              <span>•</span>
                              <span>{discussion.likes_count || 0} likes</span>
                              {discussion.quiz?.title && (
                                <>
                                  <span>•</span>
                                  <span className="truncate">Quiz: {discussion.quiz.title}</span>
                                </>
                              )}
                            </div>
                          </div>
                        </Link>
                      </div>
                    ))}
                  </div>
                  {discussionResults.length >= 6 && (
                    <div className="text-center pt-4">
                      <Link href={`/discussions?search=${encodeURIComponent(debouncedQuery)}`}>
                        <Button
                          variant="outline"
                          onClick={onClose}
                          className="bg-gradient-to-r from-violet-600 to-indigo-600 text-white border-0 hover:from-violet-700 hover:to-indigo-700"
                        >
                          View All Results
                          <ExternalLink className="ml-2 h-4 w-4" />
                        </Button>
                      </Link>
                    </div>
                  )}
                </div>
              ) : (
                <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-8 text-center">
                  <div className="p-3 rounded-xl bg-amber-500/10 w-fit mx-auto mb-3">
                    <MessageSquare className="h-8 w-8 text-amber-500" />
                  </div>
                  <p className="text-slate-500 dark:text-slate-400">
                    No discussions found matching &quot;{debouncedQuery}&quot;
                  </p>
                </div>
              )}
            </TabsContent>
          </Tabs>
        ) : (
          /* Initial Empty State */
          <div className="py-12 text-center">
            <div className="mx-auto w-20 h-20 rounded-2xl bg-gradient-to-br from-violet-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-indigo-500/30 mb-6">
              <Search className="h-10 w-10 text-white" />
            </div>
            <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 mb-2">
              What are you looking for?
            </h3>
            <p className="text-slate-500 dark:text-slate-400 mb-1">
              Find quizzes, connect with friends, or browse discussions
            </p>
            <p className="text-sm text-slate-400 dark:text-slate-500">
              Start typing to search (minimum 2 characters)
            </p>

            {/* Quick Links */}
            <div className="flex justify-center gap-3 mt-8">
              <Link href="/quizzes" onClick={onClose}>
                <Button variant="outline" className="gap-2 bg-blue-50 dark:bg-blue-900/30 border-blue-200 dark:border-blue-700 text-blue-700 dark:text-blue-300 hover:bg-blue-100 dark:hover:bg-blue-900/50">
                  <BookOpen className="h-4 w-4" />
                  Browse Quizzes
                </Button>
              </Link>
              <Link href="/friends" onClick={onClose}>
                <Button variant="outline" className="gap-2 bg-purple-50 dark:bg-purple-900/30 border-purple-200 dark:border-purple-700 text-purple-700 dark:text-purple-300 hover:bg-purple-100 dark:hover:bg-purple-900/50">
                  <Users className="h-4 w-4" />
                  Find Friends
                </Button>
              </Link>
              <Link href="/discussions" onClick={onClose}>
                <Button variant="outline" className="gap-2 bg-amber-50 dark:bg-amber-900/30 border-amber-200 dark:border-amber-700 text-amber-700 dark:text-amber-300 hover:bg-amber-100 dark:hover:bg-amber-900/50">
                  <MessageSquare className="h-4 w-4" />
                  Discussions
                </Button>
              </Link>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
