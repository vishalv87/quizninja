"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { usePreferences, useUpdatePreferences } from "@/hooks/usePreferences";
import { useCategories } from "@/hooks/useCategories";
import { preferencesSchema, type PreferencesFormData } from "@/schemas/preferences";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Checkbox } from "@/components/ui/checkbox";

export function PreferencesForm() {
  const { data: preferencesData, isLoading: isLoadingPreferences } = usePreferences();
  const { data: categoriesData } = useCategories();
  const { mutate: updatePreferences, isPending } = useUpdatePreferences();

  const preferences = preferencesData?.data;
  const categories = categoriesData || [];

  const [selectedCategories, setSelectedCategories] = useState<string[]>(
    preferences?.preferred_categories || []
  );
  const [emailNotifications, setEmailNotifications] = useState(
    preferences?.email_notifications ?? true
  );
  const [notificationFrequency, setNotificationFrequency] = useState(
    preferences?.notification_frequency || "instant"
  );
  const [preferredDifficulty, setPreferredDifficulty] = useState(
    preferences?.preferred_difficulty || "all"
  );
  const [theme, setTheme] = useState(preferences?.theme || "system");

  // Privacy settings state
  const [profileVisibility, setProfileVisibility] = useState(
    preferences?.profile_visibility || "public"
  );
  const [showAchievements, setShowAchievements] = useState(
    preferences?.show_achievements ?? true
  );
  const [showStats, setShowStats] = useState(
    preferences?.show_stats ?? true
  );
  const [allowFriendRequests, setAllowFriendRequests] = useState(
    preferences?.allow_friend_requests ?? true
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const data = {
      preferred_categories: selectedCategories,
      email_notifications: emailNotifications,
      notification_frequency: notificationFrequency as "instant" | "daily" | "weekly" | "never",
      preferred_difficulty: preferredDifficulty as "easy" | "medium" | "hard" | "all",
      theme: theme as "light" | "dark" | "system",
      // Privacy settings
      profile_visibility: profileVisibility as "public" | "friends_only" | "private",
      show_achievements: showAchievements,
      show_stats: showStats,
      allow_friend_requests: allowFriendRequests,
    };

    updatePreferences(data);
  };

  const handleCategoryToggle = (categoryId: string) => {
    setSelectedCategories((prev) =>
      prev.includes(categoryId)
        ? prev.filter((id) => id !== categoryId)
        : [...prev, categoryId]
    );
  };

  if (isLoadingPreferences) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Category Preferences */}
      <Card>
        <CardHeader>
          <CardTitle>Category Preferences</CardTitle>
          <CardDescription>
            Select categories you're interested in. We'll recommend quizzes based on your preferences.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {categories.map((category: any) => (
              <div key={category.id} className="flex items-center space-x-2">
                <Checkbox
                  id={`category-${category.id}`}
                  checked={selectedCategories.includes(category.id)}
                  onCheckedChange={() => handleCategoryToggle(category.id)}
                />
                <Label htmlFor={`category-${category.id}`} className="cursor-pointer">
                  {category.display_name || category.name}
                </Label>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Difficulty Preference */}
      <Card>
        <CardHeader>
          <CardTitle>Difficulty Level</CardTitle>
          <CardDescription>
            Choose your preferred quiz difficulty level
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <Label htmlFor="difficulty">Preferred Difficulty</Label>
            <Select value={preferredDifficulty} onValueChange={setPreferredDifficulty}>
              <SelectTrigger id="difficulty">
                <SelectValue placeholder="Select difficulty" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Levels</SelectItem>
                <SelectItem value="easy">Easy</SelectItem>
                <SelectItem value="medium">Medium</SelectItem>
                <SelectItem value="hard">Hard</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Notification Settings */}
      <Card>
        <CardHeader>
          <CardTitle>Notifications</CardTitle>
          <CardDescription>
            Manage how you receive notifications
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="email-notifications">Email Notifications</Label>
              <p className="text-sm text-muted-foreground">
                Receive notifications via email
              </p>
            </div>
            <Switch
              id="email-notifications"
              checked={emailNotifications}
              onCheckedChange={setEmailNotifications}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="notification-frequency">Notification Frequency</Label>
            <Select value={notificationFrequency} onValueChange={setNotificationFrequency}>
              <SelectTrigger id="notification-frequency">
                <SelectValue placeholder="Select frequency" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="instant">Instant</SelectItem>
                <SelectItem value="daily">Daily Digest</SelectItem>
                <SelectItem value="weekly">Weekly Digest</SelectItem>
                <SelectItem value="never">Never</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Theme Settings */}
      <Card>
        <CardHeader>
          <CardTitle>Appearance</CardTitle>
          <CardDescription>
            Customize how the app looks for you
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <Label htmlFor="theme">Theme</Label>
            <Select value={theme} onValueChange={(value) => setTheme(value as typeof theme)}>
              <SelectTrigger id="theme">
                <SelectValue placeholder="Select theme" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="light">Light</SelectItem>
                <SelectItem value="dark">Dark</SelectItem>
                <SelectItem value="system">System</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Privacy Settings */}
      <Card>
        <CardHeader>
          <CardTitle>Privacy Settings</CardTitle>
          <CardDescription>
            Control who can see your profile and activity
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="profile-visibility">Profile Visibility</Label>
            <Select value={profileVisibility} onValueChange={(value) => setProfileVisibility(value as "public" | "friends_only" | "private")}>
              <SelectTrigger id="profile-visibility">
                <SelectValue placeholder="Select visibility" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="public">Public - Anyone can view</SelectItem>
                <SelectItem value="friends_only">Friends Only - Only friends can view</SelectItem>
                <SelectItem value="private">Private - Only you can view</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="show-achievements">Show Achievements</Label>
              <p className="text-sm text-muted-foreground">
                Allow others to see your achievements
              </p>
            </div>
            <Switch
              id="show-achievements"
              checked={showAchievements}
              onCheckedChange={setShowAchievements}
            />
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="show-stats">Show Statistics</Label>
              <p className="text-sm text-muted-foreground">
                Allow others to see your quiz statistics
              </p>
            </div>
            <Switch
              id="show-stats"
              checked={showStats}
              onCheckedChange={setShowStats}
            />
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="allow-friend-requests">Allow Friend Requests</Label>
              <p className="text-sm text-muted-foreground">
                Let other users send you friend requests
              </p>
            </div>
            <Switch
              id="allow-friend-requests"
              checked={allowFriendRequests}
              onCheckedChange={setAllowFriendRequests}
            />
          </div>
        </CardContent>
      </Card>

      {/* Save Button */}
      <div className="flex justify-end">
        <Button type="submit" disabled={isPending}>
          {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Save Preferences
        </Button>
      </div>
    </form>
  );
}
