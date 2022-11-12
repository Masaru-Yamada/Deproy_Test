package grifts

import (
	"fmt"
	"strings"

	"github.com/aiit2022-pbl-okuhara/play-security/models"

	. "github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const (
	AdminDisplayID = 99999
	TestDisplayID  = 99990
)

// buffalo task list
var _ = Namespace("db", func() {

	// buffalo task db:clear
	Desc("clear", "Clear a database")
	Add("clear", func(c *Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			// all data clear
			if err := tx.TruncateAll(); err != nil {
				return err
			}
			return nil
		})
	})

	// buffalo task db:admin email=admin@example.com password=12345
	Desc("admin", "Seeds a database for admin")
	Add("admin", func(c *Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			if len(c.Args) != 2 {
				return errors.New("Invalid args")
			}

			user := &models.User{}
			for _, e := range c.Args {
				pair := strings.Split(e, "=")
				if len(pair) != 2 {
					return errors.New("Invalid args")
				}
				switch pair[0] {
				case "email":
					user.Email = pair[1]
					user.Nickname = pair[1]
				case "password":
					user.Password = pair[1]
					user.PasswordConfirmation = pair[1]
				default:
					continue
				}
			}

			organization := &models.Organization{}
			count, err := tx.Where("display_id = ?", AdminDisplayID).Count(organization)
			if err != nil {
				return err
			}

			if count < 1 {
				// create company
				company := &models.Company{
					Name: "PlaySecurityAdministrators",
				}
				if _, err := company.Create(tx); err != nil {
					return err
				}

				// create organization
				organization = &models.Organization{
					CompanyID:        company.ID,
					Name:             "Administrator",
					DisplayID:        AdminDisplayID,
					ReportSendEmails: "admin@example.com",
				}

				if _, err := organization.Create(tx); err != nil {
					return err
				}
			} else {
				if err := tx.Where("display_id = ?", AdminDisplayID).First(organization); err != nil {
					return err
				}
			}

			user.OrganizationID = organization.ID

			verrs, err := user.Create(tx)
			if err != nil {
				return err
			}
			if verrs.HasAny() {
				return verrs
			}

			administrator := &models.Administrator{
				UserID: user.ID,
			}

			verrs, err = administrator.Create(tx)
			if err != nil {
				return err
			}
			if verrs.HasAny() {
				return verrs
			}

			return nil
		})
	})

	// buffalo task db:seed
	Desc("seed", "Seeds a database")
	Add("seed", func(c *Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			organization := &models.Organization{}
			count, err := tx.Where("display_id = ?", TestDisplayID).Count(organization)
			if err != nil {
				return err
			}

			// if a test organization already exists, the process is terminated.
			if count > 0 {
				return nil
			}

			// create company
			company := &models.Company{
				Name: "PlaySecurityCompany",
			}
			if _, err := company.Create(tx); err != nil {
				return err
			}

			// create organization
			organization = &models.Organization{
				CompanyID:        company.ID,
				Name:             "PlaySecurityOrganization",
				DisplayID:        TestDisplayID,
				ReportSendEmails: "example@example.com",
			}
			if _, err := organization.Create(tx); err != nil {
				return err
			}

			// create story
			story := &models.Story{
				OrganizationID: organization.ID,
				Title:          "Story-1",
			}
			if _, err := story.Create(tx); err != nil {
				return err
			}

			// create tags
			// create story_taggings
			tags := models.Tags{
				{Name: "Tag-1"},
				{Name: "Tag-2"},
				{Name: "Tag-3"},
				{Name: "Tag-4"},
				{Name: "Tag-5"},
			}
			for _, tag := range tags {
				if _, err := tag.Create(tx); err != nil {
					return err
				}
				tagging := &models.StoryTagging{
					StoryID: story.ID,
					TagID:   tag.ID,
				}
				if _, err := tagging.Create(tx); err != nil {
					return err
				}
			}

			// create quizzes
			quiz := &models.Quiz{
				Question:       "Quiz master - Question1",
				FailureMessage: nulls.NewString("Quiz master - FailureMessage1"),
			}
			if _, err := quiz.Create(tx); err != nil {
				return err
			}

			// create quiz_options
			options := models.QuizOptions{
				{QuizID: quiz.ID, Answer: "Answer-1"},
				{QuizID: quiz.ID, Answer: "Answer-2"},
				{QuizID: quiz.ID, Answer: "Answer-3"},
				{QuizID: quiz.ID, Answer: "Answer-4"},
				{QuizID: quiz.ID, Answer: "Answer-5"},
			}
			var oids []uuid.UUID
			for _, option := range options {
				if _, err := option.Create(tx); err != nil {
					return err
				}
				oids = append(oids, option.ID)
			}

			// create roles
			// create scenarios
			roles := models.Roles{
				{OrganizationID: organization.ID, Name: "Role-1"},
				{OrganizationID: organization.ID, Name: "Role-2"},
				{OrganizationID: organization.ID, Name: "Role-3"},
			}
			for _, role := range roles {
				if _, err := role.Create(tx); err != nil {
					return err
				}
				scenario := &models.Scenario{
					StoryID:       story.ID,
					RoleID:        role.ID,
					Overview:      fmt.Sprintf("%s (%s) Overview.", story.Title, role.Name),
					Description:   fmt.Sprintf("%s (%s) Description.", story.Title, role.Name),
					HighestScore:  10,
					ResultMessage: fmt.Sprintf("%s (%s) ResultMessage.", story.Title, role.Name),
				}
				if _, err := scenario.Create(tx); err != nil {
					return err
				}

				// create scenario_quizzes
				q := &models.ScenarioQuiz{
					ScenarioID: scenario.ID,
					QuizID:     nulls.NewUUID(quiz.ID),
					First:      true,
				}
				if _, err := q.Create(tx); err != nil {
					return err
				}

				// create scenario_quiz_options
				for i, oid := range oids {
					var score int
					var status models.Status
					switch i {
					case 0:
						score = 0
						status = models.StatusFailure
					default:
						score = 10
						status = models.StatusSuccess
					}
					o := &models.ScenarioQuizOption{
						ScenarioQuizID: q.ID,
						QuizOptionID:   nulls.NewUUID(oid),
						Score:          score,
						Status:         status,
					}
					if _, err := o.Create(tx); err != nil {
						return err
					}
				}
			}
			return nil
		})
	})
})
