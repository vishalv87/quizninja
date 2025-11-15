'use client'

import { Mail, Calendar, UserPlus, UserMinus, Swords, UserCheck, Clock } from 'lucide-react'
import Link from 'next/link'
import { format } from 'date-fns'

import type { UserProfile, UserStats } from '@/types/user'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

interface UserProfileCardProps {
  profile: UserProfile
  stats?: UserStats | null
  onSendFriendRequest?: () => void
  onRemoveFriend?: () => void
  isSendingRequest?: boolean
  isRemovingFriend?: boolean
}

export function UserProfileCard({
  profile,
  stats,
  onSendFriendRequest,
  onRemoveFriend,
  isSendingRequest,
  isRemovingFriend,
}: UserProfileCardProps) {
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

  const displayName = profile.full_name || profile.name || 'User'
  const initials = getInitials(displayName, profile.email)

  // Determine friendship status
  const isFriend = profile.is_friend
  const friendRequestStatus = profile.friend_request_status

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between flex-col sm:flex-row gap-4">
          <div className="flex items-center space-x-4">
            <Avatar className="h-20 w-20">
              <AvatarImage src={profile.avatar_url} alt={displayName} />
              <AvatarFallback className="text-2xl">{initials}</AvatarFallback>
            </Avatar>
            <div>
              <CardTitle className="text-2xl">{displayName}</CardTitle>
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

          {/* Friend Status & Actions */}
          <div className="flex gap-2 flex-wrap">
            {isFriend ? (
              <>
                <Badge variant="secondary" className="h-fit">
                  <UserCheck className="mr-1 h-3 w-3" />
                  Friends
                </Badge>
                <Button variant="default" size="sm" asChild>
                  <Link href={`/challenges/create?friendId=${profile.user_id}`}>
                    <Swords className="mr-2 h-4 w-4" />
                    Challenge
                  </Link>
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onRemoveFriend}
                  disabled={isRemovingFriend}
                >
                  <UserMinus className="mr-2 h-4 w-4" />
                  {isRemovingFriend ? 'Removing...' : 'Remove Friend'}
                </Button>
              </>
            ) : friendRequestStatus === 'pending_sent' ? (
              <Badge variant="outline" className="h-fit">
                <Clock className="mr-1 h-3 w-3" />
                Friend Request Sent
              </Badge>
            ) : friendRequestStatus === 'pending_received' ? (
              <Badge variant="outline" className="h-fit">
                <Clock className="mr-1 h-3 w-3" />
                Pending Friend Request
              </Badge>
            ) : (
              <Button
                variant="default"
                size="sm"
                onClick={onSendFriendRequest}
                disabled={isSendingRequest}
              >
                <UserPlus className="mr-2 h-4 w-4" />
                {isSendingRequest ? 'Sending...' : 'Add Friend'}
              </Button>
            )}
          </div>
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

        {/* Stats - Only show if available (respecting privacy) */}
        {stats && (
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4 pt-4 border-t">
            <div className="text-center">
              <p className="text-2xl font-bold">
                {stats.total_quizzes_completed?.toLocaleString() || '0'}
              </p>
              <p className="text-xs text-muted-foreground">Quizzes Taken</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold">
                {stats.total_points?.toLocaleString() || '0'}
              </p>
              <p className="text-xs text-muted-foreground">Total Points</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold">
                {stats.achievements_unlocked?.toLocaleString() || '0'}
              </p>
              <p className="text-xs text-muted-foreground">Achievements</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold">
                {stats.rank ? `#${stats.rank}` : '—'}
              </p>
              <p className="text-xs text-muted-foreground">Rank</p>
            </div>
          </div>
        )}

        {/* Message if stats are hidden */}
        {!stats && (
          <div className="pt-4 border-t text-center">
            <p className="text-sm text-muted-foreground">
              This user has chosen to keep their statistics private.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
