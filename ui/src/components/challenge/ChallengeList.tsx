"use client";

import { ChallengeCard } from "./ChallengeCard";
import type { Challenge } from "@/types/challenge";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Swords, AlertCircle } from "lucide-react";

interface ChallengeListProps {
  challenges: Challenge[];
  isLoading?: boolean;
  error?: Error | null;
  emptyMessage?: string;
}

export function ChallengeList({
  challenges,
  isLoading = false,
  error = null,
  emptyMessage = "No challenges found",
}: ChallengeListProps) {
  // Loading state
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[...Array(6)].map((_, i) => (
          <div key={i} className="space-y-3">
            <Skeleton className="h-48 w-full rounded-lg" />
          </div>
        ))}
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          {error.message || "Failed to load challenges"}
        </AlertDescription>
      </Alert>
    );
  }

  // Empty state
  if (!challenges || challenges.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="rounded-full bg-muted p-6 mb-4">
          <Swords className="h-12 w-12 text-muted-foreground" />
        </div>
        <h3 className="text-lg font-semibold mb-2">No Challenges Yet</h3>
        <p className="text-sm text-muted-foreground max-w-md">
          {emptyMessage}
        </p>
      </div>
    );
  }

  // Render challenges grid
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {challenges.map((challenge) => (
        <ChallengeCard key={challenge.id} challenge={challenge} />
      ))}
    </div>
  );
}
