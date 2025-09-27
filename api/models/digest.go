package models

import (
	"time"

	"github.com/google/uuid"
)

// Digest represents a daily news digest
type Digest struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Date         time.Time `json:"date" db:"date"`
	Title        string    `json:"title" db:"title"`
	Summary      *string   `json:"summary" db:"summary"`
	ArticleCount int       `json:"article_count" db:"article_count"`
	IsDummy      bool      `json:"is_dummy" db:"is_dummy"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsTestData   bool      `json:"is_test_data" db:"is_test_data"`
	Articles     []Article `json:"articles,omitempty"`
}

// Article represents an article within a digest
type Article struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	DigestID         uuid.UUID  `json:"digest_id" db:"digest_id"`
	Title            string     `json:"title" db:"title"`
	Content          string     `json:"content" db:"content"`
	Summary          string     `json:"summary" db:"summary"`
	Source           *string    `json:"source" db:"source"`
	Author           *string    `json:"author" db:"author"`
	PublishedAt      *time.Time `json:"published_at" db:"published_at"`
	Category         string     `json:"category" db:"category"`
	ImageURL         *string    `json:"image_url" db:"image_url"`
	ExternalURL      *string    `json:"external_url" db:"external_url"`
	ReadTimeMinutes  *int       `json:"read_time_minutes" db:"read_time_minutes"`
	IsBreaking       bool       `json:"is_breaking" db:"is_breaking"`
	IsHot            bool       `json:"is_hot" db:"is_hot"`
	IsDummy          bool       `json:"is_dummy" db:"is_dummy"`
	IsTrending       bool       `json:"is_trending" db:"is_trending"`
	TrendingScore    float64    `json:"trending_score" db:"trending_score"`
	TrendingRank     *int       `json:"trending_rank" db:"trending_rank"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	IsTestData       bool       `json:"is_test_data" db:"is_test_data"`
}

// DigestResponse represents the API response for a single digest
type DigestResponse struct {
	Success bool   `json:"success"`
	Digest  Digest `json:"digest"`
	Error   string `json:"error,omitempty"`
}

// DigestListResponse represents the API response for multiple digests
type DigestListResponse struct {
	Success    bool     `json:"success"`
	Digests    []Digest `json:"digests"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	HasMore    bool     `json:"has_more"`
	Error      string   `json:"error,omitempty"`
}

// ArticleResponse represents the API response for a single article
type ArticleResponse struct {
	Success bool    `json:"success"`
	Article Article `json:"article"`
	Error   string  `json:"error,omitempty"`
}

// CategoriesResponse represents the API response for digest categories
type CategoriesResponse struct {
	Success    bool     `json:"success"`
	Categories []string `json:"categories"`
	Error      string   `json:"error,omitempty"`
}

// DigestRequest represents the request payload for creating/updating a digest
type DigestRequest struct {
	Date    time.Time `json:"date" binding:"required"`
	Title   string    `json:"title" binding:"required,min=1,max=255"`
	Summary *string   `json:"summary"`
	IsDummy bool      `json:"is_dummy"`
}

// ArticleRequest represents the request payload for creating/updating an article
type ArticleRequest struct {
	DigestID        uuid.UUID  `json:"digest_id" binding:"required"`
	Title           string     `json:"title" binding:"required,min=1,max=500"`
	Content         string     `json:"content" binding:"required"`
	Summary         string     `json:"summary" binding:"required"`
	Source          *string    `json:"source"`
	Author          *string    `json:"author"`
	PublishedAt     *time.Time `json:"published_at"`
	Category        string     `json:"category" binding:"required"`
	ImageURL        *string    `json:"image_url"`
	ExternalURL     *string    `json:"external_url"`
	ReadTimeMinutes *int       `json:"read_time_minutes"`
	IsBreaking      bool       `json:"is_breaking"`
	IsHot           bool       `json:"is_hot"`
	IsDummy         bool       `json:"is_dummy"`
	IsTrending      bool       `json:"is_trending"`
	TrendingScore   float64    `json:"trending_score"`
	TrendingRank    *int       `json:"trending_rank"`
}

// DigestWithStats represents a digest with additional statistics
type DigestWithStats struct {
	Digest
	TotalArticles              int       `json:"total_articles" db:"total_articles"`
	BreakingArticles           int       `json:"breaking_articles" db:"breaking_articles"`
	HotArticles                int       `json:"hot_articles" db:"hot_articles"`
	Categories                 *string   `json:"categories" db:"categories"`
	LatestArticlePublishedAt   *time.Time `json:"latest_article_published_at" db:"latest_article_published_at"`
}

// GetDigestListRequest represents query parameters for getting digest list
type GetDigestListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Category string `form:"category"`
}

// GetArticlesByDigestRequest represents query parameters for getting articles by digest
type GetArticlesByDigestRequest struct {
	Category   string `form:"category"`
	IsBreaking *bool  `form:"is_breaking"`
	IsHot      *bool  `form:"is_hot"`
	Limit      int    `form:"limit" binding:"min=1,max=100"`
	Offset     int    `form:"offset" binding:"min=0"`
}

// CreateDigestRequest represents the full request for creating a digest with articles
type CreateDigestRequest struct {
	DigestRequest
	Articles []ArticleRequest `json:"articles"`
}

// UpdateDigestRequest represents the request for updating a digest
type UpdateDigestRequest struct {
	Title   *string `json:"title"`
	Summary *string `json:"summary"`
}

// DigestStats represents statistics about the digest system
type DigestStats struct {
	TotalDigests           int      `json:"total_digests"`
	TotalArticles          int      `json:"total_articles"`
	TotalCategories        int      `json:"total_categories"`
	AverageArticlesPerDay  float64  `json:"average_articles_per_day"`
	MostActiveCategory     string   `json:"most_active_category"`
	LatestDigestDate       *time.Time `json:"latest_digest_date"`
	OldestDigestDate       *time.Time `json:"oldest_digest_date"`
}

// DigestStatsResponse represents the API response for digest statistics
type DigestStatsResponse struct {
	Success bool        `json:"success"`
	Stats   DigestStats `json:"stats"`
	Error   string      `json:"error,omitempty"`
}

// TrendingArticlesResponse represents the API response for trending articles
type TrendingArticlesResponse struct {
	Success    bool      `json:"success"`
	Articles   []Article `json:"articles"`
	TotalCount int       `json:"total_count"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	HasMore    bool      `json:"has_more"`
	Error      string    `json:"error,omitempty"`
}