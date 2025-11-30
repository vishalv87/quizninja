import type { DiscussionType } from '@/constants';

export interface Discussion {
  id: string;
  quiz_id: string;
  question_id?: string;
  user_id: string;
  title: string;
  content: string;
  type?: DiscussionType;
  likes_count: number;
  replies_count: number;
  is_liked_by_user?: boolean;
  created_at: string;
  updated_at: string;
  user: {
    id: string;
    name: string;
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
  is_liked_by_user?: boolean;
  created_at: string;
  updated_at: string;
  user: {
    id: string;
    name: string;
    avatar_url?: string;
  };
}

export interface CreateDiscussionRequest {
  quiz_id: string;
  question_id?: string;
  title: string;
  content: string;
  type?: DiscussionType;
}

export interface CreateDiscussionReplyRequest {
  content: string;
}