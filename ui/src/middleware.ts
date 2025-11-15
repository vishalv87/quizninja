import { createMiddlewareClient } from '@supabase/auth-helpers-nextjs'
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

/**
 * Next.js Middleware for Route Protection
 * Validates user sessions and redirects based on authentication status
 */

// Routes that require authentication
const protectedRoutes = [
  '/dashboard',
  '/quizzes',
  '/profile',
  '/friends',
  '/challenges',
  '/achievements',
  '/leaderboard',
  '/discussions',
  '/notifications',
  '/settings',
]

// Routes that should redirect to dashboard if user is already authenticated
const authRoutes = ['/login', '/register']

export async function middleware(req: NextRequest) {
  const res = NextResponse.next()
  const supabase = createMiddlewareClient({ req, res })

  // Get current path
  const path = req.nextUrl.pathname

  console.log(`[MIDDLEWARE] Processing request: ${path}`)

  try {
    // Get the session
    const {
      data: { session },
    } = await supabase.auth.getSession()

    const isAuthenticated = !!session

    console.log(`[MIDDLEWARE] Auth status:`, {
      path,
      isAuthenticated,
      userId: session?.user?.id,
    })

    // Check if the current path is a protected route
    const isProtectedRoute = protectedRoutes.some((route) =>
      path.startsWith(route)
    )

    // Check if the current path is an auth route
    const isAuthRoute = authRoutes.some((route) => path.startsWith(route))

    console.log(`[MIDDLEWARE] Route type:`, {
      isProtectedRoute,
      isAuthRoute,
    })

    // If accessing a protected route without authentication
    if (isProtectedRoute && !isAuthenticated) {
      console.log(`[MIDDLEWARE] Redirecting to login (protected route, not authenticated)`)
      const loginUrl = new URL('/login', req.url)
      // Store the intended destination to redirect after login
      if (path !== '/login') {
        loginUrl.searchParams.set('returnUrl', path)
      }
      return NextResponse.redirect(loginUrl)
    }

    // If accessing auth routes while already authenticated
    if (isAuthRoute && isAuthenticated) {
      console.log(`[MIDDLEWARE] Redirecting to dashboard (auth route, already authenticated)`)
      // Check if there's a return URL
      const returnUrl = req.nextUrl.searchParams.get('returnUrl')
      if (returnUrl) {
        return NextResponse.redirect(new URL(returnUrl, req.url))
      }
      // Otherwise redirect to dashboard
      return NextResponse.redirect(new URL('/dashboard', req.url))
    }

    console.log(`[MIDDLEWARE] Allowing request to proceed`)
    // Allow the request to proceed
    return res
  } catch (error) {
    console.error('[MIDDLEWARE] Error:', error)
    // On error, allow the request to proceed
    // The AuthGuard component will handle the redirect
    return res
  }
}

// Configure which routes to run middleware on
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     * - api routes (handled separately)
     */
    '/((?!_next/static|_next/image|favicon.ico|public|api).*)',
  ],
}