package actions

import (
	"log"
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
)

func (as *ActionSuite) createCompany() (*models.Company, error) {
	c := &models.Company{
		Name: "Company",
	}

	verrs, err := c.Create(as.DB)
	as.False(verrs.HasAny())
	log.Println(c)
	return c, err
}

func (as *ActionSuite) createOrganization(c *models.Company) (*models.Organization, error) {
	o := &models.Organization{
		CompanyID:        c.ID,
		DisplayID:        1,
		Name:             "organization",
		ReportSendEmails: "test@example.com",
	}

	verrs, err := o.Create(as.DB)
	as.False(verrs.HasAny())

	log.Println(o)
	return o, err
}

func (as *ActionSuite) createUser(o *models.Organization) (*models.User, error) {
	u := &models.User{
		OrganizationID:       o.ID,
		Nickname:             "nickname",
		Email:                "test@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
		DisplayID:            1,
	}

	verrs, err := u.Create(as.DB)
	as.False(verrs.HasAny())

	return u, err
}

func (as *ActionSuite) Test_Users_New() {
	res := as.HTML("/users/new").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Users_Create() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)

	c, err := as.createCompany()
	as.NoError(err)
	o, err := as.createOrganization(c)
	as.NoError(err)

	u := &models.User{
		OrganizationID:       o.ID,
		Nickname:             "nickname",
		Email:                "test@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
		DisplayID:            1,
	}

	res := as.HTML("/users").Post(u)
	as.Equal(http.StatusFound, res.Code)

	count, err = as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, count)
}
