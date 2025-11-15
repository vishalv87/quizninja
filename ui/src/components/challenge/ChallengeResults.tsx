"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { Challenge } from "@/types/challenge";
import { Trophy, Medal, Target, Clock } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { formatDistanceToNow } from "date-fns";

interface ChallengeResultsProps {
  challenge: Challenge;
}

export function ChallengeResults({ challenge }: ChallengeResultsProps) {
  const { user } = useAuth();

  // Determine if current user is challenger or opponent
  const isChallenger = user?.id === challenge.challenger_id;
  const userScore = isChallenger ? challenge.challenger_score : challenge.opponent_score;
  const opponentScore = isChallenger ? challenge.opponent_score : challenge.challenger_score;
  const opponent = isChallenger ? challenge.opponent : challenge.challenger;

  // Determine winner
  const isWinner = user?.id === challenge.winner_id;
  const isDraw = challenge.challenger_score === challenge.opponent_score;

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  return (
    <div className="space-y-6">
      {/* Winner Banner */}
      <Card className={isWinner ? "border-green-500 bg-green-50 dark:bg-green-950" : "border-red-500 bg-red-50 dark:bg-red-950"}>
        <CardContent className="pt-6">
          <div className="flex flex-col items-center text-center">
            {isDraw ? (
              <>
                <Medal className="h-16 w-16 text-yellow-500 mb-4" />
                <h2 className="text-2xl font-bold mb-2">It's a Draw!</h2>
                <p className="text-muted-foreground">
                  Both players scored {userScore} points
                </p>
              </>
            ) : isWinner ? (
              <>
                <Trophy className="h-16 w-16 text-green-500 mb-4" />
                <h2 className="text-2xl font-bold mb-2">Congratulations!</h2>
                <p className="text-muted-foreground">
                  You defeated {opponent.full_name}
                </p>
              </>
            ) : (
              <>
                <Target className="h-16 w-16 text-red-500 mb-4" />
                <h2 className="text-2xl font-bold mb-2">Better Luck Next Time!</h2>
                <p className="text-muted-foreground">
                  {opponent.full_name} won this challenge
                </p>
              </>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Score Comparison */}
      <Card>
        <CardHeader>
          <CardTitle>Final Scores</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex justify-around items-center">
            {/* User Score */}
            <div className="flex flex-col items-center space-y-2">
              <Avatar className="h-20 w-20">
                <AvatarImage src={user?.avatar_url} alt="You" />
                <AvatarFallback className="text-lg font-semibold">
                  {user?.full_name ? getInitials(user.full_name) : "ME"}
                </AvatarFallback>
              </Avatar>
              <div className="text-center">
                <p className="text-sm font-medium">You</p>
                <p className="text-3xl font-bold text-primary">{userScore ?? 0}</p>
              </div>
              {isWinner && !isDraw && (
                <Badge className="bg-green-600 hover:bg-green-700">
                  <Trophy className="mr-1 h-3 w-3" />
                  Winner
                </Badge>
              )}
            </div>

            {/* VS Divider */}
            <div className="text-4xl font-bold text-muted-foreground px-4">VS</div>

            {/* Opponent Score */}
            <div className="flex flex-col items-center space-y-2">
              <Avatar className="h-20 w-20">
                <AvatarImage src={opponent.avatar_url} alt={opponent.full_name} />
                <AvatarFallback className="text-lg font-semibold">
                  {getInitials(opponent.full_name)}
                </AvatarFallback>
              </Avatar>
              <div className="text-center">
                <p className="text-sm font-medium">{opponent.full_name}</p>
                <p className="text-3xl font-bold text-primary">{opponentScore ?? 0}</p>
              </div>
              {!isWinner && !isDraw && (
                <Badge className="bg-green-600 hover:bg-green-700">
                  <Trophy className="mr-1 h-3 w-3" />
                  Winner
                </Badge>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Challenge Details */}
      <Card>
        <CardHeader>
          <CardTitle>Challenge Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground mb-1">Quiz</p>
              <p className="font-semibold">{challenge.quiz.title}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-1">Category</p>
              <Badge variant="outline">{challenge.quiz.category}</Badge>
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-1">Difficulty</p>
              <Badge variant="outline">{challenge.quiz.difficulty}</Badge>
            </div>
            <div>
              <p className="text-sm text-muted-foreground mb-1">Completed</p>
              <p className="text-sm flex items-center">
                <Clock className="h-4 w-4 mr-1" />
                {challenge.completed_at
                  ? formatDistanceToNow(new Date(challenge.completed_at), { addSuffix: true })
                  : "N/A"}
              </p>
            </div>
          </div>

          {/* Score Difference */}
          {!isDraw && userScore !== undefined && opponentScore !== undefined && (
            <div className="pt-4 border-t">
              <p className="text-sm text-muted-foreground mb-1">Margin of Victory</p>
              <p className="text-2xl font-bold">
                {Math.abs(userScore - opponentScore)} points
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
