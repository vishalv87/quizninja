"use client";

import { Trophy, Lock } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { UserAchievement, AchievementProgress } from "@/types/achievement";

interface AchievementBadgeProps {
  achievement: UserAchievement | AchievementProgress;
  variant?: "compact" | "detailed";
}

export function AchievementBadge({ achievement, variant = "compact" }: AchievementBadgeProps) {
  // Check if it's a UserAchievement (unlocked) or AchievementProgress
  const isUnlocked = "unlocked_at" in achievement && !!achievement.unlocked_at;
  const achievementData = "achievement" in achievement ? achievement.achievement : null;

  if (!achievementData) return null;

  const progress = "progress_percentage" in achievement ? achievement.progress_percentage : 100;

  if (variant === "compact") {
    return (
      <div
        className={`flex items-center gap-3 p-3 rounded-lg border transition-all ${
          isUnlocked
            ? "bg-gradient-to-br from-yellow-50 to-amber-50 border-yellow-200 dark:from-yellow-900/10 dark:to-amber-900/10 dark:border-yellow-800"
            : "bg-muted border-muted-foreground/20 opacity-60"
        }`}
      >
        <div
          className={`flex items-center justify-center h-10 w-10 rounded-full flex-shrink-0 ${
            isUnlocked
              ? "bg-yellow-500 text-white"
              : "bg-muted-foreground/20 text-muted-foreground"
          }`}
        >
          {isUnlocked ? (
            <Trophy className="h-5 w-5" />
          ) : (
            <Lock className="h-5 w-5" />
          )}
        </div>
        <div className="flex-1 min-w-0">
          <p className={`font-medium text-sm truncate ${!isUnlocked && "text-muted-foreground"}`}>
            {achievementData.name}
          </p>
          <p className="text-xs text-muted-foreground truncate">
            {achievementData.description}
          </p>
        </div>
        <Badge variant={isUnlocked ? "default" : "secondary"} className="flex-shrink-0">
          {achievementData.points} pts
        </Badge>
      </div>
    );
  }

  // Detailed variant
  return (
    <Card
      className={`p-4 transition-all ${
        isUnlocked
          ? "bg-gradient-to-br from-yellow-50 to-amber-50 border-yellow-200 dark:from-yellow-900/10 dark:to-amber-900/10 dark:border-yellow-800"
          : "bg-muted border-muted-foreground/20 opacity-60"
      }`}
    >
      <div className="flex items-start gap-4">
        <div
          className={`flex items-center justify-center h-14 w-14 rounded-full flex-shrink-0 ${
            isUnlocked
              ? "bg-yellow-500 text-white"
              : "bg-muted-foreground/20 text-muted-foreground"
          }`}
        >
          {isUnlocked ? (
            <Trophy className="h-7 w-7" />
          ) : (
            <Lock className="h-7 w-7" />
          )}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2 mb-2">
            <div>
              <h4 className={`font-semibold ${!isUnlocked && "text-muted-foreground"}`}>
                {achievementData.name}
              </h4>
              <p className="text-sm text-muted-foreground mt-1">
                {achievementData.description}
              </p>
            </div>
            <Badge variant={isUnlocked ? "default" : "secondary"} className="flex-shrink-0">
              {achievementData.points} pts
            </Badge>
          </div>
          {!isUnlocked && "progress_percentage" in achievement && (
            <div className="mt-3">
              <div className="flex items-center justify-between text-xs text-muted-foreground mb-1">
                <span>Progress</span>
                <span>{Math.round(progress)}%</span>
              </div>
              <div className="w-full bg-muted-foreground/20 rounded-full h-2">
                <div
                  className="bg-primary rounded-full h-2 transition-all"
                  style={{ width: `${Math.min(progress, 100)}%` }}
                />
              </div>
            </div>
          )}
          {isUnlocked && "unlocked_at" in achievement && achievement.unlocked_at && (
            <p className="text-xs text-muted-foreground mt-2">
              Unlocked {new Date(achievement.unlocked_at).toLocaleDateString()}
            </p>
          )}
        </div>
      </div>
    </Card>
  );
}
