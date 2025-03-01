package question

import (
	"github.com/cleoexcel/ristek-test/app/models"
)

type QuestionService interface {
	CreateQuestion(number int, content string, tryoutID int, questionType string) error
	GetAllQuestions(tryoutID int) ([]*models.Question, error)
	EditQuestion(id int, number int, content string, questionType string) (*models.Question, error)
	DeleteQuestion(id int) error
}

type questionService struct {
	repo Repository
}

func NewQuestionService(repo Repository) QuestionService {
	return &questionService{repo}
}

func (s *questionService) CreateQuestion(number int, content string, tryoutID int, questionType string) error {
	_, err := s.repo.CreateQuestion(number, content, tryoutID, questionType)
	if err != nil {
		return err
	}
	return nil
}

func (s *questionService) GetAllQuestions(tryoutID int) ([]*models.Question, error) {
	return s.repo.GetAllQuestions(tryoutID)
}

func (s *questionService) EditQuestion(id int, number int, content string, questionType string) (*models.Question, error) {
	return s.repo.EditQuestion(id, number, content, questionType)
}

func (s *questionService) DeleteQuestion(id int) error {
	return s.repo.DeleteQuestion(id)
}
