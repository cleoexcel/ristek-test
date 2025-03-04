package question

import (
	"fmt"

	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	DB *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{DB: db}
}

func (r *QuestionRepository) GetAllQuestions(tryoutID int) ([]*models.Question, error) {
	var questions []*models.Question
	err := r.DB.
		Preload("ShortAnswer").
		Preload("TrueFalse").
		Where("tryout_id = ?", tryoutID).
		Find(&questions).
		Error
	return questions, err
}

func (r *QuestionRepository) CreateQuestion(content string, tryoutID int, questionType string, weight int) (*models.Question, error) {
	var submission models.Submission
	err := r.DB.Where("tryout_id = ?", tryoutID).First(&submission).Error
	if err == nil {
		return nil, fmt.Errorf("tryout already has a submission, cannot add or edit questions")
	}
	question := &models.Question{
		Content:      content,
		TryoutID:     tryoutID,
		QuestionType: questionType,
		Weight:       weight,
	}
	if err := r.DB.Create(question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (r *QuestionRepository) EditQuestion(id int, content string, questionType string, weight int) (*models.Question, error) {
	var question models.Question
	err := r.DB.First(&question, id).Error
	if err != nil {
		return nil, err
	}

	var submission models.Submission
	err = r.DB.Where("tryout_id = ?", question.TryoutID).First(&submission).Error
	if err == nil {
		return nil, fmt.Errorf("cannot edit question for tryout that already has a submission")
	}

	question.Content = content
	question.QuestionType = questionType
	question.Weight = weight
	if err := r.DB.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) DeleteQuestion(id int) error {
	var question models.Question
	err := r.DB.First(&question, id).Error
	if err != nil {
		return err
	}

	var submission models.Submission
	err = r.DB.Where("tryout_id = ?", question.TryoutID).First(&submission).Error
	if err == nil {
		return fmt.Errorf("cannot delete question for tryout that already has a submission")
	}

	err = r.DB.Delete(&question).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *QuestionRepository) GetQuestionByID(id int) (*models.Question, error) {
	var question models.Question
	err := r.DB.First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}