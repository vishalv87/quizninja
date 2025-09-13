package handlers

import (
	"net/http"

	"quizninja-api/config"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
)

type CategoriesHandler struct {
	categoriesRepo *repository.CategoriesRepository
	config         *config.Config
}

func NewCategoriesHandler(config *config.Config) *CategoriesHandler {
	return &CategoriesHandler{
		categoriesRepo: repository.NewCategoriesRepository(),
		config:         config,
	}
}

func (ch *CategoriesHandler) GetCategories(c *gin.Context) {
	categories, err := ch.categoriesRepo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
		"meta": gin.H{
			"total": len(categories),
		},
	})
}