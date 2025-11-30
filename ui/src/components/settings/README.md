# Settings Components

## Overview

Components for user settings and preferences management. Allows users to customize their experience, manage account details, and configure privacy options.

## Components

| Component | File | Purpose |
|-----------|------|---------|
| AccountSettingsForm | `AccountSettingsForm.tsx` | Account management |
| PreferencesForm | `PreferencesForm.tsx` | User preferences |

## PreferencesForm

Comprehensive form for managing user preferences.

### Features

Organized into multiple Card sections:

1. **Category Preferences** - Select interested quiz categories
2. **Difficulty Level** - Preferred quiz difficulty
3. **Notifications** - Email and frequency settings
4. **Appearance** - Theme selection
5. **Privacy Settings** - Visibility and social options

### Category Preferences Section

```tsx
<Card>
  <CardHeader>
    <CardTitle>Category Preferences</CardTitle>
    <CardDescription>
      Select categories you're interested in.
    </CardDescription>
  </CardHeader>
  <CardContent>
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {categories.map(category => (
        <div className="flex items-center space-x-2">
          <Checkbox
            checked={selectedCategories.includes(category.id)}
            onCheckedChange={() => handleCategoryToggle(category.id)}
          />
          <Label>{category.display_name}</Label>
        </div>
      ))}
    </div>
  </CardContent>
</Card>
```

### Difficulty Level Section

```tsx
<Select value={preferredDifficulty} onValueChange={setPreferredDifficulty}>
  <SelectItem value="all">All Levels</SelectItem>
  {DIFFICULTY_OPTIONS.map(({ value, label }) => (
    <SelectItem key={value} value={value}>{label}</SelectItem>
  ))}
</Select>
```

### Notification Settings

| Setting | Type | Options |
|---------|------|---------|
| Email Notifications | Switch | On/Off |
| Notification Frequency | Select | Instant, Daily, Weekly |

### Appearance Settings

| Setting | Type | Options |
|---------|------|---------|
| Theme | Select | Light, Dark, System |

### Privacy Settings

| Setting | Type | Description |
|---------|------|-------------|
| Profile Visibility | Select | Public, Friends Only, Private |
| Show Achievements | Switch | Allow others to see achievements |
| Show Statistics | Switch | Allow others to see quiz stats |
| Allow Friend Requests | Switch | Let users send friend requests |

### State Management

```tsx
// Uses local state for form values
const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
const [emailNotifications, setEmailNotifications] = useState(true);
const [notificationFrequency, setNotificationFrequency] = useState("instant");
const [preferredDifficulty, setPreferredDifficulty] = useState("all");
const [theme, setTheme] = useState("system");
const [profileVisibility, setProfileVisibility] = useState("public");
// ...more state
```

### Form Submission

```tsx
const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault();

  const data = {
    preferred_categories: selectedCategories,
    email_notifications: emailNotifications,
    notification_frequency: notificationFrequency,
    preferred_difficulty: preferredDifficulty,
    theme: theme,
    profile_visibility: profileVisibility,
    show_achievements: showAchievements,
    show_stats: showStats,
    allow_friend_requests: allowFriendRequests,
  };

  updatePreferences(data);
};
```

### Usage

```tsx
import { PreferencesForm } from "@/components/settings/PreferencesForm";

function SettingsPage() {
  return (
    <div className="max-w-2xl mx-auto">
      <h1>Settings</h1>
      <PreferencesForm />
    </div>
  );
}
```

---

## AccountSettingsForm

Form for managing account details.

### Features

- Change email address
- Update password
- Delete account (with confirmation)

### Sections

**Email Change:**
```tsx
<Card>
  <CardHeader>
    <CardTitle>Email Address</CardTitle>
  </CardHeader>
  <CardContent>
    <Input type="email" value={email} onChange={...} />
    <Button>Update Email</Button>
  </CardContent>
</Card>
```

