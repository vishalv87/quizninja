"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { MessageSquare, AlertCircle, Loader2 } from "lucide-react";
import { useQuizzes } from "@/hooks/useQuizzes";
import {
  useCreateDiscussion,
  useUpdateDiscussion,
} from "@/hooks/useDiscussions";
import {
  createDiscussionSchema,
  type CreateDiscussionData,
} from "@/schemas/discussion";
import type { Discussion } from "@/types/discussion";

interface CreateDiscussionDialogProps {
  trigger?: React.ReactNode;
  defaultQuizId?: string;
  discussion?: Discussion;
  open?: boolean;
  onClose?: () => void;
  onSuccess?: () => void;
}

export function CreateDiscussionDialog({
  trigger,
  defaultQuizId,
  discussion,
  open: controlledOpen,
  onClose,
  onSuccess,
}: CreateDiscussionDialogProps) {
  const [internalOpen, setInternalOpen] = useState(false);
  const isControlled = controlledOpen !== undefined;
  const open = isControlled ? controlledOpen : internalOpen;
  const setOpen = isControlled ? (value: boolean) => {
    if (!value) onClose?.();
  } : setInternalOpen;

  const { data: quizzes, isLoading: loadingQuizzes } = useQuizzes();
  const createMutation = useCreateDiscussion();
  const updateMutation = useUpdateDiscussion();

  const isEditing = !!discussion;

  const {
    register,
    setValue,
    watch,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateDiscussionData>({
    resolver: zodResolver(createDiscussionSchema),
    defaultValues: {
      quiz_id: defaultQuizId || discussion?.quiz_id || "",
      title: discussion?.title || "",
      content: discussion?.content || "",
    },
  });

  // Reset form when discussion changes
  useEffect(() => {
    if (discussion) {
      setValue("quiz_id", discussion.quiz_id);
      setValue("title", discussion.title);
      setValue("content", discussion.content);
    }
  }, [discussion, setValue]);

  const selectedQuizId = watch("quiz_id");
  const selectedQuiz = quizzes?.find((q) => q.id === selectedQuizId);

  const onSubmit = (data: CreateDiscussionData) => {
    if (isEditing && discussion) {
      updateMutation.mutate(
        {
          id: discussion.id,
          data: {
            title: data.title,
            content: data.content,
          },
        },
        {
          onSuccess: () => {
            setOpen(false);
            reset();
            onSuccess?.();
          },
        }
      );
    } else {
      createMutation.mutate(data, {
        onSuccess: () => {
          setOpen(false);
          reset();
          onSuccess?.();
        },
      });
    }
  };

  const handleClose = () => {
    setOpen(false);
    reset();
    onClose?.();
  };

  const isPending = createMutation.isPending || updateMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      {!isControlled && (
        <DialogTrigger asChild>
          {trigger || (
            <Button>
              <MessageSquare className="mr-2 h-4 w-4" />
              Start Discussion
            </Button>
          )}
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? "Edit Discussion" : "Start a Discussion"}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update your discussion details below."
              : "Share your thoughts about a quiz and start a discussion."}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Quiz Selection */}
          {!isEditing && (
            <div className="space-y-2">
              <Label htmlFor="quiz_id">
                Quiz <span className="text-destructive">*</span>
              </Label>
              {loadingQuizzes ? (
                <div className="flex items-center gap-2 p-3 border rounded-md">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span className="text-sm text-muted-foreground">
                    Loading quizzes...
                  </span>
                </div>
              ) : (
                <Select
                  value={selectedQuizId}
                  onValueChange={(value) => setValue("quiz_id", value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a quiz" />
                  </SelectTrigger>
                  <SelectContent>
                    {quizzes?.map((quiz) => (
                      <SelectItem key={quiz.id} value={quiz.id}>
                        <div className="flex items-center gap-2">
                          <span>{quiz.title}</span>
                          <Badge variant="outline" className="text-xs">
                            {quiz.difficulty}
                          </Badge>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
              {errors.quiz_id && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{errors.quiz_id.message}</AlertDescription>
                </Alert>
              )}
            </div>
          )}

          {/* Selected Quiz Display (when editing) */}
          {isEditing && discussion && (
            <div className="space-y-2">
              <Label>Quiz</Label>
              <div className="flex items-center gap-2 p-3 border rounded-md bg-muted/50">
                <span className="font-medium">{discussion.quiz.title}</span>
              </div>
            </div>
          )}

          {/* Title */}
          <div className="space-y-2">
            <Label htmlFor="title">
              Title <span className="text-destructive">*</span>
            </Label>
            <Input
              id="title"
              placeholder="Enter discussion title"
              {...register("title")}
            />
            {errors.title && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{errors.title.message}</AlertDescription>
              </Alert>
            )}
          </div>

          {/* Content */}
          <div className="space-y-2">
            <Label htmlFor="content">
              Content <span className="text-destructive">*</span>
            </Label>
            <Textarea
              id="content"
              placeholder="Share your thoughts, questions, or insights..."
              className="min-h-[150px]"
              {...register("content")}
            />
            {errors.content && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{errors.content.message}</AlertDescription>
              </Alert>
            )}
          </div>

          {/* Preview */}
          {selectedQuiz && (
            <div className="space-y-2">
              <Label>Selected Quiz</Label>
              <div className="flex items-center gap-2 p-3 border rounded-md bg-muted/50">
                <span className="font-medium">{selectedQuiz.title}</span>
                <Badge variant="outline">{selectedQuiz.difficulty}</Badge>
                <Badge variant="secondary">{selectedQuiz.category}</Badge>
              </div>
            </div>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {isPending
                ? isEditing
                  ? "Updating..."
                  : "Creating..."
                : isEditing
                ? "Update Discussion"
                : "Create Discussion"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}