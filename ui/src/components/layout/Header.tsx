'use client'

import Link from 'next/link'
import { Bell, Menu, Moon, Sun } from 'lucide-react'
import { useTheme } from 'next-themes'

import { Button } from '@/components/ui/button'
import { UserMenu } from './UserMenu'
import { useUIStore } from '@/store/uiStore'

export function Header() {
  const { theme, setTheme } = useTheme()
  const { sidebarOpen, toggleSidebar } = useUIStore()

  return (
    <header className="sticky top-0 z-40 w-full border-b bg-background">
      <div className="flex h-16 items-center px-4 sm:px-6">
        {/* Mobile Menu Toggle */}
        <Button
          variant="ghost"
          size="icon"
          className="mr-2 md:hidden"
          onClick={toggleSidebar}
        >
          <Menu className="h-5 w-5" />
          <span className="sr-only">Toggle menu</span>
        </Button>

        {/* Logo */}
        <Link href="/dashboard" className="flex items-center space-x-2">
          <span className="hidden font-bold sm:inline-block md:text-xl">
            QuizNinja
          </span>
          <span className="font-bold sm:hidden">QN</span>
        </Link>

        {/* Spacer */}
        <div className="flex-1" />

        {/* Right Side Actions */}
        <div className="flex items-center space-x-2">
          {/* Theme Toggle */}
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
          >
            <Sun className="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
            <span className="sr-only">Toggle theme</span>
          </Button>

          {/* Notifications */}
          <Button variant="ghost" size="icon" asChild>
            <Link href="/notifications">
              <Bell className="h-5 w-5" />
              <span className="sr-only">Notifications</span>
              {/* Badge for unread notifications (to be implemented) */}
              {/* <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-destructive" /> */}
            </Link>
          </Button>

          {/* User Menu */}
          <UserMenu />
        </div>
      </div>
    </header>
  )
}