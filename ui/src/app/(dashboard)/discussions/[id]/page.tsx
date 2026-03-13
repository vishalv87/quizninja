"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Separator } from "@/components/ui/separator";
import { ReplyCard } from "@/components/discussion/ReplyCard";
import { ReplyForm } from "@/components/discussion/ReplyForm";
import { CreateDiscussionDialog } from "@/components/discussion/CreateDiscussionDialog";
import {
  useDiscussion,
  useDiscussionReplies,
  useLikeDiscussion,
  useDeleteDiscussion,
} from "@/hooks/useDiscussions";
import { useAuth } from "@/hooks/useAuth";
import {
  MessageSquare,
  ArrowLeft,
  Heart,
  Clock,
  Edit,
  Trash2,
  AlertCircle,
} from "lucide-react";
import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
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

export default function DiscussionDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuth();
  const discussionId = params.id as string;

  const { data: discussion, isLoading, error } = useDiscussion(discussionId);
  const {
    data: replies = [],
    isLoading: repliesLoading,
  } = useDiscussionReplies(discussionId);

  const likeMutation = useLikeDiscussion();
  const deleteMutation = useDeleteDiscussion();

  const [showEditDialog, setShowEditDialog] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const isOwner = user?.id === discussion?.user_id;

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const handleLike = () => {
    if (discussion) {
      likeMutation.mutate(discussion.id);
    }
  };

  const handleEdit = () => {
    setShowEditDialog(true);
  };

  const handleDelete = () => {
    setShowDeleteDialog(true);
  };

  const confirmDelete = () => {
    if (discussion) {
      deleteMutation.mutate(discussion.id, {
        onSuccess: () => {
          router.push("/discussions");
        },
      });
    }
    setShowDeleteDialog(false);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <Skeleton className="h-8 w-32 mb-4" />
        <Skeleton className="h-12 w-64 mb-8" />
        <div className="space-y-4">
          <Skeleton className="h-64 w-full" />
          <Skeleton className="h-48 w-full" />
        </div>
      </div>
    );
  }

  // Error state
  if (error || !discussion) {
    return (
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>
            {error?.message || "Discussion not found"}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <>
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Discussions
        </Button>

        {/* Discussion Content */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1">
                <CardTitle className="text-3xl mb-2">
                  {discussion.title}
                </CardTitle>
                <div className="flex items-center gap-3 mt-3">
                  <Avatar className="h-10 w-10">
                    <AvatarImage
                      src={discussion.user.avatar_url}
                      alt={discussion.user.name}
                    />
                    <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                      {getInitials(discussion.user.name)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-semibold">{discussion.user.name}</p>
                    <p className="text-sm text-muted-foreground">
                      <Clock className="inline h-3 w-3 mr-1" />
                      {formatDistanceToNow(new Date(discussion.created_at), {
                        addSuffix: true,
                      })}
                      {discussion.updated_at !== discussion.created_at &&
                        " (edited)"}
                    </p>
                  </div>
                </div>
              </div>

              {/* Action Buttons */}
              {isOwner && (
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleEdit}
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleDelete}
                    disabled={deleteMutation.isPending}
                    className="text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </div>

            {/* Quiz Badge */}
            <div className="flex gap-2 mt-4">
              <Link href={`/quizzes/${discussion.quiz_id}`}>
                <Badge variant="secondary" className="cursor-pointer hover:bg-secondary/80">
                  Quiz: {discussion.quiz.title}
                </Badge>
              </Link>
            </div>
          </CardHeader>

          <CardContent>
            <p className="text-foreground whitespace-pre-wrap mb-6">
              {discussion.content}
            </p>

            <Separator className="my-4" />

            {/* Discussion Stats and Actions */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4 text-sm text-muted-foreground">
                <div className="flex items-center gap-1">
                  <MessageSquare className="h-4 w-4" />
                  <span>{discussion.replies_count} replies</span>
                </div>
                <div className="flex items-center gap-1">
                  <Heart className="h-4 w-4" />
                  <span>{discussion.likes_count} likes</span>
                </div>
              </div>

              <Button
                onClick={handleLike}
                variant="outline"
                size="sm"
                disabled={likeMutation.isPending}
              >
                <Heart className="mr-2 h-4 w-4" />
                {likeMutation.isPending ? "..." : "Like"}
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Reply Form */}
        <div className="mb-6">
          <ReplyForm discussionId={discussionId} />
        </div>

        {/* Replies */}
        <Card>
          <CardHeader>
            <CardTitle>
              Replies ({discussion.replies_count})
            </CardTitle>
            <CardDescription>
              Join the conversation and share your thoughts
            </CardDescription>
          </CardHeader>
          <CardContent>
            {repliesLoading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-32 w-full" />
                ))}
              </div>
            ) : replies.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <MessageSquare className="h-12 w-12 mx-auto mb-3 opacity-50" />
                <p>No replies yet. Be the first to reply!</p>
              </div>
            ) : (
              <div className="space-y-4">
                {replies.map((reply) => (
                  <ReplyCard
                    key={reply.id}
                    reply={reply}
                    discussionId={discussionId}
                  />
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Edit Dialog */}
      {discussion && (
        <CreateDiscussionDialog
          open={showEditDialog}
          onClose={() => setShowEditDialog(false)}
          discussion={discussion}
        />
      )}

      {/* Delete Confirmation */}
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
