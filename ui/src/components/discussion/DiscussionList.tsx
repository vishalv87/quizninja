"use client";

import type { Discussion } from "@/types/discussion";
import { DiscussionCard } from "./DiscussionCard";
import { EmptyState } from "@/components/common/EmptyState";
import { MessageSquare } from "lucide-react";
import { useState } from "react";
import { CreateDiscussionDialog } from "./CreateDiscussionDialog";

interface DiscussionListProps {
  discussions: Discussion[];
  isLoading?: boolean;
  onEditDiscussion?: (discussion: Discussion) => void;
}

export function DiscussionList({
  discussions,
  isLoading,
  onEditDiscussion,
}: DiscussionListProps) {
  const [discussionToEdit, setDiscussionToEdit] = useState<Discussion | null>(
    null
  );

  const handleEdit = (discussion: Discussion) => {
    setDiscussionToEdit(discussion);
    onEditDiscussion?.(discussion);
  };

  const handleCloseEdit = () => {
    setDiscussionToEdit(null);
  };

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="h-64 rounded-lg bg-muted animate-pulse"
          ></div>
        ))}
      </div>
    );
  }

  if (!discussions || discussions.length === 0) {
    return (
      <EmptyState
        icon={MessageSquare}
        title="No discussions yet"
        description="Be the first to start a discussion about this quiz!"
      />
    );
  }

  return (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {discussions.map((discussion) => (
          <DiscussionCard
            key={discussion.id}
            discussion={discussion}
            onEdit={handleEdit}
          />
        ))}
      </div>

      {discussionToEdit && (
        <CreateDiscussionDialog
          open={!!discussionToEdit}
          onClose={handleCloseEdit}
          discussion={discussionToEdit}
        />
      )}
    </>
  );
}
