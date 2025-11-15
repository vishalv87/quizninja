"use client";

import Link from "next/link";
import { Trophy, ArrowRight } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { AchievementBadge } from "@/components/achievement/AchievementBadge";
import { useUserAchievements } from "@/hooks/useUserAchievements";

/**
 * AchievementBadges Component
 * Displays user's recent unlocked achievements on the profile page
 * Shows up to 6 most recent achievements with a link to view all
 */
export function AchievementBadges() {
  const { data: achievements = [], isLoading, error } = useUserAchievements();

  // Sort by unlock date (most recent first) and take first 6
  const recentAchievements = [...achievements]
    .sort((a, b) => new Date(b.unlocked_at).getTime() - new Date(a.unlocked_at).getTime())
    .slice(0, 6);

  // Calculate total points from all unlocked achievements
  const totalPoints = achievements.reduce((sum, ach) => sum + (ach.achievement.points || 0), 0);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Trophy className="h-5 w-5 text-yellow-500" />
              Achievements
            </CardTitle>
            <CardDescription>
              {achievements.length} unlocked · {totalPoints} points earned
            </CardDescription>
          </div>
          <Button variant="ghost" size="sm" asChild>
            <Link href="/achievements">
              View All
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center gap-3">
                <Skeleton className="h-10 w-10 rounded-full" />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-3 w-1/2" />
                </div>
                <Skeleton className="h-6 w-16" />
              </div>
            ))}
          </div>
        ) : error ? (
          <div className="text-center py-8">
            <p className="text-sm text-muted-foreground">
              Failed to load achievements. Please try again later.
            </p>
          </div>
        ) : achievements.length === 0 ? (
          <div className="text-center py-8">
            <div className="flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100 text-yellow-600 dark:bg-yellow-900/20 dark:text-yellow-400 mx-auto mb-4">
              <Trophy className="h-6 w-6" />
            </div>
            <p className="text-sm text-muted-foreground mb-1">No achievements yet</p>
            <p className="text-xs text-muted-foreground">
              Complete quizzes and challenges to unlock achievements!
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            {recentAchievements.map((achievement) => (
              <AchievementBadge
                key={achievement.id}
                achievement={achievement}
                variant="compact"
              />
            ))}
            {achievements.length > 6 && (
              <div className="pt-2">
                <Button variant="outline" size="sm" asChild className="w-full">
                  <Link href="/achievements">
                    View {achievements.length - 6} more achievement{achievements.length - 6 !== 1 ? 's' : ''}
                  </Link>
                </Button>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}