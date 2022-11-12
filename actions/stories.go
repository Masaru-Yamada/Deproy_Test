package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func StoriesNew(c buffalo.Context) error {

	return c.Render(http.StatusOK, r.HTML("admin/organizations/stories/new.plush.html"))
}

func StoriesCreate(c buffalo.Context) error {

	return nil
}

func StoriesList(c buffalo.Context) error {

	return c.Render(http.StatusOK, r.HTML("admin/organizations/stories/index.plush.html"))
}
