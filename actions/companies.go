package actions

import (
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func CompaniesNew(c buffalo.Context) error {
	company := models.Company{}
	c.Set("company", company)
	return c.Render(http.StatusOK, r.HTML("admin/companies/new.plush.html"))
}

func CompaniesCreate(c buffalo.Context) error {
	company := &models.Company{}
	if err := c.Bind(company); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := company.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("company", company)
		c.Flash().Add("error", verrs.String())
		return c.Render(http.StatusOK, r.HTML("admin/companies/new.plush.html"))
	}

	return c.Redirect(http.StatusFound, "/admin")
}
