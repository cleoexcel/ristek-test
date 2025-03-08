package question

import (
	"net/http"
	"strconv"

	"github.com/cleoexcel/ristek-test/app/models"
	"github.com/cleoexcel/ristek-test/middleware"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	Service *QuestionService
}

func NewQuestionHandler(service *QuestionService) *QuestionHandler {
	return &QuestionHandler{Service: service}
}

func (h *QuestionHandler) GetAllQuestionsByTryoutID(c *gin.Context) {
	tryoutID, err := strconv.Atoi(c.Param("tryoutid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID"})
		return
	}

	questions, err := h.Service.GetAllQuestionsByTryoutID(tryoutID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	userID, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tryoutID, err := strconv.Atoi(c.Param("tryoutid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tryout ID"})
		return
	}

	var input struct {
		Content      string      `json:"content"`
		QuestionType string      `json:"question_type"`
		Weight       int         `json:"weight"`
		ExpectAnswer interface{} `json:"expectanswer"`
		Options      []models.MultipleChoiceOption    `json:"options,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tryout models.Tryout
	if err := h.Service.Repo.DB.First(&tryout, tryoutID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tryout not found"})
		return
	}

	if tryout.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to create this question"})
		return
	}

	question, err := h.Service.CreateQuestion(input.Content, tryoutID, input.QuestionType, input.Weight, input.ExpectAnswer)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	
	answer, err := h.Service.AnswerService.CreateAnswer(question.ID, input.QuestionType, input.ExpectAnswer, input.Options)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Question created, but failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Question and answer created successfully",
		"question": question,
		"answer":   answer,
	})
}

func (h *QuestionHandler) EditQuestionByQuestionID(c *gin.Context) {
	userID, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var question models.Question
	if err := h.Service.Repo.DB.Preload("Tryout").Preload("ShortAnswer").Preload("TrueFalse").Preload("MultipleChoice").First(&question, questionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if question.Tryout.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this question"})
		return
	}

	var input struct {
		Content      string      `json:"content"`
		Weight       int         `json:"weight"`
		ExpectAnswer interface{} `json:"expectanswer"`
		Options      []models.MultipleChoiceOption    `json:"options,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.EditQuestionByQuestionID(questionID, input.Content, input.Weight, input.ExpectAnswer)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.Service.AnswerService.UpdateAnswer(questionID, input.ExpectAnswer, input.Options)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Question updated, but failed to update answer"})
		return
	}
	if mc, ok := answer.(*models.MultipleChoice); ok {
		if err := h.Service.Repo.DB.Preload("Options").First(&mc, mc.ID).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Failed to load options"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Question and answer updated successfully",
		"answer":  answer,
	})
}

func (h *QuestionHandler) DeleteQuestionByQuestionID(c *gin.Context) {
	userID, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var question models.Question
	if err := h.Service.Repo.DB.Preload("Tryout").First(&question, questionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if question.Tryout.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this question"})
		return
	}

	questionType := question.QuestionType
	err = h.Service.AnswerService.DeleteAnswer(questionID, questionType)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to delete answer"})
		return
	}
	err = h.Service.DeleteQuestionByQuestionID(questionID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Question and answer deleted successfully"})
}
