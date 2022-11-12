package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// StoryTagging is used by pop to map your story_taggings database table to your go code.
type StoryTagging struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StoryID   uuid.UUID `json:"story_id" db:"story_id"`
	TagID     uuid.UUID `json:"tag_id" db:"tag_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Story     *Story    `json:"story,omitempty" belongs_to:"story"`
	Tag       *Tag      `json:"tag,omitempty" belongs_to:"tag"`
}

func (s *StoryTagging) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(s)
}

// String is not required by pop and may be deleted
func (s StoryTagging) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// StoryTaggings is not required by pop and may be deleted
type StoryTaggings []StoryTagging

// String is not required by pop and may be deleted
func (s StoryTaggings) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *StoryTagging) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *StoryTagging) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *StoryTagging) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
