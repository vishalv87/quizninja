"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent } from "@/components/ui/card";
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
    <div className="w-full max-w-4xl mx-auto space-y-4">
      {/* Search Input */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
        <Input
          placeholder="Search for quizzes, users, or discussions..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-11 h-12 text-lg"
          autoFocus
        />
        {isLoading && (
          <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground animate-spin" />
        )}
      </div>

      {/* Search Results */}
      {showResults ? (
        <Tabs defaultValue="quizzes" className="w-full">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="quizzes" className="gap-2">
              <BookOpen className="h-4 w-4" />
              Quizzes
              {quizResults.length > 0 && (
                <Badge variant="secondary" className="ml-1">
                  {quizResults.length}
                </Badge>
              )}
            </TabsTrigger>
            <TabsTrigger value="users" className="gap-2">
              <Users className="h-4 w-4" />
              Users
              {userResults.length > 0 && (
                <Badge variant="secondary" className="ml-1">
                  {userResults.length}
                </Badge>
              )}
            </TabsTrigger>
            <TabsTrigger value="discussions" className="gap-2">
              <MessageSquare className="h-4 w-4" />
              Discussions
              {discussionResults.length > 0 && (
                <Badge variant="secondary" className="ml-1">
                  {discussionResults.length}
                </Badge>
              )}
            </TabsTrigger>
          </TabsList>

          {/* Quizzes Tab */}
          <TabsContent value="quizzes" className="mt-6">
            {quizzesLoading ? (
              <Card>
                <CardContent className="p-8 flex items-center justify-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin" />
                  <p className="text-muted-foreground">Searching quizzes...</p>
                </CardContent>
              </Card>
            ) : quizResults.length > 0 ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground">
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
                      <Button variant="outline" onClick={onClose}>
                        View All Results
                        <ExternalLink className="ml-2 h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                )}
              </div>
            ) : (
              <Card>
                <CardContent className="p-8 text-center">
                  <BookOpen className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
                  <p className="text-muted-foreground">
                    No quizzes found matching &quot;{debouncedQuery}&quot;
                  </p>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {/* Users Tab */}
          <TabsContent value="users" className="mt-6">
            {usersLoading ? (
              <Card>
                <CardContent className="p-8 flex items-center justify-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin" />
                  <p className="text-muted-foreground">Searching users...</p>
                </CardContent>
              </Card>
            ) : userResults.length > 0 ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground">
                  Found {userResults.length} {userResults.length === 1 ? "user" : "users"}
                </p>
                <div className="space-y-2">
                  {userResults.map((user: any) => (
                    <Card key={user.id} className="hover:shadow-md transition-shadow">
                      <CardContent className="p-4">
                        <div className="flex items-center gap-3">
                          <Avatar className="h-12 w-12">
                            <AvatarImage src={user.avatar_url} alt={user.full_name} />
                            <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                              {getInitials(user.full_name)}
                            </AvatarFallback>
                          </Avatar>
                          <div className="flex-1 min-w-0">
                            <h4 className="font-semibold truncate">{user.full_name}</h4>
                            <p className="text-sm text-muted-foreground truncate">
                              {user.email}
                            </p>
                          </div>
                          <Link href="/friends">
                            <Button variant="outline" size="sm" onClick={onClose}>
                              View Profile
                            </Button>
                          </Link>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
            ) : (
              <Card>
                <CardContent className="p-8 text-center">
                  <Users className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
                  <p className="text-muted-foreground">
                    No users found matching &quot;{debouncedQuery}&quot;
                  </p>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {/* Discussions Tab */}
          <TabsContent value="discussions" className="mt-6">
            {discussionsLoading ? (
              <Card>
                <CardContent className="p-8 flex items-center justify-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin" />
                  <p className="text-muted-foreground">Searching discussions...</p>
                </CardContent>
              </Card>
            ) : discussionResults.length > 0 ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground">
                  Found {discussionResults.length} {discussionResults.length === 1 ? "discussion" : "discussions"}
                </p>
                <div className="space-y-2">
                  {discussionResults.map((discussion: Discussion) => (
                    <Card key={discussion.id} className="hover:shadow-md transition-shadow">
                      <CardContent className="p-4">
                        <Link href={`/discussions/${discussion.id}`} onClick={onClose}>
                          <div className="space-y-2">
                            <div className="flex items-start justify-between gap-2">
                              <h4 className="font-semibold line-clamp-2 flex-1">
                                {discussion.title}
                              </h4>
                              <Badge variant="outline" className="shrink-0">
                                <MessageSquare className="h-3 w-3 mr-1" />
                                {discussion.replies_count || 0}
                              </Badge>
                            </div>
                            {discussion.content && (
                              <p className="text-sm text-muted-foreground line-clamp-2">
                                {discussion.content}
                              </p>
                            )}
                            <div className="flex items-center gap-3 text-xs text-muted-foreground">
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
                      </CardContent>
                    </Card>
                  ))}
                </div>
                {discussionResults.length >= 6 && (
                  <div className="text-center pt-4">
                    <Link href={`/discussions?search=${encodeURIComponent(debouncedQuery)}`}>
                      <Button variant="outline" onClick={onClose}>
                        View All Results
                        <ExternalLink className="ml-2 h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                )}
              </div>
            ) : (
              <Card>
                <CardContent className="p-8 text-center">
                  <MessageSquare className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
                  <p className="text-muted-foreground">
                    No discussions found matching &quot;{debouncedQuery}&quot;
                  </p>
                </CardContent>
              </Card>
            )}
          </TabsContent>
        </Tabs>
      ) : (
        <Card>
          <CardContent className="p-12 text-center">
            <Search className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">Search QuizNinja</h3>
            <p className="text-muted-foreground">
              Find quizzes, connect with friends, or browse discussions
            </p>
            <p className="text-sm text-muted-foreground mt-2">
              Start typing to search (minimum 2 characters)
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
