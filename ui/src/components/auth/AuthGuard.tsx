'use client'

import { useEffect, ReactNode } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { authLogger } from '@/lib/logger'
import { LoadingPage } from '@/components/common/LoadingSpinner'

interface AuthGuardProps {
  children: ReactNode
  requireAuth?: boolean
  redirectTo?: string
}

/**
 * AuthGuard Component
 * Protects routes by checking authentication status
 *
 * @param requireAuth - If true, redirects unauthenticated users to login
 * @param redirectTo - Custom redirect path (defaults to /login or /dashboard)
 */
export function AuthGuard({
  children,
  requireAuth = true,
  redirectTo,
}: AuthGuardProps) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, isLoading } = useAuth()

  useEffect(() => {
    // Wait for auth state to load
    if (isLoading) {
      return
    }

    // If authentication is required but user is not authenticated
    if (requireAuth && !user) {
      authLogger.warn('AuthGuard: Auth required but no user, redirecting to login')
      const loginPath = redirectTo || '/login'
      // Store the intended destination to redirect after login
      const returnUrl = pathname !== '/login' && pathname !== '/register'
        ? `?returnUrl=${encodeURIComponent(pathname)}`
        : ''
      router.push(`${loginPath}${returnUrl}`)
      return
    }

    // If authentication is NOT required but user IS authenticated
    // (e.g., login/register pages when already logged in)
    if (!requireAuth && user) {
      const dashboardPath = redirectTo || '/dashboard'
      router.push(dashboardPath)
      return
    }
  }, [user, isLoading, requireAuth, redirectTo, router, pathname])

  // Show loading state while checking authentication
  if (isLoading) {
    return <LoadingPage text="Loading..." />
  }

  // If requireAuth is true and no user, show loading (will redirect via useEffect)
  if (requireAuth && !user) {
    return <LoadingPage text="Redirecting to login..." />
  }

  // If requireAuth is false and user exists, show loading (will redirect via useEffect)
  if (!requireAuth && user) {
    return <LoadingPage text="Redirecting..." />
  }

  // Render children if auth requirements are met
  return <>{children}</>
}

/**
 * Protect a component with authentication
 * Usage: const ProtectedComponent = withAuth(MyComponent)
 */
export function withAuth<P extends object>(
  Component: React.ComponentType<P>,
  options?: Omit<AuthGuardProps, 'children'>
) {
  return function WithAuthComponent(props: P) {
    return (
      <AuthGuard {...options}>
        <Component {...props} />
      </AuthGuard>
    )
  }
}