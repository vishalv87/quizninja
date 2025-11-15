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
  X,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { useUIStore } from '@/store/uiStore'

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

export function MobileNav() {
  const pathname = usePathname()
  const { sidebarOpen, toggleSidebar } = useUIStore()

  const handleNavigation = () => {
    toggleSidebar()
  }

  return (
    <Sheet open={sidebarOpen} onOpenChange={toggleSidebar}>
      <SheetContent side="left" className="w-64 p-0">
        <SheetHeader className="border-b p-4">
          <SheetTitle>Navigation</SheetTitle>
        </SheetHeader>
        <ScrollArea className="h-[calc(100vh-5rem)]">
          <nav className="space-y-1 p-3">
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
                  onClick={handleNavigation}
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
      </SheetContent>
    </Sheet>
  )
}