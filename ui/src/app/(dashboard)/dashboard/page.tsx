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
import { FeaturedQuizzesDashboard } from '@/components/dashboard/FeaturedQuizzesDashboard'

export default function DashboardPage() {
  const { user, isLoading } = useAuth()
  const { data: statsData, isLoading: isLoadingStats } = useUserStats()

  const stats = statsData?.data

  if (isLoading) {
    return (
      <div className="space-y-8">
        <Skeleton className="h-48 w-full rounded-3xl" />
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-32 rounded-xl" />
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
      bgColor: 'bg-blue-100',
    },
    {
      title: 'Your Rank',
      value: isLoadingStats ? '...' : stats?.rank ? `#${stats.rank}` : '—',
      description: 'On the leaderboard',
      icon: Trophy,
      href: '/leaderboard',
      color: 'text-yellow-600',
      bgColor: 'bg-yellow-100',
    },
    {
      title: 'Total Points',
      value: isLoadingStats ? '...' : stats?.total_points?.toLocaleString() || '0',
      description: 'Points earned',
      icon: Trophy,
      href: '/profile',
      color: 'text-purple-600',
      bgColor: 'bg-purple-100',
    },
    {
      title: 'Achievements',
      value: isLoadingStats ? '...' : stats?.achievements_unlocked?.toString() || '0',
      description: 'Unlocked achievements',
      icon: Trophy,
      href: '/achievements',
      color: 'text-green-600',
      bgColor: 'bg-green-100',
    },
  ]

  const quickActions = [
    {
      title: 'Browse Quizzes',
      description: 'Find and take quizzes on various topics',
      href: '/quizzes',
      icon: FileQuestion,
      gradient: 'from-blue-500 to-cyan-500',
    },
    {
      title: 'Challenge Friends',
      description: 'Compete with your friends in quiz battles',
      href: '/challenges',
      icon: Swords,
      gradient: 'from-purple-500 to-pink-500',
    },
    {
      title: 'View Achievements',
      description: 'Track your progress and unlock rewards',
      href: '/achievements',
      icon: Trophy,
      gradient: 'from-orange-500 to-yellow-500',
    },
  ]

  return (
    <div className="space-y-10 pb-10">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-r from-violet-600 to-indigo-600 p-8 text-white shadow-xl lg:p-12">
        <div className="relative z-10 max-w-2xl">
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Welcome back, {user?.name || 'Quiz Ninja'}! 👋
          </h1>
          <p className="text-lg text-indigo-100 mb-6">
            Ready to test your knowledge today? "Knowledge is power, but enthusiasm pulls the switch."
          </p>
          <Button 
            size="lg" 
            className="bg-white text-indigo-600 hover:bg-indigo-50 border-0 font-semibold"
            asChild
          >
            <Link href="/quizzes">
              Start a Quiz
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
          </Button>
        </div>
        
        {/* Decorative background elements */}
        <div className="absolute right-0 top-0 -mt-10 -mr-10 h-64 w-64 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute bottom-0 right-20 -mb-10 h-40 w-40 rounded-full bg-indigo-400/20 blur-2xl" />
      </div>

      {/* Quick Stats */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {quickStats.map((stat) => {
          const Icon = stat.icon
          return (
            <Card key={stat.title} className="group overflow-hidden border border-gray-200/60 shadow-sm hover:shadow-lg transition-all duration-300 bg-white hover:-translate-y-1">
              <CardContent className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <div className={`p-3 rounded-2xl ${stat.bgColor} group-hover:scale-110 transition-transform duration-300`}>
                    <Icon className={`h-6 w-6 ${stat.color}`} />
                  </div>
                  {stat.value !== '...' && (
                     <span className="text-xs font-medium text-muted-foreground bg-gray-100 px-2.5 py-1 rounded-full border border-gray-100">
                       {stat.title}
                     </span>
                  )}
                </div>
                <div className="space-y-1">
                  <h3 className="text-3xl font-bold tracking-tight text-gray-900">
                    {stat.value}
                  </h3>
                  <p className="text-sm text-muted-foreground font-medium">
                    {stat.description}
                  </p>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>

      <div className="grid gap-8 lg:grid-cols-3">
        {/* Main Content Area (2 columns) */}
        <div className="lg:col-span-2 space-y-8">
          {/* Quick Actions */}
          <div>
            <h2 className="text-2xl font-bold tracking-tight mb-6 flex items-center gap-2">
              <span className="bg-primary/10 p-2 rounded-lg text-primary">⚡</span>
              Quick Actions
            </h2>
            <div className="grid gap-4 sm:grid-cols-3">
              {quickActions.map((action) => {
                const Icon = action.icon
                return (
                  <Link key={action.title} href={action.href} className="group block h-full">
                    <div className="relative h-full overflow-hidden rounded-2xl border bg-card p-6 shadow-sm transition-all duration-300 hover:shadow-lg hover:-translate-y-1">
                      <div className={`absolute inset-0 opacity-0 group-hover:opacity-5 transition-opacity duration-300 bg-gradient-to-br ${action.gradient}`} />
                      <div className="relative z-10 flex flex-col h-full">
                        <div className={`mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br ${action.gradient} text-white shadow-md`}>
                          <Icon className="h-6 w-6" />
                        </div>
                        <h3 className="mb-2 font-semibold tracking-tight text-lg">{action.title}</h3>
                        <p className="text-sm text-muted-foreground">{action.description}</p>
                      </div>
                    </div>
                  </Link>
                )
              })}
            </div>
          </div>

          {/* Recent Activity */}
          <RecentActivity />
        </div>

        {/* Sidebar Area (1 column) */}
        <div className="space-y-8">
          {/* Featured Quizzes */}
          <FeaturedQuizzesDashboard />
        </div>
      </div>
    </div>
  )
}