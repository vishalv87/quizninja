"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import type { AchievementProgress } from "@/types/achievement";
import { AchievementCategory, ACHIEVEMENT_CATEGORY_LABELS } from "@/constants";
import { Lock, Trophy, Star, Award, Users, Flame, BookOpen, Target } from "lucide-react";
import { formatDistanceToNow } from "date-fns";

interface AchievementCardProps {
  achievementProgress: AchievementProgress;
}

export function AchievementCard({ achievementProgress }: AchievementCardProps) {
  const { achievement, current_value, target_value, progress_percentage, is_unlocked, unlocked_at } = achievementProgress;

  // Get icon based on achievement category
  const getIcon = () => {
    const iconClass = "h-12 w-12";

    switch (achievement.category) {
      case AchievementCategory.QUIZ_MASTER:
        return <Trophy className={iconClass} />;
      case AchievementCategory.SOCIAL:
        return <Users className={iconClass} />;
      case AchievementCategory.STREAK:
        return <Flame className={iconClass} />;
      case AchievementCategory.KNOWLEDGE:
        return <BookOpen className={iconClass} />;
      case AchievementCategory.COMPETITOR:
        return <Target className={iconClass} />;
      default:
        return <Award className={iconClass} />;
    }
  };

  return (
    <Card className={`hover:shadow-lg transition-all duration-300 ${is_unlocked ? "bg-gradient-to-br from-yellow-50 to-amber-50 dark:from-yellow-950/20 dark:to-amber-950/20" : "opacity-75"}`}>
      <CardHeader className="space-y-3">
        {/* Icon and Status */}
        <div className="flex items-start justify-between gap-3">
          <div className={`p-3 rounded-full ${is_unlocked ? "bg-yellow-100 dark:bg-yellow-900/30" : "bg-muted"}`}>
            {is_unlocked ? getIcon() : <Lock className="h-12 w-12" />}
          </div>
          <div className="flex flex-col items-end gap-1">
            {is_unlocked ? (
              <Badge className="bg-yellow-600 hover:bg-yellow-700">
                <Trophy className="mr-1 h-3 w-3" />
                Unlocked
              </Badge>
            ) : (
              <Badge variant="secondary">
                <Lock className="mr-1 h-3 w-3" />
                Locked
              </Badge>
            )}
            <Badge variant="outline">{achievement.points} pts</Badge>
          </div>
        </div>

        {/* Title and Description */}
        <div>
          <h3 className="text-lg font-bold mb-1">{achievement.name}</h3>
          <p className="text-sm text-muted-foreground line-clamp-2">
            {achievement.is_secret && !is_unlocked
              ? "This is a secret achievement. Keep playing to unlock it!"
              : achievement.description}
          </p>
        </div>

        {/* Category Badge */}
        <div>
          <Badge variant="outline">{ACHIEVEMENT_CATEGORY_LABELS[achievement.category]}</Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {/* Progress Bar (for locked achievements) */}
        {!is_unlocked && (
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Progress</span>
              <span className="font-medium">
                {current_value} / {target_value}
              </span>
            </div>
            <Progress value={progress_percentage} className="h-2" />
            <p className="text-xs text-muted-foreground text-center">
              {progress_percentage.toFixed(0)}% Complete
            </p>
          </div>
        )}

        {/* Unlock Date (for unlocked achievements) */}
        {is_unlocked && unlocked_at && (
          <div className="text-sm text-muted-foreground text-center pt-2 border-t">
            Unlocked {formatDistanceToNow(new Date(unlocked_at), { addSuffix: true })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}