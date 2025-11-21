import { ReactNode } from 'react'
import { AuthGuard } from '@/components/auth/AuthGuard'
import { Header } from '@/components/layout/Header'
import { Sidebar } from '@/components/layout/Sidebar'
import { MobileNav } from '@/components/layout/MobileNav'
import { DashboardNotificationListener } from '@/components/notification/DashboardNotificationListener'

export const metadata = {
  title: 'Dashboard | QuizNinja',
  description: 'Your QuizNinja dashboard',
}

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <AuthGuard requireAuth={true}>
      {/* Notification Toast Listener - Shows toasts when new notifications arrive */}
      <DashboardNotificationListener />

      <div className="flex min-h-screen flex-col">
        {/* Header */}
        <Header />

        <div className="flex flex-1">
          {/* Desktop Sidebar */}
          <Sidebar />

          {/* Mobile Navigation */}
          <MobileNav />

          {/* Main Content */}
          <main className="flex-1 overflow-y-auto">
            <div className="container mx-auto px-4 py-6 sm:px-6">
              {children}
            </div>
          </main>
        </div>
      </div>
    </AuthGuard>
  )
}