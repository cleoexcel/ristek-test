package submission

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type SubmissionService interface {
	CreateSubmission(tryoutID int, userID int) (*models.Submission, error)
	GetSubmissionByTryoutID(tryoutID int) ([]models.Submission, error)
	CreateSubmissionAnswer(submissionID int, questionID int, submitted_answer interface{}) (interface{}, error)
	GetAllAnswersBySubmissionID(submissionID int) ([]interface{}, error)
	CalculateScoreBySubmissionID(submissionID int) (float64, error)
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

func (s *submissionService) CreateSubmissionAnswer(submissionID int, questionID int, submitted_answer interface{}) (interface{}, error) {
	return s.repo.CreateSubmissionAnswer(submissionID, questionID, submitted_answer)
}

func (s *submissionService) GetSubmissionByTryoutID(tryoutID int) ([]models.Submission, error) {
	return s.repo.GetSubmissionByTryoutID(tryoutID)
}

func (s *submissionService) GetAllAnswersBySubmissionID(submissionID int) ([]interface{}, error) {
	return s.repo.GetAllAnswersBySubmissionID(submissionID)
}

func (s *submissionService) CalculateScoreBySubmissionID(submissionID int) (float64, error) {
	return s.repo.CalculateScoreBySubmissionID(submissionID)
}