'use client'

import Link from 'next/link'
import { ArrowRight, FileQuestion, Users, Trophy } from 'lucide-react'
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
      color: 'text-blue-500',
      bgColor: 'bg-blue-500/10',
      borderColor: 'border-blue-200/20',
    },
    {
      title: 'Your Rank',
      value: isLoadingStats ? '...' : stats?.rank ? `#${stats.rank}` : '—',
      description: 'On the leaderboard',
      icon: Trophy,
      href: '/leaderboard',
      color: 'text-yellow-500',
      bgColor: 'bg-yellow-500/10',
      borderColor: 'border-yellow-200/20',
    },
    {
      title: 'Total Points',
      value: isLoadingStats ? '...' : stats?.total_points?.toLocaleString() || '0',
      description: 'Points earned',
      icon: Trophy,
      href: '/profile',
      color: 'text-purple-500',
      bgColor: 'bg-purple-500/10',
      borderColor: 'border-purple-200/20',
    },
    {
      title: 'Achievements',
      value: isLoadingStats ? '...' : stats?.achievements_unlocked?.toString() || '0',
      description: 'Unlocked achievements',
      icon: Trophy,
      href: '/achievements',
      color: 'text-green-500',
      bgColor: 'bg-green-500/10',
      borderColor: 'border-green-200/20',
    },
  ]

  const quickActions = [
    {
      title: 'Browse Quizzes',
      description: 'Find and take quizzes on various topics',
      href: '/quizzes',
      icon: FileQuestion,
      gradient: 'from-blue-500 to-cyan-500',
      shadow: 'shadow-blue-500/20',
    },
    {
      title: 'Find Friends',
      description: 'Connect with other quiz enthusiasts',
      href: '/friends',
      icon: Users,
      gradient: 'from-purple-500 to-pink-500',
      shadow: 'shadow-purple-500/20',
    },
    {
      title: 'View Achievements',
      description: 'Track your progress and unlock rewards',
      href: '/achievements',
      icon: Trophy,
      gradient: 'from-orange-500 to-yellow-500',
      shadow: 'shadow-orange-500/20',
    },
  ]

  return (
    <div className="space-y-10 pb-10">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-violet-600 via-indigo-600 to-purple-700 p-8 text-white shadow-2xl shadow-indigo-500/30 lg:p-12 border border-white/10">
        <div className="relative z-10 max-w-2xl">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl mb-6 drop-shadow-sm">
            Welcome back, {user?.name || 'Quiz Ninja'}! 👋
          </h1>
          <p className="text-xl text-indigo-100 mb-8 font-medium leading-relaxed">
            Ready to test your knowledge today? "Knowledge is power, but enthusiasm pulls the switch."
          </p>
          <Button 
            size="lg" 
            className="bg-white text-indigo-600 hover:bg-indigo-50 border-0 font-bold h-12 px-8 rounded-xl shadow-lg shadow-black/10 transition-all hover:scale-105 hover:shadow-xl"
            asChild
          >
            <Link href="/quizzes">
              Start a Quiz
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
          </Button>
        </div>
        
        {/* Decorative background elements */}
        <div className="absolute right-0 top-0 -mt-20 -mr-20 h-96 w-96 rounded-full bg-white/10 blur-3xl" />
        <div className="absolute bottom-0 right-20 -mb-20 h-64 w-64 rounded-full bg-indigo-400/20 blur-3xl" />
        <div className="absolute left-10 bottom-10 h-32 w-32 rounded-full bg-purple-400/20 blur-2xl" />
      </div>

      {/* Quick Stats */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {quickStats.map((stat) => {
          const Icon = stat.icon
          return (
            <Card key={stat.title} className={`group overflow-hidden border border-white/20 dark:border-white/10 shadow-lg shadow-black/5 hover:shadow-xl transition-all duration-300 bg-white/40 dark:bg-black/40 backdrop-blur-md hover:-translate-y-1 ${stat.borderColor}`}>
              <CardContent className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <div className={`p-3 rounded-2xl ${stat.bgColor} group-hover:scale-110 transition-transform duration-300 ring-1 ring-white/20`}>
                    <Icon className={`h-6 w-6 ${stat.color}`} />
                  </div>
                  {stat.value !== '...' && (
                     <span className="text-xs font-bold text-slate-500 dark:text-slate-400 bg-white/50 dark:bg-white/10 px-2.5 py-1 rounded-full border border-white/20 backdrop-blur-sm">
                       {stat.title}
                     </span>
                  )}
                </div>
                <div className="space-y-1">
                  <h3 className="text-3xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
                    {stat.value}
                  </h3>
                  <p className="text-sm text-slate-500 dark:text-slate-400 font-medium">
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
            <h2 className="text-2xl font-bold tracking-tight mb-6 flex items-center gap-2 text-slate-800 dark:text-slate-100">
              <span className="bg-gradient-to-br from-amber-400 to-orange-500 text-white p-1.5 rounded-lg shadow-sm">⚡</span>
              Quick Actions
            </h2>
            <div className="grid gap-4 sm:grid-cols-3">
              {quickActions.map((action) => {
                const Icon = action.icon
                return (
                  <Link key={action.title} href={action.href} className="group block h-full">
                    <div className="relative h-full overflow-hidden rounded-2xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md p-6 shadow-sm transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
                      <div className={`absolute inset-0 opacity-0 group-hover:opacity-10 transition-opacity duration-500 bg-gradient-to-br ${action.gradient}`} />
                      <div className="relative z-10 flex flex-col h-full">
                        <div className={`mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br ${action.gradient} text-white shadow-lg ${action.shadow} group-hover:scale-110 transition-transform duration-300`}>
                          <Icon className="h-6 w-6" />
                        </div>
                        <h3 className="mb-2 font-bold tracking-tight text-lg text-slate-800 dark:text-slate-100 group-hover:text-slate-900 dark:group-hover:text-white transition-colors">{action.title}</h3>
                        <p className="text-sm text-slate-500 dark:text-slate-400 group-hover:text-slate-700 dark:group-hover:text-slate-200 transition-colors">{action.description}</p>
                      </div>
                    </div>
                  </Link>
                )
              })}
            </div>
          </div>

          {/* Recent Activity */}
          <div className="rounded-3xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md shadow-sm overflow-hidden">
             <RecentActivity />
          </div>
        </div>

        {/* Sidebar Area (1 column) */}
        <div className="space-y-8">
          {/* Featured Quizzes */}
          <div className="rounded-3xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md shadow-sm overflow-hidden">
             <FeaturedQuizzesDashboard />
          </div>
        </div>
      </div>
    </div>
  )
}