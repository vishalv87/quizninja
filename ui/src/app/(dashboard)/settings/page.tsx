'use client'

import { Settings as SettingsIcon } from 'lucide-react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AccountSettingsForm } from '@/components/settings/AccountSettingsForm'
import { PreferencesForm } from '@/components/settings/PreferencesForm'

export default function SettingsPage() {
  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
          <SettingsIcon className="h-8 w-8" />
          Settings
        </h1>
        <p className="text-muted-foreground mt-2">
          Manage your account, preferences, and customize your quiz experience
        </p>
      </div>

      <Tabs defaultValue="account" className="w-full">
        <TabsList className="grid w-full max-w-md grid-cols-2">
          <TabsTrigger value="account">Account</TabsTrigger>
          <TabsTrigger value="preferences">Preferences</TabsTrigger>
        </TabsList>

        <TabsContent value="account" className="mt-6">
          <AccountSettingsForm />
        </TabsContent>

        <TabsContent value="preferences" className="mt-6">
          <PreferencesForm />
        </TabsContent>
      </Tabs>
    </div>
  )
}
