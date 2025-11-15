"use client";

import { useParams, useRouter } from "next/navigation";
import { QuizResults } from "@/components/quiz/QuizResults";
import { useQuizAttemptResults } from "@/hooks/useQuizAttempt";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

/**
 * Quiz Results Page
 * Displays the results of a completed quiz attempt
 */
export default function QuizResultsPage() {
  const params = useParams();
  const router = useRouter();
  const quizId = params.id as string;
  const attemptId = params.attemptId as string;

  const { data: results, isLoading, error } = useQuizAttemptResults(attemptId);

  const handleRetry = () => {
    router.push(`/quizzes/${quizId}`);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="container max-w-4xl mx-auto py-8 px-4">
        <div className="space-y-6">
          <Card>
            <CardContent className="pt-6 space-y-4">
              <div className="flex justify-center">
                <Skeleton className="w-20 h-20 rounded-full" />
              </div>
              <Skeleton className="h-12 w-64 mx-auto" />
              <Skeleton className="h-8 w-32 mx-auto" />
              <Skeleton className="h-3 w-full max-w-md mx-auto" />
            </CardContent>
          </Card>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {[1, 2, 3, 4].map((i) => (
              <Card key={i}>
                <CardContent className="pt-6">
                  <Skeleton className="h-20 w-full" />
                </CardContent>
              </Card>
            ))}
          </div>
          <Card>
            <CardContent className="pt-6">
              <Skeleton className="h-64 w-full" />
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="container max-w-4xl mx-auto py-8 px-4">
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error Loading Results</AlertTitle>
          <AlertDescription>
            {error instanceof Error
              ? error.message
              : "Failed to load quiz results. Please try again."}
          </AlertDescription>
        </Alert>
        <div className="mt-6 flex gap-4">
          <Link href={`/quizzes/${quizId}`}>
            <Button variant="outline">Back to Quiz</Button>
          </Link>
          <Link href="/quizzes">
            <Button>Browse Quizzes</Button>
          </Link>
        </div>
      </div>
    );
  }

  // Results loaded successfully
  if (!results) {
    return (
      <div className="container max-w-4xl mx-auto py-8 px-4">
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>No Results Found</AlertTitle>
          <AlertDescription>
            The quiz attempt results could not be found.
          </AlertDescription>
        </Alert>
        <div className="mt-6">
          <Link href="/quizzes">
            <Button>Browse Quizzes</Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="container max-w-4xl mx-auto py-8 px-4">
      <QuizResults results={results} onRetry={handleRetry} />
    </div>
  );
}
