package repository

import (
	"database/sql"
	"fmt"
	"time"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
)

type DigestRepository struct {
	db *sql.DB
}

func NewDigestRepository() *DigestRepository {
	return &DigestRepository{
		db: database.DB,
	}
}

// GetTodaysDigest retrieves today's digest with articles
func (dr *DigestRepository) GetTodaysDigest() (*models.Digest, error) {
	today := time.Now().Format("2006-01-02")
	return dr.GetDigestByDate(today)
}

// GetDigestByDate retrieves a digest for a specific date with articles
func (dr *DigestRepository) GetDigestByDate(dateStr string) (*models.Digest, error) {
	// Parse the date string
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// Get the digest
	digest := &models.Digest{}
	query := `
		SELECT id, date, title, summary, article_count, is_dummy, created_at, updated_at, is_test_data
		FROM digests
		WHERE date = $1
	`
	err = dr.db.QueryRow(query, date).Scan(&digest.ID, &digest.Date, &digest.Title, &digest.Summary, &digest.ArticleCount, &digest.IsDummy, &digest.CreatedAt, &digest.UpdatedAt, &digest.IsTestData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("digest not found for date %s", dateStr)
		}
		return nil, fmt.Errorf("failed to get digest: %v", err)
	}

	// Get articles for this digest
	articles, err := dr.GetArticlesByDigestID(digest.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles for digest: %v", err)
	}

	digest.Articles = articles
	return digest, nil
}

// GetDigestByID retrieves a digest by ID with articles
func (dr *DigestRepository) GetDigestByID(digestID uuid.UUID) (*models.Digest, error) {
	digest := &models.Digest{}
	query := `
		SELECT id, date, title, summary, article_count, is_dummy, created_at, updated_at, is_test_data
		FROM digests
		WHERE id = $1
	`
	err := dr.db.QueryRow(query, digestID).Scan(&digest.ID, &digest.Date, &digest.Title, &digest.Summary, &digest.ArticleCount, &digest.IsDummy, &digest.CreatedAt, &digest.UpdatedAt, &digest.IsTestData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("digest not found")
		}
		return nil, fmt.Errorf("failed to get digest: %v", err)
	}

	// Get articles for this digest
	articles, err := dr.GetArticlesByDigestID(digest.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles for digest: %v", err)
	}

	digest.Articles = articles
	return digest, nil
}

