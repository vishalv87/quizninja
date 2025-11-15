"use client";

import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import type { Friend } from "@/types/user";
import { MoreVertical, UserMinus, Swords, User } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { useRemoveFriend } from "@/hooks/useFriends";

interface FriendCardProps {
  friend: Friend;
}

export function FriendCard({ friend }: FriendCardProps) {
  const [showRemoveDialog, setShowRemoveDialog] = useState(false);
  const removeFriendMutation = useRemoveFriend();

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const handleRemoveFriend = () => {
    removeFriendMutation.mutate(friend.id, {
      onSuccess: () => {
        setShowRemoveDialog(false);
      },
    });
  };

  return (
    <>
      <Card className="hover:shadow-lg transition-shadow duration-300 flex flex-col h-full">
        <CardHeader className="space-y-3">
          <div className="flex items-start justify-between gap-2">
            <div className="flex items-center gap-3 flex-1 min-w-0">
              <Avatar className="h-12 w-12 flex-shrink-0">
                <AvatarImage src={friend.avatar_url} alt={friend.name} />
                <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                  {getInitials(friend.name)}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <h3 className="text-lg font-bold truncate">{friend.name}</h3>
                {friend.email && (
                  <p className="text-sm text-muted-foreground truncate">
                    {friend.email}
                  </p>
                )}
              </div>
            </div>

            {/* Actions Menu */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="flex-shrink-0">
                  <MoreVertical className="h-4 w-4" />
                  <span className="sr-only">Open menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem asChild>
                  <Link href={`/profile/${friend.id}`} className="cursor-pointer">
                    <User className="mr-2 h-4 w-4" />
                    View Profile
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href={`/challenges/create?friendId=${friend.id}`} className="cursor-pointer">
                    <Swords className="mr-2 h-4 w-4" />
                    Challenge Friend
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setShowRemoveDialog(true)}
                  className="text-destructive focus:text-destructive cursor-pointer"
                >
                  <UserMinus className="mr-2 h-4 w-4" />
                  Remove Friend
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {/* Friend Since Badge */}
          <div>
            <Badge variant="outline" className="text-xs">
              Friends since {new Date(friend.friends_since).toLocaleDateString()}
            </Badge>
          </div>
        </CardHeader>

        <CardContent className="flex-1">
          <div className="text-sm text-muted-foreground">
            <p>Click to view profile or challenge your friend to a quiz!</p>
          </div>
        </CardContent>

        <CardFooter className="flex gap-2">
          <Link href={`/profile/${friend.id}`} className="flex-1">
            <Button variant="outline" className="w-full">
              <User className="mr-2 h-4 w-4" />
              View Profile
            </Button>
          </Link>
          <Link href={`/challenges/create?friendId=${friend.id}`} className="flex-1">
            <Button className="w-full">
              <Swords className="mr-2 h-4 w-4" />
              Challenge
            </Button>
          </Link>
        </CardFooter>
      </Card>

      {/* Remove Friend Confirmation Dialog */}
      <AlertDialog open={showRemoveDialog} onOpenChange={setShowRemoveDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Remove Friend?</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to remove <strong>{friend.name}</strong> from
              your friends list? You can always send them a friend request again later.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={removeFriendMutation.isPending}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleRemoveFriend}
              disabled={removeFriendMutation.isPending}
              className="bg-destructive hover:bg-destructive/90"
            >
              {removeFriendMutation.isPending ? "Removing..." : "Remove Friend"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
