'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Menu, Moon, Sun, Search, Bell } from 'lucide-react'
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
  const [scrolled, setScrolled] = useState(false)

  // Handle scroll effect
  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10)
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

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
    <header 
      className={`sticky top-0 z-40 w-full transition-all duration-300 ${
        scrolled 
          ? 'bg-white/70 dark:bg-black/70 backdrop-blur-xl border-b border-white/20 dark:border-white/10 shadow-sm' 
          : 'bg-transparent border-b border-transparent'
      }`}
    >
      <div className="grid grid-cols-[auto_1fr_auto] h-20 items-center px-4 sm:px-6 gap-4">
        {/* Left Section: Mobile Menu + Logo */}
        <div className="flex items-center">
          <Button
            variant="ghost"
            size="icon"
            className="mr-2 md:hidden rounded-xl hover:bg-white/20 transition-all duration-300"
            onClick={toggleSidebar}
          >
            <Menu className="h-5 w-5" />
            <span className="sr-only">Toggle menu</span>
          </Button>

          <Link href="/dashboard" className="flex items-center space-x-2 group relative">
            <div className="absolute -inset-2 bg-gradient-to-r from-violet-600/20 to-indigo-600/20 rounded-xl blur-lg opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
            <span className="relative z-10 hidden font-bold sm:inline-block md:text-2xl bg-gradient-to-r from-violet-600 to-indigo-600 bg-clip-text text-transparent tracking-tight drop-shadow-sm">
              QuizNinja
            </span>
            <span className="relative z-10 font-bold sm:hidden bg-gradient-to-r from-violet-600 to-indigo-600 bg-clip-text text-transparent text-xl">QN</span>
          </Link>
        </div>

        {/* Search Button */}
        <div className="flex justify-center max-w-2xl mx-auto w-full">
          <Button
            variant="outline"
            className="w-full justify-start text-muted-foreground hidden sm:flex rounded-2xl border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/20 shadow-sm hover:shadow-md hover:bg-white/60 dark:hover:bg-black/40 hover:border-violet-200/50 transition-all duration-300 h-11 backdrop-blur-sm group"
            onClick={() => setSearchOpen(true)}
          >
            <Search className="mr-3 h-4 w-4 text-violet-500 group-hover:scale-110 transition-transform" />
            <span className="text-slate-600 dark:text-slate-400">Search quizzes, users, discussions...</span>
            <kbd className="ml-auto pointer-events-none inline-flex h-6 select-none items-center gap-1 rounded-lg border border-white/20 bg-white/50 px-2 font-mono text-[10px] font-medium text-slate-500">
              <span className="text-xs">⌘</span>K
            </kbd>
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="sm:hidden rounded-xl hover:bg-white/20 transition-all duration-300"
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
            className="rounded-xl hover:bg-white/40 dark:hover:bg-white/10 transition-all duration-300 hover:scale-105 hover:text-violet-600"
            onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
          >
            <Sun className="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
            <span className="sr-only">Toggle theme</span>
          </Button>

          {/* Notifications */}
          <div className="relative">
             <NotificationDropdown />
          </div>

          {/* User Menu */}
          <div className="pl-2 border-l border-gray-200/30 dark:border-gray-700/30">
            <UserMenu />
          </div>
        </div>
      </div>

      {/* Global Search Dialog */}
      <Dialog open={searchOpen} onOpenChange={setSearchOpen}>
        <DialogContent className="max-w-5xl max-h-[85vh] overflow-y-auto p-0 border-0 bg-transparent shadow-2xl">
          <div className="bg-white/90 dark:bg-slate-950/90 backdrop-blur-xl rounded-lg border border-white/20 overflow-hidden">
             <DialogHeader className="p-6 pb-0">
              <DialogTitle className="sr-only">Search QuizNinja</DialogTitle>
            </DialogHeader>
            <div className="p-6 pt-4">
              <GlobalSearch onClose={() => setSearchOpen(false)} />
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </header>
  )
}