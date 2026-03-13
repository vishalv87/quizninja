'use client'

import { Loader2 } from 'lucide-react'
import { useProfile } from '@/hooks/useProfile'
import { useUserStats } from '@/hooks/useUserStats'
import { useAuth } from '@/hooks/useAuth'
import { ProfileCard } from '@/components/profile/ProfileCard'
import { AttemptHistory } from '@/components/profile/AttemptHistory'
import { AchievementBadges } from '@/components/profile/AchievementBadges'
import { DetailedStatistics } from '@/components/profile/DetailedStatistics'
import { ErrorDisplay } from '@/components/common/ErrorBoundary'
import { GlassCard } from '@/components/common/GlassCard'
import { Skeleton } from '@/components/ui/skeleton'

export default function ProfilePage() {
  const { user } = useAuth()
  const { data: profile, isLoading, error } = useProfile()
  const { data: statsData, isLoading: isLoadingStats } = useUserStats()

  const stats = statsData?.data

  if (isLoading) {
    return (
      <div className="space-y-10 pb-10">
        <Skeleton className="h-48 w-full rounded-3xl" />
        <Skeleton className="h-64 w-full rounded-2xl" />
        <Skeleton className="h-48 w-full rounded-2xl" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-10 pb-10">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
            Profile
          </h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">
            View and manage your personal information
          </p>
        </div>
        <div className="container px-0 md:px-4">
          <GlassCard>
            <ErrorDisplay
              error={error as Error}
              onRetry={() => window.location.reload()}
            />
          </GlassCard>
        </div>
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="space-y-10 pb-10">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
            Profile
          </h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">
            View and manage your personal information
          </p>
        </div>
        <div className="container px-0 md:px-4">
          <GlassCard>
            <p className="text-center text-muted-foreground py-8">
              Profile not found
            </p>
          </GlassCard>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-10 pb-10">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
          Profile
        </h1>
        <p className="text-slate-500 dark:text-slate-400 mt-1">
          View your profile, track your progress, and see your achievements
        </p>
      </div>

      <div className="container px-0 md:px-4 space-y-8">
        {/* Profile Card */}
        <GlassCard padding="none" rounded="2xl">
          <ProfileCard profile={profile} stats={stats} isLoadingStats={isLoadingStats} />
        </GlassCard>

        {/* Detailed Statistics */}
        <GlassCard padding="none" rounded="2xl">
          <div className="p-6 border-b border-white/10">
            <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
              <span className="bg-gradient-to-br from-blue-500 to-cyan-500 text-white p-1.5 rounded-lg shadow-sm">
                📊
              </span>
              Statistics
            </h2>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
              Your detailed quiz performance metrics
            </p>
          </div>
          <div className="p-6">
            <DetailedStatistics stats={stats} isLoading={isLoadingStats} />
          </div>
        </GlassCard>

        {/* Achievements */}
        <GlassCard padding="none" rounded="2xl">
          <div className="p-6 border-b border-white/10">
            <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
              <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">
                🏆
              </span>
              Achievements
            </h2>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
              Your unlocked badges and rewards
            </p>
          </div>
          <div className="p-6">
            <AchievementBadges />
          </div>
        </GlassCard>

        {/* Attempt History */}
        <GlassCard padding="none" rounded="2xl">
          <div className="p-6 border-b border-white/10">
            <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
              <span className="bg-gradient-to-br from-violet-500 to-purple-600 text-white p-1.5 rounded-lg shadow-sm">
                📝
              </span>
              Quiz History
            </h2>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
              Your recent quiz attempts and scores
            </p>
          </div>
          <div className="p-6">
            <AttemptHistory />
          </div>
        </GlassCard>
      </div>
    </div>
  )
}
