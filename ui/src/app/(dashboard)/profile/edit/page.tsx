'use client'

import { Loader2 } from 'lucide-react'
import { useProfile } from '@/hooks/useProfile'
import { ProfileEditForm } from '@/components/profile/ProfileEditForm'
import { ErrorDisplay } from '@/components/common/ErrorBoundary'
import { Card, CardContent } from '@/components/ui/card'

export default function ProfileEditPage() {
  const { data: profile, isLoading, error } = useProfile()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto">
        <ErrorDisplay
          error={error as Error}
          onRetry={() => window.location.reload()}
        />
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-muted-foreground">
              Profile not found
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Edit Profile</h1>
        <p className="text-muted-foreground mt-2">
          Update your personal information
        </p>
      </div>

      <ProfileEditForm profile={profile} />
    </div>
  )
}