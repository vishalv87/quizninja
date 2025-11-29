'use client'

import { Settings as SettingsIcon, User, Sliders } from 'lucide-react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AccountSettingsForm } from '@/components/settings/AccountSettingsForm'
import { PreferencesForm } from '@/components/settings/PreferencesForm'
import { PageHero } from '@/components/common/PageHero'
import { GlassCard } from '@/components/common/GlassCard'

export default function SettingsPage() {
  return (
    <div className="space-y-10 pb-10">
      {/* Hero Section */}
      <PageHero
        title="Settings"
        icon="⚙️"
        description="Manage your account, preferences, and customize your quiz experience to suit your style."
      />

      <div className="container px-0 md:px-4">
        <GlassCard padding="none" rounded="2xl">
          <Tabs defaultValue="account" className="w-full">
            <div className="p-6 border-b border-white/10">
              <TabsList className="grid w-full max-w-md grid-cols-2 bg-white/60 dark:bg-black/40 backdrop-blur-md border border-white/20 dark:border-white/10 p-1 rounded-xl shadow-sm">
                <TabsTrigger
                  value="account"
                  className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  <User className="h-4 w-4" />
                  Account
                </TabsTrigger>
                <TabsTrigger
                  value="preferences"
                  className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  <Sliders className="h-4 w-4" />
                  Preferences
                </TabsTrigger>
              </TabsList>
            </div>

            <TabsContent value="account" className="mt-0">
              <div className="p-6 border-b border-white/10">
                <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
                  <span className="bg-gradient-to-br from-blue-500 to-cyan-500 text-white p-1.5 rounded-lg shadow-sm">
                    <User className="h-4 w-4" />
                  </span>
                  Account Settings
                </h2>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  Manage your account information and security settings
                </p>
              </div>
              <div className="p-6">
                <AccountSettingsForm />
              </div>
            </TabsContent>

            <TabsContent value="preferences" className="mt-0">
              <div className="p-6 border-b border-white/10">
                <h2 className="text-xl font-bold tracking-tight flex items-center gap-2 text-slate-800 dark:text-slate-100">
                  <span className="bg-gradient-to-br from-violet-500 to-purple-600 text-white p-1.5 rounded-lg shadow-sm">
                    <Sliders className="h-4 w-4" />
                  </span>
                  Preferences
                </h2>
                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                  Customize your quiz experience and notification settings
                </p>
              </div>
              <div className="p-6">
                <PreferencesForm />
              </div>
            </TabsContent>
          </Tabs>
        </GlassCard>
      </div>
    </div>
  )
}
