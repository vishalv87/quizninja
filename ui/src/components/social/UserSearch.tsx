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
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search for users by name or email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
        {isLoading && (
          <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground animate-spin" />
        )}
      </div>

      {/* Search Results */}
      {showResults && (
        <div className="space-y-2">
          {isLoading ? (
            <Card>
              <CardContent className="p-4">
                <div className="flex items-center gap-3">
                  <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                  <p className="text-sm text-muted-foreground">Searching...</p>
                </div>
              </CardContent>
            </Card>
          ) : hasResults ? (
            <>
              <p className="text-sm text-muted-foreground">
                Found {searchResults.length} {searchResults.length === 1 ? "user" : "users"}
              </p>
              <div className="space-y-2">
                {searchResults.map((user) => (
                  <Card key={user.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3">
                        {/* Avatar */}
                        <Avatar className="h-12 w-12 flex-shrink-0">
                          <AvatarImage src={user.avatar_url} alt={user.full_name} />
                          <AvatarFallback className="bg-primary/10 text-primary font-semibold">
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
                            <Badge variant="secondary" className="gap-1">
                              <Check className="h-3 w-3" />
                              Friends
                            </Badge>
                          ) : user.has_pending_request ? (
                            <Badge variant="outline" className="gap-1">
                              <Clock className="h-3 w-3" />
                              Pending
                            </Badge>
                          ) : user.is_request_sent ? (
                            <Badge variant="outline" className="gap-1">
                              <Clock className="h-3 w-3" />
                              Sent
                            </Badge>
                          ) : (
                            <Button
                              size="sm"
                              onClick={() => handleSendRequest(user.id)}
                              disabled={sendRequestMutation.isPending}
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
                    </CardContent>
                  </Card>
                ))}
              </div>
            </>
          ) : (
            <Card>
              <CardContent className="p-8 text-center">
                <p className="text-muted-foreground">
                  No users found matching &quot;{debouncedQuery}&quot;
                </p>
                <p className="text-sm text-muted-foreground mt-2">
                  Try searching with a different name or email
                </p>
              </CardContent>
            </Card>
          )}
        </div>
      )}

      {/* Initial State */}
      {!showResults && (
        <Card>
          <CardContent className="p-8 text-center">
            <Search className="h-12 w-12 mx-auto text-muted-foreground mb-3" />
            <p className="text-muted-foreground">
              Start typing to search for users
            </p>
            <p className="text-sm text-muted-foreground mt-2">
              Search by name or email address
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
