'use client'

import { Loader2 } from 'lucide-react'
import { useProfile } from '@/hooks/useProfile'
import { useUserStats } from '@/hooks/useUserStats'
import { ProfileCard } from '@/components/profile/ProfileCard'
import { AttemptHistory } from '@/components/profile/AttemptHistory'
import { AchievementBadges } from '@/components/profile/AchievementBadges'
import { DetailedStatistics } from '@/components/profile/DetailedStatistics'
import { ErrorDisplay } from '@/components/common/ErrorBoundary'
import { Card, CardContent } from '@/components/ui/card'

export default function ProfilePage() {
  const { data: profile, isLoading, error } = useProfile()
  const { data: statsData, isLoading: isLoadingStats } = useUserStats()

  const stats = statsData?.data

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto">
        <ErrorDisplay
          error={error as Error}
          onRetry={() => window.location.reload()}
        />
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-muted-foreground">
              Profile not found
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Profile</h1>
        <p className="text-muted-foreground mt-2">
          Manage your personal information and settings
        </p>
      </div>

      <ProfileCard profile={profile} stats={stats} isLoadingStats={isLoadingStats} />

      {/* Detailed Statistics */}
      <DetailedStatistics stats={stats} isLoading={isLoadingStats} />

      {/* Achievements */}
      <AchievementBadges />

      {/* Attempt History */}
      <AttemptHistory />
    </div>
  )
}