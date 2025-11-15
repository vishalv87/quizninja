'use client'

import { useEffect } from 'react'
import { useAuthStore } from '@/store/authStore'
import { getSession } from '@/lib/supabase/client'
import { authLogger } from '@/lib/logger'

/**
 * SessionValidator Component
 * Validates the stored session on mount and clears it if expired
 * Used on auth pages to ensure stale sessions don't cause issues
 */
export function SessionValidator() {
  const { session, clearAuth } = useAuthStore()

  useEffect(() => {
    const validateSession = async () => {
      // If there's a stored session, validate it with Supabase
      if (session) {
        authLogger.info('[SESSION VALIDATOR] Found stored session, validating...', {
          userId: session.user?.id,
        })

        try {
          const currentSession = await getSession()

          if (!currentSession) {
            authLogger.warn('[SESSION VALIDATOR] Session expired, clearing auth state')
            clearAuth()
          } else {
            authLogger.info('[SESSION VALIDATOR] Session is valid', {
              userId: currentSession.user?.id,
            })
          }
        } catch (error) {
          authLogger.error('[SESSION VALIDATOR] Error validating session, clearing', error)
          clearAuth()
        }
      } else {
        authLogger.debug('[SESSION VALIDATOR] No stored session found')
      }
    }

    validateSession()
  }, [session, clearAuth])

  // This component doesn't render anything
  return null
}
