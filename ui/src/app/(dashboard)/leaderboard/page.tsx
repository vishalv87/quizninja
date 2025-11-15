"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LeaderboardTable } from "@/components/leaderboard/LeaderboardTable";
import { UserRankCard } from "@/components/leaderboard/UserRankCard";
import { useLeaderboard, useUserRank } from "@/hooks/useLeaderboard";
import { Trophy, Users, Award, TrendingUp } from "lucide-react";
import { useLeaderboardStats } from "@/hooks/useLeaderboard";

export default function LeaderboardPage() {
  const {
    data: leaderboard = [],
    isLoading: leaderboardLoading,
    error: leaderboardError,
  } = useLeaderboard(50); // Top 50 users

  const {
    data: userRank,
    isLoading: rankLoading,
  } = useUserRank();

  const { data: stats } = useLeaderboardStats();

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Leaderboard</h1>
        <p className="text-muted-foreground">
          See how you rank among the top quiz masters!
        </p>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Users</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.total_users?.toLocaleString() ?? '0'}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Points</CardTitle>
              <Trophy className="h-4 w-4 text-yellow-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.total_points_distributed?.toLocaleString() ?? '0'}</div>
              <p className="text-xs text-muted-foreground mt-1">
                Points distributed
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Average Points</CardTitle>
              <TrendingUp className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{Math.round(stats.average_points ?? 0).toLocaleString()}</div>
              <p className="text-xs text-muted-foreground mt-1">
                Per user
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Top User</CardTitle>
              <Award className="h-4 w-4 text-purple-600" />
            </CardHeader>
            <CardContent>
              {stats.top_user ? (
                <>
                  <div className="text-lg font-bold truncate">{stats.top_user.full_name}</div>
                  <p className="text-xs text-muted-foreground mt-1">
                    {stats.top_user.total_points.toLocaleString()} points
                  </p>
                </>
              ) : (
                <div className="text-sm text-muted-foreground">No data</div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* Main Content */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* User Rank Card - Sticky on desktop */}
        <div className="lg:col-span-1">
          <div className="sticky top-8">
            <UserRankCard rankData={userRank} isLoading={rankLoading} />
          </div>
        </div>

        {/* Leaderboard Table */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Trophy className="h-5 w-5" />
                Top Players
              </CardTitle>
              <CardDescription>
                The top 50 players ranked by total points
              </CardDescription>
            </CardHeader>
            <CardContent>
              <LeaderboardTable
                entries={leaderboard}
                isLoading={leaderboardLoading}
                error={leaderboardError}
              />
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}