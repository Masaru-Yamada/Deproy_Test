package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID                     uuid.UUID               `json:"id" db:"id"`
	OrganizationID         uuid.UUID               `json:"organization_id" db:"organization_id"`
	CreatedAt              time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at" db:"updated_at"`
	Nickname               string                  `json:"nickname" db:"nickname"`
	Email                  string                  `json:"email" db:"email"`
	PasswordHash           string                  `json:"password_hash" db:"password_hash"`
	Organization           *Organization           `json:"organization,omitempty" belongs_to:"organization"`
	UserAuthenticationLogs []UserAuthenticationLog `json:"user_authentication_logs,omitempty" has_many:"user_authentication_logs"`
	UserScenarioHistories  []UserScenarioHistory   `json:"user_scenario_histories,omitempty" has_many:"user_scenario_histories"`
	UserQuizHistories      []UserQuizHistory       `json:"user_quiz_histories,omitempty" has_many:"user_quiz_histories"`

	DisplayID            int    `json:"-" db:"-"`
	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndCreate(u)
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.OrganizationID, Name: "OrganizationID"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email).Where("organization_id = ?", u.OrganizationID)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
