"use client";

import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { Discussion } from "@/types/discussion";
import { MessageSquare, Heart, Clock, Trash2, Edit } from "lucide-react";
import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";
import { useLikeDiscussion, useDeleteDiscussion } from "@/hooks/useDiscussions";
import { formatDistanceToNow } from "date-fns";
import { useState } from "react";
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

interface DiscussionCardProps {
  discussion: Discussion;
  onEdit?: (discussion: Discussion) => void;
}

export function DiscussionCard({ discussion, onEdit }: DiscussionCardProps) {
  const { user } = useAuth();
  const likeMutation = useLikeDiscussion();
  const deleteMutation = useDeleteDiscussion();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const isOwner = user?.id === discussion.user_id;

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const handleLike = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    likeMutation.mutate(discussion.id);
  };

  const handleEdit = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onEdit?.(discussion);
  };

  const handleDelete = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowDeleteDialog(true);
  };

  const confirmDelete = () => {
    deleteMutation.mutate(discussion.id);
    setShowDeleteDialog(false);
  };

  return (
    <>
      <Link href={`/discussions/${discussion.id}`}>
        <Card className="hover:shadow-lg transition-shadow duration-300 cursor-pointer">
          <CardHeader className="space-y-3">
            {/* Header with user info */}
            <div className="flex items-start justify-between gap-2">
              <div className="flex items-center gap-3 flex-1 min-w-0">
                <Avatar className="h-10 w-10 flex-shrink-0">
                  <AvatarImage
                    src={discussion.user.avatar_url}
                    alt={discussion.user.name}
                  />
                  <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                    {getInitials(discussion.user.name)}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                  <h3 className="text-lg font-bold truncate">
                    {discussion.title}
                  </h3>
                  <p className="text-sm text-muted-foreground truncate">
                    by {discussion.user.name}
                  </p>
                </div>
              </div>
            </div>

            {/* Quiz badge */}
            <div className="flex gap-2 flex-wrap">
              <Badge variant="outline">{discussion.quiz.title}</Badge>
            </div>
          </CardHeader>

          <CardContent className="space-y-3">
            {/* Discussion content preview */}
            <p className="text-sm text-muted-foreground line-clamp-3">
              {discussion.content}
            </p>

            {/* Discussion info */}
            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <div className="flex items-center gap-1">
                <MessageSquare className="h-4 w-4" />
                <span>{discussion.replies_count} replies</span>
              </div>
              <div className="flex items-center gap-1">
                <Heart className="h-4 w-4" />
                <span>{discussion.likes_count} likes</span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="h-4 w-4" />
                <span>
                  {formatDistanceToNow(new Date(discussion.created_at), {
                    addSuffix: true,
                  })}
                </span>
              </div>
            </div>
          </CardContent>

          <CardFooter className="flex gap-2">
            <Button
              onClick={handleLike}
              variant="outline"
              size="sm"
              className="flex-1"
              disabled={likeMutation.isPending}
            >
              <Heart className="mr-2 h-4 w-4" />
              {likeMutation.isPending ? "..." : "Like"}
            </Button>

            {isOwner && (
              <>
                <Button
                  onClick={handleEdit}
                  variant="outline"
                  size="sm"
                  className="flex-shrink-0"
                >
                  <Edit className="h-4 w-4" />
                </Button>
                <Button
                  onClick={handleDelete}
                  variant="outline"
                  size="sm"
                  className="flex-shrink-0 text-destructive hover:text-destructive"
                  disabled={deleteMutation.isPending}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </>
            )}
          </CardFooter>
        </Card>
      </Link>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Discussion</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this discussion? This action
              cannot be undone and will also delete all replies.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}