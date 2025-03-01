package submission

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type SubmissionService interface {
	CreateSubmission(tryoutID int, userID int) (*models.Submission, error)
	GetSubmissionByTryoutID(tryoutID int) (*models.Submission, error)
}

type submissionService struct {
	repo SubmissionRepository
}

func NewSubmissionService(repo SubmissionRepository) SubmissionService {
	return &submissionService{repo}
}

func (s *submissionService) CreateSubmission(tryoutID int, userID int) (*models.Submission, error) {
	return s.repo.CreateSubmission(tryoutID, userID)
}

func (s *submissionService) GetSubmissionByTryoutID(tryoutID int) (*models.Submission, error) {
	return s.repo.GetSubmissionByTryoutID(tryoutID)
}
