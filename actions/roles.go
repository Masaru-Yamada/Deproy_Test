package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func RolesNew(c buffalo.Context) error {

	return c.Render(http.StatusOK, r.HTML("admin/organizations/tags/new.plush.html"))
}

func RolesCreate(c buffalo.Context) error {

	return nil
}

func RolesList(c buffalo.Context) error {

	return c.Render(http.StatusOK, r.HTML("admin/organizations/roles/index.plush.html"))
}
