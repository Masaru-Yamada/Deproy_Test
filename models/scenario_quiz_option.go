package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type Status int

const (
	StatusFailure Status = iota
	StatusSuccess
	StatusNextQuiz
)

// ScenarioQuizOption is used by pop to map your scenario_quiz_options database table to your go code.
type ScenarioQuizOption struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	ScenarioQuizID     uuid.UUID         `json:"scenario_quiz_id" db:"scenario_quiz_id"`
	QuizOptionID       nulls.UUID        `json:"quiz_option_id" db:"quiz_option_id"`
	Answer             nulls.String      `json:"answer" db:"answer"`
	Score              int               `json:"score" db:"score"`
	NextScenarioQuizID nulls.UUID        `json:"next_scenario_quiz_id" db:"next_scenario_quiz_id"`
	Status             Status            `json:"status" db:"status"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at"`
	ScenarioQuiz       *ScenarioQuiz     `json:"scenario_quiz,omitempty" belongs_to:"scenario_quiz"`
	QuizOption         *QuizOption       `json:"quiz_option,omitempty" belongs_to:"quiz_option"`
	UserQuizHistories  []UserQuizHistory `json:"user_quiz_histories,omitempty" has_many:"user_quiz_histories"`
}

func (s *ScenarioQuizOption) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(s)
}

// String is not required by pop and may be deleted
func (s ScenarioQuizOption) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ScenarioQuizOptions is not required by pop and may be deleted
type ScenarioQuizOptions []ScenarioQuizOption

// String is not required by pop and may be deleted
func (s ScenarioQuizOptions) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ScenarioQuizOption) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ScenarioQuizOption) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ScenarioQuizOption) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
