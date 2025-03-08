package models

type MultipleChoice struct {
	ID         int                    `json:"id" gorm:"primaryKey;autoIncrement;index"`
	QuestionID int                    `json:"question_id"`
	Options    []MultipleChoiceOption `json:"options" gorm:"foreignKey:MultipleChoiceID;constraint:OnDelete:CASCADE;"`
}


type MultipleChoiceOption struct {
	ID               int    `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	MultipleChoiceID int    `json:"multiple_choice_id,omitempty"` 
	OptionText       string `json:"option_text" gorm:"not null"`
	IsCorrect        bool   `json:"is_correct"`
}
