package actions

import (
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// AdminHandler is a default handler to serve up
// a admin root page.
func AdminHandler(c buffalo.Context) error {
	// helper function to handle bad attempts
	bad := func() error {
		c.Flash().Add("danger", "Internal Errors")

		// TODO: 適切なエラーページを作る
		return c.Render(http.StatusBadRequest, r.HTML("home/index.plush.html"))
	}

	tx := c.Value("tx").(*pop.Connection)
	o := &models.Organizations{}
	if err := tx.Eager("Company").All(o); err != nil {
		return bad()
	}

	c.Set("output", o)

	return c.Render(http.StatusOK, r.HTML("admin/index.plush.html"))
}
