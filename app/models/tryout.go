package models

import (
	"time"
)

type Tryout struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"not null"`
	UserID      int        `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	Category    string     `json:"category"`
}
