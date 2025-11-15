"use client";

import { useEffect } from "react";
import { useAuthStore } from "@/store/authStore";
import { getSession, onAuthStateChange } from "@/lib/supabase/client";
import { authLogger } from "@/lib/logger";

export function useAuth() {
  const {
    user,
    session,
    isLoading,
    isAuthenticated,
    setSession,
    setLoading,
    clearAuth,
  } = useAuthStore();

  useEffect(() => {
    authLogger.info('useAuth: Initializing auth hook');

    // Check the actual session from Supabase on mount
    const initializeAuth = async () => {
      try {
        authLogger.debug('useAuth: Fetching current session from Supabase');
        const currentSession = await getSession();

        if (currentSession) {
          authLogger.info('useAuth: Valid session found', {
            userId: currentSession.user?.id,
          });
          setSession(currentSession as any);
        } else {
          authLogger.warn('useAuth: No valid session found, clearing auth state');
          clearAuth();
        }
      } catch (error) {
        authLogger.error('useAuth: Error fetching session', error);
        // If there's an error, clear the auth state to be safe
        clearAuth();
      } finally {
        setLoading(false);
      }
    };

    // Initialize auth state
    initializeAuth();

    // Listen to auth state changes
    authLogger.debug('useAuth: Setting up auth state change listener');
    const { data: authListener } = onAuthStateChange((event, session) => {
      authLogger.info('useAuth: Auth state changed', { event, hasSession: !!session });

      if (session) {
        setSession(session);
      } else {
        clearAuth();
      }
    });

    return () => {
      authLogger.debug('useAuth: Cleaning up auth listener');
      authListener.subscription.unsubscribe();
    };
  }, [setSession, setLoading, clearAuth]);

  return {
    user,
    session,
    isLoading,
    isAuthenticated,
  };
}