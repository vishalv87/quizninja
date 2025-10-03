package handlers

import (
	"net/http"
	"strconv"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DigestHandler struct {
	digestRepo *repository.DigestRepository
	config     *config.Config
}

func NewDigestHandler(config *config.Config) *DigestHandler {
	return &DigestHandler{
		digestRepo: repository.NewDigestRepository(),
		config:     config,
	}
}

// GetTodaysDigest handles GET /digest/today
func (dh *DigestHandler) GetTodaysDigest(c *gin.Context) {
	digest, err := dh.digestRepo.GetTodaysDigest()
	if err != nil {
		c.JSON(http.StatusNotFound, models.DigestResponse{
			Success: false,
			Error:   "Today's digest not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.DigestResponse{
		Success: true,
		Digest:  *digest,
	})
}

// GetDigestByDate handles GET /digest/:date
func (dh *DigestHandler) GetDigestByDate(c *gin.Context) {
	dateStr := c.Param("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, models.DigestResponse{
			Success: false,
			Error:   "Date parameter is required",
		})
		return
	}

	digest, err := dh.digestRepo.GetDigestByDate(dateStr)
	if err != nil {
		c.JSON(http.StatusNotFound, models.DigestResponse{
			Success: false,
			Error:   "Digest not found for the specified date",
		})
		return
	}

	c.JSON(http.StatusOK, models.DigestResponse{
		Success: true,
		Digest:  *digest,
	})
}

// GetDigestList handles GET /digest
func (dh *DigestHandler) GetDigestList(c *gin.Context) {
	// Parse query parameters
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	digests, totalCount, err := dh.digestRepo.GetDigestList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DigestListResponse{
			Success: false,
			Error:   "Failed to fetch digest list",
		})
		return
	}

	hasMore := page*pageSize < totalCount

	c.JSON(http.StatusOK, models.DigestListResponse{
		Success:    true,
		Digests:    digests,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	})
}

// GetArticleByID handles GET /digest/articles/:id
func (dh *DigestHandler) GetArticleByID(c *gin.Context) {
	articleIDStr := c.Param("id")
	if articleIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ArticleResponse{
			Success: false,
			Error:   "Article ID is required",
		})
		return
	}

	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ArticleResponse{
			Success: false,
			Error:   "Invalid article ID format",
		})
		return
	}

	article, err := dh.digestRepo.GetArticleByID(articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ArticleResponse{
			Success: false,
			Error:   "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ArticleResponse{
		Success: true,
		Article: *article,
	})
}

// GetDigestCategories handles GET /digest/categories
func (dh *DigestHandler) GetDigestCategories(c *gin.Context) {
	categories, err := dh.digestRepo.GetDigestCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.CategoriesResponse{
			Success: false,
			Error:   "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, models.CategoriesResponse{
		Success:    true,
		Categories: categories,
	})
}

// CreateDigest handles POST /digest
func (dh *DigestHandler) CreateDigest(c *gin.Context) {
	var req models.DigestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.DigestResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	digest, err := dh.digestRepo.CreateDigest(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DigestResponse{
			Success: false,
			Error:   "Failed to create digest",
		})
		return
	}

	// Initialize empty articles array
	digest.Articles = []models.Article{}

	c.JSON(http.StatusCreated, models.DigestResponse{
		Success: true,
		Digest:  *digest,
	})
}

// CreateArticle handles POST /digest/:digestId/articles
func (dh *DigestHandler) CreateArticle(c *gin.Context) {
	digestIDStr := c.Param("digestId")
	if digestIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ArticleResponse{
			Success: false,
			Error:   "Digest ID is required",
		})
		return
	}

	digestID, err := uuid.Parse(digestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ArticleResponse{
			Success: false,
			Error:   "Invalid digest ID format",
		})
		return
	}

	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ArticleResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Set the digest ID from the URL parameter
	req.DigestID = digestID

	article, err := dh.digestRepo.CreateArticle(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ArticleResponse{
			Success: false,
			Error:   "Failed to create article",
		})
		return
	}

	c.JSON(http.StatusCreated, models.ArticleResponse{
		Success: true,
		Article: *article,
	})
}

// UpdateDigest handles PUT /digest/:id
func (dh *DigestHandler) UpdateDigest(c *gin.Context) {
	digestIDStr := c.Param("id")
	if digestIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Digest ID is required",
		})
		return
	}

	digestID, err := uuid.Parse(digestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid digest ID format",
		})
		return
	}

	var req models.UpdateDigestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	err = dh.digestRepo.UpdateDigest(digestID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update digest",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Digest updated successfully",
	})
}

// DeleteDigest handles DELETE /digest/:id
func (dh *DigestHandler) DeleteDigest(c *gin.Context) {
	digestIDStr := c.Param("id")
	if digestIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Digest ID is required",
		})
		return
	}

	digestID, err := uuid.Parse(digestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid digest ID format",
		})
		return
	}

	err = dh.digestRepo.DeleteDigest(digestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete digest",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Digest deleted successfully",
	})
}

// DeleteArticle handles DELETE /digest/articles/:id
func (dh *DigestHandler) DeleteArticle(c *gin.Context) {
	articleIDStr := c.Param("id")
	if articleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Article ID is required",
		})
		return
	}

	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid article ID format",
		})
		return
	}

	err = dh.digestRepo.DeleteArticle(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete article",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Article deleted successfully",
	})
}

// GetDigestStats handles GET /digest/stats
func (dh *DigestHandler) GetDigestStats(c *gin.Context) {
	stats, err := dh.digestRepo.GetDigestStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DigestStatsResponse{
			Success: false,
			Error:   "Failed to fetch digest statistics",
		})
		return
	}

	c.JSON(http.StatusOK, models.DigestStatsResponse{
		Success: true,
		Stats:   *stats,
	})
}

// GetOrCreateTodaysDigest handles GET /digest/today/ensure
func (dh *DigestHandler) GetOrCreateTodaysDigest(c *gin.Context) {
	digest, err := dh.digestRepo.GetOrCreateTodaysDigest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DigestResponse{
			Success: false,
			Error:   "Failed to get or create today's digest",
		})
		return
	}

	c.JSON(http.StatusOK, models.DigestResponse{
		Success: true,
		Digest:  *digest,
	})
}

// GetTrendingArticles handles GET /digest/trending
func (dh *DigestHandler) GetTrendingArticles(c *gin.Context) {
	// Parse query parameters
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	articles, totalCount, err := dh.digestRepo.GetTrendingArticles(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TrendingArticlesResponse{
			Success: false,
			Error:   "Failed to fetch trending articles",
		})
		return
	}

	hasMore := page*pageSize < totalCount

	c.JSON(http.StatusOK, models.TrendingArticlesResponse{
		Success:    true,
		Articles:   articles,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	})
}

// GetTrendingArticlesByCategory handles GET /digest/trending/:category
func (dh *DigestHandler) GetTrendingArticlesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, models.TrendingArticlesResponse{
			Success: false,
			Error:   "Category parameter is required",
		})
		return
	}

	// Parse query parameters
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	articles, totalCount, err := dh.digestRepo.GetTrendingArticlesByCategory(category, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TrendingArticlesResponse{
			Success: false,
			Error:   "Failed to fetch trending articles for category",
		})
		return
	}

	hasMore := page*pageSize < totalCount

	c.JSON(http.StatusOK, models.TrendingArticlesResponse{
		Success:    true,
		Articles:   articles,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    hasMore,
	})
}
