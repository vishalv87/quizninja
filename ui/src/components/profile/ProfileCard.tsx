'use client'

import { Edit, Mail, Calendar, Trophy } from 'lucide-react'
import Link from 'next/link'
import { format } from 'date-fns'

import type { Profile } from '@/types/auth'
import type { UserStats } from '@/types/user'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

interface ProfileCardProps {
  profile: Profile
  stats?: UserStats
  isLoadingStats?: boolean
}

export function ProfileCard({ profile, stats, isLoadingStats }: ProfileCardProps) {
  const getInitials = (name?: string, email?: string) => {
    if (name) {
      return name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    }
    if (email) {
      return email[0].toUpperCase()
    }
    return 'U'
  }

  const initials = getInitials(profile.name, profile.email)

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-4">
            <Avatar className="h-20 w-20">
              <AvatarImage src={profile.avatar_url} alt={profile.name || profile.email} />
              <AvatarFallback className="text-2xl">{initials}</AvatarFallback>
            </Avatar>
            <div>
              <CardTitle className="text-2xl">{profile.name || 'User'}</CardTitle>
              <CardDescription className="flex items-center mt-1">
                <Mail className="mr-1 h-3 w-3" />
                {profile.email}
              </CardDescription>
              {profile.created_at && (
                <CardDescription className="flex items-center mt-1">
                  <Calendar className="mr-1 h-3 w-3" />
                  Joined {format(new Date(profile.created_at), 'MMMM yyyy')}
                </CardDescription>
              )}
            </div>
          </div>
          <Button variant="outline" asChild>
            <Link href="/profile/edit">
              <Edit className="mr-2 h-4 w-4" />
              Edit Profile
            </Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Bio */}
        {profile.bio && (
          <div>
            <h3 className="text-sm font-medium mb-1">Bio</h3>
            <p className="text-sm text-muted-foreground">{profile.bio}</p>
          </div>
        )}

        {/* Stats */}
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4 pt-4 border-t">
          <div className="text-center">
            <p className="text-2xl font-bold">
              {isLoadingStats ? '...' : stats?.total_quizzes_completed?.toLocaleString() || '0'}
            </p>
            <p className="text-xs text-muted-foreground">Quizzes Taken</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold">
              {isLoadingStats ? '...' : stats?.total_points?.toLocaleString() || '0'}
            </p>
            <p className="text-xs text-muted-foreground">Total Points</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold">
              {isLoadingStats ? '...' : stats?.achievements_unlocked?.toLocaleString() || '0'}
            </p>
            <p className="text-xs text-muted-foreground">Achievements</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold">
              {isLoadingStats ? '...' : stats?.rank ? `#${stats.rank}` : '—'}
            </p>
            <p className="text-xs text-muted-foreground">Rank</p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}