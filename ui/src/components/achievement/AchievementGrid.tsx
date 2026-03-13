"use client";

import { AchievementCard } from "./AchievementCard";
import type { AchievementProgress } from "@/types/achievement";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Trophy, AlertCircle } from "lucide-react";

interface AchievementGridProps {
  achievements: AchievementProgress[];
  isLoading?: boolean;
  error?: Error | null;
  emptyMessage?: string;
  filter?: "all" | "unlocked" | "locked";
}

export function AchievementGrid({
  achievements,
  isLoading = false,
  error = null,
  emptyMessage = "No achievements found",
  filter = "all",
}: AchievementGridProps) {
  // Filter achievements based on filter prop
  const filteredAchievements = achievements.filter((achievement) => {
    if (filter === "unlocked") {
      return achievement.is_unlocked;
    } else if (filter === "locked") {
      return !achievement.is_unlocked;
    }
    return true; // "all"
  });

  // Sort achievements: unlocked first, then by progress percentage
  const sortedAchievements = [...filteredAchievements].sort((a, b) => {
    // Unlocked achievements come first
    if (a.is_unlocked && !b.is_unlocked) return -1;
    if (!a.is_unlocked && b.is_unlocked) return 1;

    // If both unlocked, sort by unlock date (most recent first)
    if (a.is_unlocked && b.is_unlocked) {
      if (a.unlocked_at && b.unlocked_at) {
        return new Date(b.unlocked_at).getTime() - new Date(a.unlocked_at).getTime();
      }
      return 0;
    }

    // If both locked, sort by progress percentage (highest first)
    return b.progress_percentage - a.progress_percentage;
  });

  // Loading state
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[...Array(6)].map((_, i) => (
          <div key={i} className="space-y-3">
            <Skeleton className="h-64 w-full rounded-lg" />
          </div>
        ))}
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          {error.message || "Failed to load achievements"}
        </AlertDescription>
      </Alert>
    );
  }

  // Empty state
  if (!sortedAchievements || sortedAchievements.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="rounded-full bg-muted p-6 mb-4">
          <Trophy className="h-12 w-12 text-muted-foreground" />
        </div>
        <h3 className="text-lg font-semibold mb-2">No Achievements Yet</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          {emptyMessage}
        </p>
      </div>
    );
  }

  // Render achievements grid
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {sortedAchievements.map((achievement) => (
        <AchievementCard key={achievement.achievement_id} achievementProgress={achievement} />
      ))}
    </div>
  );
}
