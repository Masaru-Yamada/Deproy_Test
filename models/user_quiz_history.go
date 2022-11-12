package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// UserQuizHistory is used by pop to map your user_quiz_histories database table to your go code.
type UserQuizHistory struct {
	ID                    uuid.UUID           `json:"id" db:"id"`
	UserScenarioHistoryID uuid.UUID           `json:"user_scenario_history_id" db:"user_scenario_history_id"`
	UserID                uuid.UUID           `json:"user_id" db:"user_id"`
	ScenarioID            uuid.UUID           `json:"scenario_id" db:"scenario_id"`
	ScenarioQuizID        uuid.UUID           `json:"scenario_quiz_id" db:"scenario_quiz_id"`
	ScenarioQuizOptionID  uuid.UUID           `json:"scenario_quiz_option_id" db:"scenario_quiz_option_id"`
	Score                 int                 `json:"score" db:"score"`
	CreatedAt             time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at" db:"updated_at"`
	User                  *User               `json:"user,omitempty" belongs_to:"user"`
	Scenario              *Scenario           `json:"scenario,omitempty" belongs_to:"scenario"`
	ScenarioQuiz          *ScenarioQuiz       `json:"scenario_quiz,omitempty" belongs_to:"scenario_quiz"`
	ScenarioQuizOption    *ScenarioQuizOption `json:"scenario_quiz_option,omitempty" belongs_to:"scenario_quiz_option"`
}

func (u *UserQuizHistory) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(u)
}

// String is not required by pop and may be deleted
func (u UserQuizHistory) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// UserQuizHistories is not required by pop and may be deleted
type UserQuizHistories []UserQuizHistory

// String is not required by pop and may be deleted
func (u UserQuizHistories) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UserQuizHistory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UserQuizHistory) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UserQuizHistory) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u *UserQuizHistories) DeleteByUserScenarioHistoryID(tx *pop.Connection, userScenarioHistoryID uuid.UUID) error {
	return tx.Where("user_scenario_history_id = ?", userScenarioHistoryID).Delete(u)
}
