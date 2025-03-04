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
	CreateSubmissionAnswer(submissionID int, questionID int, submittedAnswer interface{}) (interface{}, error)
	GetAllAnswersBySubmissionID(submissionID int) ([]interface{}, error)
	CalculateScoreBySubmissionID(submissionID int) (int, error)
}

type submissionRepository struct {
	DB *gorm.DB
	QuestionRepo question.QuestionRepository
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{
		DB: db, 
		QuestionRepo: *question.NewQuestionRepository(db),
	}
}

func (r *submissionRepository) CreateSubmission(tryoutID int, userID int) (*models.Submission, error) {
	var submission *models.Submission
	var numberOfAttempt int64

	r.DB.Model(&models.Submission{}).Where("tryout_id = ?", tryoutID).Count(&numberOfAttempt)

	submission = &models.Submission{
		TryoutID: tryoutID,
		UserID: userID,
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
		fmt.Println("Error fetching question:", err)
		return nil, err
	}
	fmt.Println("Question found:", question)

	questionType := question.QuestionType
	fmt.Println("Question Type:", questionType)

	if questionType == "TrueFalse" {
		submissionAnswerTF := &models.SubmissionAnswerTrueFalse{
			SubmissionId:    submissionID,
			QuestionID:      questionID,
			AnswerSubmitted: submittedAnswer.(bool),
		}

		if err := r.DB.Create(submissionAnswerTF).Error; err != nil {
			fmt.Println("Error inserting TrueFalse Answer:", err)
			return nil, err
		}
		fmt.Println("Inserted TrueFalse Answer:", submissionAnswerTF)

		if err := r.DB.Preload("Question").First(submissionAnswerTF, submissionAnswerTF.ID).Error; err != nil {
			fmt.Println("Error preloading TrueFalse Answer:", err)
			return nil, err
		}

		submissionAnswer = submissionAnswerTF
	} else {
		submissionAnswerShortAns := &models.SubmissionAnswerShortAnswer{
			SubmissionId:    submissionID,
			QuestionID:      questionID,
			AnswerSubmitted: submittedAnswer.(string),
		}

		if err := r.DB.Create(submissionAnswerShortAns).Error; err != nil {
			fmt.Println("Error inserting Short Answer:", err)
			return nil, err
		}
		fmt.Println("Inserted Short Answer:", submissionAnswerShortAns)

		if err := r.DB.Preload("Question").First(submissionAnswerShortAns, submissionAnswerShortAns.ID).Error; err != nil {
			fmt.Println("Error preloading Short Answer:", err)
			return nil, err
		}

		submissionAnswer = submissionAnswerShortAns
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

func (r *submissionRepository) GetAllAnswersBySubmissionID(submissionID int) ([]interface{}, error) {
	var trueFalseAnswers []models.SubmissionAnswerTrueFalse
	var shortAnswers []models.SubmissionAnswerShortAnswer
	var answers []interface{}

	if err := r.DB.Preload("Question").Where("submission_id = ?", submissionID).Find(&trueFalseAnswers).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Question").Where("submission_id = ?", submissionID).Find(&shortAnswers).Error; err != nil {
		return nil, err
	}

	for _, answer := range trueFalseAnswers {
		answers = append(answers, answer)
	}
	for _, answer := range shortAnswers {
		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *submissionRepository) CalculateScoreBySubmissionID(submissionID int) (int, error) {
	var submission models.Submission
	if err := r.DB.First(&submission, submissionID).Error; err != nil {
		return 0, err
	}

	answers, err := r.GetAllAnswersBySubmissionID(submissionID)
	if err != nil {
		return 0, err
	}

	totalScore := 0

	for _, answer := range answers {
		switch ans := answer.(type) {
		case models.SubmissionAnswerTrueFalse:
			var correctAnswer models.TrueFalse
			if err := r.DB.Where("question_id = ?", ans.QuestionID).First(&correctAnswer).Error; err != nil {
				continue
			}
			if ans.AnswerSubmitted == correctAnswer.ExpectAnswer {
				totalScore += ans.Question.Weight
			}

		case models.SubmissionAnswerShortAnswer:
			var correctAnswer models.ShortAnswer
			if err := r.DB.Where("question_id = ?", ans.QuestionID).First(&correctAnswer).Error; err != nil {
				continue
			}
			if ans.AnswerSubmitted == correctAnswer.ExpectAnswer {
				totalScore += ans.Question.Weight
			}
		}
	}

	submission.TotalScore = totalScore
	if err := r.DB.Save(&submission).Error; err != nil {
		return 0, err
	}

	return totalScore, nil
}