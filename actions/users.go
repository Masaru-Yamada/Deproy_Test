package actions

import (
	"fmt"
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	o := &models.Organization{}
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find an organization by display_id
	if err := tx.Eager("Company").Where("display_id = ?", user.DisplayID).First(o); err != nil {
		return errors.WithStack(err)
	}
	user.OrganizationID = o.ID

	verrs, err := user.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		users := &models.Users{}
		if err := tx.Where("organization_id = ?", o.ID).All(users); err != nil {
			return errors.WithStack(err)
		}

		c.Set("organization", o)
		c.Set("user", user)
		c.Set("users", users)

		// TODO: error を表示する
		c.Set("errors", verrs)

		return c.Render(http.StatusOK, r.HTML("admin/organizations/users/index.plush.html"))
	}

	c.Flash().Add("success", "User を作成しました")
	return c.Redirect(http.StatusFound, fmt.Sprintf("/admin/organizations/%v/users", o.ID))
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			o := &models.Organization{}
			tx := c.Value("tx").(*pop.Connection)

			if err := tx.Find(u, uid); err != nil {
				// session inconsistency
				c.Session().Clear()
				return c.Redirect(http.StatusFound, "/")
			}

			if err := tx.Find(o, u.OrganizationID); err != nil {
				// session inconsistency
				c.Session().Clear()
				return c.Redirect(http.StatusFound, "/")
			}

			u.DisplayID = o.DisplayID
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			if err := c.Session().Save(); err != nil {
				c.Session().Clear()
				return c.Redirect(http.StatusFound, "/")
			}

			c.Flash().Add("danger", "PlaySecurity にログインしてください")
			return c.Redirect(http.StatusFound, "/signin")
		}
		return next(c)
	}
}

func AdminAuthorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		auth := true
		uid := c.Session().Get("current_user_id")
		if uid == nil {
			auth = false
		} else {
			admin := &models.Administrator{}
			tx := c.Value("tx").(*pop.Connection)
			count, err := tx.Where("user_id = ?", uid).Count(admin)
			if count < 1 || err != nil {
				auth = false
			}
		}

		if !auth {
			c.Session().Set("redirectURL", c.Request().URL.String())

			if err := c.Session().Save(); err != nil {
				c.Session().Clear()
				return c.Redirect(http.StatusFound, "/")
			}

			// TODO: admin が存在することを示唆するメッセージなのであとで消す (debug 用)
			c.Flash().Add("danger", "admin へのアクセスができません")
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}

func UsersList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	organizationID, err := uuid.FromString(c.Param("organization_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	o := &models.Organization{}
	if err := tx.Eager("Company").Find(o, organizationID); err != nil {
		return errors.WithStack(err)
	}
	user := &models.User{}
	user.DisplayID = o.DisplayID

	users := &models.Users{}
	if err := tx.Where("organization_id = ?", organizationID).All(users); err != nil {
		return errors.WithStack(err)
	}

	c.Set("organization", o)
	c.Set("user", user)
	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("admin/organizations/users/index.plush.html"))
}
