package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Story is used by pop to map your stories database table to your go code.
type Story struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
	Title          string         `json:"title" db:"title"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	Organization   *Organization  `json:"organization,omitempty" belongs_to:"organization"`
	StoryTaggings  []StoryTagging `json:"story_taggings,omitempty" has_many:"story_taggings"`
	Scenarios      []Scenario     `json:"scenarios,omitempty" has_many:"scenarios"`
}

func (s *Story) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(s)
}

// String is not required by pop and may be deleted
func (s Story) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Stories is not required by pop and may be deleted
type Stories []Story

// String is not required by pop and may be deleted
func (s Stories) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Story) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Story) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Story) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
