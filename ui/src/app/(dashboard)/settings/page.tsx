'use client'

import { Settings as SettingsIcon } from 'lucide-react'
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
          Manage your preferences and customize your quiz experience
        </p>
      </div>

      <PreferencesForm />
    </div>
  )
}
