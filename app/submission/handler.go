package submission

import (
	"net/http"
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
		TryoutID int `json:"tryout_id"`
		UserID   int `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submission, err := h.service.CreateSubmission(input.TryoutID, input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission created successfully", "submission": submission})
}
