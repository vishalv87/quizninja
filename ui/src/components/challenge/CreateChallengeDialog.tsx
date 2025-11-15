"use client";

import { useState } from "react";
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
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Swords, AlertCircle, Loader2 } from "lucide-react";
import { useFriends } from "@/hooks/useFriends";
import { useQuizzes } from "@/hooks/useQuizzes";
import { useCreateChallenge } from "@/hooks/useCreateChallenge";
import { createChallengeSchema, type CreateChallengeData } from "@/schemas/challenge";

interface CreateChallengeDialogProps {
  trigger?: React.ReactNode;
  defaultFriendId?: string;
  defaultQuizId?: string;
  onSuccess?: () => void;
}

export function CreateChallengeDialog({
  trigger,
  defaultFriendId,
  defaultQuizId,
  onSuccess,
}: CreateChallengeDialogProps) {
  const [open, setOpen] = useState(false);

  const { data: friends, isLoading: loadingFriends } = useFriends();
  const { data: quizzes, isLoading: loadingQuizzes } = useQuizzes();
  const createMutation = useCreateChallenge();

  const {
    setValue,
    watch,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateChallengeData>({
    resolver: zodResolver(createChallengeSchema),
    defaultValues: {
      opponent_id: defaultFriendId || "",
      quiz_id: defaultQuizId || "",
    },
  });

  const selectedFriendId = watch("opponent_id");
  const selectedQuizId = watch("quiz_id");

  const selectedFriend = friends?.find((f) => f.id === selectedFriendId);
  const selectedQuiz = quizzes?.find((q) => q.id === selectedQuizId);

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    const parts = name.split(" ");
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const onSubmit = (data: CreateChallengeData) => {
    createMutation.mutate(data, {
      onSuccess: () => {
        setOpen(false);
        reset();
        onSuccess?.();
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button>
            <Swords className="mr-2 h-4 w-4" />
            Create Challenge
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>Create Challenge</DialogTitle>
            <DialogDescription>
              Challenge a friend to compete on a quiz. Winner takes the glory!
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-6 py-6">
            {/* Friend Selection */}
            <div className="space-y-2">
              <Label htmlFor="opponent">Select Friend</Label>
              {loadingFriends ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !friends || friends.length === 0 ? (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    You need friends to create challenges. Add some friends first!
                  </AlertDescription>
                </Alert>
              ) : (
                <>
                  <Select
                    value={selectedFriendId}
                    onValueChange={(value) => setValue("opponent_id", value)}
                  >
                    <SelectTrigger id="opponent" className={errors.opponent_id ? "border-destructive" : ""}>
                      <SelectValue placeholder="Choose a friend..." />
                    </SelectTrigger>
                    <SelectContent>
                      {friends.map((friend) => (
                        <SelectItem key={friend.id} value={friend.id}>
                          <div className="flex items-center gap-2">
                            <span>{friend.name}</span>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {errors.opponent_id && (
                    <p className="text-sm text-destructive">{errors.opponent_id.message}</p>
                  )}
                  {selectedFriend && (
                    <div className="flex items-center gap-3 p-3 bg-muted rounded-lg">
                      <Avatar className="h-10 w-10">
                        <AvatarImage src={selectedFriend.avatar_url} alt={selectedFriend.name} />
                        <AvatarFallback>{getInitials(selectedFriend.name)}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium">{selectedFriend.name}</p>
                        <p className="text-sm text-muted-foreground">{selectedFriend.email}</p>
                      </div>
                    </div>
                  )}
                </>
              )}
            </div>

            {/* Quiz Selection */}
            <div className="space-y-2">
              <Label htmlFor="quiz">Select Quiz</Label>
              {loadingQuizzes ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !quizzes || quizzes.length === 0 ? (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    No quizzes available at the moment.
                  </AlertDescription>
                </Alert>
              ) : (
                <>
                  <Select
                    value={selectedQuizId}
                    onValueChange={(value) => setValue("quiz_id", value)}
                  >
                    <SelectTrigger id="quiz" className={errors.quiz_id ? "border-destructive" : ""}>
                      <SelectValue placeholder="Choose a quiz..." />
                    </SelectTrigger>
                    <SelectContent>
                      {quizzes.map((quiz) => (
                        <SelectItem key={quiz.id} value={quiz.id}>
                          {quiz.title}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {errors.quiz_id && (
                    <p className="text-sm text-destructive">{errors.quiz_id.message}</p>
                  )}
                  {selectedQuiz && (
                    <div className="p-3 bg-muted rounded-lg space-y-2">
                      <p className="font-medium">{selectedQuiz.title}</p>
                      <p className="text-sm text-muted-foreground line-clamp-2">
                        {selectedQuiz.description}
                      </p>
                      <div className="flex gap-2">
                        <Badge variant="outline">{selectedQuiz.category}</Badge>
                        <Badge variant="outline">{selectedQuiz.difficulty}</Badge>
                        <Badge variant="outline">{selectedQuiz.question_count} questions</Badge>
                      </div>
                    </div>
                  )}
                </>
              )}
            </div>

            {/* Info Alert */}
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Both you and your friend will take the quiz. The player with the highest score wins!
              </AlertDescription>
            </Alert>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={createMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={
                createMutation.isPending ||
                !selectedFriendId ||
                !selectedQuizId ||
                loadingFriends ||
                loadingQuizzes
              }
            >
              {createMutation.isPending ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Sending...
                </>
              ) : (
                <>
                  <Swords className="mr-2 h-4 w-4" />
                  Send Challenge
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
