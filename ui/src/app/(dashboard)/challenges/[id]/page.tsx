"use client";

import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ChallengeResults } from "@/components/challenge/ChallengeResults";
import { useChallenge } from "@/hooks/useChallenges";
import { useAcceptChallenge, useDeclineChallenge } from "@/hooks/useChallengeActions";
import { useAuth } from "@/hooks/useAuth";
import {
  Swords,
  ArrowLeft,
  Clock,
  Trophy,
  Play,
  AlertCircle,
  CheckCircle2,
  XCircle,
  Timer,
} from "lucide-react";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";

export default function ChallengeDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuth();
  const challengeId = params.id as string;

  const { data: challenge, isLoading, error } = useChallenge(challengeId);
  const acceptMutation = useAcceptChallenge();
  const declineMutation = useDeclineChallenge();

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <Skeleton className="h-8 w-32 mb-4" />
        <Skeleton className="h-12 w-64 mb-8" />
        <div className="space-y-4">
          <Skeleton className="h-64 w-full" />
          <Skeleton className="h-48 w-full" />
        </div>
      </div>
    );
  }

  // Error state
  if (error || !challenge) {
    return (
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            {error?.message || "Challenge not found"}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  // Determine if current user is challenger or opponent
  const isChallenger = user?.id === challenge.challenger_id;
  const opponent = isChallenger ? challenge.opponent : challenge.challenger;
  const userHasCompleted = isChallenger
    ? !!challenge.challenger_attempt_id
    : !!challenge.opponent_attempt_id;
  const opponentHasCompleted = isChallenger
    ? !!challenge.opponent_attempt_id
    : !!challenge.challenger_attempt_id;

  // Get status badge
  const getStatusBadge = () => {
    switch (challenge.status) {
      case "pending":
        return (
          <Badge variant="secondary" className="text-base">
            <Clock className="mr-1 h-4 w-4" />
            Pending
          </Badge>
        );
      case "accepted":
        return (
          <Badge variant="default" className="text-base">
            <Swords className="mr-1 h-4 w-4" />
            Active
          </Badge>
        );
      case "completed":
        return (
          <Badge variant="outline" className="text-base">
            <CheckCircle2 className="mr-1 h-4 w-4" />
            Completed
          </Badge>
        );
      case "declined":
        return (
          <Badge variant="destructive" className="text-base">
            <XCircle className="mr-1 h-4 w-4" />
            Declined
          </Badge>
        );
      case "expired":
        return (
          <Badge variant="secondary" className="text-base">
            <Timer className="mr-1 h-4 w-4" />
            Expired
          </Badge>
        );
      default:
        return null;
    }
  };

  const handleAccept = () => {
    acceptMutation.mutate(challengeId, {
      onSuccess: () => {
        // Optionally redirect or show success message
      },
    });
  };

  const handleDecline = () => {
    declineMutation.mutate(challengeId, {
      onSuccess: () => {
        router.push("/challenges");
      },
    });
  };

  const handleTakeQuiz = () => {
    router.push(`/quizzes/${challenge.quiz_id}/take?challengeId=${challengeId}`);
  };

  return (
    <div className="container mx-auto py-8 px-4 max-w-4xl">
      {/* Back Button */}
      <Button variant="ghost" onClick={() => router.back()} className="mb-6">
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to Challenges
      </Button>

      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h1 className="text-4xl font-bold mb-2">Challenge Details</h1>
            <p className="text-muted-foreground">
              {isChallenger ? `Your challenge to ${opponent.full_name}` : `Challenge from ${opponent.full_name}`}
            </p>
          </div>
          {getStatusBadge()}
        </div>
      </div>

      {/* Completed Challenge - Show Results */}
      {challenge.status === "completed" && (
        <ChallengeResults challenge={challenge} />
      )}

      {/* Active or Pending Challenge */}
      {(challenge.status === "pending" || challenge.status === "accepted") && (
        <div className="space-y-6">
          {/* Challenge Info Card */}
          <Card>
            <CardHeader>
              <CardTitle>Challenge Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Participants */}
              <div className="grid grid-cols-2 gap-4">
                <div className="flex flex-col items-center text-center space-y-2">
                  <Avatar className="h-20 w-20">
                    <AvatarImage
                      src={challenge.challenger.avatar_url}
                      alt={challenge.challenger.full_name}
                    />
                    <AvatarFallback className="text-lg font-semibold">
                      {getInitials(challenge.challenger.full_name)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-semibold">{challenge.challenger.full_name}</p>
                    <p className="text-sm text-muted-foreground">Challenger</p>
                  </div>
                  {challenge.challenger_attempt_id && (
                    <Badge variant="outline">
                      <CheckCircle2 className="mr-1 h-3 w-3" />
                      Completed
                    </Badge>
                  )}
                </div>

                <div className="flex flex-col items-center text-center space-y-2">
                  <Avatar className="h-20 w-20">
                    <AvatarImage
                      src={challenge.opponent.avatar_url}
                      alt={challenge.opponent.full_name}
                    />
                    <AvatarFallback className="text-lg font-semibold">
                      {getInitials(challenge.opponent.full_name)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-semibold">{challenge.opponent.full_name}</p>
                    <p className="text-sm text-muted-foreground">Opponent</p>
                  </div>
                  {challenge.opponent_attempt_id && (
                    <Badge variant="outline">
                      <CheckCircle2 className="mr-1 h-3 w-3" />
                      Completed
                    </Badge>
                  )}
                </div>
              </div>

              {/* Quiz Details */}
              <div className="border-t pt-4">
                <h3 className="font-semibold mb-3">Quiz</h3>
                <div className="space-y-2">
                  <p className="text-lg font-medium">{challenge.quiz.title}</p>
                  <div className="flex gap-2">
                    <Badge variant="outline">{challenge.quiz.category}</Badge>
                    <Badge variant="outline">{challenge.quiz.difficulty}</Badge>
                  </div>
                </div>
              </div>

              {/* Timeline */}
              <div className="border-t pt-4 space-y-2 text-sm">
                <p className="flex items-center text-muted-foreground">
                  <Clock className="mr-2 h-4 w-4" />
                  Created {formatDistanceToNow(new Date(challenge.created_at), { addSuffix: true })}
                </p>
                {challenge.status === "pending" && (
                  <p className="flex items-center text-muted-foreground">
                    <Timer className="mr-2 h-4 w-4" />
                    Expires {formatDistanceToNow(new Date(challenge.expires_at), { addSuffix: true })}
                  </p>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Actions Card */}
          <Card>
            <CardHeader>
              <CardTitle>Next Steps</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Pending - Opponent needs to accept */}
              {challenge.status === "pending" && !isChallenger && (
                <div className="space-y-4">
                  <Alert>
                    <Swords className="h-4 w-4" />
                    <AlertDescription>
                      {opponent.full_name} has challenged you! Accept to compete on this quiz.
                    </AlertDescription>
                  </Alert>
                  <div className="flex gap-3">
                    <Button
                      onClick={handleAccept}
                      disabled={acceptMutation.isPending}
                      className="flex-1"
                    >
                      <Trophy className="mr-2 h-4 w-4" />
                      {acceptMutation.isPending ? "Accepting..." : "Accept Challenge"}
                    </Button>
                    <Button
                      onClick={handleDecline}
                      disabled={declineMutation.isPending}
                      variant="outline"
                      className="flex-1"
                    >
                      {declineMutation.isPending ? "Declining..." : "Decline"}
                    </Button>
                  </div>
                </div>
              )}

              {/* Pending - Challenger waiting */}
              {challenge.status === "pending" && isChallenger && (
                <Alert>
                  <Clock className="h-4 w-4" />
                  <AlertDescription>
                    Waiting for {opponent.full_name} to accept your challenge...
                  </AlertDescription>
                </Alert>
              )}

              {/* Active - User needs to complete quiz */}
              {challenge.status === "accepted" && !userHasCompleted && (
                <div className="space-y-4">
                  <Alert>
                    <Play className="h-4 w-4" />
                    <AlertDescription>
                      The challenge is active! Take the quiz to submit your score.
                    </AlertDescription>
                  </Alert>
                  <Button onClick={handleTakeQuiz} className="w-full">
                    <Play className="mr-2 h-4 w-4" />
                    Take Quiz
                  </Button>
                </div>
              )}

              {/* Active - User completed, waiting for opponent */}
              {challenge.status === "accepted" && userHasCompleted && !opponentHasCompleted && (
                <Alert>
                  <CheckCircle2 className="h-4 w-4" />
                  <AlertDescription>
                    You've completed the quiz! Waiting for {opponent.full_name} to complete theirs...
                  </AlertDescription>
                </Alert>
              )}

              {/* Active - Both completed, waiting for system to finalize */}
              {challenge.status === "accepted" && userHasCompleted && opponentHasCompleted && (
                <Alert>
                  <Trophy className="h-4 w-4" />
                  <AlertDescription>
                    Both players have completed the quiz! Results will be available shortly.
                  </AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* Declined or Expired */}
      {(challenge.status === "declined" || challenge.status === "expired") && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              {challenge.status === "declined" ? (
                <>
                  <XCircle className="h-16 w-16 mx-auto text-red-500 mb-4" />
                  <h2 className="text-2xl font-bold mb-2">Challenge Declined</h2>
                  <p className="text-muted-foreground">
                    This challenge was declined and is no longer active.
                  </p>
                </>
              ) : (
                <>
                  <Timer className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
                  <h2 className="text-2xl font-bold mb-2">Challenge Expired</h2>
                  <p className="text-muted-foreground">
                    This challenge expired before it was accepted.
                  </p>
                </>
              )}
              <Button onClick={() => router.push("/challenges")} className="mt-6">
                Back to Challenges
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
