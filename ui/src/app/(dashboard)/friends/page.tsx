"use client";

import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { FriendList } from "@/components/social/FriendList";
import { FriendRequestCard } from "@/components/social/FriendRequestCard";
import { UserSearch } from "@/components/social/UserSearch";
import { useFriends } from "@/hooks/useFriends";
import { useFriendRequests } from "@/hooks/useFriendRequests";
import { Users, UserPlus, Inbox, AlertCircle } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useAuth } from "@/hooks/useAuth";

export default function FriendsPage() {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState("friends");

  const {
    data: friends = [],
    isLoading: friendsLoading,
    error: friendsError,
  } = useFriends();

  const {
    data: requests = [],
    isLoading: requestsLoading,
    error: requestsError,
  } = useFriendRequests();

  // Filter incoming and outgoing requests
  const incomingRequests = requests.filter((req) => req.requested_id === user?.id);
  const outgoingRequests = requests.filter((req) => req.requester_id === user?.id);

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Friends</h1>
        <p className="text-muted-foreground">
          Connect with other quiz enthusiasts and challenge your friends!
        </p>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-auto lg:inline-grid">
          <TabsTrigger value="friends" className="gap-2">
            <Users className="h-4 w-4" />
            Friends
            {friends.length > 0 && (
              <span className="ml-1 rounded-full bg-primary px-2 py-0.5 text-xs text-primary-foreground">
                {friends.length}
              </span>
            )}
          </TabsTrigger>
          <TabsTrigger value="requests" className="gap-2">
            <Inbox className="h-4 w-4" />
            Requests
            {incomingRequests.length > 0 && (
              <span className="ml-1 rounded-full bg-destructive px-2 py-0.5 text-xs text-destructive-foreground">
                {incomingRequests.length}
              </span>
            )}
          </TabsTrigger>
          <TabsTrigger value="search" className="gap-2">
            <UserPlus className="h-4 w-4" />
            Add Friends
          </TabsTrigger>
        </TabsList>

        {/* Friends Tab */}
        <TabsContent value="friends" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Your Friends</CardTitle>
              <CardDescription>
                View and manage your connections
              </CardDescription>
            </CardHeader>
            <CardContent>
              {friendsError ? (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Failed to load friends. Please try again later.
                  </AlertDescription>
                </Alert>
              ) : friendsLoading ? (
                <div className="space-y-4">
                  <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {[1, 2, 3].map((i) => (
                      <Card key={i}>
                        <CardHeader className="space-y-3">
                          <div className="flex items-center gap-3">
                            <Skeleton className="h-12 w-12 rounded-full" />
                            <div className="space-y-2 flex-1">
                              <Skeleton className="h-4 w-32" />
                              <Skeleton className="h-3 w-48" />
                            </div>
                          </div>
                        </CardHeader>
                        <CardContent>
                          <Skeleton className="h-4 w-full" />
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </div>
              ) : (
                <FriendList friends={friends} />
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Requests Tab */}
        <TabsContent value="requests" className="space-y-4">
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
              {requestsError ? (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Failed to load friend requests. Please try again later.
                  </AlertDescription>
                </Alert>
              ) : requestsLoading ? (
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
                <div className="grid gap-3 sm:grid-cols-2">
                  {incomingRequests.map((request) => (
                    <FriendRequestCard key={request.id} request={request} />
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Inbox className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
                  <p className="text-muted-foreground">No incoming friend requests</p>
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
              {requestsError ? (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Failed to load friend requests. Please try again later.
                  </AlertDescription>
                </Alert>
              ) : requestsLoading ? (
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
                <div className="grid gap-3 sm:grid-cols-2">
                  {outgoingRequests.map((request) => (
                    <FriendRequestCard key={request.id} request={request} />
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Inbox className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
                  <p className="text-muted-foreground">No pending sent requests</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Add Friends Tab */}
        <TabsContent value="search" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Find Friends</CardTitle>
              <CardDescription>
                Search for users and send friend requests
              </CardDescription>
            </CardHeader>
            <CardContent>
              <UserSearch />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
