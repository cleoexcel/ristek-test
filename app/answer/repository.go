package answer

import (
	"errors"
	"fmt"

	"github.com/cleoexcel/ristek-test/app/models"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	DB *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{DB: db}
}

func (r *AnswerRepository) GetAllAnswers() (map[string]interface{}, error) {
	var TrueFalseAnswers []models.TrueFalse
	var ShortAnswerAnswers []models.ShortAnswer

	if err := r.DB.Find(&TrueFalseAnswers).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Find(&ShortAnswerAnswers).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TrueFalse":   TrueFalseAnswers,
		"ShortAnswer": ShortAnswerAnswers,
	}, nil
}

func (r *AnswerRepository) CreateAnswer(questionID int, questionType string, expectAnswer interface{}) (interface{}, error) {
	if questionID <= 0 {
		return nil, errors.New("invalid question ID")
	}

	switch questionType {
	case "ShortAnswer":
		answer := &models.ShortAnswer{
			QuestionID:   questionID,
			ExpectAnswer: expectAnswer.(string),
		}
		if err := r.DB.Create(answer).Error; err != nil {
			return nil, err
		}
		return answer, nil

	case "TrueFalse":
		answer := &models.TrueFalse{
			QuestionID:   questionID,
			ExpectAnswer: expectAnswer.(bool),
		}
		if err := r.DB.Create(answer).Error; err != nil {
			return nil, err
		}
		return answer, nil

	default:
		return nil, errors.New("invalid question type")
	}
}

func (r *AnswerRepository) GetAnswer(questionID int, questionType string) (interface{}, error) {
	if questionID <= 0 {
		return nil, errors.New("invalid question ID")
	}

	switch questionType {
	case "ShortAnswer":
		var answer models.ShortAnswer
		if err := r.DB.Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil

	case "TrueFalse":
		var answer models.TrueFalse
		if err := r.DB.Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil

	default:
		return nil, errors.New("invalid question type")
	}
}

func (r *AnswerRepository) UpdateAnswer(questionID int, questionType string, expectAnswer interface{}) (interface{}, error) {
	if questionID <= 0 {
		return nil, errors.New("invalid question ID")
	}

	switch questionType {
	case "ShortAnswer":
		var answer models.ShortAnswer
		if err := r.DB.Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		answer.ExpectAnswer = expectAnswer.(string)
		if err := r.DB.Save(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil

	case "TrueFalse":
		var answer models.TrueFalse
		if err := r.DB.Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		answer.ExpectAnswer = expectAnswer.(bool)
		if err := r.DB.Save(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil

	default:
		return nil, errors.New("invalid question type")
	}
}

func (r *AnswerRepository) DeleteAnswer(questionID int, questionType string) error {
	if questionID <= 0 {
		return errors.New("invalid question ID")
	}
	fmt.Println("Question Type:", questionType)

	switch questionType {
	case "ShortAnswer":
		
		if err := r.DB.Where("question_id = ?", questionID).Delete(&models.ShortAnswer{}).Error; err != nil {
			return err
		}
	case "TrueFalse":
		
		if err := r.DB.Where("question_id = ?", questionID).Delete(&models.TrueFalse{}).Error; err != nil {
			return err
		}
	default:
		fmt.Println("Invalid question type:", questionType)
		fmt.Println("salaaahhh")
		return errors.New("invalid question type")
	}
	return nil
}