**Password Change:**
```tsx
<Card>
  <CardHeader>
    <CardTitle>Change Password</CardTitle>
  </CardHeader>
  <CardContent>
    <Input type="password" placeholder="Current Password" />
    <Input type="password" placeholder="New Password" />
    <Input type="password" placeholder="Confirm New Password" />
    <Button>Update Password</Button>
  </CardContent>
</Card>
```

**Danger Zone:**
```tsx
<Card className="border-destructive">
  <CardHeader>
    <CardTitle>Danger Zone</CardTitle>
  </CardHeader>
  <CardContent>
    <p>Permanently delete your account and all data.</p>
    <Button variant="destructive">Delete Account</Button>
  </CardContent>
</Card>
```

### Usage

```tsx
import { AccountSettingsForm } from "@/components/settings/AccountSettingsForm";

<Tabs defaultValue="preferences">
  <TabsList>
    <TabsTrigger value="preferences">Preferences</TabsTrigger>
    <TabsTrigger value="account">Account</TabsTrigger>
  </TabsList>

  <TabsContent value="preferences">
    <PreferencesForm />
  </TabsContent>

  <TabsContent value="account">
    <AccountSettingsForm />
  </TabsContent>
</Tabs>
```

---

## Settings Page Structure

```tsx
// /settings page
<div className="space-y-8">
  {/* Page Header */}
  <div>
    <h1>Settings</h1>
    <p>Manage your account and preferences</p>
  </div>

  <Tabs defaultValue="preferences">
    <TabsList className="grid w-full grid-cols-2">
      <TabsTrigger value="preferences">Preferences</TabsTrigger>
      <TabsTrigger value="account">Account</TabsTrigger>
    </TabsList>

    <TabsContent value="preferences" className="mt-6">
      <PreferencesForm />
    </TabsContent>

    <TabsContent value="account" className="mt-6">
      <AccountSettingsForm />
    </TabsContent>
  </Tabs>
</div>
```

## Constants Used

From `@/constants`:

```typescript
// Difficulty options
DIFFICULTY_OPTIONS: [
  { value: "beginner", label: "Beginner" },
  { value: "intermediate", label: "Intermediate" },
  { value: "advanced", label: "Advanced" },
]

// Theme options
THEME_OPTIONS: [
  { value: "light", label: "Light" },
  { value: "dark", label: "Dark" },
  { value: "system", label: "System" },
]

// Profile visibility options
PROFILE_VISIBILITY_OPTIONS: [
  { value: "public", label: "Public" },
  { value: "friends", label: "Friends Only" },
  { value: "private", label: "Private" },
]

// Notification frequency options
NOTIFICATION_FREQUENCY_OPTIONS: [
  { value: "instant", label: "Instant" },
  { value: "daily", label: "Daily Digest" },
  { value: "weekly", label: "Weekly Summary" },
]
```

## Hooks Used

```typescript
// Get current preferences
const { data: preferences, isLoading } = usePreferences();

// Update preferences
const { mutate: updatePreferences, isPending } = useUpdatePreferences();

// Get categories for selection
const { data: categories } = useCategories();
```

## Data Types

```typescript
interface UserPreferences {
  preferred_categories: string[];
  preferred_difficulty: QuizDifficulty | "all";
  email_notifications: boolean;
  notification_frequency: NotificationFrequency;
  theme: Theme;
  profile_visibility: ProfileVisibility;
  show_achievements: boolean;
  show_stats: boolean;
  allow_friend_requests: boolean;
}

type Theme = "light" | "dark" | "system";
type ProfileVisibility = "public" | "friends" | "private";
type NotificationFrequency = "instant" | "daily" | "weekly";
```

## Related Documentation

- [Parent: Components Overview](../README.md)
- [Constants](../../constants/README.md) - Option constants
- [Preferences Schema](../../schemas/README.md)
- [usePreferences Hook](../../hooks/README.md)

