'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'

import { profileUpdateSchema, type ProfileUpdateFormData } from '@/schemas/auth'
import { useUpdateProfile } from '@/hooks/useProfile'
import type { Profile } from '@/types/auth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

interface ProfileEditFormProps {
  profile: Profile
}

export function ProfileEditForm({ profile }: ProfileEditFormProps) {
  const router = useRouter()
  const { mutate: updateProfile, isPending } = useUpdateProfile()

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<ProfileUpdateFormData>({
    resolver: zodResolver(profileUpdateSchema),
    defaultValues: {
      name: profile.name || '',
      email: profile.email || '',
      avatar_url: profile.avatar_url || '',
      bio: profile.bio || '',
    },
  })

  // Update form when profile changes
  useEffect(() => {
    reset({
      name: profile.name || '',
      email: profile.email || '',
      avatar_url: profile.avatar_url || '',
      bio: profile.bio || '',
    })
  }, [profile, reset])

  const onSubmit = async (data: ProfileUpdateFormData) => {
    // Filter out empty optional fields
    const updateData: Partial<Profile> = {}
    if (data.name) updateData.name = data.name
    if (data.email) updateData.email = data.email
    if (data.avatar_url) updateData.avatar_url = data.avatar_url
    if (data.bio) updateData.bio = data.bio

    updateProfile(updateData, {
      onSuccess: () => {
        // Redirect to profile page after successful update
        router.push('/profile')
      },
    })
  }

  const handleCancel = () => {
    router.back()
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Edit Profile</CardTitle>
        <CardDescription>
          Update your personal information and profile details
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Name Field */}
          <div className="space-y-2">
            <Label htmlFor="name">Full Name</Label>
            <Input
              id="name"
              type="text"
              placeholder="John Doe"
              {...register('name')}
              disabled={isPending}
              className={errors.name ? 'border-destructive' : ''}
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          {/* Email Field */}
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="john.doe@example.com"
              {...register('email')}
              disabled={isPending}
              className={errors.email ? 'border-destructive' : ''}
            />
            {errors.email && (
              <p className="text-sm text-destructive">{errors.email.message}</p>
            )}
            <p className="text-xs text-muted-foreground">
              Changing your email may require verification
            </p>
          </div>

          {/* Avatar URL Field */}
          <div className="space-y-2">
            <Label htmlFor="avatar_url">Avatar URL</Label>
            <Input
              id="avatar_url"
              type="url"
              placeholder="https://example.com/avatar.jpg"
              {...register('avatar_url')}
              disabled={isPending}
              className={errors.avatar_url ? 'border-destructive' : ''}
            />
            {errors.avatar_url && (
              <p className="text-sm text-destructive">{errors.avatar_url.message}</p>
            )}
            <p className="text-xs text-muted-foreground">
              Provide a URL to your profile picture
            </p>
          </div>

          {/* Bio Field */}
          <div className="space-y-2">
            <Label htmlFor="bio">Bio</Label>
            <Textarea
              id="bio"
              placeholder="Tell us about yourself..."
              rows={4}
              {...register('bio')}
              disabled={isPending}
              className={errors.bio ? 'border-destructive' : ''}
            />
            {errors.bio && (
              <p className="text-sm text-destructive">{errors.bio.message}</p>
            )}
            <p className="text-xs text-muted-foreground">
              Maximum 500 characters
            </p>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-2">
            <Button type="submit" disabled={isPending}>
              {isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Saving...
                </>
              ) : (
                'Save Changes'
              )}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={handleCancel}
              disabled={isPending}
            >
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}