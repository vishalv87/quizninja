package models

import (
	"time"

	"github.com/google/uuid"
)

// Discussion represents a discussion/comment about a quiz or question
type Discussion struct {
	ID            uuid.UUID           `json:"id" db:"id"`
	QuizID        uuid.UUID           `json:"quiz_id" db:"quiz_id"`
	QuestionID    *uuid.UUID          `json:"question_id,omitempty" db:"question_id"`
	UserID        uuid.UUID           `json:"user_id" db:"user_id"`
	Content       string              `json:"content" db:"content"`
	LikesCount    int                 `json:"likes_count" db:"likes_count"`
	RepliesCount  int                 `json:"replies_count" db:"replies_count"`
	Type          string              `json:"type" db:"type"` // general, question, explanation, help
	CreatedAt     time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" db:"updated_at"`
	User          *User               `json:"user,omitempty"`
	Quiz          *QuizSummary        `json:"quiz,omitempty"`
	Question      *Question           `json:"question,omitempty"`
	TopReplies    []DiscussionReply   `json:"top_replies,omitempty"`
	IsLikedByUser bool                `json:"is_liked_by_user"`
}

// DiscussionReply represents a reply to a discussion
type DiscussionReply struct {
	ID            uuid.UUID `json:"id" db:"id"`
	DiscussionID  uuid.UUID `json:"discussion_id" db:"discussion_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Content       string    `json:"content" db:"content"`
	LikesCount    int       `json:"likes_count" db:"likes_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	User          *User     `json:"user,omitempty"`
	IsLikedByUser bool      `json:"is_liked_by_user"`
}

// DiscussionLike represents a like on a discussion
type DiscussionLike struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DiscussionID uuid.UUID `json:"discussion_id" db:"discussion_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DiscussionReplyLike represents a like on a discussion reply
type DiscussionReplyLike struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ReplyID   uuid.UUID `json:"reply_id" db:"reply_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Discussion DTOs

// CreateDiscussionRequest represents the request to create a new discussion
type CreateDiscussionRequest struct {
	QuizID     uuid.UUID  `json:"quiz_id" binding:"required"`
	QuestionID *uuid.UUID `json:"question_id,omitempty"`
	Content    string     `json:"content" binding:"required,min=1,max=2000"`
	Type       string     `json:"type" binding:"omitempty,oneof=general question explanation help"`
}

// CreateDiscussionReplyRequest represents the request to create a new reply
type CreateDiscussionReplyRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}

// UpdateDiscussionRequest represents the request to update a discussion
type UpdateDiscussionRequest struct {
	Content *string `json:"content,omitempty" binding:"omitempty,min=1,max=2000"`
	Type    *string `json:"type,omitempty" binding:"omitempty,oneof=general question explanation help"`
}

// UpdateDiscussionReplyRequest represents the request to update a reply
type UpdateDiscussionReplyRequest struct {
	Content *string `json:"content,omitempty" binding:"omitempty,min=1,max=2000"`
}

// DiscussionFilters represents filters for discussion queries
type DiscussionFilters struct {
	QuizID     *uuid.UUID `form:"quiz_id"`
	QuestionID *uuid.UUID `form:"question_id"`
	UserID     *uuid.UUID `form:"user_id"`
	Type       string     `form:"type" binding:"omitempty,oneof=general question explanation help"`
	Search     string     `form:"search"`
	Page       int        `form:"page,default=1" binding:"min=1"`
	PageSize   int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy     string     `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at likes_count replies_count"`
	SortOrder  string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// DiscussionListResponse represents the response for discussion list
type DiscussionListResponse struct {
	Discussions []Discussion `json:"discussions"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
	TotalPages  int          `json:"total_pages"`
}

// DiscussionDetailResponse represents the response for a single discussion detail
type DiscussionDetailResponse struct {
	Discussion Discussion `json:"discussion"`
}

// DiscussionRepliesResponse represents the response for discussion replies
type DiscussionRepliesResponse struct {
	Replies    []DiscussionReply `json:"replies"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// LikeResponse represents the response for like/unlike actions
type LikeResponse struct {
	IsLiked    bool `json:"is_liked"`
	LikesCount int  `json:"likes_count"`
	Message    string `json:"message"`
}

// DiscussionWithDetails represents a discussion with full user and quiz details
type DiscussionWithDetails struct {
	Discussion
	UserName      string `json:"user_name"`
	UserAvatar    string `json:"user_avatar"`
	QuizTitle     string `json:"quiz_title"`
	QuizCategory  string `json:"quiz_category"`
	QuestionText  string `json:"question_text,omitempty"`
}

// DiscussionReplyWithDetails represents a reply with full user details
type DiscussionReplyWithDetails struct {
	DiscussionReply
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
}

// DiscussionStatsResponse represents discussion statistics
type DiscussionStatsResponse struct {
	TotalDiscussions      int     `json:"total_discussions"`
	TotalReplies          int     `json:"total_replies"`
	TotalLikes            int     `json:"total_likes"`
	AverageRepliesPerPost float64 `json:"average_replies_per_post"`
	AverageLikesPerPost   float64 `json:"average_likes_per_post"`
	MostActiveQuiz        *struct {
		QuizID           uuid.UUID `json:"quiz_id"`
		QuizTitle        string    `json:"quiz_title"`
		DiscussionCount  int       `json:"discussion_count"`
	} `json:"most_active_quiz,omitempty"`
}