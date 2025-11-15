export interface Discussion {
  id: string;
  quiz_id: string;
  user_id: string;
  title: string;
  content: string;
  likes_count: number;
  replies_count: number;
  created_at: string;
  updated_at: string;
  user: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  quiz: {
    id: string;
    title: string;
  };
}

export interface DiscussionReply {
  id: string;
  discussion_id: string;
  user_id: string;
  content: string;
  likes_count: number;
  created_at: string;
  updated_at: string;
  user: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
}

export interface CreateDiscussionRequest {
  quiz_id: string;
  title: string;
  content: string;
}

export interface CreateDiscussionReplyRequest {
  content: string;
}