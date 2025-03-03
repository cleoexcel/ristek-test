package question

import (
	
	"net/http"
	"strconv"
	"github.com/cleoexcel/ristek-test/app/models"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	Service *QuestionService
}

func NewQuestionHandler(service *QuestionService) *QuestionHandler {
	return &QuestionHandler{Service: service}
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	tryoutID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tryout ID"})
		return
	}

	questions, err := h.Service.GetAllQuestions(tryoutID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var input struct {
		Content      string      `json:"content"`
		TryoutID     int         `json:"tryout_id"`
		QuestionType string      `json:"question_type"`
		Weight       int         `json:"weight"`
		ExpectAnswer interface{} `json:"expectanswer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.Service.CreateQuestion(input.Content, input.TryoutID, input.QuestionType, input.Weight, input.ExpectAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.Service.AnswerService.CreateAnswer(question.ID, input.QuestionType, input.ExpectAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Question created, but failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Question and answer created successfully",
		"question": question,
		"answer": answer,
	})
}

func (h *QuestionHandler) EditQuestion(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var input struct {
		Content      string      `json:"content"`
		Weight       int         `json:"weight"`
		QuestionType string      `json:"question_type"`
		ExpectAnswer interface{} `json:"expectanswer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.EditQuestion(questionID, input.Content, input.QuestionType, input.Weight, input.ExpectAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.Service.AnswerService.UpdateAnswer(questionID, input.QuestionType, input.ExpectAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Question updated, but failed to update answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question and answer updated successfully",
		"answer": answer,
	})
}


func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var question models.Question
	if err := h.Service.Repo.DB.First(&question, questionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	questionType := question.QuestionType
	
	err = h.Service.AnswerService.DeleteAnswer(questionID, questionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answer"})
		return
	}
	err = h.Service.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question and answer deleted successfully"})
}



