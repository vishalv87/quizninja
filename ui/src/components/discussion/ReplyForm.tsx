"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent } from "@/components/ui/card";
import { AlertCircle, Loader2, Send } from "lucide-react";
import { useCreateReply } from "@/hooks/useDiscussions";
import {
  createDiscussionReplySchema,
  type CreateDiscussionReplyData,
} from "@/schemas/discussion";

interface ReplyFormProps {
  discussionId: string;
  onSuccess?: () => void;
}

export function ReplyForm({ discussionId, onSuccess }: ReplyFormProps) {
  const createReplyMutation = useCreateReply();
  const [isFocused, setIsFocused] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<CreateDiscussionReplyData>({
    resolver: zodResolver(createDiscussionReplySchema),
    defaultValues: {
      content: "",
    },
  });

  const content = watch("content");

  const onSubmit = (data: CreateDiscussionReplyData) => {
    createReplyMutation.mutate(
      {
        discussionId,
        data,
      },
      {
        onSuccess: () => {
          reset();
          setIsFocused(false);
          onSuccess?.();
        },
      }
    );
  };

  const handleCancel = () => {
    reset();
    setIsFocused(false);
  };

  return (
    <Card>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="reply-content">
              {isFocused ? "Write your reply" : "Add a reply"}
            </Label>
            <Textarea
              id="reply-content"
              placeholder="Share your thoughts..."
              className="min-h-[100px]"
              {...register("content")}
              onFocus={() => setIsFocused(true)}
            />
            {errors.content && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{errors.content.message}</AlertDescription>
              </Alert>
            )}
          </div>

          {isFocused && (
            <div className="flex items-center justify-between gap-2">
              <p className="text-sm text-muted-foreground">
                {content?.length || 0} / 2000 characters
              </p>
              <div className="flex gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleCancel}
                  disabled={createReplyMutation.isPending}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={
                    !content?.trim() || createReplyMutation.isPending
                  }
                >
                  {createReplyMutation.isPending ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Posting...
                    </>
                  ) : (
                    <>
                      <Send className="mr-2 h-4 w-4" />
                      Post Reply
                    </>
                  )}
                </Button>
              </div>
            </div>
          )}
        </form>
      </CardContent>
    </Card>
  );
}
