'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  LayoutDashboard,
  FileQuestion,
  Users,
  Trophy,
  Award,
  MessageSquare,
  Bell,
  Settings,
  Swords,
  BarChart3,
  FolderOpen,
  Heart,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { ScrollArea } from '@/components/ui/scroll-area'

const navigation = [
  {
    name: 'Dashboard',
    href: '/dashboard',
    icon: LayoutDashboard,
  },
  {
    name: 'Quizzes',
    href: '/quizzes',
    icon: FileQuestion,
  },
  {
    name: 'Categories',
    href: '/categories',
    icon: FolderOpen,
  },
  {
    name: 'Favorites',
    href: '/favorites',
    icon: Heart,
  },
  {
    name: 'Challenges',
    href: '/challenges',
    icon: Swords,
  },
  {
    name: 'Friends',
    href: '/friends',
    icon: Users,
  },
  {
    name: 'Leaderboard',
    href: '/leaderboard',
    icon: BarChart3,
  },
  {
    name: 'Achievements',
    href: '/achievements',
    icon: Trophy,
  },
  {
    name: 'Discussions',
    href: '/discussions',
    icon: MessageSquare,
  },
  {
    name: 'Notifications',
    href: '/notifications',
    icon: Bell,
  },
  {
    name: 'Settings',
    href: '/settings',
    icon: Settings,
  },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <div className="hidden border-r border-gray-200/60 dark:border-gray-800/60 bg-white/80 dark:bg-background/95 backdrop-blur-md md:block w-64">
      <ScrollArea className="h-[calc(100vh-4rem)] py-6">
        <nav className="space-y-1.5 px-4">
          {navigation.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
            const Icon = item.icon

            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all duration-300',
                  isActive
                    ? 'bg-gradient-to-r from-violet-500/10 to-indigo-500/10 text-violet-700 dark:text-violet-400 shadow-sm border border-violet-200/50 dark:border-violet-500/20'
                    : 'text-muted-foreground hover:bg-primary/10 hover:text-foreground'
                )}
              >
                <Icon className={cn(
                  'h-5 w-5 transition-colors duration-300',
                  isActive ? 'text-violet-600 dark:text-violet-400' : ''
                )} />
                {item.name}
              </Link>
            )
          })}
        </nav>
      </ScrollArea>
    </div>
  )
}