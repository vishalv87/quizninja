'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  LayoutDashboard,
  FileQuestion,
  Users,
  Trophy,
  MessageSquare,
  Bell,
  Settings,
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
    <div className="hidden border-r border-white/20 dark:border-white/10 bg-white/30 dark:bg-black/30 backdrop-blur-xl md:block w-72 shadow-[4px_0_24px_-12px_rgba(0,0,0,0.1)] z-30 h-full overflow-hidden">
      <ScrollArea className="h-full py-6">
        <nav className="space-y-2 px-4">
          {navigation.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
            const Icon = item.icon

            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'group flex items-center gap-3 px-4 py-3 rounded-2xl text-sm font-medium transition-all duration-300 ease-out',
                  isActive
                    ? 'bg-gradient-to-r from-violet-600/90 to-indigo-600/90 text-white shadow-lg shadow-indigo-500/20 scale-[1.02]'
                    : 'text-slate-600 dark:text-slate-300 hover:bg-white/50 dark:hover:bg-white/10 hover:text-violet-600 dark:hover:text-violet-300 hover:shadow-sm hover:scale-[1.01]'
                )}
              >
                <div className={cn(
                  "p-1 rounded-lg transition-all duration-300",
                  isActive ? "bg-white/20" : "bg-transparent group-hover:bg-violet-100/50 dark:group-hover:bg-violet-900/30"
                )}>
                  <Icon className={cn(
                    'h-5 w-5 transition-transform duration-300',
                    isActive ? 'text-white' : 'text-slate-500 dark:text-slate-400 group-hover:text-violet-600 dark:group-hover:text-violet-300 group-hover:scale-110'
                  )} />
                </div>
                <span className="tracking-wide">{item.name}</span>
                {isActive && (
                  <div className="ml-auto w-1.5 h-1.5 rounded-full bg-white shadow-[0_0_8px_rgba(255,255,255,0.8)]" />
                )}
              </Link>
            )
          })}
        </nav>
      </ScrollArea>
    </div>
  )
}