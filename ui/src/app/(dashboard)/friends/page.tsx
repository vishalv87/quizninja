"use client";

import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { FriendList } from "@/components/social/FriendList";
import { FriendRequestCard } from "@/components/social/FriendRequestCard";
import { UserSearch } from "@/components/social/UserSearch";
import { useFriends } from "@/hooks/useFriends";
import { useFriendRequests } from "@/hooks/useFriendRequests";
import { Users, UserPlus, Inbox, AlertCircle, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
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
    <div className="space-y-10 pb-10">
      {/* Hero Section */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-violet-600 via-indigo-600 to-purple-700 p-8 text-white shadow-2xl shadow-indigo-500/30 border border-white/10 lg:p-12">
        <div className="relative z-10 max-w-2xl">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl mb-6 drop-shadow-sm">
            Hey{user?.name ? `, ${user.name.split(' ')[0]}` : ''}! 👋
          </h1>
          <p className="text-xl text-indigo-100 mb-8 font-medium leading-relaxed">
            Connect with quiz enthusiasts and challenge your friends to epic battles!
          </p>
          <Button
            size="lg"
            className="bg-white text-indigo-600 hover:bg-indigo-50 border-0 font-bold h-12 px-8 rounded-xl shadow-lg shadow-black/10 transition-all hover:scale-105 hover:shadow-xl active:scale-95"
            onClick={() => setActiveTab("search")}
          >
            <Search className="mr-2 h-5 w-5" />
            Find New Friends
          </Button>
        </div>

        {/* Decorative background elements with subtle animation */}
        <div className="absolute right-0 top-0 -mt-20 -mr-20 h-96 w-96 rounded-full bg-white/10 blur-3xl animate-pulse" style={{ animationDuration: '4s' }} />
        <div className="absolute bottom-0 right-20 -mb-20 h-64 w-64 rounded-full bg-indigo-400/20 blur-3xl animate-pulse" style={{ animationDuration: '6s' }} />
        <div className="absolute left-10 bottom-10 h-32 w-32 rounded-full bg-purple-400/20 blur-2xl animate-pulse" style={{ animationDuration: '5s' }} />
      </div>

      <div className="container px-0 md:px-4">
        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full grid-cols-3 lg:w-auto lg:inline-grid bg-white/60 dark:bg-black/40 backdrop-blur-md border border-white/20 dark:border-white/10 p-1 rounded-xl shadow-sm">
            <TabsTrigger
              value="friends"
              className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
            >
              <Users className="h-4 w-4" />
              Friends
              {friends.length > 0 && (
                <span className="ml-1 rounded-full bg-violet-600 px-2 py-0.5 text-xs text-white">
                  {friends.length}
                </span>
              )}
            </TabsTrigger>
            <TabsTrigger
              value="requests"
              className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
            >
              <Inbox className="h-4 w-4" />
              Requests
              {incomingRequests.length > 0 && (
                <span className="ml-1 rounded-full bg-destructive px-2 py-0.5 text-xs text-destructive-foreground animate-pulse shadow-sm shadow-destructive/50">
                  {incomingRequests.length}
                </span>
              )}
            </TabsTrigger>
            <TabsTrigger
              value="search"
              className="gap-2 rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
            >
              <UserPlus className="h-4 w-4" />
              Add Friends
            </TabsTrigger>
          </TabsList>

        {/* Friends Tab */}
        <TabsContent value="friends" className="space-y-4">
          <div className="rounded-2xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md shadow-sm overflow-hidden">
            <div className="p-6 pb-4">
              <h2 className="text-xl font-semibold">Your Friends</h2>
              <p className="text-sm text-muted-foreground mt-1">
                View and manage your connections
              </p>
            </div>
            <div className="px-6 pb-6">
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
                      <div key={i} className="border border-white/20 dark:border-white/10 rounded-2xl p-6 space-y-4 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5">
                        <div className="h-1.5 w-full rounded-full bg-gradient-to-r from-violet-400/50 to-indigo-500/50 -mt-6 -mx-6 mb-4" style={{ width: 'calc(100% + 3rem)' }} />
                        <div className="flex items-center gap-3">
                          <Skeleton className="h-12 w-12 rounded-full" />
                          <div className="space-y-2 flex-1">
                            <Skeleton className="h-4 w-32 rounded-lg" />
                            <Skeleton className="h-3 w-48 rounded-lg" />
                          </div>
                        </div>
                        <Skeleton className="h-4 w-full rounded-lg" />
                        <div className="flex gap-2 pt-2">
                          <Skeleton className="h-9 flex-1 rounded-lg" />
                          <Skeleton className="h-9 flex-1 rounded-lg" />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : (
                <FriendList friends={friends} onSearchClick={() => setActiveTab("search")} />
              )}
            </div>
          </div>
        </TabsContent>

        {/* Requests Tab */}
        <TabsContent value="requests" className="space-y-6">
          {/* Incoming Requests */}
          <div className="rounded-2xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md shadow-sm overflow-hidden">
            <div className="p-6 pb-4">
              <h2 className="text-xl font-semibold flex items-center gap-2">
                Incoming Requests
                {incomingRequests.length > 0 && (
                  <span className="text-sm font-normal text-muted-foreground">
                    ({incomingRequests.length})
                  </span>
                )}
              </h2>
              <p className="text-sm text-muted-foreground mt-1">
                Friend requests you&apos;ve received
              </p>
            </div>
            <div className="px-6 pb-6">
              {requestsError ? (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Failed to load friend requests. Please try again later.
                  </AlertDescription>
                </Alert>
              ) : requestsLoading ? (
                <div className="grid gap-3 sm:grid-cols-2">
                  {[1, 2].map((i) => (
                    <div key={i} className="border border-white/20 dark:border-white/10 rounded-2xl p-4 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5">
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-12 w-12 rounded-full" />
                        <div className="space-y-2 flex-1">
                          <Skeleton className="h-4 w-32 rounded-lg" />
                          <Skeleton className="h-3 w-24 rounded-lg" />
                        </div>
                      </div>
                    </div>
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
                  <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-br from-amber-100 to-orange-100 dark:from-amber-900/30 dark:to-orange-900/30 mb-4">
                    <Inbox className="h-8 w-8 text-amber-500" />
                  </div>
                  <p className="font-medium text-foreground mb-1">No incoming requests</p>
                  <p className="text-sm text-muted-foreground">When someone sends you a friend request, it will appear here</p>
                </div>
              )}
            </div>
          </div>

          {/* Outgoing Requests */}
          <div className="rounded-2xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md shadow-sm overflow-hidden">
            <div className="p-6 pb-4">
              <h2 className="text-xl font-semibold flex items-center gap-2">
                Sent Requests
                {outgoingRequests.length > 0 && (
                  <span className="text-sm font-normal text-muted-foreground">
                    ({outgoingRequests.length})
                  </span>
                )}
              </h2>
              <p className="text-sm text-muted-foreground mt-1">
                Friend requests you&apos;ve sent
              </p>
            </div>
            <div className="px-6 pb-6">
              {requestsError ? (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Failed to load friend requests. Please try again later.
                  </AlertDescription>
                </Alert>
              ) : requestsLoading ? (
                <div className="grid gap-3 sm:grid-cols-2">
                  {[1].map((i) => (
                    <div key={i} className="border border-white/20 dark:border-white/10 rounded-2xl p-4 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5">
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-12 w-12 rounded-full" />
                        <div className="space-y-2 flex-1">
                          <Skeleton className="h-4 w-32 rounded-lg" />
                          <Skeleton className="h-3 w-24 rounded-lg" />
                        </div>
                      </div>
                    </div>
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
                  <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-br from-violet-100 to-indigo-100 dark:from-violet-900/30 dark:to-indigo-900/30 mb-4">
                    <UserPlus className="h-8 w-8 text-violet-500" />
                  </div>
                  <p className="font-medium text-foreground mb-1">No pending requests</p>
                  <p className="text-sm text-muted-foreground">Friend requests you send will appear here until accepted</p>
                </div>
              )}
            </div>
          </div>
        </TabsContent>

        {/* Add Friends Tab */}
        <TabsContent value="search" className="space-y-4">
          <div className="rounded-2xl border border-white/20 dark:border-white/10 bg-white/40 dark:bg-black/40 backdrop-blur-md shadow-sm overflow-hidden">
            <div className="p-6 pb-4">
              <h2 className="text-xl font-semibold">Find Friends</h2>
              <p className="text-sm text-muted-foreground mt-1">
                Search for users and send friend requests
              </p>
            </div>
            <div className="px-6 pb-6">
              <UserSearch />
            </div>
          </div>
        </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
