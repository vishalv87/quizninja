"use client";

import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import type { FriendRequest } from "@/types/user";
import { Check, X, Clock } from "lucide-react";
import { useAcceptFriendRequest, useDeclineFriendRequest, useCancelFriendRequest } from "@/hooks/useFriendRequests";
import { useAuth } from "@/hooks/useAuth";
import { formatDistanceToNow } from "date-fns";

interface FriendRequestCardProps {
  request: FriendRequest;
}

export function FriendRequestCard({ request }: FriendRequestCardProps) {
  const { user } = useAuth();
  const acceptMutation = useAcceptFriendRequest();
  const declineMutation = useDeclineFriendRequest();
  const cancelMutation = useCancelFriendRequest();

  // Determine if this is an incoming or outgoing request
  const isIncoming = user?.id === request.requested_id;
  const requestUser = request.requester;

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const handleAccept = () => {
    acceptMutation.mutate(request.id);
  };

  const handleDecline = () => {
    declineMutation.mutate(request.id);
  };

  const handleCancel = () => {
    cancelMutation.mutate(request.id);
  };

  const timeAgo = formatDistanceToNow(new Date(request.created_at), {
    addSuffix: true,
  });

  return (
    <Card className="group relative overflow-hidden border border-white/20 dark:border-white/10 shadow-md shadow-black/5 transition-all duration-300 hover:shadow-xl hover:shadow-black/10 hover:-translate-y-1 hover:scale-[1.02] bg-white/90 dark:bg-background/90 backdrop-blur-sm rounded-2xl">
      {/* Top Decoration Bar - different color for incoming vs outgoing */}
      <div className={`h-1.5 w-full bg-gradient-to-r ${isIncoming ? "from-amber-400 to-orange-500 animate-pulse" : "from-violet-400 to-indigo-500"}`} />

      <CardHeader className="space-y-3">
        <div className="flex items-start gap-3">
          <Avatar className={`h-12 w-12 flex-shrink-0 ring-2 transition-all duration-300 group-hover:ring-4 ${isIncoming ? "ring-amber-200/50 dark:ring-amber-800/50 group-hover:ring-amber-300/50 dark:group-hover:ring-amber-700/50" : "ring-violet-200/50 dark:ring-violet-800/50 group-hover:ring-violet-300/50 dark:group-hover:ring-violet-700/50"}`}>
            <AvatarImage src={requestUser.avatar_url} alt={requestUser.full_name} />
            <AvatarFallback className={`font-semibold text-white ${isIncoming ? "bg-gradient-to-br from-amber-500 to-orange-500" : "bg-gradient-to-br from-violet-500 to-indigo-500"}`}>
              {getInitials(requestUser.full_name)}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-bold truncate">{requestUser.full_name}</h3>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant="outline" className={`text-xs ${isIncoming ? "border-amber-200/50 dark:border-amber-800/50 bg-amber-50/50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-400" : "border-violet-200/50 dark:border-violet-800/50 bg-violet-50/50 dark:bg-violet-900/20 text-violet-700 dark:text-violet-400"}`}>
                {isIncoming ? "Incoming" : "Sent"}
              </Badge>
              <div className="flex items-center gap-1 text-xs text-muted-foreground">
                <Clock className="h-3 w-3" />
                <span>{timeAgo}</span>
              </div>
            </div>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        <p className="text-sm text-muted-foreground border-t border-gray-200/30 dark:border-gray-700/30 pt-3">
          {isIncoming
            ? `${requestUser.full_name} wants to be your friend`
            : `Waiting for ${requestUser.full_name} to accept your request`}
        </p>
      </CardContent>

      <CardFooter className="flex gap-2">
        {isIncoming ? (
          <>
            <Button
              variant="outline"
              className="flex-1 rounded-xl border-gray-200/50 dark:border-gray-700/50 hover:border-red-400/50 dark:hover:border-red-600/50 hover:bg-red-50/50 dark:hover:bg-red-900/10 transition-all duration-200 active:scale-95 hover:shadow-md hover:shadow-red-500/10"
              onClick={handleDecline}
              disabled={declineMutation.isPending || acceptMutation.isPending}
            >
              {declineMutation.isPending ? (
                <>Declining...</>
              ) : (
                <>
                  <X className="mr-2 h-4 w-4" />
                  Decline
                </>
              )}
            </Button>
            <Button
              className="flex-1 rounded-xl bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600 shadow-md shadow-green-500/20 transition-all duration-200 active:scale-95 hover:shadow-lg hover:shadow-green-500/30"
              onClick={handleAccept}
              disabled={acceptMutation.isPending || declineMutation.isPending}
            >
              {acceptMutation.isPending ? (
                <>Accepting...</>
              ) : (
                <>
                  <Check className="mr-2 h-4 w-4" />
                  Accept
                </>
              )}
            </Button>
          </>
        ) : (
          <Button
            variant="outline"
            className="w-full rounded-xl border-gray-200/50 dark:border-gray-700/50 hover:border-red-400/50 dark:hover:border-red-600/50 hover:bg-red-50/50 dark:hover:bg-red-900/10 transition-all duration-200 active:scale-95 hover:shadow-md hover:shadow-red-500/10"
            onClick={handleCancel}
            disabled={cancelMutation.isPending}
          >
            {cancelMutation.isPending ? (
              <>Canceling...</>
            ) : (
              <>
                <X className="mr-2 h-4 w-4" />
                Cancel Request
              </>
            )}
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}
