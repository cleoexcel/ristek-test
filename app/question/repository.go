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

func (r *QuestionRepository) GetAllQuestionsByTryoutID(tryoutID int) ([]*models.Question, error) {
	var questions []*models.Question
	err := r.DB.
	Preload("Tryout").
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

	var lastQuestion models.Question
	var newNumber int

	err = r.DB.Where("tryout_id = ?", tryoutID).Order("number DESC").First(&lastQuestion).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newNumber = 1
		} else {
			return nil, err
		}
	} else {
		newNumber = lastQuestion.Number + 1
	}

	question := &models.Question{
		Number:       newNumber,
		Content:      content,
		TryoutID:     tryoutID,
		QuestionType: questionType,
		Weight:       weight,
	}
	if err := r.DB.Create(question).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Tryout").First(&question, question.ID).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (r *QuestionRepository) EditQuestionByQuestionID(id int, content string, weight int) (*models.Question, error) {
	var question models.Question
	err := r.DB.Preload("Tryout").First(&question, id).Error
	if err != nil {
		return nil, err
	}

	var submission models.Submission
	err = r.DB.Where("tryout_id = ?", question.TryoutID).First(&submission).Error
	if err == nil {
		return nil, fmt.Errorf("cannot edit question for tryout that already has a submission")
	}

	question.Content = content
	question.Weight = weight
	if err := r.DB.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) DeleteQuestionByQuestionID(id int) error {
	var question models.Question
	err := r.DB.Preload("Tryout").First(&question, id).Error 
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
	err := r.DB.
		Preload("Tryout"). 
		Preload("ShortAnswer").
		Preload("TrueFalse").
		First(&question, id).
		Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}