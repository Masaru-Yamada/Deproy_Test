package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

const (
	MinDisplayID = 10000
	MaxDisplayID = 99990
)

// Organization is used by pop to map your organizations database table to your go code.
type Organization struct {
	ID               uuid.UUID `json:"id" db:"id"`
	CompanyID        uuid.UUID `json:"-" db:"company_id"`
	DisplayID        int       `json:"display_id" db:"display_id"`
	Name             string    `json:"name" db:"name"`
	ReportSendEmails string    `json:"report_send_emails" db:"report_send_emails"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Company          *Company  `json:"company,omitempty" belongs_to:"company"`
	Users            []User    `json:"users,omitempty" has_many:"users"`
	Roles            []Role    `json:"roles,omitempty" has_many:"roles"`
	Stories          []Story   `json:"stories,omitempty" has_many:"stories"`
}

func (o *Organization) Create(tx *pop.Connection) (*validate.Errors, error) {
	o.Name = strings.Trim(o.Name, " ")
	o.ReportSendEmails = strings.Trim(o.ReportSendEmails, " ")
	return tx.ValidateAndCreate(o)
}

// String is not required by pop and may be deleted
func (o Organization) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Organizations is not required by pop and may be deleted
type Organizations []Organization

// String is not required by pop and may be deleted
func (o Organizations) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *Organization) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *Organization) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *Organization) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
