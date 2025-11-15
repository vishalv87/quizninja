import {
  UserPlus,
  UserCheck,
  Swords,
  Trophy,
  Bell,
  MessageSquare,
  CheckCircle2,
  Info,
  type LucideIcon
} from "lucide-react";
import type { NotificationType } from "@/types/notification";
import { formatDistanceToNow } from "date-fns";

/**
 * Get the icon component for a notification type
 */
export function getNotificationIcon(type: string): LucideIcon {
  const iconMap: Record<NotificationType, LucideIcon> = {
    friend_request: UserPlus,
    friend_accepted: UserCheck,
    challenge_received: Swords,
    challenge_accepted: Swords,
    challenge_completed: CheckCircle2,
    achievement_unlocked: Trophy,
    quiz_reminder: Bell,
    discussion_reply: MessageSquare,
    system: Info,
  };

  return iconMap[type as NotificationType] || Info;
}

/**
 * Get the color scheme for a notification type
 * Returns Tailwind CSS classes for icon and background
 */
export function getNotificationColor(type: string): {
  icon: string;
  bg: string;
  text: string;
} {
  const colorMap: Record<NotificationType, { icon: string; bg: string; text: string }> = {
    friend_request: {
      icon: "text-blue-600 dark:text-blue-400",
      bg: "bg-blue-50 dark:bg-blue-950/50",
      text: "text-blue-900 dark:text-blue-100",
    },
    friend_accepted: {
      icon: "text-green-600 dark:text-green-400",
      bg: "bg-green-50 dark:bg-green-950/50",
      text: "text-green-900 dark:text-green-100",
    },
    challenge_received: {
      icon: "text-orange-600 dark:text-orange-400",
      bg: "bg-orange-50 dark:bg-orange-950/50",
      text: "text-orange-900 dark:text-orange-100",
    },
    challenge_accepted: {
      icon: "text-orange-600 dark:text-orange-400",
      bg: "bg-orange-50 dark:bg-orange-950/50",
      text: "text-orange-900 dark:text-orange-100",
    },
    challenge_completed: {
      icon: "text-purple-600 dark:text-purple-400",
      bg: "bg-purple-50 dark:bg-purple-950/50",
      text: "text-purple-900 dark:text-purple-100",
    },
    achievement_unlocked: {
      icon: "text-yellow-600 dark:text-yellow-400",
      bg: "bg-yellow-50 dark:bg-yellow-950/50",
      text: "text-yellow-900 dark:text-yellow-100",
    },
    quiz_reminder: {
      icon: "text-indigo-600 dark:text-indigo-400",
      bg: "bg-indigo-50 dark:bg-indigo-950/50",
      text: "text-indigo-900 dark:text-indigo-100",
    },
    discussion_reply: {
      icon: "text-pink-600 dark:text-pink-400",
      bg: "bg-pink-50 dark:bg-pink-950/50",
      text: "text-pink-900 dark:text-pink-100",
    },
    system: {
      icon: "text-gray-600 dark:text-gray-400",
      bg: "bg-gray-50 dark:bg-gray-950/50",
      text: "text-gray-900 dark:text-gray-100",
    },
  };

  return colorMap[type as NotificationType] || colorMap.system;
}

/**
 * Format notification timestamp to relative time
 */
export function formatNotificationTime(timestamp: string): string {
  try {
    return formatDistanceToNow(new Date(timestamp), { addSuffix: true });
  } catch (error) {
    return "Unknown time";
  }
}

/**
 * Get user-friendly label for notification type
 */
export function getNotificationTypeLabel(type: string): string {
  const labelMap: Record<NotificationType, string> = {
    friend_request: "Friend Request",
    friend_accepted: "Friend Accepted",
    challenge_received: "Challenge Received",
    challenge_accepted: "Challenge Accepted",
    challenge_completed: "Challenge Completed",
    achievement_unlocked: "Achievement Unlocked",
    quiz_reminder: "Quiz Reminder",
    discussion_reply: "Discussion Reply",
    system: "System Notification",
  };

  return labelMap[type as NotificationType] || "Notification";
}

/**
 * Get action URL based on notification type and data
 */
export function getNotificationActionUrl(type: string, data?: Record<string, any>): string | null {
  if (!data) return null;

  const urlMap: Record<NotificationType, (data: Record<string, any>) => string | null> = {
    friend_request: (d) => d.user_id ? `/profile/${d.user_id}` : "/friends",
    friend_accepted: (d) => d.user_id ? `/profile/${d.user_id}` : "/friends",
    challenge_received: (d) => d.challenge_id ? `/challenges/${d.challenge_id}` : "/challenges",
    challenge_accepted: (d) => d.challenge_id ? `/challenges/${d.challenge_id}` : "/challenges",
    challenge_completed: (d) => d.challenge_id ? `/challenges/${d.challenge_id}` : "/challenges",
    achievement_unlocked: (d) => d.achievement_id ? `/achievements` : "/achievements",
    quiz_reminder: (d) => d.quiz_id ? `/quizzes/${d.quiz_id}` : "/quizzes",
    discussion_reply: (d) => d.discussion_id ? `/discussions/${d.discussion_id}` : "/discussions",
    system: () => null,
  };

  const urlGenerator = urlMap[type as NotificationType];
  return urlGenerator ? urlGenerator(data) : null;
}

/**
 * Check if notification is expired
 */
export function isNotificationExpired(expiresAt?: string): boolean {
  if (!expiresAt) return false;
  try {
    return new Date(expiresAt) < new Date();
  } catch (error) {
    return false;
  }
}

/**
 * Group notifications by date (Today, Yesterday, This Week, Earlier)
 */
export function groupNotificationsByDate(notifications: Array<{ created_at: string; [key: string]: any }>) {
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  const thisWeek = new Date(today);
  thisWeek.setDate(thisWeek.getDate() - 7);

  const groups = {
    today: [] as typeof notifications,
    yesterday: [] as typeof notifications,
    thisWeek: [] as typeof notifications,
    earlier: [] as typeof notifications,
  };

  notifications.forEach((notification) => {
    const date = new Date(notification.created_at);

    if (date >= today) {
      groups.today.push(notification);
    } else if (date >= yesterday) {
      groups.yesterday.push(notification);
    } else if (date >= thisWeek) {
      groups.thisWeek.push(notification);
    } else {
      groups.earlier.push(notification);
    }
  });

  return groups;
}
