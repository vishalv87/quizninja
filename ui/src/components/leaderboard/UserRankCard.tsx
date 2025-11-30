"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import type { UserRankResponse } from "@/lib/api/leaderboard";
import { Trophy, Target, Award, TrendingUp } from "lucide-react";
import { cn } from "@/lib/utils";

interface UserRankCardProps {
  rankData: UserRankResponse | null | undefined;
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

  // Glassmorphic card base styles
  const glassCardStyles = cn(
    "bg-white/40 dark:bg-black/40",
    "backdrop-blur-md",
    "border border-white/20 dark:border-white/10",
    "shadow-lg shadow-black/5",
    "rounded-2xl",
    "overflow-hidden"
  );

  // Loading state
  if (isLoading) {
    return (
      <div className={glassCardStyles}>
        <div className="p-6 border-b border-white/10">
          <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 flex items-center gap-2">
            <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
              <Trophy className="h-4 w-4" />
            </span>
            Your Rank
          </h3>
        </div>
        <div className="p-6 space-y-4">
          <div className="flex items-center gap-4">
            <Skeleton className="h-20 w-20 rounded-full" />
            <div className="flex-1 space-y-2">
              <Skeleton className="h-6 w-32" />
              <Skeleton className="h-4 w-24" />
            </div>
          </div>
          <div className="grid grid-cols-3 gap-4">
            <Skeleton className="h-20 rounded-lg" />
            <Skeleton className="h-20 rounded-lg" />
            <Skeleton className="h-20 rounded-lg" />
          </div>
        </div>
      </div>
    );
  }

  // No data state
  if (!rankData) {
    return (
      <div className={glassCardStyles}>
        <div className="p-6 border-b border-white/10">
          <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 flex items-center gap-2">
            <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
              <Trophy className="h-4 w-4" />
            </span>
            Your Rank
          </h3>
        </div>
        <div className="p-6">
          <div className="text-center py-8">
            <Trophy className="h-12 w-12 mx-auto text-slate-400 dark:text-slate-500 mb-3" />
            <p className="text-sm text-slate-500 dark:text-slate-400">
              Complete quizzes to appear on the leaderboard!
            </p>
          </div>
        </div>
      </div>
    );
  }

  // Safe percentile calculation - handle edge cases like 0 users or division by zero
  const totalUsers = rankData.total_users || 1;
  const rank = rankData.rank || 1;
  const percentile = totalUsers > 0
    ? Math.max(0, Math.min(100, ((totalUsers - rank + 1) / totalUsers) * 100))
    : 0;

  return (
    <div className={glassCardStyles}>
      {/* Header */}
      <div className="p-6 border-b border-white/10">
        <h3 className="text-lg font-semibold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
            <Trophy className="h-4 w-4" />
          </span>
          Your Rank
        </h3>
      </div>

      {/* Content */}
      <div className="p-6 space-y-6">
        {/* User Info and Rank */}
        <div className="flex items-center gap-4">
          <Avatar className="h-20 w-20 ring-2 ring-amber-400/30">
            <AvatarImage src={rankData.user?.avatar_url} alt={rankData.user?.full_name ?? 'User'} />
            <AvatarFallback className="text-lg font-semibold bg-gradient-to-br from-amber-100 to-orange-100 dark:from-amber-900/50 dark:to-orange-900/50 text-amber-700 dark:text-amber-300">
              {getInitials(rankData.user?.full_name ?? 'User')}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <h3 className="text-xl font-bold mb-1 text-slate-800 dark:text-slate-100">{rankData.user?.full_name ?? 'User'}</h3>
            <div className="flex items-center gap-2">
              <Badge variant={getRankVariant(rank)} className="text-base px-3 py-1">
                #{rank}
              </Badge>
              <span className="text-sm text-slate-500 dark:text-slate-400">
                out of {totalUsers.toLocaleString()} {totalUsers === 1 ? 'user' : 'users'}
              </span>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-3 gap-3">
          {/* Total Points */}
          <div className="text-center p-3 bg-white/50 dark:bg-black/30 rounded-xl border border-white/20 dark:border-white/5">
            <div className="flex justify-center mb-2">
              <Trophy className="h-5 w-5 text-yellow-500" />
            </div>
            <p className="text-2xl font-bold text-slate-800 dark:text-slate-100">{rankData.total_points?.toLocaleString() ?? '0'}</p>
            <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Points</p>
          </div>

          {/* Quizzes Completed */}
          <div className="text-center p-3 bg-white/50 dark:bg-black/30 rounded-xl border border-white/20 dark:border-white/5">
            <div className="flex justify-center mb-2">
              <Target className="h-5 w-5 text-blue-500" />
            </div>
            <p className="text-2xl font-bold text-slate-800 dark:text-slate-100">{rankData.quizzes_completed ?? 0}</p>
            <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Quizzes</p>
          </div>

          {/* Achievements */}
          <div className="text-center p-3 bg-white/50 dark:bg-black/30 rounded-xl border border-white/20 dark:border-white/5">
            <div className="flex justify-center mb-2">
              <Award className="h-5 w-5 text-purple-500" />
            </div>
            <p className="text-2xl font-bold text-slate-800 dark:text-slate-100">{rankData.achievements_unlocked ?? 0}</p>
            <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Achievements</p>
          </div>
        </div>

        {/* Percentile */}
        <div className="flex items-center justify-between p-3 bg-white/50 dark:bg-black/30 rounded-xl border border-white/20 dark:border-white/5">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-green-500" />
            <span className="text-sm font-medium text-slate-700 dark:text-slate-200">Top {percentile.toFixed(1)}%</span>
          </div>
          <span className="text-xs text-slate-500 dark:text-slate-400">
            Better than {percentile.toFixed(0)}% of users
          </span>
        </div>
      </div>
    </div>
  );
}