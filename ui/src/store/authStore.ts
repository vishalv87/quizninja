import { create } from "zustand";
import { persist } from "zustand/middleware";
import { storeLogger } from "@/lib/logger";
import type { User, Session } from "@/types/auth";

interface AuthState {
  user: User | null;
  session: Session | null;
  isLoading: boolean;
  isAuthenticated: boolean;

  setUser: (user: User | null) => void;
  setSession: (session: Session | null) => void;
  setLoading: (loading: boolean) => void;
  clearAuth: () => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      session: null,
      isLoading: true,
      isAuthenticated: false,

      setUser: (user) => {
        storeLogger.debug('Setting user', { userId: user?.id, email: user?.email });
        set({
          user,
          isAuthenticated: !!user,
        });
      },

      setSession: (session) => {
        storeLogger.info('Setting session', {
          hasSession: !!session,
          userId: session?.user?.id,
        });
        set({
          session,
          user: session?.user || null,
          isAuthenticated: !!session,
        });
      },

      setLoading: (isLoading) => {
        storeLogger.debug('Setting loading state', { isLoading });
        set({ isLoading });
      },

      clearAuth: () => {
        storeLogger.warn('Clearing auth state (invalid/expired session)');
        set({
          user: null,
          session: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      logout: () => {
        storeLogger.info('Logging out user');
        set({
          user: null,
          session: null,
          isAuthenticated: false,
        });
      },
    }),
    {
      name: "auth-storage",
      // Only persist session data, not loading or auth state
      // These should be computed fresh on each page load
      partialize: (state) => ({
        session: state.session,
      }),
      onRehydrateStorage: () => {
        storeLogger.info('Rehydrating auth state from localStorage');
        return (state, error) => {
          if (error) {
            storeLogger.error('Failed to rehydrate auth state', error);
          } else if (state?.session) {
            storeLogger.info('Auth state rehydrated', {
              hasSession: !!state.session,
              userId: state.session.user?.id,
            });
          } else {
            storeLogger.debug('No stored session found');
          }
        };
      },
    }
  )
);