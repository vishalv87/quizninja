'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Menu, Moon, Sun, Search } from 'lucide-react'
import { useTheme } from 'next-themes'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { UserMenu } from './UserMenu'
import { NotificationDropdown } from '@/components/notification/NotificationDropdown'
import { GlobalSearch } from '@/components/search/GlobalSearch'
import { useUIStore } from '@/store/uiStore'

export function Header() {
  const { theme, setTheme } = useTheme()
  const { sidebarOpen, toggleSidebar } = useUIStore()
  const [searchOpen, setSearchOpen] = useState(false)

  // Keyboard shortcut for search (⌘K or Ctrl+K)
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Check for Cmd+K (Mac) or Ctrl+K (Windows/Linux)
      if ((event.metaKey || event.ctrlKey) && event.key === 'k') {
        event.preventDefault()
        setSearchOpen(true)
      }
    }

    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [])

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

        {/* Search Button */}
        <div className="flex-1 flex justify-center px-4 max-w-2xl">
          <Button
            variant="outline"
            className="w-full max-w-md justify-start text-muted-foreground hidden sm:flex"
            onClick={() => setSearchOpen(true)}
          >
            <Search className="mr-2 h-4 w-4" />
            <span>Search quizzes, users, discussions...</span>
            <kbd className="ml-auto pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
              <span className="text-xs">⌘</span>K
            </kbd>
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="sm:hidden"
            onClick={() => setSearchOpen(true)}
          >
            <Search className="h-5 w-5" />
            <span className="sr-only">Search</span>
          </Button>
        </div>

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
          <NotificationDropdown />

          {/* User Menu */}
          <UserMenu />
        </div>
      </div>

      {/* Global Search Dialog */}
      <Dialog open={searchOpen} onOpenChange={setSearchOpen}>
        <DialogContent className="max-w-5xl max-h-[85vh] overflow-y-auto p-0">
          <DialogHeader className="p-6 pb-0">
            <DialogTitle className="sr-only">Search QuizNinja</DialogTitle>
          </DialogHeader>
          <div className="p-6 pt-4">
            <GlobalSearch onClose={() => setSearchOpen(false)} />
          </div>
        </DialogContent>
      </Dialog>
    </header>
  )
}