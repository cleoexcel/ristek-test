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
	var input struct {
		TryoutID         int `json:"tryout_id"`
		SubmittedAnswers []struct {
			QuestionID      int         `json:"question_id"`
			SubmittedAnswer interface{} `json:"submitted_answer"` // interface spy bisa bool (T/F) atau string (short answer):D
		} `json:"answers_submitted"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(middleware.ExtractUserID(c))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	submission, err := h.service.CreateSubmission(input.TryoutID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
		return
	}

	fmt.Println("WOI ADA GA ", input.SubmittedAnswers)

	for id2, submissionQuestion := range input.SubmittedAnswers {
		submissionAnswer, err := h.service.CreateSubmissionAnswer(submission.ID, submissionQuestion.QuestionID, submissionQuestion.SubmittedAnswer)
		fmt.Println("SUBANS", id2, submissionAnswer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission answer"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission created successfully", "submission": submission})
}

func (h *SubmissionHandler) GetSubmissionByTryoutID(c *gin.Context) {
	tryoutID, err := strconv.Atoi(c.Param("id"))
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
	submissionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Submission ID"})
		return
	}

	answers, err := h.service.GetAllAnswersBySubmissionID(submissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answers not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answers": answers})
}