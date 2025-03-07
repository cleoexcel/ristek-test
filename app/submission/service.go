package submission

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type SubmissionService interface {
	CreateSubmission(tryoutID int, userID int) (*models.Submission, error)
	GetSubmissionByTryoutID(tryoutID int) ([]models.Submission, error)
	CreateSubmissionAnswer(SubmissionID int, questionID int, submitted_answer interface{}) (interface{}, error)
	GetAllAnswersBySubmissionID(SubmissionID int) ([]interface{}, error)
	CalculateScoreBySubmissionID(SubmissionID int) (float64, error)
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

func (s *submissionService) CreateSubmissionAnswer(SubmissionID int, questionID int, submitted_answer interface{}) (interface{}, error) {
	return s.repo.CreateSubmissionAnswer(SubmissionID, questionID, submitted_answer)
}

func (s *submissionService) GetSubmissionByTryoutID(tryoutID int) ([]models.Submission, error) {
	return s.repo.GetSubmissionByTryoutID(tryoutID)
}

func (s *submissionService) GetAllAnswersBySubmissionID(SubmissionID int) ([]interface{}, error) {
	return s.repo.GetAllAnswersBySubmissionID(SubmissionID)
}

func (s *submissionService) CalculateScoreBySubmissionID(SubmissionID int) (float64, error) {
	return s.repo.CalculateScoreBySubmissionID(SubmissionID)
}
