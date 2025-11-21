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
        set({
          user,
          isAuthenticated: !!user,
        });
      },

      setSession: (session) => {
        set({
          session,
          user: session?.user || null,
          isAuthenticated: !!session,
        });
      },

      setLoading: (isLoading) => {
        set({ isLoading });
      },

      clearAuth: () => {
        storeLogger.warn("Clearing auth state (invalid/expired session)");
        set({
          user: null,
          session: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      logout: () => {
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
        return (state, error) => {
          if (error) {
            storeLogger.error("Failed to rehydrate auth state", error);
          }
        };
      },
    }
  )
);
