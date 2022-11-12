package actions

import (
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func TagsCreate(c buffalo.Context) error {
	tag := &models.Tag{}
	if err := c.Bind(tag); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tag.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("tag", tag)
		c.Flash().Add("error", verrs.String())
		return c.Render(http.StatusOK, r.HTML("admin/tags/index.plush.html"))
	}

	return c.Redirect(http.StatusFound, "/admin/tags/")
}

func TagsList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	tags := &models.Tags{}
	tag := models.Tag{}

	if err := tx.All(tags); err != nil {
		return errors.WithStack(err)
	}

	c.Set("tag", tag)
	c.Set("tags", tags)
	return c.Render(http.StatusOK, r.HTML("admin/tags/index.plush.html"))
}
