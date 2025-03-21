package models

import (
	"time"
)

type SubmissionAnswerShortAnswer struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SubmissionID    int       `json:"submission_id" gorm:"index"`
	QuestionID      int       `json:"question_id" gorm:"index"`
	AnswerSubmitted string    `json:"answer_submitted"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	Question Question `json:"question" gorm:"foreignKey:QuestionID;references:ID"`
}
