"use client";

import { Trophy, Users, Award, TrendingUp } from "lucide-react";
import { LeaderboardTable } from "@/components/leaderboard/LeaderboardTable";
import { UserRankCard } from "@/components/leaderboard/UserRankCard";
import { useLeaderboard, useUserRank, useLeaderboardStats } from "@/hooks/useLeaderboard";
import { PageHero } from "@/components/common/PageHero";
import { GlassCard } from "@/components/common/GlassCard";
import { StatsCard } from "@/components/common/StatsCard";
import { StatsGrid } from "@/components/common/StatsGrid";

export default function LeaderboardPage() {
  const {
    data: leaderboard = [],
    isLoading: leaderboardLoading,
    error: leaderboardError,
  } = useLeaderboard(50);

  const {
    data: userRank,
    isLoading: rankLoading,
  } = useUserRank();

  const { data: stats, isLoading: statsLoading } = useLeaderboardStats();

  return (
    <div className="space-y-10 pb-10">
      {/* Hero Section */}
      <PageHero
        title="Leaderboard"
        icon="🏆"
        description="See how you rank among the top quiz masters! Climb the ranks by completing quizzes and earning points."
      />

      {/* Stats Cards */}
      <StatsGrid columns={4}>
        <StatsCard
          title="Total Users"
          value={stats?.total_users?.toLocaleString() ?? "0"}
          description="Active players"
          icon={Users}
          color="blue"
          loading={statsLoading}
        />
        <StatsCard
          title="Total Points"
          value={stats?.total_points_distributed?.toLocaleString() ?? "0"}
          description="Points distributed"
          icon={Trophy}
          color="yellow"
          loading={statsLoading}
        />
        <StatsCard
          title="Average Points"
          value={Math.round(stats?.average_points ?? 0).toLocaleString()}
          description="Per user"
          icon={TrendingUp}
          color="green"
          loading={statsLoading}
        />
        <StatsCard
          title="Top User"
          value={stats?.top_user?.full_name ?? "—"}
          description={stats?.top_user ? `${stats.top_user.total_points.toLocaleString()} points` : "No data"}
          icon={Award}
          color="purple"
          loading={statsLoading}
        />
      </StatsGrid>

      {/* Main Content */}
      <div className="container px-0 md:px-4">
        <div className="grid gap-8 lg:grid-cols-3">
          {/* User Rank Card - Sticky on desktop */}
          <div className="lg:col-span-1">
            <div className="sticky top-24">
              <GlassCard padding="none" rounded="2xl">
                <UserRankCard rankData={userRank} isLoading={rankLoading} />
              </GlassCard>
            </div>
          </div>

          {/* Leaderboard Table */}
          <div className="lg:col-span-2">
            <GlassCard padding="none" rounded="2xl">
              <div className="p-6 border-b border-white/10">
                <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
                  <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
                    <Trophy className="h-4 w-4" />
                  </span>
                  Top Players
                </h2>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  The top 50 players ranked by total points
                </p>
              </div>
              <div className="p-6 pt-4">
                <LeaderboardTable
                  entries={leaderboard}
                  isLoading={leaderboardLoading}
                  error={leaderboardError}
                />
              </div>
            </GlassCard>
          </div>
        </div>
      </div>
    </div>
  );
}