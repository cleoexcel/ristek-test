package submission

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cleoexcel/ristek-test/middleware"
	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	service SubmissionService
}

func NewSubmissionHandler(service SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	tryoutID, err := strconv.Atoi(c.Param("tryoutid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tryout ID"})
		return
	}

	var input struct {
		SubmittedAnswers []struct {
			QuestionID      int         `json:"question_id"`
			SubmittedAnswer interface{} `json:"submitted_answer"`
		} `json:"submitted_answers"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.SubmittedAnswers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No submitted answers provided"})
		return
	}

	userID, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	submission, err := h.service.CreateSubmission(tryoutID, userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create submission"})
		return
	}

	for i, submissionQuestion := range input.SubmittedAnswers {
		fmt.Println("Answer", i, "for Question ID:", submissionQuestion.QuestionID)

		submissionAnswer, err := h.service.CreateSubmissionAnswer(submission.ID, submissionQuestion.QuestionID, submissionQuestion.SubmittedAnswer)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Failed to create submission answer"})
			return
		}

		fmt.Println(" Submission Answer Created:", submissionAnswer)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Submission created successfully", "submission": submission})
}


func (h *SubmissionHandler) GetSubmissionByTryoutID(c *gin.Context) {
	tryoutID, err := strconv.Atoi(c.Param("tryoutid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tryout ID"})
		return
	}

	submission, err := h.service.GetSubmissionByTryoutID(tryoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submission": submission})
}

func (h *SubmissionHandler) GetAllAnswerBySubmissionID(c *gin.Context) {
	SubmissionID, err := strconv.Atoi(c.Param("submissionid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Submission ID"})
		return
	}

	answers, err := h.service.GetAllAnswersBySubmissionID(SubmissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answers not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answers": answers})
}

func (h *SubmissionHandler) CalculateScoreBySubmissionID(c *gin.Context) {
	SubmissionID, err := strconv.Atoi(c.Param("submissionid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Submission ID"})
		return
	}

	score, err := h.service.CalculateScoreBySubmissionID(SubmissionID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to calculate score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Score calculated successfully",
		"total_score": score,
	})
}
