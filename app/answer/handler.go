package answer

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AnswerHandler struct {
	service Service
}

func NewAnswerHandler(service Service) *AnswerHandler {
	return &AnswerHandler{service: service}
}

func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	var input struct {
		QuestionID   int         `json:"question_id"`
		ExpectAnswer interface{} `json:"expectanswer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var answer interface{}
	var err error
	switch input.ExpectAnswer.(type) {
	case bool:
		answer, err = h.service.CreateTruefalseAnswer(input.QuestionID, input.ExpectAnswer.(bool))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Truefalse answer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Truefalse Answer created", "answer": answer})

	case string:
		answer, err = h.service.CreateShortanswerAnswer(input.QuestionID, input.ExpectAnswer.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Shortanswer answer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Shortanswer Answer created", "answer": answer})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer type"})
		return
	}
}

func (h *AnswerHandler) GetAllAnswers(c *gin.Context) {
	truefalseAnswers, shortanswerAnswers, err := h.service.GetAllAnswers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch answers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"truefalse_answers": truefalseAnswers,
		"shortanswer_answers": shortanswerAnswers,
	})
}

func (h *AnswerHandler) EditAnswer(c *gin.Context) {
	var input struct {
		ExpectAnswer interface{} `json:"expectanswer"` 
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedAnswer interface{}
	switch input.ExpectAnswer.(type) {
	case bool:
		updatedAnswer, err = h.service.EditTruefalseAnswer(id, input.ExpectAnswer.(bool))
	case string:
		updatedAnswer, err = h.service.EditShortanswerAnswer(id, input.ExpectAnswer.(string))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Answer updated successfully", "answer": updatedAnswer})
}

func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	if err := h.service.DeleteAnswer(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully"})
}