// GetDigestList retrieves a paginated list of digests
func (dr *DigestRepository) GetDigestList(page, pageSize int) ([]models.Digest, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM digests"
	err := dr.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get digest count: %v", err)
	}

	// Get digests
	var digests []models.Digest
	query := `
		SELECT id, date, title, summary, article_count, is_dummy, created_at, updated_at, is_test_data
		FROM digests
		ORDER BY date DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := dr.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get digests: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var digest models.Digest
		err := rows.Scan(&digest.ID, &digest.Date, &digest.Title, &digest.Summary, &digest.ArticleCount, &digest.IsDummy, &digest.CreatedAt, &digest.UpdatedAt, &digest.IsTestData)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan digest: %v", err)
		}
		digests = append(digests, digest)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate digests: %v", err)
	}

	// Optionally load articles for each digest (can be disabled for performance)
	for i := range digests {
		articles, err := dr.GetArticlesByDigestID(digests[i].ID)
		if err != nil {
			// Log error but continue
			continue
		}
		digests[i].Articles = articles
	}

	return digests, totalCount, nil
}

// GetArticlesByDigestID retrieves all articles for a specific digest
func (dr *DigestRepository) GetArticlesByDigestID(digestID uuid.UUID) ([]models.Article, error) {
	var articles []models.Article
	query := `
		SELECT id, digest_id, title, content, summary, source, author, published_at,
		       category, image_url, external_url, read_time_minutes, is_breaking,
		       is_hot, is_dummy, is_trending, trending_score, trending_rank, created_at, is_test_data
		FROM digest_articles
		WHERE digest_id = $1
		ORDER BY
			CASE WHEN is_breaking = true THEN 0 ELSE 1 END,
			CASE WHEN is_hot = true THEN 0 ELSE 1 END,
			published_at DESC NULLS LAST,
			created_at DESC
	`
	rows, err := dr.db.Query(query, digestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.DigestID, &article.Title, &article.Content, &article.Summary, &article.Source, &article.Author, &article.PublishedAt, &article.Category, &article.ImageURL, &article.ExternalURL, &article.ReadTimeMinutes, &article.IsBreaking, &article.IsHot, &article.IsDummy, &article.IsTrending, &article.TrendingScore, &article.TrendingRank, &article.CreatedAt, &article.IsTestData)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %v", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate articles: %v", err)
	}

	return articles, nil
}

// GetArticleByID retrieves a single article by ID
func (dr *DigestRepository) GetArticleByID(articleID uuid.UUID) (*models.Article, error) {
	article := &models.Article{}
	query := `
		SELECT id, digest_id, title, content, summary, source, author, published_at,
		       category, image_url, external_url, read_time_minutes, is_breaking,
		       is_hot, is_dummy, is_trending, trending_score, trending_rank, created_at, is_test_data
		FROM digest_articles
		WHERE id = $1
	`
	err := dr.db.QueryRow(query, articleID).Scan(&article.ID, &article.DigestID, &article.Title, &article.Content, &article.Summary, &article.Source, &article.Author, &article.PublishedAt, &article.Category, &article.ImageURL, &article.ExternalURL, &article.ReadTimeMinutes, &article.IsBreaking, &article.IsHot, &article.IsDummy, &article.IsTrending, &article.TrendingScore, &article.TrendingRank, &article.CreatedAt, &article.IsTestData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to get article: %v", err)
	}

	return article, nil
}

// GetDigestCategories retrieves all unique categories from articles
func (dr *DigestRepository) GetDigestCategories() ([]string, error) {
	var categories []string
	query := `
		SELECT DISTINCT category
		FROM digest_articles
		WHERE is_dummy = false
		ORDER BY category
	`
	rows, err := dr.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %v", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate categories: %v", err)
	}

	return categories, nil
}

// CreateDigest creates a new digest
func (dr *DigestRepository) CreateDigest(digest *models.DigestRequest) (*models.Digest, error) {
	newDigest := &models.Digest{
		ID:      uuid.New(),
		Date:    digest.Date,
		Title:   digest.Title,
		Summary: digest.Summary,
		IsDummy: digest.IsDummy,
	}

	query := `
		INSERT INTO digests (id, date, title, summary, is_dummy, is_test_data)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING created_at, updated_at
	`
	err := dr.db.QueryRow(query, newDigest.ID, newDigest.Date, newDigest.Title,
		newDigest.Summary, newDigest.IsDummy).Scan(&newDigest.CreatedAt, &newDigest.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create digest: %v", err)
	}

	return newDigest, nil
}

// CreateArticle creates a new article
func (dr *DigestRepository) CreateArticle(article *models.ArticleRequest) (*models.Article, error) {
	newArticle := &models.Article{
		ID:               uuid.New(),
		DigestID:         article.DigestID,
		Title:            article.Title,
		Content:          article.Content,
		Summary:          article.Summary,
		Source:           article.Source,
		Author:           article.Author,
		PublishedAt:      article.PublishedAt,
		Category:         article.Category,
		ImageURL:         article.ImageURL,
		ExternalURL:      article.ExternalURL,
		ReadTimeMinutes:  article.ReadTimeMinutes,
		IsBreaking:       article.IsBreaking,
		IsHot:            article.IsHot,
		IsDummy:          article.IsDummy,
		IsTrending:       article.IsTrending,
		TrendingScore:    article.TrendingScore,
		TrendingRank:     article.TrendingRank,
	}

	query := `
		INSERT INTO digest_articles (id, digest_id, title, content, summary, source, author,
		                            published_at, category, image_url, external_url,
		                            read_time_minutes, is_breaking, is_hot, is_dummy,
		                            is_trending, trending_score, trending_rank, is_test_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, true)
		RETURNING created_at
	`
	err := dr.db.QueryRow(query, newArticle.ID, newArticle.DigestID, newArticle.Title,
		newArticle.Content, newArticle.Summary, newArticle.Source, newArticle.Author,
		newArticle.PublishedAt, newArticle.Category, newArticle.ImageURL, newArticle.ExternalURL,
		newArticle.ReadTimeMinutes, newArticle.IsBreaking, newArticle.IsHot,
		newArticle.IsDummy, newArticle.IsTrending, newArticle.TrendingScore,
		newArticle.TrendingRank).Scan(&newArticle.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %v", err)
	}

	return newArticle, nil
}

// UpdateDigest updates an existing digest
func (dr *DigestRepository) UpdateDigest(digestID uuid.UUID, update *models.UpdateDigestRequest) error {
	query := `
		UPDATE digests
		SET title = COALESCE($2, title),
		    summary = COALESCE($3, summary),
		    updated_at = NOW()
		WHERE id = $1
	`
	result, err := dr.db.Exec(query, digestID, update.Title, update.Summary)
	if err != nil {
		return fmt.Errorf("failed to update digest: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("digest not found")
	}

	return nil
}

// DeleteDigest deletes a digest and all its articles
func (dr *DigestRepository) DeleteDigest(digestID uuid.UUID) error {
	query := "DELETE FROM digests WHERE id = $1"
	result, err := dr.db.Exec(query, digestID)
	if err != nil {
		return fmt.Errorf("failed to delete digest: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("digest not found")
	}

	return nil
}

// DeleteArticle deletes an article
func (dr *DigestRepository) DeleteArticle(articleID uuid.UUID) error {
	query := "DELETE FROM digest_articles WHERE id = $1"
	result, err := dr.db.Exec(query, articleID)
	if err != nil {
		return fmt.Errorf("failed to delete article: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

// GetDigestStats retrieves statistics about the digest system
func (dr *DigestRepository) GetDigestStats() (*models.DigestStats, error) {
	stats := &models.DigestStats{}

	// Get basic counts
	query := `
		SELECT
			(SELECT COUNT(*) FROM digests) as total_digests,
			(SELECT COUNT(*) FROM digest_articles) as total_articles,
			(SELECT COUNT(DISTINCT category) FROM digest_articles) as total_categories
	`
	err := dr.db.QueryRow(query).Scan(&stats.TotalDigests, &stats.TotalArticles, &stats.TotalCategories)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %v", err)
	}

	// Calculate average articles per day
	if stats.TotalDigests > 0 {
		stats.AverageArticlesPerDay = float64(stats.TotalArticles) / float64(stats.TotalDigests)
	}

	// Get most active category
	categoryQuery := `
		SELECT category
		FROM digest_articles
		GROUP BY category
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`
	err = dr.db.QueryRow(categoryQuery).Scan(&stats.MostActiveCategory)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get most active category: %v", err)
	}

	// Get date ranges
	dateQuery := `
		SELECT MIN(date) as oldest_date, MAX(date) as latest_date
		FROM digests
	`
	err = dr.db.QueryRow(dateQuery).Scan(&stats.OldestDigestDate, &stats.LatestDigestDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get date ranges: %v", err)
	}

	return stats, nil
}

// GetOrCreateTodaysDigest gets today's digest or creates it if it doesn't exist
func (dr *DigestRepository) GetOrCreateTodaysDigest() (*models.Digest, error) {
	// Try to get today's digest first
	digest, err := dr.GetTodaysDigest()
	if err == nil {
		return digest, nil
	}

	// If not found, create today's digest
	today := time.Now()
	createRequest := &models.DigestRequest{
		Date:    today,
		Title:   "Daily News Digest",
		Summary: stringPtr("Stay informed with the most important news from around the world."),
		IsDummy: false,
	}

	newDigest, err := dr.CreateDigest(createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create today's digest: %v", err)
	}

	// Return the new digest with empty articles array
	newDigest.Articles = []models.Article{}
	return newDigest, nil
}

// GetTrendingArticles retrieves a paginated list of trending articles
func (dr *DigestRepository) GetTrendingArticles(page, pageSize int) ([]models.Article, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count of trending articles
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM digest_articles WHERE is_trending = true"
	err := dr.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get trending articles count: %v", err)
	}

	// Get trending articles
	var articles []models.Article
	query := `
		SELECT id, digest_id, title, content, summary, source, author, published_at,
		       category, image_url, external_url, read_time_minutes, is_breaking,
		       is_hot, is_dummy, is_trending, trending_score, trending_rank, created_at, is_test_data
		FROM digest_articles
		WHERE is_trending = true
		ORDER BY trending_rank ASC NULLS LAST, trending_score DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := dr.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get trending articles: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.DigestID, &article.Title, &article.Content, &article.Summary, &article.Source, &article.Author, &article.PublishedAt, &article.Category, &article.ImageURL, &article.ExternalURL, &article.ReadTimeMinutes, &article.IsBreaking, &article.IsHot, &article.IsDummy, &article.IsTrending, &article.TrendingScore, &article.TrendingRank, &article.CreatedAt, &article.IsTestData)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan trending article: %v", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate trending articles: %v", err)
	}

	return articles, totalCount, nil
}

