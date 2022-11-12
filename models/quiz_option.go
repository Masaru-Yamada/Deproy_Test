package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// QuizOption is used by pop to map your quiz_options database table to your go code.
type QuizOption struct {
	ID                  uuid.UUID            `json:"id" db:"id"`
	QuizID              uuid.UUID            `json:"quiz_id" db:"quiz_id"`
	Answer              string               `json:"answer" db:"answer"`
	CreatedAt           time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" db:"updated_at"`
	Quiz                *Quiz                `json:"quiz,omitempty" belongs_to:"quiz"`
	ScenarioQuizOptions []ScenarioQuizOption `json:"scenario_quiz_options,omitempty" has_many:"scenario_quiz_options"`
}

func (q *QuizOption) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(q)
}

// String is not required by pop and may be deleted
func (q QuizOption) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// QuizOptions is not required by pop and may be deleted
type QuizOptions []QuizOption

// String is not required by pop and may be deleted
func (q QuizOptions) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *QuizOption) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *QuizOption) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *QuizOption) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
