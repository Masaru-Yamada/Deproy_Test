package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// UserScenarioHistory is used by pop to map your user_scenario_histories database table to your go code.
type UserScenarioHistory struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	ScenarioID uuid.UUID  `json:"scenario_id" db:"scenario_id"`
	TotalScore int        `json:"total_score" db:"total_score"`
	PlayedAt   nulls.Time `json:"played_at" db:"played_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	User       *User      `json:"user,omitempty" belongs_to:"user"`
	Scenario   *Scenario  `json:"scenario,omitempty" belongs_to:"scenario"`
}

func (u *UserScenarioHistory) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(u)
}

func (u *UserScenarioHistory) Update(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(u)
}

func (u *UserScenarioHistory) Delete(tx *pop.Connection) error {
	return tx.Where("id = ?", u.ID).Delete(u)
}

// String is not required by pop and may be deleted
func (u UserScenarioHistory) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// UserScenarioHistories is not required by pop and may be deleted
type UserScenarioHistories []UserScenarioHistory

// String is not required by pop and may be deleted
func (u UserScenarioHistories) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UserScenarioHistory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UserScenarioHistory) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UserScenarioHistory) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u *UserScenarioHistories) ListByUserID(tx *pop.Connection, userID uuid.UUID) error {

	tx.Where(`played_at is not null and user_id = ?`, userID).Order("scenario_id, total_score desc").All(u)
	if err := tx.All(u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
