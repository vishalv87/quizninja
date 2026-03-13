"use client";

import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import type { DiscussionReply } from "@/types/discussion";
import { Heart, Clock, Trash2, Edit } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import {
  useLikeReply,
  useDeleteReply,
  useUpdateReply,
} from "@/hooks/useDiscussions";
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
import { Textarea } from "@/components/ui/textarea";

interface ReplyCardProps {
  reply: DiscussionReply;
  discussionId: string;
}

export function ReplyCard({ reply, discussionId }: ReplyCardProps) {
  const { user } = useAuth();
  const likeMutation = useLikeReply();
  const deleteMutation = useDeleteReply();
  const updateMutation = useUpdateReply();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(reply.content);

  const isOwner = user?.id === reply.user_id;

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const handleLike = () => {
    likeMutation.mutate({ replyId: reply.id, discussionId });
  };

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
    setEditContent(reply.content);
  };

  const handleSaveEdit = () => {
    updateMutation.mutate(
      {
        replyId: reply.id,
        discussionId,
        data: { content: editContent },
      },
      {
        onSuccess: () => {
          setIsEditing(false);
        },
      }
    );
  };

  const handleDelete = () => {
    setShowDeleteDialog(true);
  };

  const confirmDelete = () => {
    deleteMutation.mutate({ replyId: reply.id, discussionId });
    setShowDeleteDialog(false);
  };

  return (
    <>
      <Card className="hover:shadow-md transition-shadow duration-200">
        <CardContent className="pt-6">
          {/* Header with user info */}
          <div className="flex items-start gap-3 mb-3">
            <Avatar className="h-8 w-8 flex-shrink-0">
              <AvatarImage
                src={reply.user.avatar_url}
                alt={reply.user.name}
              />
              <AvatarFallback className="bg-primary/10 text-primary text-xs font-semibold">
                {getInitials(reply.user.name)}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-semibold">{reply.user.name}</p>
              <p className="text-xs text-muted-foreground">
                {formatDistanceToNow(new Date(reply.created_at), {
                  addSuffix: true,
                })}
                {reply.updated_at !== reply.created_at && " (edited)"}
              </p>
            </div>
          </div>

          {/* Reply content */}
          {isEditing ? (
            <div className="space-y-2">
              <Textarea
                value={editContent}
                onChange={(e) => setEditContent(e.target.value)}
                className="min-h-[80px]"
                placeholder="Write your reply..."
              />
              <div className="flex gap-2 justify-end">
                <Button
                  onClick={handleCancelEdit}
                  variant="outline"
                  size="sm"
                >
                  Cancel
                </Button>
                <Button
                  onClick={handleSaveEdit}
                  size="sm"
                  disabled={
                    !editContent.trim() || updateMutation.isPending
                  }
                >
                  {updateMutation.isPending ? "Saving..." : "Save"}
                </Button>
              </div>
            </div>
          ) : (
            <p className="text-sm text-foreground whitespace-pre-wrap">
              {reply.content}
            </p>
          )}
        </CardContent>

        {!isEditing && (
          <CardFooter className="flex gap-2 pt-0">
            <Button
              onClick={handleLike}
              variant="ghost"
              size="sm"
              className="flex items-center gap-1"
              disabled={likeMutation.isPending}
            >
              <Heart className="h-4 w-4" />
              <span className="text-xs">{reply.likes_count}</span>
            </Button>

            {isOwner && (
              <>
                <Button
                  onClick={handleEdit}
                  variant="ghost"
                  size="sm"
                  className="ml-auto"
                >
                  <Edit className="h-4 w-4" />
                </Button>
                <Button
                  onClick={handleDelete}
                  variant="ghost"
                  size="sm"
                  className="text-destructive hover:text-destructive"
                  disabled={deleteMutation.isPending}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </>
            )}
          </CardFooter>
        )}
      </Card>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Reply</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this reply? This action cannot be
              undone.
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
