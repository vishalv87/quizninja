'use client'

import { useParams, useRouter } from 'next/navigation'
import { Loader2, ArrowLeft, UserPlus, UserMinus, Swords, Shield } from 'lucide-react'
import { useUserProfile } from '@/hooks/useProfile'
import { useSendFriendRequest, useAcceptFriendRequest, useDeclineFriendRequest } from '@/hooks/useFriendRequests'
import { useRemoveFriend } from '@/hooks/useFriends'
import { UserProfileCard } from '@/components/profile/UserProfileCard'
import { AchievementBadges } from '@/components/profile/AchievementBadges'
import { DetailedStatistics } from '@/components/profile/DetailedStatistics'
import { ErrorDisplay } from '@/components/common/ErrorBoundary'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import Link from 'next/link'
import { toast } from 'sonner'

export default function UserProfilePage() {
  const params = useParams()
  const router = useRouter()
  const userId = params.userId as string

  const { data: userProfile, isLoading, error } = useUserProfile(userId)
  const sendFriendRequest = useSendFriendRequest()
  const removeFriend = useRemoveFriend()

  const handleSendFriendRequest = () => {
    sendFriendRequest.mutate(userId, {
      onSuccess: () => {
        toast.success('Friend request sent!')
      },
    })
  }

  const handleRemoveFriend = () => {
    removeFriend.mutate(userId, {
      onSuccess: () => {
        toast.success('Friend removed')
      },
    })
  }

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
        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="mb-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Go Back
        </Button>
        <ErrorDisplay
          error={error as Error}
          onRetry={() => window.location.reload()}
        />
      </div>
    )
  }

  if (!userProfile) {
    return (
      <div className="max-w-2xl mx-auto">
        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="mb-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Go Back
        </Button>
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-muted-foreground">
              User profile not found
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  // Determine if profile is private
  const isPrivate = userProfile.preferences?.profile_visibility === 'private'
  const isFriendsOnly = userProfile.preferences?.profile_visibility === 'friends_only'
  const canView = !isPrivate && (!isFriendsOnly || userProfile.is_friend)

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Back Button */}
      <Button
        variant="ghost"
        onClick={() => router.back()}
      >
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back
      </Button>

      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          {userProfile.full_name || userProfile.name || 'User'}'s Profile
        </h1>
        <p className="text-muted-foreground mt-2">
          View profile information and statistics
        </p>
      </div>

      {/* Privacy Message */}
      {!canView && (
        <Card>
          <CardContent className="pt-6 text-center space-y-4">
            <Shield className="h-12 w-12 mx-auto text-muted-foreground" />
            <div>
              <h3 className="text-lg font-semibold">Private Profile</h3>
              <p className="text-sm text-muted-foreground mt-2">
                {isPrivate
                  ? 'This user has set their profile to private.'
                  : 'This profile is only visible to friends.'}
              </p>
            </div>
            {!userProfile.is_friend && userProfile.friend_request_status === 'none' && (
              <Button onClick={handleSendFriendRequest} disabled={sendFriendRequest.isPending}>
                <UserPlus className="mr-2 h-4 w-4" />
                {sendFriendRequest.isPending ? 'Sending...' : 'Send Friend Request'}
              </Button>
            )}
          </CardContent>
        </Card>
      )}

      {/* Profile Content - Only show if user can view */}
      {canView && (
        <>
          {/* User Profile Card */}
          <UserProfileCard
            profile={userProfile}
            stats={userProfile.stats}
            onSendFriendRequest={handleSendFriendRequest}
            onRemoveFriend={handleRemoveFriend}
            isSendingRequest={sendFriendRequest.isPending}
            isRemovingFriend={removeFriend.isPending}
          />

          {/* Detailed Statistics - Only show if allowed by privacy settings */}
          {userProfile.preferences?.show_stats !== false && userProfile.stats && (
            <DetailedStatistics stats={userProfile.stats} isLoading={false} />
          )}

          {/* Achievements - Only show if allowed by privacy settings */}
          {/* Note: AchievementBadges will be updated in Phase 7 to support viewing other users' achievements */}
          {userProfile.preferences?.show_achievements !== false && (
            <AchievementBadges />
          )}
        </>
      )}
    </div>
  )
}
