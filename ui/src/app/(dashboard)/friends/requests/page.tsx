"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { FriendRequestCard } from "@/components/social/FriendRequestCard";
import { useFriendRequests } from "@/hooks/useFriendRequests";
import { Inbox, AlertCircle, ArrowLeft } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function FriendRequestsPage() {
  const { user } = useAuth();

  const {
    data: requests = [],
    isLoading,
    error,
  } = useFriendRequests();

  // Filter incoming and outgoing requests
  const incomingRequests = requests.filter((req) => req.requested_id === user?.id);
  const outgoingRequests = requests.filter((req) => req.requester_id === user?.id);

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <Link href="/friends">
          <Button variant="ghost" size="sm" className="mb-4">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Friends
          </Button>
        </Link>
        <h1 className="text-4xl font-bold mb-2">Friend Requests</h1>
        <p className="text-muted-foreground">
          Manage your incoming and outgoing friend requests
        </p>
      </div>

      <div className="space-y-6">
        {/* Incoming Requests */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              Incoming Requests
              {incomingRequests.length > 0 && (
                <span className="text-sm font-normal text-muted-foreground">
                  ({incomingRequests.length})
                </span>
              )}
            </CardTitle>
            <CardDescription>
              Friend requests you&apos;ve received
            </CardDescription>
          </CardHeader>
          <CardContent>
            {error ? (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  Failed to load friend requests. Please try again later.
                </AlertDescription>
              </Alert>
            ) : isLoading ? (
              <div className="space-y-2">
                {[1, 2].map((i) => (
                  <Card key={i}>
                    <CardHeader className="space-y-3">
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-12 w-12 rounded-full" />
                        <div className="space-y-2 flex-1">
                          <Skeleton className="h-4 w-32" />
                          <Skeleton className="h-3 w-24" />
                        </div>
                      </div>
                    </CardHeader>
                  </Card>
                ))}
              </div>
            ) : incomingRequests.length > 0 ? (
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {incomingRequests.map((request) => (
                  <FriendRequestCard key={request.id} request={request} />
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <Inbox className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">No Incoming Requests</h3>
                <p className="text-muted-foreground mb-4">
                  You don&apos;t have any pending friend requests
                </p>
                <Link href="/friends?tab=search">
                  <Button>Find Friends</Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Outgoing Requests */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              Sent Requests
              {outgoingRequests.length > 0 && (
                <span className="text-sm font-normal text-muted-foreground">
                  ({outgoingRequests.length})
                </span>
              )}
            </CardTitle>
            <CardDescription>
              Friend requests you&apos;ve sent
            </CardDescription>
          </CardHeader>
          <CardContent>
            {error ? (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  Failed to load friend requests. Please try again later.
                </AlertDescription>
              </Alert>
            ) : isLoading ? (
              <div className="space-y-2">
                {[1].map((i) => (
                  <Card key={i}>
                    <CardHeader className="space-y-3">
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-12 w-12 rounded-full" />
                        <div className="space-y-2 flex-1">
                          <Skeleton className="h-4 w-32" />
                          <Skeleton className="h-3 w-24" />
                        </div>
                      </div>
                    </CardHeader>
                  </Card>
                ))}
              </div>
            ) : outgoingRequests.length > 0 ? (
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {outgoingRequests.map((request) => (
                  <FriendRequestCard key={request.id} request={request} />
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <Inbox className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">No Pending Requests</h3>
                <p className="text-muted-foreground mb-4">
                  You haven&apos;t sent any friend requests yet
                </p>
                <Link href="/friends?tab=search">
                  <Button>Find Friends</Button>
                </Link>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
