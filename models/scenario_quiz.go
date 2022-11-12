package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// ScenarioQuiz is used by pop to map your scenario_quizzes database table to your go code.
type ScenarioQuiz struct {
	ID                  uuid.UUID            `json:"id" db:"id"`
	ScenarioID          uuid.UUID            `json:"scenario_id" db:"scenario_id"`
	QuizID              nulls.UUID           `json:"quiz_id" db:"quiz_id"`
	First               bool                 `json:"first" db:"first"`
	Question            nulls.String         `json:"question" db:"question"`
	FailureMessage      nulls.String         `json:"failure_message" db:"failure_message"`
	CreatedAt           time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" db:"updated_at"`
	Scenario            *Scenario            `json:"scenario,omitempty" belongs_to:"scenario"`
	Quiz                *Quiz                `json:"quiz,omitempty" belongs_to:"quiz"`
	ScenarioQuizOptions []ScenarioQuizOption `json:"scenario_quiz_options,omitempty" has_many:"scenario_quiz_options"`
	UserQuizHistories   []UserQuizHistory    `json:"user_quiz_histories,omitempty" has_many:"user_quiz_histories"`
}

func (s *ScenarioQuiz) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(s)
}

// String is not required by pop and may be deleted
func (s ScenarioQuiz) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ScenarioQuizzes is not required by pop and may be deleted
type ScenarioQuizzes []ScenarioQuiz

// String is not required by pop and may be deleted
func (s ScenarioQuizzes) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ScenarioQuiz) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ScenarioQuiz) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ScenarioQuiz) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