// GetTrendingArticlesByCategory retrieves a paginated list of trending articles filtered by category
func (dr *DigestRepository) GetTrendingArticlesByCategory(category string, page, pageSize int) ([]models.Article, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count of trending articles in this category
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM digest_articles WHERE is_trending = true AND LOWER(category) = LOWER($1)"
	err := dr.db.QueryRow(countQuery, category).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get trending articles count for category: %v", err)
	}

	// Get trending articles for this category
	var articles []models.Article
	query := `
		SELECT id, digest_id, title, content, summary, source, author, published_at,
		       category, image_url, external_url, read_time_minutes, is_breaking,
		       is_hot, is_dummy, is_trending, trending_score, trending_rank, created_at, is_test_data
		FROM digest_articles
		WHERE is_trending = true AND LOWER(category) = LOWER($1)
		ORDER BY trending_rank ASC NULLS LAST, trending_score DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := dr.db.Query(query, category, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get trending articles by category: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.DigestID, &article.Title, &article.Content, &article.Summary, &article.Source, &article.Author, &article.PublishedAt, &article.Category, &article.ImageURL, &article.ExternalURL, &article.ReadTimeMinutes, &article.IsBreaking, &article.IsHot, &article.IsDummy, &article.IsTrending, &article.TrendingScore, &article.TrendingRank, &article.CreatedAt, &article.IsTestData)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan trending article: %v", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate trending articles by category: %v", err)
	}

	return articles, totalCount, nil
}

// UpdateTrendingRankings updates the trending rankings for all articles
func (dr *DigestRepository) UpdateTrendingRankings() error {
	_, err := dr.db.Exec("SELECT update_trending_rankings()")
	if err != nil {
		return fmt.Errorf("failed to update trending rankings: %v", err)
	}
	return nil
}


