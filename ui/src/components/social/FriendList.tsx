"use client";

import { FriendCard } from "./FriendCard";
import type { Friend } from "@/types/user";
import { Users } from "lucide-react";

interface FriendListProps {
  friends: Friend[];
}

export function FriendList({ friends }: FriendListProps) {
  // Handle empty state
  if (friends.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="rounded-full bg-muted p-6 mb-4">
          <Users className="h-12 w-12 text-muted-foreground" />
        </div>
        <h3 className="text-xl font-semibold mb-2">No Friends Yet</h3>
        <p className="text-muted-foreground max-w-md">
          Start building your network! Search for users and send friend requests to
          connect with other quiz enthusiasts.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Friends Count */}
      <div className="flex items-center gap-2">
        <Users className="h-5 w-5 text-muted-foreground" />
        <p className="text-sm text-muted-foreground">
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
