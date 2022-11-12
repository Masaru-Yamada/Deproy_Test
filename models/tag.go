package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Tag is used by pop to map your tags database table to your go code.
type Tag struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	StoryTaggings []StoryTagging `json:"story_taggings,omitempty" has_many:"story_taggings"`
}

func (t *Tag) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(t)
}

// String is not required by pop and may be deleted
func (t Tag) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tags is not required by pop and may be deleted
type Tags []Tag

// String is not required by pop and may be deleted
func (t Tags) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tag) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tag) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tag) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
