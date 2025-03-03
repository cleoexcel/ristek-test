package models

type Shortanswer struct {
	ID           int  `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionID   int  `json:"question_id"`
	ExpectAnswer string `json:"expectanswer"`
}

