"use client";

import { Trophy, Sparkles } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import type { UserAchievement } from "@/types/achievement";

interface AchievementToastContentProps {
  achievement: UserAchievement;
}

/**
 * Custom toast content for achievement unlocks
 * Displays rich achievement information with animations
 */
function AchievementToastContent({ achievement }: AchievementToastContentProps) {
  const achievementData = achievement.achievement;

  return (
    <div className="flex items-start gap-3 w-full">
      {/* Icon with animation */}
      <div className="relative flex-shrink-0">
        <div className="flex items-center justify-center h-12 w-12 rounded-full bg-gradient-to-br from-yellow-400 to-amber-500 animate-pulse">
          <Trophy className="h-6 w-6 text-white" />
        </div>
        {/* Sparkle effect */}
        <Sparkles className="absolute -top-1 -right-1 h-4 w-4 text-yellow-400 animate-bounce" />
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <p className="font-semibold text-sm">Achievement Unlocked!</p>
          <Badge variant="secondary" className="text-xs bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400">
            +{achievementData.points} pts
          </Badge>
        </div>
        <p className="font-medium text-base">{achievementData.name}</p>
        <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
          {achievementData.description}
        </p>
        {achievementData.category && (
          <Badge variant="outline" className="mt-2 text-xs">
            {achievementData.category}
          </Badge>
        )}
      </div>
    </div>
  );
}

/**
 * Shows a custom toast notification for achievement unlocks
 * Includes celebration animations, achievement details, and category badge
 *
 * @param achievement - The unlocked achievement to display
 */
export function showAchievementToast(achievement: UserAchievement) {
  toast.custom(
    (t) => (
      <div className="bg-background border border-border rounded-lg shadow-lg p-4 max-w-md w-full">
        <AchievementToastContent achievement={achievement} />
      </div>
    ),
    {
      duration: 6000,
      position: "top-center",
    }
  );
}

/**
 * Shows multiple achievement unlock toasts
 * Staggers the display to avoid overwhelming the user
 *
 * @param achievements - Array of unlocked achievements
 */
export function showMultipleAchievementToasts(achievements: UserAchievement[]) {
  achievements.forEach((achievement, index) => {
    setTimeout(() => {
      showAchievementToast(achievement);
    }, index * 1000); // Stagger by 1 second
  });
}