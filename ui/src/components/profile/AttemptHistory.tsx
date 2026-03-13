"use client";

import { useState } from "react";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { Clock, Award, CheckCircle2, XCircle, Filter, Calendar, Trophy } from "lucide-react";
import { useUserAttempts } from "@/hooks/useUserStats";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";

export function AttemptHistory() {
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [limit, setLimit] = useState<number>(10);

  // Build filters object
  const filters: any = { limit };
  if (statusFilter !== "all") {
    filters.status = statusFilter;
  }

  const { data: attemptsData, isLoading } = useUserAttempts(filters);
  // Note: Each user has only ONE attempt per quiz
  const attempts = attemptsData || [];

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "completed":
        return (
          <Badge variant="default" className="bg-green-500">
            <CheckCircle2 className="mr-1 h-3 w-3" />
            Completed
          </Badge>
        );
      case "in_progress":
        return (
          <Badge variant="secondary">
            <Clock className="mr-1 h-3 w-3" />
            In Progress
          </Badge>
        );
      case "abandoned":
        return (
          <Badge variant="destructive">
            <XCircle className="mr-1 h-3 w-3" />
            Abandoned
          </Badge>
        );
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const formatTime = (seconds?: number) => {
    if (!seconds) return "N/A";
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes}m ${secs}s`;
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Attempt History</CardTitle>
          <CardDescription>All your quiz attempts</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="flex items-center space-x-4 p-4 rounded-lg border">
              <Skeleton className="h-12 w-12 rounded-full" />
              <div className="space-y-2 flex-1">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-3 w-1/2" />
              </div>
              <Skeleton className="h-6 w-20" />
            </div>
          ))}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between flex-wrap gap-4">
          <div>
            <CardTitle>Attempt History</CardTitle>
            <CardDescription>All your quiz attempts</CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[150px]">
                <Filter className="mr-2 h-4 w-4" />
                <SelectValue placeholder="Filter" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="in_progress">In Progress</SelectItem>
                <SelectItem value="abandoned">Abandoned</SelectItem>
              </SelectContent>
            </Select>
            <Select value={limit.toString()} onValueChange={(val) => setLimit(parseInt(val))}>
              <SelectTrigger className="w-[120px]">
                <SelectValue placeholder="Show" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="10">Show 10</SelectItem>
                <SelectItem value="25">Show 25</SelectItem>
                <SelectItem value="50">Show 50</SelectItem>
                <SelectItem value="100">Show 100</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {!attempts || attempts.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12">
            <Trophy className="h-12 w-12 text-muted-foreground mb-4" />
            <p className="text-muted-foreground text-center mb-2">
              No attempts found
            </p>
            <p className="text-sm text-muted-foreground text-center">
              {statusFilter !== "all"
                ? `Try adjusting your filters to see more results`
                : `Start taking quizzes to build your attempt history!`}
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            {attempts.map((attempt: any) => {
              const correctAnswers = attempt.score ?? 0;
              const totalQuestions = attempt.total_points ?? 0;
              const percentage = totalQuestions > 0
                ? Math.round((correctAnswers / totalQuestions) * 100)
                : 0;
              const passed = percentage >= 60;

              return (
                <div
                  key={attempt.id}
                  className="flex items-center justify-between p-4 rounded-lg border hover:bg-accent transition-colors"
                >
                  <div className="flex items-center space-x-4 flex-1">
                    <div
                      className={`flex items-center justify-center h-12 w-12 rounded-full flex-shrink-0 ${
                        attempt.status === "completed"
                          ? passed
                            ? "bg-green-100 text-green-600 dark:bg-green-900/20 dark:text-green-400"
                            : "bg-yellow-100 text-yellow-600 dark:bg-yellow-900/20 dark:text-yellow-400"
                          : attempt.status === "in_progress"
                          ? "bg-blue-100 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
                          : "bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400"
                      }`}
                    >
                      <Award className="h-6 w-6" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <p className="font-medium truncate">
                          {attempt.quiz_title || "Quiz Attempt"}
                        </p>
                        {attempt.category && (
                          <Badge variant="outline" className="text-xs">
                            {attempt.category}
                          </Badge>
                        )}
                      </div>
                      <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground flex-wrap">
                        <span className="flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {attempt.started_at &&
                            formatDistanceToNow(new Date(attempt.started_at), {
                              addSuffix: true,
                            })}
                        </span>
                        {attempt.time_spent && (
                          <span className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            {formatTime(attempt.time_spent)}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-4 flex-shrink-0 ml-4">
                    {attempt.status === "completed" && (
                      <div className="text-right">
                        <p className="text-lg font-bold">{percentage}%</p>
                        <p className="text-xs text-muted-foreground">
                          {correctAnswers}/{totalQuestions} correct
                        </p>
                      </div>
                    )}
                    {getStatusBadge(attempt.status)}
                  </div>
                </div>
              );
            })}

            {/* Show more button if there might be more results */}
            {attempts.length >= limit && (
              <div className="flex justify-center pt-4">
                <Button
                  variant="outline"
                  onClick={() => setLimit(limit + 10)}
                >
                  Load More
                </Button>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
