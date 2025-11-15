"use client";

import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { Challenge } from "@/types/challenge";
import { Swords, Trophy, Clock, CheckCircle2, XCircle, Timer } from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";
import { useAcceptChallenge, useDeclineChallenge } from "@/hooks/useChallengeActions";
import { formatDistanceToNow } from "date-fns";

interface ChallengeCardProps {
  challenge: Challenge;
}

export function ChallengeCard({ challenge }: ChallengeCardProps) {
  const { user } = useAuth();
  const acceptMutation = useAcceptChallenge();
  const declineMutation = useDeclineChallenge();

  // Determine if current user is challenger or opponent
  const isChallenger = user?.id === challenge.challenger_id;
  const opponent = isChallenger ? challenge.opponent : challenge.challenger;

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  // Get status badge
  const getStatusBadge = () => {
    switch (challenge.status) {
      case "pending":
        return <Badge variant="secondary"><Clock className="mr-1 h-3 w-3" />Pending</Badge>;
      case "accepted":
        return <Badge variant="default"><Swords className="mr-1 h-3 w-3" />Active</Badge>;
      case "completed":
        return <Badge variant="outline"><CheckCircle2 className="mr-1 h-3 w-3" />Completed</Badge>;
      case "declined":
        return <Badge variant="destructive"><XCircle className="mr-1 h-3 w-3" />Declined</Badge>;
      case "expired":
        return <Badge variant="secondary"><Timer className="mr-1 h-3 w-3" />Expired</Badge>;
      default:
        return null;
    }
  };

  // Get winner badge
  const getWinnerBadge = () => {
    if (challenge.status === "completed" && challenge.winner_id) {
      const isWinner = user?.id === challenge.winner_id;
      return isWinner ? (
        <Badge className="bg-green-600 hover:bg-green-700">
          <Trophy className="mr-1 h-3 w-3" />You Won!
        </Badge>
      ) : (
        <Badge variant="destructive">You Lost</Badge>
      );
    }
    return null;
  };

  const handleAccept = () => {
    acceptMutation.mutate(challenge.id);
  };

  const handleDecline = () => {
    declineMutation.mutate(challenge.id);
  };

  return (
    <Card className="hover:shadow-lg transition-shadow duration-300">
      <CardHeader className="space-y-3">
        {/* Header with opponent and status */}
        <div className="flex items-start justify-between gap-2">
          <div className="flex items-center gap-3 flex-1 min-w-0">
            <Avatar className="h-12 w-12 flex-shrink-0">
              <AvatarImage src={opponent.avatar_url} alt={opponent.full_name} />
              <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                {getInitials(opponent.full_name)}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <h3 className="text-lg font-bold truncate">
                {isChallenger ? `Challenge to ${opponent.full_name}` : `Challenge from ${opponent.full_name}`}
              </h3>
              <p className="text-sm text-muted-foreground truncate">
                {challenge.quiz.title}
              </p>
            </div>
          </div>
        </div>

        {/* Badges */}
        <div className="flex gap-2 flex-wrap">
          {getStatusBadge()}
          {getWinnerBadge()}
          <Badge variant="outline">{challenge.quiz.difficulty}</Badge>
          <Badge variant="outline">{challenge.quiz.category}</Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {/* Scores */}
        {challenge.status === "completed" && (
          <div className="flex justify-around items-center p-4 bg-muted/50 rounded-lg">
            <div className="text-center">
              <p className="text-sm text-muted-foreground mb-1">
                {isChallenger ? "Your Score" : `${challenge.challenger.full_name}`}
              </p>
              <p className="text-2xl font-bold">
                {challenge.challenger_score ?? "-"}
              </p>
            </div>
            <div className="text-2xl font-bold text-muted-foreground">VS</div>
            <div className="text-center">
              <p className="text-sm text-muted-foreground mb-1">
                {isChallenger ? opponent.full_name : "Your Score"}
              </p>
              <p className="text-2xl font-bold">
                {challenge.opponent_score ?? "-"}
              </p>
            </div>
          </div>
        )}

        {/* Challenge info */}
        <div className="text-sm text-muted-foreground space-y-1">
          <p>
            <Clock className="inline h-4 w-4 mr-1" />
            Created {formatDistanceToNow(new Date(challenge.created_at), { addSuffix: true })}
          </p>
          {challenge.status === "pending" && (
            <p>
              <Timer className="inline h-4 w-4 mr-1" />
              Expires {formatDistanceToNow(new Date(challenge.expires_at), { addSuffix: true })}
            </p>
          )}
          {challenge.status === "completed" && challenge.completed_at && (
            <p>
              <CheckCircle2 className="inline h-4 w-4 mr-1" />
              Completed {formatDistanceToNow(new Date(challenge.completed_at), { addSuffix: true })}
            </p>
          )}
        </div>
      </CardContent>

      <CardFooter className="flex gap-2">
        {/* Actions based on status and user role */}
        {challenge.status === "pending" && !isChallenger && (
          <>
            <Button
              onClick={handleAccept}
              disabled={acceptMutation.isPending}
              className="flex-1"
            >
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
          </>
        )}

        {challenge.status === "pending" && isChallenger && (
          <Button disabled variant="outline" className="w-full">
            Waiting for Response...
          </Button>
        )}

        {challenge.status === "accepted" && (
          <Link href={`/challenges/${challenge.id}`} className="flex-1">
            <Button className="w-full">
              <Swords className="mr-2 h-4 w-4" />
              View Challenge
            </Button>
          </Link>
        )}

        {challenge.status === "completed" && (
          <Link href={`/challenges/${challenge.id}`} className="flex-1">
            <Button variant="outline" className="w-full">
              View Results
            </Button>
          </Link>
        )}

        {(challenge.status === "declined" || challenge.status === "expired") && (
          <Button disabled variant="outline" className="w-full">
            {challenge.status === "declined" ? "Declined" : "Expired"}
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}
