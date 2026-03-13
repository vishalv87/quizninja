# Profile Components

## Overview

Components for displaying and editing user profiles, statistics, and quiz history.

## Components

| Component | File | Purpose |
|-----------|------|---------|
| ProfileCard | `ProfileCard.tsx` | Current user's profile summary |
| UserProfileCard | `UserProfileCard.tsx` | Other user's profile view |
| ProfileEditForm | `ProfileEditForm.tsx` | Profile editing form |
| DetailedStatistics | `DetailedStatistics.tsx` | Comprehensive stats breakdown |
| AchievementBadges | `AchievementBadges.tsx` | Profile achievement display |
| AttemptHistory | `AttemptHistory.tsx` | Quiz attempt history list |

## ProfileCard

Displays the current user's profile summary.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `user` | `User` | User data |
| `stats` | `UserStats` | User statistics |
| `className` | `string` | Additional classes |

### Display

- Avatar
- Name and email
- Bio
- Join date
- Quick stats (quizzes, score, rank)
- Edit profile button

### Usage

```tsx
import { ProfileCard } from "@/components/profile/ProfileCard";
import { useAuth } from "@/hooks/useAuth";
import { useUserStats } from "@/hooks/useUserStats";

function ProfilePage() {
  const { user } = useAuth();
  const { data: stats } = useUserStats();

  return <ProfileCard user={user} stats={stats} />;
}
```

---

## UserProfileCard

View another user's profile (read-only).

### Props

| Prop | Type | Description |
|------|------|-------------|
| `profile` | `UserProfile` | User profile data |
| `friendshipStatus` | `FriendshipStatus` | Relationship status |
| `onAddFriend` | `() => void` | Add friend handler |

### Features

- Avatar and name
- Bio (if public)
- Statistics (if visible)
- Friend action button
- Privacy-aware display

### Usage

```tsx
<UserProfileCard
  profile={profile}
  friendshipStatus={friendship}
  onAddFriend={handleAddFriend}
/>
```

---

## ProfileEditForm

Form for editing user profile.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `profile` | `Profile` | Current profile data |
| `onSubmit` | `(data) => void` | Form submission |
| `isSubmitting` | `boolean` | Loading state |

### Editable Fields

| Field | Type | Validation |
|-------|------|------------|
| Name | `input` | 2-50 characters |
| Bio | `textarea` | Max 500 characters |
| Avatar URL | `input[url]` | Valid URL |

### Usage

```tsx
import { ProfileEditForm } from "@/components/profile/ProfileEditForm";
import { useUpdateProfile } from "@/hooks/useProfile";

function EditProfilePage() {
  const { updateProfile, isUpdating } = useUpdateProfile();

  return (
    <ProfileEditForm
      profile={currentProfile}
      onSubmit={updateProfile}
      isSubmitting={isUpdating}
    />
  );
}
```

---

## DetailedStatistics

Comprehensive statistics breakdown.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `userId` | `string` | User ID |
| `stats` | `UserStats` | Statistics data |

### Statistics Categories

**Quiz Performance:**
- Total quizzes completed
- Total points earned
- Average score
- Best score
- Pass rate

**Progress:**
- Current streak
- Longest streak
- Categories completed
- Achievements earned

**Rankings:**
- Global rank
- Category rankings

### Usage

```tsx
<DetailedStatistics userId={userId} stats={stats} />
```

---

## AttemptHistory

List of user's quiz attempts.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `userId` | `string` | User ID |
| `limit` | `number` | Max items |
| `showAll` | `boolean` | Show pagination |

### Display per Attempt

- Quiz title
- Score (X/Y)
- Percentage
- Pass/fail badge
- Date completed
- Time taken
- Link to results

### Usage

```tsx
<AttemptHistory userId={userId} limit={10} showAll={true} />
```

---

## AchievementBadges

Mini achievement display for profile.

### Props

| Prop | Type | Description |
|------|------|-------------|
| `achievements` | `UserAchievement[]` | Earned achievements |
| `limit` | `number` | Max to show |
| `showMore` | `boolean` | Show "view all" link |

### Usage

```tsx
<AchievementBadges
  achievements={userAchievements}
  limit={6}
  showMore={true}
/>
```

---

## Profile Page Structure

```tsx
// Own profile (/profile)
<div className="space-y-8">
  <ProfileCard user={user} stats={stats} />
  <DetailedStatistics userId={user.id} stats={stats} />
  <AchievementBadges achievements={achievements} />
  <AttemptHistory userId={user.id} />
</div>

// Other user's profile (/profile/[userId])
<div className="space-y-8">
  <UserProfileCard
    profile={profile}
    friendshipStatus={friendship}
    onAddFriend={handleAddFriend}
  />
  {profile.visibility !== "private" && (
    <>
      <DetailedStatistics userId={userId} />
      <AttemptHistory userId={userId} />
    </>
  )}
</div>
```

## Related Documentation

- [Parent: Components Overview](../README.md)
- [Profile Routes](../../app/(dashboard)/README.md)
- [Auth Types](../../types/README.md)
- [useProfile Hook](../../hooks/README.md)
