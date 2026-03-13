"use client";

import { BarChart3, Target, Zap, Calendar, Award } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import type { UserStats } from "@/types/user";

interface DetailedStatisticsProps {
  stats?: UserStats;
  isLoading?: boolean;
}

export function DetailedStatistics({ stats, isLoading }: DetailedStatisticsProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Detailed Statistics</CardTitle>
          <CardDescription>In-depth analysis of your quiz performance</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="space-y-2">
              <Skeleton className="h-4 w-1/3" />
              <Skeleton className="h-2 w-full" />
            </div>
          ))}
        </CardContent>
      </Card>
    );
  }

  if (!stats) {
    return null;
  }

  const completionRate = stats.total_quizzes_taken > 0
    ? Math.round((stats.total_quizzes_completed / stats.total_quizzes_taken) * 100)
    : 0;

  const averageTimePerQuiz = stats.total_quizzes_completed > 0
    ? Math.round(stats.total_time_spent_minutes / stats.total_quizzes_completed)
    : 0;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <BarChart3 className="h-5 w-5 text-primary" />
          Detailed Statistics
        </CardTitle>
        <CardDescription>In-depth analysis of your quiz performance</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Performance Metrics */}
        <div className="grid gap-6 md:grid-cols-2">
          {/* Completion Rate */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Target className="h-4 w-4 text-blue-500" />
                <span className="text-sm font-medium">Completion Rate</span>
              </div>
              <span className="text-sm font-bold">{completionRate}%</span>
            </div>
            <Progress value={completionRate} className="h-2" />
            <p className="text-xs text-muted-foreground">
              {stats.total_quizzes_completed} of {stats.total_quizzes_taken} quizzes completed
            </p>
          </div>

          {/* Average Score */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Award className="h-4 w-4 text-yellow-500" />
                <span className="text-sm font-medium">Average Score</span>
              </div>
              <span className="text-sm font-bold">{stats.average_score.toFixed(1)}%</span>
            </div>
            <Progress value={stats.average_score} className="h-2" />
            <p className="text-xs text-muted-foreground">
              Across all completed quizzes
            </p>
          </div>
        </div>

        {/* Activity Metrics */}
        <div className="grid gap-4 md:grid-cols-3">
          {/* Current Streak */}
          <div className="p-4 rounded-lg border bg-gradient-to-br from-orange-50 to-red-50 dark:from-orange-900/10 dark:to-red-900/10 border-orange-200 dark:border-orange-800">
            <div className="flex items-center gap-2 mb-2">
              <Zap className="h-4 w-4 text-orange-600 dark:text-orange-400" />
              <span className="text-xs font-medium text-orange-900 dark:text-orange-100">Current Streak</span>
            </div>
            <p className="text-2xl font-bold text-orange-900 dark:text-orange-100">
              {stats.current_streak} days
            </p>
            <p className="text-xs text-orange-700 dark:text-orange-300 mt-1">
              Longest: {stats.longest_streak} days
            </p>
          </div>

          {/* Total Time */}
          <div className="p-4 rounded-lg border bg-gradient-to-br from-blue-50 to-cyan-50 dark:from-blue-900/10 dark:to-cyan-900/10 border-blue-200 dark:border-blue-800">
            <div className="flex items-center gap-2 mb-2">
              <Calendar className="h-4 w-4 text-blue-600 dark:text-blue-400" />
              <span className="text-xs font-medium text-blue-900 dark:text-blue-100">Time Spent</span>
            </div>
            <p className="text-2xl font-bold text-blue-900 dark:text-blue-100">
              {Math.round(stats.total_time_spent_minutes / 60)} hrs
            </p>
            <p className="text-xs text-blue-700 dark:text-blue-300 mt-1">
              Avg: {averageTimePerQuiz} min/quiz
            </p>
          </div>

        </div>

        {/* Point and Achievement Summary */}
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="p-4 rounded-lg border">
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm font-medium">Total Points</span>
              <Award className="h-5 w-5 text-primary" />
            </div>
            <div className="space-y-2">
              <p className="text-3xl font-bold text-primary">{stats.total_points.toLocaleString()}</p>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <span>Rank #{stats.rank}</span>
                <span>•</span>
                <span>{stats.achievements_unlocked} achievements</span>
              </div>
            </div>
          </div>

          <div className="p-4 rounded-lg border">
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm font-medium">Quiz Performance</span>
              <BarChart3 className="h-5 w-5 text-primary" />
            </div>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Taken</span>
                <span className="font-semibold">{stats.total_quizzes_taken}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Completed</span>
                <span className="font-semibold text-green-600 dark:text-green-400">
                  {stats.total_quizzes_completed}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Completion Rate</span>
                <span className="font-semibold">{completionRate}%</span>
              </div>
            </div>
          </div>
        </div>

        {/* Future Enhancement Placeholder */}
        <div className="p-4 rounded-lg border border-dashed bg-muted/50">
          <p className="text-sm text-muted-foreground text-center">
            More detailed analytics coming soon: category performance, difficulty breakdown, and time-based trends
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
