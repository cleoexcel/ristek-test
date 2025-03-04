package submission

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/cleoexcel/ristek-test/middleware"
)

type SubmissionHandler struct {
	service SubmissionService
}

func NewSubmissionHandler(service SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	var input struct {
		TryoutID int `json:"tryout_id"`
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