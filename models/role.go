package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Role is used by pop to map your roles database table to your go code.
type Role struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	OrganizationID uuid.UUID     `json:"organization_id" db:"organization_id"`
	Name           string        `json:"name" db:"name"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	Organization   *Organization `json:"organization,omitempty" belongs_to:"organization"`
	Scenarios      []Scenario    `json:"scenarios,omitempty" has_many:"scenarios"`
}

func (r *Role) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(r)
}

// String is not required by pop and may be deleted
func (r Role) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Roles is not required by pop and may be deleted
type Roles []Role

// String is not required by pop and may be deleted
func (r Roles) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Role) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Role) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Role) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
