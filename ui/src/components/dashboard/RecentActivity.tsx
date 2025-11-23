"use client";

import Link from "next/link";
import { ArrowRight, Clock, Award, CheckCircle2, XCircle } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { useUserAttempts } from "@/hooks/useUserStats";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";

export function RecentActivity() {
  const { data: attemptsData, isLoading } = useUserAttempts({ limit: 5 });

  // Note: Each user has only ONE attempt per quiz
  const attempts = attemptsData || [];

  if (isLoading) {
    return (
      <Card className="border-none shadow-sm bg-white/50 backdrop-blur-sm">
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Your latest quiz attempts</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="flex items-center space-x-4">
              <Skeleton className="h-12 w-12 rounded-full" />
              <div className="space-y-2 flex-1">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-3 w-1/2" />
              </div>
            </div>
          ))}
        </CardContent>
      </Card>
    );
  }

  if (!attempts || attempts.length === 0) {
    return (
      <Card className="border-none shadow-sm bg-white/50 backdrop-blur-sm">
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Your latest quiz attempts</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center py-12 text-center">
          <div className="bg-muted/50 p-4 rounded-full mb-4">
            <Clock className="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 className="font-semibold text-lg mb-2">No activity yet</h3>
          <p className="text-muted-foreground mb-6 max-w-xs">
            Start taking quizzes to see your progress and achievements here!
          </p>
          <Button asChild>
            <Link href="/quizzes">
              Browse Quizzes
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </CardContent>
      </Card>
    );
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "completed":
        return (
          <Badge variant="default" className="bg-green-100 text-green-700 hover:bg-green-200 border-green-200">
            <CheckCircle2 className="mr-1 h-3 w-3" />
            Completed
          </Badge>
        );
      case "in_progress":
        return (
          <Badge variant="secondary" className="bg-blue-100 text-blue-700 hover:bg-blue-200 border-blue-200">
            <Clock className="mr-1 h-3 w-3" />
            In Progress
          </Badge>
        );
      case "abandoned":
        return (
          <Badge variant="destructive" className="bg-red-100 text-red-700 hover:bg-red-200 border-red-200">
            <XCircle className="mr-1 h-3 w-3" />
            Abandoned
          </Badge>
        );
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  return (
    <Card className="border-none shadow-sm bg-white/50 backdrop-blur-sm">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5 text-primary" />
              Recent Activity
            </CardTitle>
            <CardDescription>Your latest quiz attempts</CardDescription>
          </div>
          <Button variant="ghost" size="sm" className="text-muted-foreground hover:text-primary" asChild>
            <Link href="/profile">
              View All
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
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
                className="group flex items-center justify-between p-4 rounded-xl border bg-card hover:shadow-md transition-all duration-300 hover:border-primary/20"
              >
                <div className="flex items-center space-x-4 flex-1">
                  <div
                    className={`flex items-center justify-center h-12 w-12 rounded-full shadow-sm ${
                      attempt.status === "completed"
                        ? passed
                          ? "bg-gradient-to-br from-green-100 to-green-200 text-green-600"
                          : "bg-gradient-to-br from-yellow-100 to-yellow-200 text-yellow-600"
                        : "bg-gradient-to-br from-gray-100 to-gray-200 text-gray-600"
                    }`}
                  >
                    <Award className="h-6 w-6" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold truncate">Quiz Attempt</p>
                    <div className="flex items-center text-xs text-muted-foreground mt-0.5">
                      <span>
                        {attempt.started_at &&
                          formatDistanceToNow(new Date(attempt.started_at), {
                            addSuffix: true,
                          })}
                      </span>
                      <span className="mx-2">•</span>
                      <span>{attempt.quiz_title || 'General Knowledge'}</span>
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  {attempt.status === "completed" && (
                    <div className="text-right hidden sm:block">
                      <p className="text-sm font-bold">{percentage}%</p>
                      <p className="text-xs text-muted-foreground">
                        {correctAnswers}/{totalQuestions} correct
                      </p>
                    </div>
                  )}
                  <div className="flex flex-col items-end gap-2">
                    {getStatusBadge(attempt.status)}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}