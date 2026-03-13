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
import { MoreVertical, UserMinus, User, Calendar } from "lucide-react";
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
      <Card className="group relative flex flex-col h-full overflow-hidden border border-white/20 dark:border-white/10 shadow-md shadow-black/5 transition-all duration-300 hover:shadow-xl hover:shadow-violet-500/20 hover:-translate-y-1 hover:scale-[1.02] bg-white/90 dark:bg-background/90 backdrop-blur-sm rounded-2xl">
        {/* Top Decoration Bar */}
        <div className="h-1.5 w-full bg-gradient-to-r from-violet-400 to-indigo-500" />

        <CardHeader className="space-y-3">
          <div className="flex items-start justify-between gap-2">
            <div className="flex items-center gap-3 flex-1 min-w-0">
              <Avatar className="h-12 w-12 flex-shrink-0 ring-2 ring-violet-200/50 dark:ring-violet-800/50 transition-all duration-300 group-hover:ring-4 group-hover:ring-violet-300/50 dark:group-hover:ring-violet-700/50">
                <AvatarImage src={friend.avatar_url} alt={friend.name} />
                <AvatarFallback className="bg-gradient-to-br from-violet-500 to-indigo-500 text-white font-semibold">
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
                <Button variant="ghost" size="icon" className="flex-shrink-0 hover:bg-violet-100/50 dark:hover:bg-violet-900/30">
                  <MoreVertical className="h-4 w-4" />
                  <span className="sr-only">Open menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="rounded-xl">
                <DropdownMenuItem asChild>
                  <Link href={`/profile/${friend.id}`} className="cursor-pointer">
                    <User className="mr-2 h-4 w-4" />
                    View Profile
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
            <Badge variant="outline" className="text-xs border-violet-200/50 dark:border-violet-800/50 bg-violet-50/50 dark:bg-violet-900/20 gap-1">
              <Calendar className="h-3 w-3" />
              Friends since {new Date(friend.friends_since).toLocaleDateString()}
            </Badge>
          </div>
        </CardHeader>

        <CardContent className="flex-1">
          <div className="text-sm text-muted-foreground border-t border-gray-200/30 dark:border-gray-700/30 pt-3">
            <p>Click to view your friend&apos;s profile!</p>
          </div>
        </CardContent>

        <CardFooter className="flex gap-2">
          <Link href={`/profile/${friend.id}`} className="flex-1">
            <Button className="w-full rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 shadow-md shadow-indigo-500/20 transition-all duration-200 active:scale-95 hover:shadow-lg hover:shadow-indigo-500/30">
              <User className="mr-2 h-4 w-4" />
              View Profile
            </Button>
          </Link>
        </CardFooter>
      </Card>

      {/* Remove Friend Confirmation Dialog */}
      <AlertDialog open={showRemoveDialog} onOpenChange={setShowRemoveDialog}>
        <AlertDialogContent className="bg-white/95 dark:bg-slate-950/95 backdrop-blur-xl border border-white/20 dark:border-white/10 rounded-2xl shadow-2xl">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-xl">Remove Friend?</AlertDialogTitle>
            <AlertDialogDescription className="text-muted-foreground">
              Are you sure you want to remove <strong className="text-foreground">{friend.name}</strong> from
              your friends list? You can always send them a friend request again later.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              disabled={removeFriendMutation.isPending}
              className="rounded-xl transition-all duration-200 active:scale-95"
            >
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleRemoveFriend}
              disabled={removeFriendMutation.isPending}
              className="bg-destructive hover:bg-destructive/90 rounded-xl transition-all duration-200 active:scale-95 shadow-md shadow-destructive/20 hover:shadow-lg hover:shadow-destructive/30"
            >
              {removeFriendMutation.isPending ? "Removing..." : "Remove Friend"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
