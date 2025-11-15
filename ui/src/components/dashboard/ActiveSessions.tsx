"use client";

import Link from "next/link";
import { ArrowRight, Clock, PlayCircle } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { useActiveSessions } from "@/hooks/useUserStats";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";

export function ActiveSessions() {
  const { data: sessionsData, isLoading } = useActiveSessions();

  // Note: Each user has only ONE attempt per quiz
  // Ensure sessions is always an array
  const sessions = Array.isArray(sessionsData) ? sessionsData : [];

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Active Sessions</CardTitle>
          <CardDescription>Continue your in-progress quizzes</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {[...Array(2)].map((_, i) => (
            <div key={i} className="space-y-2">
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-2 w-full" />
              <Skeleton className="h-3 w-1/2" />
            </div>
          ))}
        </CardContent>
      </Card>
    );
  }

  if (!sessions || sessions.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Active Sessions</CardTitle>
          <CardDescription>Continue your in-progress quizzes</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <p className="text-muted-foreground text-center">
            No active quiz sessions
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Active Sessions</CardTitle>
        <CardDescription>Continue your in-progress quizzes</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {sessions.map((session: any) => {
            const progress = (session.questions_answered / session.total_questions) * 100;

            return (
              <div
                key={session.id}
                className="p-4 rounded-lg border hover:bg-accent transition-colors"
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <h4 className="font-medium">{session.quiz_title}</h4>
                    <div className="flex items-center gap-2 mt-1">
                      <Badge variant="outline">{session.category}</Badge>
                      <Badge variant="outline">{session.difficulty}</Badge>
                      {session.status === "paused" && (
                        <Badge variant="secondary">Paused</Badge>
                      )}
                    </div>
                  </div>
                  <Button size="sm" asChild>
                    <Link href={`/quizzes/${session.quiz_id}/take`}>
                      <PlayCircle className="mr-2 h-4 w-4" />
                      Resume
                    </Link>
                  </Button>
                </div>

                <div className="space-y-2">
                  <Progress value={progress} className="h-2" />
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <span>
                      {session.questions_answered} of {session.total_questions} questions
                    </span>
                    <span className="flex items-center">
                      <Clock className="mr-1 h-3 w-3" />
                      Started{" "}
                      {formatDistanceToNow(new Date(session.started_at), {
                        addSuffix: true,
                      })}
                    </span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}