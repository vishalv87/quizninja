"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Search, UserPlus, Check, Clock, Loader2 } from "lucide-react";
import { useSearchUsers } from "@/hooks/useSearchUsers";
import { useSendFriendRequest } from "@/hooks/useFriendRequests";
import { useDebounce } from "@/hooks/useDebounce";

export function UserSearch() {
  const [searchQuery, setSearchQuery] = useState("");
  const debouncedQuery = useDebounce(searchQuery, 500); // 500ms debounce

  const { data: searchResults = [], isLoading } = useSearchUsers(debouncedQuery);
  const sendRequestMutation = useSendFriendRequest();

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const handleSendRequest = (userId: string) => {
    sendRequestMutation.mutate(userId);
  };

  const showResults = debouncedQuery.length >= 2;
  const hasResults = searchResults.length > 0;

  return (
    <div className="space-y-4">
      {/* Search Input */}
      <div className="relative group">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground transition-colors duration-200 group-focus-within:text-violet-500" />
        <Input
          placeholder="Search for users by name or email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10 bg-white/90 dark:bg-background/90 backdrop-blur-sm border-gray-200/50 dark:border-gray-700/50 rounded-xl transition-all duration-300 hover:border-violet-400/50 dark:hover:border-violet-600/50 focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 focus:shadow-lg focus:shadow-violet-500/10"
        />
        {isLoading && (
          <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-violet-500 animate-spin" />
        )}
      </div>

      {/* Search Results */}
      {showResults && (
        <div className="space-y-3">
          {isLoading ? (
            <div className="border border-white/20 dark:border-white/10 rounded-2xl p-4 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5">
              <div className="flex items-center gap-3">
                <Loader2 className="h-5 w-5 animate-spin text-violet-500" />
                <p className="text-sm text-muted-foreground">Searching...</p>
              </div>
            </div>
          ) : hasResults ? (
            <>
              <p className="text-sm text-muted-foreground">
                Found {searchResults.length} {searchResults.length === 1 ? "user" : "users"}
              </p>
              <div className="space-y-3">
                {searchResults.map((user, index) => (
                  <div
                    key={user.id}
                    className="group border border-white/20 dark:border-white/10 rounded-2xl p-4 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5 transition-all duration-300 hover:shadow-xl hover:shadow-violet-500/10 hover:-translate-y-0.5 hover:scale-[1.01]"
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
                    <div className="flex items-center gap-3">
                      {/* Avatar */}
                      <Avatar className="h-12 w-12 flex-shrink-0 ring-2 ring-violet-200/50 dark:ring-violet-800/50 transition-all duration-300 group-hover:ring-4 group-hover:ring-violet-300/50">
                        <AvatarImage src={user.avatar_url} alt={user.full_name} />
                        <AvatarFallback className="bg-gradient-to-br from-violet-500 to-indigo-500 text-white font-semibold">
                          {getInitials(user.full_name)}
                        </AvatarFallback>
                      </Avatar>

                      {/* User Info */}
                      <div className="flex-1 min-w-0">
                        <h4 className="font-semibold truncate">{user.full_name}</h4>
                        <p className="text-sm text-muted-foreground truncate">
                          {user.email}
                        </p>
                      </div>

                      {/* Action Button / Status */}
                      <div className="flex-shrink-0">
                        {user.is_friend ? (
                          <Badge variant="secondary" className="gap-1 bg-green-100/50 dark:bg-green-900/20 text-green-700 dark:text-green-400 border border-green-200/50 dark:border-green-800/50">
                            <Check className="h-3 w-3" />
                            Friends
                          </Badge>
                        ) : user.has_pending_request ? (
                          <Badge variant="outline" className="gap-1 border-amber-200/50 dark:border-amber-800/50 bg-amber-50/50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-400">
                            <Clock className="h-3 w-3" />
                            Pending
                          </Badge>
                        ) : user.is_request_sent ? (
                          <Badge variant="outline" className="gap-1 border-violet-200/50 dark:border-violet-800/50 bg-violet-50/50 dark:bg-violet-900/20 text-violet-700 dark:text-violet-400">
                            <Clock className="h-3 w-3" />
                            Sent
                          </Badge>
                        ) : (
                          <Button
                            size="sm"
                            onClick={() => handleSendRequest(user.id)}
                            disabled={sendRequestMutation.isPending}
                            className="rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 shadow-md shadow-indigo-500/20 transition-all duration-200 active:scale-95 hover:shadow-lg hover:shadow-indigo-500/30 hover:scale-105"
                          >
                            {sendRequestMutation.isPending ? (
                              <Loader2 className="h-4 w-4 animate-spin" />
                            ) : (
                              <>
                                <UserPlus className="mr-2 h-4 w-4" />
                                Add Friend
                              </>
                            )}
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </>
          ) : (
            <div className="border border-white/20 dark:border-white/10 rounded-2xl p-8 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5 text-center">
              <p className="text-muted-foreground">
                No users found matching &quot;{debouncedQuery}&quot;
              </p>
              <p className="text-sm text-muted-foreground mt-2">
                Try searching with a different name or email
              </p>
            </div>
          )}
        </div>
      )}

      {/* Initial State */}
      {!showResults && (
        <div className="border border-white/20 dark:border-white/10 rounded-2xl p-8 bg-white/90 dark:bg-background/90 backdrop-blur-sm shadow-md shadow-black/5 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-br from-violet-100 to-indigo-100 dark:from-violet-900/30 dark:to-indigo-900/30 mb-4 animate-pulse">
            <Search className="h-8 w-8 text-violet-500" />
          </div>
          <p className="font-medium text-foreground">
            Start typing to search for users
          </p>
          <p className="text-sm text-muted-foreground mt-2">
            Search by name or email address
          </p>
        </div>
      )}
    </div>
  );
}
