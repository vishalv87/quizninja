"use client";

import Link from "next/link";
import { Trophy, ArrowRight, Lock } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

/**
 * AchievementBadges Component
 * Displays user achievements on the profile page
 *
 * Note: This is a Phase 4 placeholder component.
 * Full implementation will be completed in Phase 7: Achievements & Gamification
 * which includes:
 * - Achievement hooks (useAchievements, useUserAchievements)
 * - Achievement API integration
 * - Achievement progress tracking
 * - Achievement unlock animations
 */
export function AchievementBadges() {
  // Placeholder: In Phase 7, this will use useUserAchievements() hook
  const isLoading = false;
  const achievements: any[] = [];

  // Show placeholder UI since full achievement system is Phase 7
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Trophy className="h-5 w-5 text-yellow-500" />
              Achievements
            </CardTitle>
            <CardDescription>Your unlocked achievements and progress</CardDescription>
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
        {/* Placeholder content - will be replaced in Phase 7 */}
        <div className="flex flex-col items-center justify-center py-12 px-4">
          <div className="flex items-center gap-3 mb-4">
            <div className="flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100 text-yellow-600 dark:bg-yellow-900/20 dark:text-yellow-400">
              <Trophy className="h-6 w-6" />
            </div>
            <div className="flex items-center justify-center h-12 w-12 rounded-full bg-gray-100 text-gray-400 dark:bg-gray-800">
              <Lock className="h-6 w-6" />
            </div>
            <div className="flex items-center justify-center h-12 w-12 rounded-full bg-gray-100 text-gray-400 dark:bg-gray-800">
              <Lock className="h-6 w-6" />
            </div>
          </div>
          <p className="text-muted-foreground text-center mb-2">
            Achievement system coming soon!
          </p>
          <p className="text-sm text-muted-foreground text-center max-w-md">
            Full achievement tracking, progress monitoring, and unlock notifications will be
            available in Phase 7 of the migration.
          </p>
          <div className="mt-6 flex gap-2">
            <Badge variant="outline" className="text-xs">
              Phase 7: Achievements & Gamification
            </Badge>
          </div>
        </div>

        {/*
        Future Phase 7 implementation will include:
        - Grid of unlocked achievement badges
        - Progress bars for locked achievements
        - Filter by category
        - Achievement statistics
        - Recent unlocks section

        Example code structure:

        {isLoading ? (
          <LoadingSkeleton />
        ) : achievements.length > 0 ? (
          <div className="grid gap-3 md:grid-cols-2">
            {achievements.slice(0, 6).map((achievement) => (
              <AchievementBadge
                key={achievement.id}
                achievement={achievement}
                variant="compact"
              />
            ))}
          </div>
        ) : (
          <EmptyState />
        )}
        */}
      </CardContent>
    </Card>
  );
}