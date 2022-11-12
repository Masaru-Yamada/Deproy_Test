package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Quiz is used by pop to map your quizzes database table to your go code.
type Quiz struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	Question        string         `json:"question" db:"question"`
	FailureMessage  nulls.String   `json:"failure_message" db:"failure_message"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	QuizOptions     []QuizOption   `json:"quiz_options,omitempty" has_many:"quiz_options"`
	ScenarioQuizzes []ScenarioQuiz `json:"scenario_quizzes,omitempty" has_many:"scenario_quizzes"`
}

func (q *Quiz) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(q)
}

// String is not required by pop and may be deleted
func (q Quiz) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// quizzes is not required by pop and may be deleted
type Quizzes []Quiz

// String is not required by pop and may be deleted
func (q Quizzes) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *Quiz) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *Quiz) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *Quiz) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
