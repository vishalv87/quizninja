"use client";

import { useState, useMemo } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AchievementGrid } from "@/components/achievement/AchievementGrid";
import { useAchievementProgress } from "@/hooks/useAchievementProgress";
import { useAchievementStats } from "@/hooks/useAchievementProgress";
import { Trophy, Lock, Star, TrendingUp } from "lucide-react";
import { Badge } from "@/components/ui/badge";

export default function AchievementsPage() {
  const [activeTab, setActiveTab] = useState("all");

  const {
    data: achievementsProgress = [],
    isLoading: progressLoading,
    error: progressError,
  } = useAchievementProgress();

  const { data: stats } = useAchievementStats();

  // Calculate counts for each category
  const counts = useMemo(() => {
    const unlockedCount = achievementsProgress.filter(a => a.is_unlocked).length;
    const lockedCount = achievementsProgress.filter(a => !a.is_unlocked).length;
    return {
      all: achievementsProgress.length,
      unlocked: unlockedCount,
      locked: lockedCount,
    };
  }, [achievementsProgress]);

  // Get unique categories
  const categories = useMemo(() => {
    const categorySet = new Set(achievementsProgress.map(a => a.achievement.category));
    return Array.from(categorySet);
  }, [achievementsProgress]);

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Achievements</h1>
        <p className="text-muted-foreground">
          Track your progress and unlock all achievements!
        </p>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Achievements</CardTitle>
              <Trophy className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.total_achievements}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Unlocked</CardTitle>
              <Star className="h-4 w-4 text-yellow-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-yellow-600">{stats.unlocked_achievements}</div>
              <p className="text-xs text-muted-foreground mt-1">
                {stats.completion_percentage.toFixed(1)}% Complete
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Points Earned</CardTitle>
              <TrendingUp className="h-4 w-4 text-green-600" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">{stats.points_earned}</div>
              <p className="text-xs text-muted-foreground mt-1">
                out of {stats.total_points} points
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Locked</CardTitle>
              <Lock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {stats.total_achievements - stats.unlocked_achievements}
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                Keep playing to unlock!
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-auto lg:inline-grid">
          <TabsTrigger value="all" className="gap-2">
            All
            {counts.all > 0 && (
              <Badge variant="outline" className="ml-1 px-2 py-0.5 text-xs">
                {counts.all}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="unlocked" className="gap-2">
            <Star className="h-4 w-4" />
            Unlocked
            {counts.unlocked > 0 && (
              <Badge variant="default" className="ml-1 px-2 py-0.5 text-xs">
                {counts.unlocked}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="locked" className="gap-2">
            <Lock className="h-4 w-4" />
            Locked
            {counts.locked > 0 && (
              <Badge variant="secondary" className="ml-1 px-2 py-0.5 text-xs">
                {counts.locked}
              </Badge>
            )}
          </TabsTrigger>
        </TabsList>

        {/* All Achievements Tab */}
        <TabsContent value="all" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>All Achievements</CardTitle>
              <CardDescription>
                View all available achievements and track your progress
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AchievementGrid
                achievements={achievementsProgress}
                isLoading={progressLoading}
                error={progressError}
                filter="all"
                emptyMessage="No achievements available yet. Check back soon!"
              />
            </CardContent>
          </Card>

          {/* Categories */}
          {categories.length > 0 && (
            <div className="space-y-6">
              {categories.map((category) => {
                const categoryAchievements = achievementsProgress.filter(
                  (a) => a.achievement.category === category
                );
                const unlockedInCategory = categoryAchievements.filter((a) => a.is_unlocked).length;

                return (
                  <Card key={category}>
                    <CardHeader>
                      <div className="flex items-center justify-between">
                        <CardTitle>{category}</CardTitle>
                        <Badge variant="outline">
                          {unlockedInCategory} / {categoryAchievements.length}
                        </Badge>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <AchievementGrid
                        achievements={categoryAchievements}
                        filter="all"
                      />
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          )}
        </TabsContent>

        {/* Unlocked Achievements Tab */}
        <TabsContent value="unlocked" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Unlocked Achievements</CardTitle>
              <CardDescription>
                Achievements you've successfully completed
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AchievementGrid
                achievements={achievementsProgress}
                isLoading={progressLoading}
                error={progressError}
                filter="unlocked"
                emptyMessage="No unlocked achievements yet. Complete quizzes to earn achievements!"
              />
            </CardContent>
          </Card>
        </TabsContent>

        {/* Locked Achievements Tab */}
        <TabsContent value="locked" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Locked Achievements</CardTitle>
              <CardDescription>
                Achievements you're working towards unlocking
              </CardDescription>
            </CardHeader>
            <CardContent>
              <AchievementGrid
                achievements={achievementsProgress}
                isLoading={progressLoading}
                error={progressError}
                filter="locked"
                emptyMessage="Great job! You've unlocked all achievements!"
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
