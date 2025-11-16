'use client';

import { ReactNode } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { SessionValidator } from '@/components/auth/SessionValidator'

/**
 * Auth Layout
 *
 * NOTE: We don't use AuthGuard here because the middleware already handles
 * redirecting authenticated users away from auth pages. Using AuthGuard here
 * causes a redirect loop due to client/server state mismatch.
 *
 * We use SessionValidator to detect and clear expired sessions from localStorage.
 */
export default function AuthLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const isFullWidthPage = pathname === '/welcome' || pathname === '/preferences';

  return (
    <div className="min-h-screen flex flex-col">
      {/* Validates and clears expired sessions */}
      <SessionValidator />

      {/* Header - Hidden on full-width pages (they have their own) */}
      {!isFullWidthPage && (
        <header className="border-b">
          <div className="container flex h-16 items-center px-4">
            <Link href="/" className="flex items-center space-x-2">
              <span className="text-2xl font-bold text-primary">QuizNinja</span>
            </Link>
          </div>
        </header>
      )}

      {/* Main Content */}
      <main className={isFullWidthPage ? "flex-1" : "flex-1 flex items-center justify-center p-4"}>
        <div className={isFullWidthPage ? "w-full" : "w-full max-w-md"}>
          {children}
        </div>
      </main>

      {/* Footer - Hidden on full-width pages (they have their own) */}
      {!isFullWidthPage && (
        <footer className="border-t py-6 md:py-0">
          <div className="container flex h-16 items-center justify-center px-4">
            <p className="text-sm text-muted-foreground">
              © {new Date().getFullYear()} QuizNinja. All rights reserved.
            </p>
          </div>
        </footer>
      )}
    </div>
  )
}