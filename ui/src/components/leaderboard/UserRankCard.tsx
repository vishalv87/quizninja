"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import type { UserRankResponse } from "@/lib/api/leaderboard";
import { Trophy, Target, Award, TrendingUp } from "lucide-react";

interface UserRankCardProps {
  rankData: UserRankResponse | undefined;
  isLoading?: boolean;
}

export function UserRankCard({ rankData, isLoading = false }: UserRankCardProps) {
  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  // Get rank badge variant
  const getRankVariant = (rank: number): "default" | "secondary" | "outline" => {
    if (rank === 1) return "default";
    if (rank <= 3) return "secondary";
    if (rank <= 10) return "outline";
    return "outline";
  };

  // Loading state
  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Your Rank</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-4">
            <Skeleton className="h-20 w-20 rounded-full" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-6 w-32" />
              <Skeleton className="h-4 w-24" />
            </div>
          </div>
          <div className="grid grid-cols-3 gap-4">
            <Skeleton className="h-20" />
            <Skeleton className="h-20" />
            <Skeleton className="h-20" />
          </div>
        </CardContent>
      </Card>
    );
  }

  // No data state
  if (!rankData) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Your Rank</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <Trophy className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
            <p className="text-sm text-muted-foreground">
              Complete quizzes to appear on the leaderboard!
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const percentile = ((rankData.total_users ?? 1) - (rankData.rank ?? 0)) / (rankData.total_users ?? 1) * 100;

  return (
    <Card className="bg-gradient-to-br from-primary/5 to-primary/10">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Trophy className="h-5 w-5" />
          Your Rank
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* User Info and Rank */}
        <div className="flex items-center gap-4">
          <Avatar className="h-20 w-20 ring-2 ring-primary/20">
            <AvatarImage src={rankData.user?.avatar_url} alt={rankData.user?.full_name ?? 'User'} />
            <AvatarFallback className="text-lg font-semibold">
              {getInitials(rankData.user?.full_name ?? 'User')}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <h3 className="text-xl font-bold mb-1">{rankData.user?.full_name ?? 'User'}</h3>
            <div className="flex items-center gap-2">
              <Badge variant={getRankVariant(rankData.rank ?? 0)} className="text-base px-3 py-1">
                #{rankData.rank ?? 0}
              </Badge>
              <span className="text-sm text-muted-foreground">
                out of {rankData.total_users?.toLocaleString() ?? '0'} users
              </span>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-3 gap-4">
          {/* Total Points */}
          <div className="text-center p-3 bg-background rounded-lg">
            <div className="flex justify-center mb-2">
              <Trophy className="h-5 w-5 text-yellow-500" />
            </div>
            <p className="text-2xl font-bold">{rankData.total_points?.toLocaleString() ?? '0'}</p>
            <p className="text-xs text-muted-foreground mt-1">Points</p>
          </div>

          {/* Quizzes Completed */}
          <div className="text-center p-3 bg-background rounded-lg">
            <div className="flex justify-center mb-2">
              <Target className="h-5 w-5 text-blue-500" />
            </div>
            <p className="text-2xl font-bold">{rankData.quizzes_completed ?? 0}</p>
            <p className="text-xs text-muted-foreground mt-1">Quizzes</p>
          </div>

          {/* Achievements */}
          <div className="text-center p-3 bg-background rounded-lg">
            <div className="flex justify-center mb-2">
              <Award className="h-5 w-5 text-purple-500" />
            </div>
            <p className="text-2xl font-bold">{rankData.achievements_unlocked ?? 0}</p>
            <p className="text-xs text-muted-foreground mt-1">Achievements</p>
          </div>
        </div>

        {/* Percentile */}
        <div className="flex items-center justify-between p-3 bg-background rounded-lg">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-green-500" />
            <span className="text-sm font-medium">Top {percentile.toFixed(1)}%</span>
          </div>
          <span className="text-xs text-muted-foreground">
            Better than {percentile.toFixed(0)}% of users
          </span>
        </div>
      </CardContent>
    </Card>
  );
}