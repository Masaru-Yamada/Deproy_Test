package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Scenario is used by pop to map your scenarios database table to your go code.
type Scenario struct {
	ID                    uuid.UUID             `json:"id" db:"id"`
	StoryID               uuid.UUID             `json:"story_id" db:"story_id"`
	RoleID                uuid.UUID             `json:"role_id" db:"role_id"`
	Overview              string                `json:"overview" db:"overview"`
	Description           string                `json:"description" db:"description"`
	HighestScore          int                   `json:"highest_score" db:"highest_score"`
	ResultMessage         string                `json:"result_message" db:"result_message"`
	CreatedAt             time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at" db:"updated_at"`
	Story                 *Story                `json:"story,omitempty" belongs_to:"story"`
	Role                  *Role                 `json:"role,omitempty" belongs_to:"role"`
	ScenarioQuizzes       []ScenarioQuiz        `json:"scenario_quizzes,omitempty" has_many:"scenario_quizzes"`
	UserScenarioHistories []UserScenarioHistory `json:"user_scenario_histories,omitempty" has_many:"user_scenario_histories"`
	UserQuizHistories     []UserQuizHistory     `json:"user_quiz_histories,omitempty" has_many:"user_quiz_histories"`
}

func (s *Scenario) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(s)
}

// String is not required by pop and may be deleted
func (s Scenario) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Scenarios is not required by pop and may be deleted
type Scenarios []Scenario

// String is not required by pop and may be deleted
func (s Scenarios) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Scenario) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Scenario) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Scenario) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
