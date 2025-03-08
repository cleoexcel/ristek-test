package answer

import (
	"errors"
	

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
	var MultipleChoiceAnswers []models.MultipleChoice

	if err := r.DB.Find(&TrueFalseAnswers).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Find(&ShortAnswerAnswers).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Preload("Options").Find(&MultipleChoiceAnswers).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TrueFalse":   TrueFalseAnswers,
		"ShortAnswer": ShortAnswerAnswers,
		"MultipleChoice": MultipleChoiceAnswers,
	}, nil
}

func (r *AnswerRepository) CreateAnswer(questionID int, questionType string, expectAnswer interface{}, options []models.MultipleChoiceOption) (interface{}, error) {
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
		
		case "MultipleChoice":
			if len(options) <= 1 {
				return nil, errors.New("there should be atleast 2 options")
			}
			answer := &models.MultipleChoice{
				QuestionID: questionID,
			}

			if err := r.DB.Create(answer).Error; err != nil {
				return nil, err
			}

			for _, option := range options {
				multipleChoiceOption := &models.MultipleChoiceOption{
					MultipleChoiceID: answer.ID,
					OptionText: option.OptionText,
					IsCorrect: option.IsCorrect,
				}
	
				if err := r.DB.Create(multipleChoiceOption).Error; err != nil {
					return nil, err
				}
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

	case "MultipleChoice":
		var answer models.MultipleChoice
		if err := r.DB.Preload("Options").Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil

	default:
		return nil, errors.New("invalid question type")
	}
}

func (r *AnswerRepository) UpdateAnswer(questionID int, expectAnswer interface{}, options []models.MultipleChoiceOption) (interface{}, error) {
	if questionID <= 0 {
		return nil, errors.New("invalid question ID")
	}

	var question models.Question
	if err := r.DB.First(&question, questionID).Error; err != nil {
		return nil, errors.New("question not found")
	}

	switch question.QuestionType { 
	case "ShortAnswer":
		var answer models.ShortAnswer
		if err := r.DB.Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		if val, ok := expectAnswer.(string); ok {
			answer.ExpectAnswer = val
		} else {
			return nil, errors.New("invalid type for ShortAnswer")
		}
		if err := r.DB.Save(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil

	case "TrueFalse":
		var answer models.TrueFalse
		if err := r.DB.Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
		if val, ok := expectAnswer.(bool); ok {
			answer.ExpectAnswer = val
		} else {
			return nil, errors.New("invalid type for TrueFalse")
		}
		if err := r.DB.Save(&answer).Error; err != nil {
			return nil, err
		}
		return &answer, nil
	case "MultipleChoice":
		var answer models.MultipleChoice
		if err := r.DB.Preload("Options").Where("question_id = ?", questionID).First(&answer).Error; err != nil {
			return nil, err
		}
	
		existingOptions := make(map[int]*models.MultipleChoiceOption)
		for i := range answer.Options {
			existingOptions[answer.Options[i].ID] = &answer.Options[i]
		}
	
		updatedOptionIDs := make(map[int]bool)
	
		for _, option := range options {
			if option.ID != 0 { 
				if existingOpt, exists := existingOptions[option.ID]; exists {
					existingOpt.OptionText = option.OptionText
					existingOpt.IsCorrect = option.IsCorrect
					if err := r.DB.Save(existingOpt).Error; err != nil {
						return nil, err
					}
					updatedOptionIDs[option.ID] = true
				}
			}
		}
	
		if err := r.DB.Preload("Options").First(&answer, answer.ID).Error; err != nil {
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
	switch questionType {
	case "ShortAnswer":
		
		if err := r.DB.Where("question_id = ?", questionID).Delete(&models.ShortAnswer{}).Error; err != nil {
			return err
		}
	case "TrueFalse":
		
		if err := r.DB.Where("question_id = ?", questionID).Delete(&models.TrueFalse{}).Error; err != nil {
			return err
		}
	case "MultipleChoice":
		var multipleChoice models.MultipleChoice
		
		if err := r.DB.Where("question_id = ?", questionID).First(&multipleChoice).Error; err != nil {
			return err
		}

		if err := r.DB.Where("multiple_choice_id = ?", multipleChoice.ID).Delete(&models.MultipleChoiceOption{}).Error; err != nil {
			return err
		}
		if err := r.DB.Delete(&multipleChoice).Error; err != nil {
			return err
		}
	default:
		return errors.New("invalid question type")
	}
	return nil
}
