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
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
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
    <div className="hidden border-r bg-background md:block w-64">
      <ScrollArea className="h-[calc(100vh-4rem)] py-4">
        <nav className="space-y-1 px-3">
          {navigation.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/')
            const Icon = item.icon

            return (
              <Button
                key={item.name}
                variant={isActive ? 'secondary' : 'ghost'}
                className={cn(
                  'w-full justify-start',
                  isActive && 'bg-secondary'
                )}
                asChild
              >
                <Link href={item.href}>
                  <Icon className="mr-2 h-4 w-4" />
                  {item.name}
                </Link>
              </Button>
            )
          })}
        </nav>
      </ScrollArea>
    </div>
  )
}