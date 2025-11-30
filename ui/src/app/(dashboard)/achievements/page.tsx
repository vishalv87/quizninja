"use client";

import { useState, useMemo } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AchievementGrid } from "@/components/achievement/AchievementGrid";
import { useAchievementProgress, useAchievementStats } from "@/hooks/useAchievementProgress";
import { Trophy, Lock, Star, TrendingUp } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { GlassCard } from "@/components/common/GlassCard";
import { StatsCard } from "@/components/common/StatsCard";
import { StatsGrid } from "@/components/common/StatsGrid";

export default function AchievementsPage() {
  const [activeTab, setActiveTab] = useState("all");

  const {
    data: achievementsProgress = [],
    isLoading: progressLoading,
    error: progressError,
  } = useAchievementProgress();

  const { data: stats, isLoading: statsLoading } = useAchievementStats();

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
    <div className="space-y-10 pb-10">
      {/* Stats Cards */}
      <StatsGrid columns={4}>
        <StatsCard
          title="Total"
          value={stats?.total_achievements ?? 0}
          description="Achievements available"
          icon={Trophy}
          color="blue"
          loading={statsLoading}
        />
        <StatsCard
          title="Unlocked"
          value={stats?.unlocked_achievements ?? 0}
          description={`${stats?.completion_percentage?.toFixed(1) ?? 0}% Complete`}
          icon={Star}
          color="yellow"
          loading={statsLoading}
        />
        <StatsCard
          title="Points Earned"
          value={stats?.points_earned ?? 0}
          description={`out of ${stats?.total_points ?? 0} points`}
          icon={TrendingUp}
          color="green"
          loading={statsLoading}
        />
        <StatsCard
          title="Locked"
          value={(stats?.total_achievements ?? 0) - (stats?.unlocked_achievements ?? 0)}
          description="Keep playing to unlock!"
          icon={Lock}
          color="purple"
          loading={statsLoading}
        />
      </StatsGrid>

      {/* Tabs */}
      <div className="container px-0 md:px-4">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full max-w-md grid-cols-3 bg-white/60 dark:bg-black/40 backdrop-blur-md border border-white/20 dark:border-white/10 p-1 rounded-xl shadow-sm">
            <TabsTrigger
              value="all"
              className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
            >
              All
              {counts.all > 0 && (
                <Badge variant="outline" className="ml-1 px-2 py-0.5 text-xs">
                  {counts.all}
                </Badge>
              )}
            </TabsTrigger>
            <TabsTrigger
              value="unlocked"
              className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
            >
              <Star className="h-4 w-4" />
              Unlocked
              {counts.unlocked > 0 && (
                <Badge variant="default" className="ml-1 px-2 py-0.5 text-xs">
                  {counts.unlocked}
                </Badge>
              )}
            </TabsTrigger>
            <TabsTrigger
              value="locked"
              className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
            >
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
          <TabsContent value="all" className="space-y-6">
            <GlassCard padding="none" rounded="2xl">
              <div className="p-6 border-b border-white/10">
                <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
                  <span className="bg-gradient-to-br from-violet-500 to-purple-600 text-white p-1.5 rounded-lg shadow-sm">
                    <Trophy className="h-4 w-4" />
                  </span>
                  All Achievements
                </h2>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  View all available achievements and track your progress
                </p>
              </div>
              <div className="p-6">
                <AchievementGrid
                  achievements={achievementsProgress}
                  isLoading={progressLoading}
                  error={progressError}
                  filter="all"
                  emptyMessage="No achievements available yet. Check back soon!"
                />
              </div>
            </GlassCard>

            {/* Categories */}
            {categories.length > 0 && (
              <div className="space-y-6">
                {categories.map((category) => {
                  const categoryAchievements = achievementsProgress.filter(
                    (a) => a.achievement.category === category
                  );
                  const unlockedInCategory = categoryAchievements.filter((a) => a.is_unlocked).length;

                  return (
                    <GlassCard key={category} padding="none" rounded="2xl">
                      <div className="p-6 border-b border-white/10">
                        <div className="flex items-center justify-between">
                          <h3 className="text-lg font-bold tracking-tight text-slate-800 dark:text-slate-100">
                            {category}
                          </h3>
                          <Badge variant="outline" className="bg-white/50 dark:bg-white/10">
                            {unlockedInCategory} / {categoryAchievements.length}
                          </Badge>
                        </div>
                      </div>
                      <div className="p-6">
                        <AchievementGrid
                          achievements={categoryAchievements}
                          filter="all"
                        />
                      </div>
                    </GlassCard>
                  );
                })}
              </div>
            )}
          </TabsContent>

          {/* Unlocked Achievements Tab */}
          <TabsContent value="unlocked" className="space-y-6">
            <GlassCard padding="none" rounded="2xl">
              <div className="p-6 border-b border-white/10">
                <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
                  <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
                    <Star className="h-4 w-4" />
                  </span>
                  Unlocked Achievements
                </h2>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  Achievements you&apos;ve successfully completed
                </p>
              </div>
              <div className="p-6">
                <AchievementGrid
                  achievements={achievementsProgress}
                  isLoading={progressLoading}
                  error={progressError}
                  filter="unlocked"
                  emptyMessage="No unlocked achievements yet. Complete quizzes to earn achievements!"
                />
              </div>
            </GlassCard>
          </TabsContent>

          {/* Locked Achievements Tab */}
          <TabsContent value="locked" className="space-y-6">
            <GlassCard padding="none" rounded="2xl">
              <div className="p-6 border-b border-white/10">
                <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
                  <span className="bg-gradient-to-br from-slate-400 to-slate-600 text-white p-1.5 rounded-lg shadow-sm">
                    <Lock className="h-4 w-4" />
                  </span>
                  Locked Achievements
                </h2>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  Achievements you&apos;re working towards unlocking
                </p>
              </div>
              <div className="p-6">
                <AchievementGrid
                  achievements={achievementsProgress}
                  isLoading={progressLoading}
                  error={progressError}
                  filter="locked"
                  emptyMessage="Great job! You've unlocked all achievements!"
                />
              </div>
            </GlassCard>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}