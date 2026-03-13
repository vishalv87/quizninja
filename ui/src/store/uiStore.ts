import { create } from "zustand";
import { persist } from "zustand/middleware";

interface UIState {
  sidebarOpen: boolean;
  notificationCount: number;

  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setNotificationCount: (count: number) => void;
  incrementNotificationCount: () => void;
  decrementNotificationCount: () => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      sidebarOpen: true,
      notificationCount: 0,

      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),

      setSidebarOpen: (sidebarOpen) => set({ sidebarOpen }),

      setNotificationCount: (notificationCount) => set({ notificationCount }),

      incrementNotificationCount: () =>
        set((state) => ({ notificationCount: state.notificationCount + 1 })),

      decrementNotificationCount: () =>
        set((state) => ({
          notificationCount: Math.max(0, state.notificationCount - 1),
        })),
    }),
    {
      name: "ui-storage",
    }
  )
);