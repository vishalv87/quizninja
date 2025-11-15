import { ReactNode } from 'react'
import Link from 'next/link'
import { SessionValidator } from '@/components/auth/SessionValidator'

export const metadata = {
  title: 'Authentication | QuizNinja',
  description: 'Sign in or create an account to access QuizNinja',
}

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
  return (
    <div className="min-h-screen flex flex-col">
      {/* Validates and clears expired sessions */}
      <SessionValidator />
      {/* Header */}
      <header className="border-b">
        <div className="container flex h-16 items-center px-4">
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-bold text-primary">QuizNinja</span>
          </Link>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          {children}
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t py-6 md:py-0">
        <div className="container flex h-16 items-center justify-center px-4">
          <p className="text-sm text-muted-foreground">
            © {new Date().getFullYear()} QuizNinja. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  )
}