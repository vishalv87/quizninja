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
    <Card className="hover:shadow-md transition-shadow duration-300">
      <CardHeader className="space-y-3">
        <div className="flex items-start gap-3">
          <Avatar className="h-12 w-12 flex-shrink-0">
            <AvatarImage src={requestUser.avatar_url} alt={requestUser.full_name} />
            <AvatarFallback className="bg-primary/10 text-primary font-semibold">
              {getInitials(requestUser.full_name)}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-bold truncate">{requestUser.full_name}</h3>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant="outline" className="text-xs">
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
        <p className="text-sm text-muted-foreground">
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
              className="flex-1"
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
              className="flex-1"
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
            className="w-full"
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
