"use client";

import { FriendCard } from "./FriendCard";
import type { Friend } from "@/types/user";
import { Users, Search } from "lucide-react";
import { Button } from "@/components/ui/button";

interface FriendListProps {
  friends: Friend[];
  onSearchClick?: () => void;
}

export function FriendList({ friends, onSearchClick }: FriendListProps) {
  // Handle empty state
  if (friends.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gradient-to-br from-violet-100 to-indigo-100 dark:from-violet-900/30 dark:to-indigo-900/30 mb-4 animate-pulse">
          <Users className="h-10 w-10 text-violet-500" />
        </div>
        <h3 className="text-xl font-semibold mb-2">No Friends Yet</h3>
        <p className="text-muted-foreground max-w-md mb-6">
          Start building your network! Search for users and send friend requests to
          connect with other quiz enthusiasts.
        </p>
        {onSearchClick && (
          <Button
            onClick={onSearchClick}
            className="rounded-xl bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 shadow-md shadow-indigo-500/20 transition-all duration-200 active:scale-95 hover:shadow-lg hover:shadow-indigo-500/30 hover:scale-105"
          >
            <Search className="mr-2 h-4 w-4" />
            Search for Friends
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Friends Count */}
      <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-violet-100/50 dark:bg-violet-900/20 border border-violet-200/50 dark:border-violet-800/50">
        <Users className="h-4 w-4 text-violet-500" />
        <p className="text-sm font-medium text-violet-700 dark:text-violet-400">
          {friends.length} {friends.length === 1 ? "friend" : "friends"}
        </p>
      </div>

      {/* Friends Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {friends.map((friend) => (
          <FriendCard key={friend.id} friend={friend} />
        ))}
      </div>
    </div>
  );
}
