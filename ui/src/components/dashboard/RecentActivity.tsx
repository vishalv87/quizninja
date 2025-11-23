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
      <Card>
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
      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Your latest quiz attempts</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <p className="text-muted-foreground text-center mb-4">
            No recent activity yet. Start taking quizzes to see your progress here!
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

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>Your latest quiz attempts</CardDescription>
          </div>
          <Button variant="ghost" size="sm" asChild>
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
                className="flex items-center justify-between p-4 rounded-lg border hover:bg-accent transition-colors"
              >
                <div className="flex items-center space-x-4 flex-1">
                  <div
                    className={`flex items-center justify-center h-10 w-10 rounded-full ${
                      attempt.status === "completed"
                        ? passed
                          ? "bg-green-100 text-green-600"
                          : "bg-yellow-100 text-yellow-600"
                        : "bg-gray-100 text-gray-600"
                    }`}
                  >
                    <Award className="h-5 w-5" />
                  </div>
                  <div className="flex-1">
                    <p className="font-medium">Quiz Attempt</p>
                    <p className="text-sm text-muted-foreground">
                      {attempt.started_at &&
                        formatDistanceToNow(new Date(attempt.started_at), {
                          addSuffix: true,
                        })}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  {attempt.status === "completed" && (
                    <div className="text-right">
                      <p className="text-sm font-semibold">{percentage}%</p>
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
        </div>
      </CardContent>
    </Card>
  );
}