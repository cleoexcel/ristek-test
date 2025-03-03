package models

import (
	"time"
)

type Question struct {
	ID           int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Number       int        `json:"number" gorm:"autoIncrement"`
	Content      string     `json:"content" gorm:"not null"`
	Weight       int        `json:"weight" gorm:"not null"`
	TryoutID     int        `json:"tryout_id"`
	QuestionType string     `json:"question_type"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
