"use client";

import { useEffect } from "react";
import { useAuthStore } from "@/store/authStore";
import { onAuthStateChange } from "@/lib/supabase/client";

export function useAuth() {
  const { user, session, isLoading, isAuthenticated, setSession, setLoading } = useAuthStore();

  useEffect(() => {
    // Listen to auth state changes
    const { data: authListener } = onAuthStateChange((_event, session) => {
      setSession(session);
      setLoading(false);
    });

    return () => {
      authListener.subscription.unsubscribe();
    };
  }, [setSession, setLoading]);

  return {
    user,
    session,
    isLoading,
    isAuthenticated,
  };
}