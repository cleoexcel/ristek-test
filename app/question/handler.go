package question

import (
	"net/http"
	"strconv"

	"github.com/cleoexcel/ristek-test/app/models"
	"github.com/cleoexcel/ristek-test/app/utils/enum"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	service QuestionService
	repo    Repository
}

func NewQuestionHandler(service QuestionService, repo Repository) *QuestionHandler {
	return &QuestionHandler{service: service, repo: repo}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var input struct {
		Content      string `json:"content"`
		TryoutID     int    `json:"tryout_id"`
		QuestionType string `json:"question_type"`
		Weight       int    `json:"weight"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !enum.IsValidQuestionType(input.QuestionType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question type"})
		return
	}

	var tryout models.Tryout
	if err := h.repo.(*repository).DB.Where("id = ?", input.TryoutID).First(&tryout).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
		return
	}

	err := h.service.CreateQuestion( input.Content, input.TryoutID, input.QuestionType, input.Weight)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question created successfully"})
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	tryoutIDStr := c.Param("id")
	tryoutID, err := strconv.Atoi(tryoutIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID"})
		return
	}

	var tryout models.Tryout
	if err := h.repo.(*repository).DB.Where("id = ?", tryoutID).First(&tryout).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
		return
	}

	questions, err := h.service.GetAllQuestions(tryoutID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *QuestionHandler) EditQuestion(c *gin.Context) {
	var input struct {
		Content      string `json:"content"`
		QuestionType string `json:"question_type"`
		Weight       int    `json:"weight"`
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.service.EditQuestion(id, input.Content, input.QuestionType, input.Weight)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully", "question": question})
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := h.service.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
