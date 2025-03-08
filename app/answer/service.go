package answer

import (
	"log"

	"github.com/cleoexcel/ristek-test/app/models"
)

type AnswerService struct {
	Repo *AnswerRepository
}

func NewAnswerService(repo *AnswerRepository) *AnswerService {
	return &AnswerService{Repo: repo}
}

func (s *AnswerService) GetAllAnswers() (interface{}, error) {
	return s.Repo.GetAllAnswers()
}

func (s *AnswerService) CreateAnswer(questionID int, questionType string, expectAnswer interface{}, options []models.MultipleChoiceOption) (interface{}, error) {
	answer, err := s.Repo.CreateAnswer(questionID, questionType, expectAnswer, options)
	if err != nil {
		log.Printf("Error creating answer: %v", err)
		return nil, err
	}
	return answer, nil
}

func (s *AnswerService) UpdateAnswer(questionID int, expectAnswer interface{}, options []models.MultipleChoiceOption) (interface{}, error) {
	answer, err := s.Repo.UpdateAnswer(questionID, expectAnswer, options)
	if err != nil {
		log.Printf("Error updating answer: %v", err)
		return nil, err
	}
	return answer, nil
}

func (s *AnswerService) DeleteAnswer(questionID int, questionType string) error {
	err := s.Repo.DeleteAnswer(questionID, questionType)
	if err != nil {
		log.Printf("Error deleting answer: %v", err)
		return err
	}
	return nil
}
