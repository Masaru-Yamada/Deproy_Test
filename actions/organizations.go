package actions

import (
	"math/rand"
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func OrganizationsNew(c buffalo.Context) error {
	companies := &models.Companies{}
	organization := &models.Organization{}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.All(companies); err != nil {
		// TODO: error log
	}

	companyMap := make(map[string]interface{}, len(*companies))
	var companySelectedValue interface{}
	for i, company := range *companies {
		companyMap[company.Name] = company.ID
		if i == 0 {
			companySelectedValue = company.ID
		}
	}

	c.Set("companies", companyMap)
	c.Set("companyValue", companySelectedValue)
	// DisplayID をランダムに生成する (10000 〜 99990 までの間)
	c.Set("displayIDValue", rand.Intn(models.MaxDisplayID-models.MinDisplayID)+models.MinDisplayID)
	c.Set("organization", organization)
	return c.Render(http.StatusOK, r.HTML("admin/organizations/new.plush.html"))
}

func OrganizationsCreate(c buffalo.Context) error {
	organization := &models.Organization{}
	if err := c.Bind(organization); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := organization.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO: validation
	// DisplayID unique
	// CompanyID exists
	// ReportSendEmails format

	if verrs.HasAny() {
		c.Set("organization", organization)
		c.Flash().Add("error", verrs.String())
		return c.Render(http.StatusOK, r.HTML("admin/organizations/new.plush.html"))
	}

	return c.Redirect(http.StatusFound, "/admin")
}
