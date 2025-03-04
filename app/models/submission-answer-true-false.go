package models

import (
	"time"
)

type SubmissionAnswerTrueFalse struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SubmissionId    int       `json:"submission_id" gorm:"index"`
	QuestionID      int       `json:"question_id" gorm:"index"`
	AnswerSubmitted bool      `json:"answer_submitted"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`

	Submission Submission `json:"submission" gorm:"foreignKey:SubmissionId;references:ID"`
	Question   Question   `json:"question" gorm:"foreignKey:QuestionID;references:ID"`
}