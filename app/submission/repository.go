package submission

import (
	"fmt"

	"github.com/cleoexcel/ristek-test/app/models"
	"github.com/cleoexcel/ristek-test/app/question"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	CreateSubmission(tryoutID int, userID int) (*models.Submission, error)
	GetSubmissionByTryoutID(tryoutID int) ([]models.Submission, error)
	CreateSubmissionAnswer(SubmissionID int, questionID int, submittedAnswer interface{}) (interface{}, error)
	GetAllAnswersBySubmissionID(SubmissionID int) ([]interface{}, error)
	CalculateScoreBySubmissionID(SubmissionID int) (float64, error)
}

type submissionRepository struct {
	DB           *gorm.DB
	QuestionRepo question.QuestionRepository
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{
		DB:           db,
		QuestionRepo: *question.NewQuestionRepository(db),
	}
}

func (r *submissionRepository) CreateSubmission(tryoutID int, userID int) (*models.Submission, error) {
	var submission *models.Submission
	var numberOfAttempt int64

	r.DB.Model(&models.Submission{}).Where("tryout_id = ?", tryoutID).Count(&numberOfAttempt)

	submission = &models.Submission{
		TryoutID:        tryoutID,
		UserID:          userID,
		NumberOfAttempt: int(numberOfAttempt) + 1,
	}

	if err := r.DB.Create(submission).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Tryout").First(submission, submission.ID).Error; err != nil {
		return nil, err
	}

	return submission, nil
}

func (r *submissionRepository) CreateSubmissionAnswer(submissionID int, questionID int, submittedAnswer interface{}) (interface{}, error) {
	var submissionAnswer interface{}

	question, err := r.QuestionRepo.GetQuestionByID(questionID)
	if err != nil {
		return nil, err
	}
	questionType := question.QuestionType

	switch questionType {
	case "TrueFalse":
		answer, ok := submittedAnswer.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool for TrueFalse question")
		}

		submissionAnswerTF := &models.SubmissionAnswerTrueFalse{
			SubmissionID:    submissionID,
			QuestionID:      questionID,
			AnswerSubmitted: answer,
		}

		if err := r.DB.Create(submissionAnswerTF).Error; err != nil {
			return nil, err
		}
		submissionAnswer = submissionAnswerTF

	case "ShortAnswer":
		answer, ok := submittedAnswer.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for ShortAnswer question")
		}

		submissionAnswerShortAns := &models.SubmissionAnswerShortAnswer{
			SubmissionID:    submissionID,
			QuestionID:      questionID,
			AnswerSubmitted: answer,
		}

		if err := r.DB.Create(submissionAnswerShortAns).Error; err != nil {
			return nil, err
		}
		submissionAnswer = submissionAnswerShortAns

	case "MultipleChoice":
		optionIDFloat, ok := submittedAnswer.(float64) 
		if !ok {
			return nil, fmt.Errorf("expected number for MultipleChoice question")
		}
		optionID := int(optionIDFloat)

		submissionMCAns := &models.SubmissionAnswerMultipleChoice{
			SubmissionID:           submissionID,
			QuestionID:             questionID,
			MultipleChoiceOptionID: optionID,
		}
		if err := r.DB.Create(submissionMCAns).Error; err != nil {
			return nil, err
		}
		submissionAnswer = submissionMCAns

	default:
		return nil, fmt.Errorf("invalid question type")
	}

	return submissionAnswer, nil
}


func (r *submissionRepository) GetSubmissionByTryoutID(tryoutID int) ([]models.Submission, error) {
	var submissions []models.Submission
	err := r.DB.Preload("Tryout").Where("tryout_id = ?", tryoutID).Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *submissionRepository) GetAllAnswersBySubmissionID(SubmissionID int) ([]interface{}, error) {
	var trueFalseAnswers []models.SubmissionAnswerTrueFalse
	var shortAnswers []models.SubmissionAnswerShortAnswer
	var multipleChoiceAnswers []models.SubmissionAnswerMultipleChoice
	var answers []interface{}

	var submission models.Submission
	if err := r.DB.First(&submission, SubmissionID).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Question").Preload("Question.Tryout").
		Where("submission_id = ?", SubmissionID).Find(&trueFalseAnswers).Error; err != nil {
	} else {
		answers = append(answers, trueFalseAnswers)
	}

	if err := r.DB.Preload("Question").Preload("Question.Tryout").
		Where("submission_id = ?", SubmissionID).Find(&shortAnswers).Error; err != nil {
	} else {
		answers = append(answers, shortAnswers)
	}

	if err := r.DB.Preload("Question").Preload("Question.Tryout").Preload("MultipleChoiceOption").
		Where("submission_id = ?", SubmissionID).Find(&multipleChoiceAnswers).Error; err != nil {
	} else {
		answers = append(answers, multipleChoiceAnswers)
	}

	return answers, nil
}


func (r *submissionRepository) CalculateScoreBySubmissionID(submissionID int) (float64, error) {
	var submission models.Submission
	if err := r.DB.First(&submission, submissionID).Error; err != nil {
		return 0, err
	}

	answers, err := r.GetAllAnswersBySubmissionID(submissionID)
	if err != nil {
		return 0, err
	}

	if len(answers) == 0 {
		return 0, fmt.Errorf("no answers found")
	}

	totalScore := 0
	totalWeight := 0

	for _, answer := range answers {
		switch ans := answer.(type) {

		case []models.SubmissionAnswerTrueFalse:
			var correctAnswer models.TrueFalse
			if err := r.DB.Where("question_id = ?", ans[0].QuestionID).First(&correctAnswer).Error; err != nil {
				continue
			}

			var question models.Question
			if err := r.DB.Where("id = ?", ans[0].QuestionID).First(&question).Error; err != nil {
				continue
			}

			totalWeight += question.Weight
			if ans[0].AnswerSubmitted == correctAnswer.ExpectAnswer {
				totalScore += question.Weight
			}

		case []models.SubmissionAnswerShortAnswer:
			var correctAnswer models.ShortAnswer
			if err := r.DB.Where("question_id = ?", ans[0].QuestionID).First(&correctAnswer).Error; err != nil {
				continue
			}

			var question models.Question
			if err := r.DB.Where("id = ?", ans[0].QuestionID).First(&question).Error; err != nil {
				continue
			}


			totalWeight += question.Weight
			if ans[0].AnswerSubmitted == correctAnswer.ExpectAnswer {
				totalScore += question.Weight
			}

		case []models.SubmissionAnswerMultipleChoice:
			var correctOption models.MultipleChoiceOption
			if err := r.DB.Where("id = ?", ans[0].MultipleChoiceOptionID).First(&correctOption).Error; err != nil {
				continue
			}

			var question models.Question
			if err := r.DB.Where("id = ?", ans[0].QuestionID).First(&question).Error; err != nil {
				continue
			}


			totalWeight += question.Weight
			if correctOption.IsCorrect {
				totalScore += question.Weight
			}
		}
	}

	if totalWeight == 0 {
		return 0, fmt.Errorf("total weight is zero, invalid scoring")
	}

	finalScore := (float64(totalScore) / float64(totalWeight)) * 100

	submission.TotalScore = finalScore
	if err := r.DB.Save(&submission).Error; err != nil {
		return 0, err
	}

	return finalScore, nil
}
