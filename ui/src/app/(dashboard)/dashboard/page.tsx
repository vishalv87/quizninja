'use client'

import Link from 'next/link'
import { ArrowRight, FileQuestion, Swords, Trophy } from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { useUserStats } from '@/hooks/useUserStats'
import { authLogger } from '@/lib/logger'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { RecentActivity } from '@/components/dashboard/RecentActivity'
import { ActiveSessions } from '@/components/dashboard/ActiveSessions'
import { FeaturedQuizzesDashboard } from '@/components/dashboard/FeaturedQuizzesDashboard'

export default function DashboardPage() {
  const { user, isLoading } = useAuth()
  const { data: statsData, isLoading: isLoadingStats } = useUserStats()

  const stats = statsData?.data

  authLogger.info('DashboardPage rendering', {
    isLoading,
    hasUser: !!user,
    userEmail: user?.email,
    hasStats: !!stats,
  })

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-64" />
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
      </div>
    )
  }

  const quickStats = [
    {
      title: 'Total Quizzes',
      value: isLoadingStats ? '...' : stats?.total_quizzes_completed?.toString() || '0',
      description: 'Quizzes completed',
      icon: FileQuestion,
      href: '/quizzes',
      color: 'text-blue-600',
    },
    {
      title: 'Your Rank',
      value: isLoadingStats ? '...' : stats?.rank ? `#${stats.rank}` : '—',
      description: 'On the leaderboard',
      icon: Trophy,
      href: '/leaderboard',
      color: 'text-yellow-600',
    },
    {
      title: 'Total Points',
      value: isLoadingStats ? '...' : stats?.total_points?.toLocaleString() || '0',
      description: 'Points earned',
      icon: Trophy,
      href: '/profile',
      color: 'text-purple-600',
    },
    {
      title: 'Achievements',
      value: isLoadingStats ? '...' : stats?.achievements_unlocked?.toString() || '0',
      description: 'Unlocked achievements',
      icon: Trophy,
      href: '/achievements',
      color: 'text-green-600',
    },
  ]

  const quickActions = [
    {
      title: 'Browse Quizzes',
      description: 'Find and take quizzes on various topics',
      href: '/quizzes',
      icon: FileQuestion,
    },
    {
      title: 'Challenge Friends',
      description: 'Compete with your friends in quiz battles',
      href: '/challenges',
      icon: Swords,
    },
    {
      title: 'View Achievements',
      description: 'Track your progress and unlock rewards',
      href: '/achievements',
      icon: Trophy,
    },
  ]

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          Welcome back, {user?.name || 'Quiz Ninja'}! 👋
        </h1>
        <p className="text-muted-foreground mt-2">
          Ready to test your knowledge? Here's your dashboard overview.
        </p>
      </div>

      {/* Quick Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {quickStats.map((stat) => {
          const Icon = stat.icon
          return (
            <Card key={stat.title} className="hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {stat.title}
                </CardTitle>
                <Icon className={`h-4 w-4 ${stat.color}`} />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className="text-xs text-muted-foreground mt-1">
                  {stat.description}
                </p>
                <Button
                  variant="ghost"
                  size="sm"
                  className="mt-2 h-8 px-2"
                  asChild
                >
                  <Link href={stat.href}>
                    View details
                    <ArrowRight className="ml-1 h-3 w-3" />
                  </Link>
                </Button>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Quick Actions */}
      <div>
        <h2 className="text-2xl font-bold tracking-tight mb-4">Quick Actions</h2>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {quickActions.map((action) => {
            const Icon = action.icon
            return (
              <Card key={action.title} className="hover:shadow-md transition-shadow">
                <CardHeader>
                  <div className="flex items-center space-x-2">
                    <Icon className="h-5 w-5 text-primary" />
                    <CardTitle className="text-lg">{action.title}</CardTitle>
                  </div>
                  <CardDescription>{action.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <Button className="w-full" asChild>
                    <Link href={action.href}>
                      Get Started
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Link>
                  </Button>
                </CardContent>
              </Card>
            )
          })}
        </div>
      </div>

      {/* Active Sessions */}
      <ActiveSessions />

      {/* Recent Activity */}
      <RecentActivity />

      {/* Featured Quizzes */}
      <FeaturedQuizzesDashboard />
    </div>
  )
}