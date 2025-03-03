package tryout

import (
	"net/http"
	"strconv"
	
	"github.com/cleoexcel/ristek-test/app/utils/enum"
	"github.com/cleoexcel/ristek-test/middleware"
	"github.com/gin-gonic/gin"
)

type TryoutHandler struct {
	service TryoutService
}

func NewTryoutHandler(service TryoutService) *TryoutHandler {
	return &TryoutHandler{service: service}
}

func (h *TryoutHandler) CreateTryout(c *gin.Context) {
	userId, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ! enum.IsValidCategory(input.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is invalid"})
		return
	}

	tryout, err := h.service.CreateTryout(input.Title, input.Description, userId, input.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tryout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tryout created successfully", "tryout": tryout})
}

func (h *TryoutHandler) GetAllTryout(c *gin.Context) {
	userId, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    title := c.DefaultQuery("title", "")
    category := c.DefaultQuery("category", "")
    createdAt := c.DefaultQuery("date", "")

    isByUser := c.DefaultQuery("is_by_user", "false")

    if isByUser == "false" {
         userId = 0
    }

    tryouts, err := h.service.GetAllTryout(title, category, createdAt, userId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tryouts"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"tryouts": tryouts})
}



func (h *TryoutHandler) GetDetailTryout(c *gin.Context) {
	idString := c.Param("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID"})
		return
	}
	tryout, err := h.service.GetDetailTryout(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tryout": tryout})
}

func (h *TryoutHandler) EditTryout(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID"})
		return
	}

	userIdString := middleware.ExtractUserID(c)
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tryout, err := h.service.EditTryout(id, input.Title, input.Description, userId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to edit this tryout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tryout updated successfully", "tryout": tryout})
}

func (h *TryoutHandler) DeleteTryoutById(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id"})
		return
	}

	err = h.service.DeleteTryoutById(int(idUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tryout not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tryout deleted successfully"})
}
