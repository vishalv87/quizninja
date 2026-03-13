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
    // Check the actual session from Supabase on mount
    const initializeAuth = async () => {
      try {
        const currentSession = await getSession();

        if (currentSession) {
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
    const { data: authListener } = onAuthStateChange((event, session) => {
      if (session) {
        setSession(session);
      } else {
        clearAuth();
      }
    });

    return () => {
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