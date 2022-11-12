package actions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const RedirectURLAfterAuth = "/mypage"

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {
	if uid := c.Session().Get("current_user_id"); uid != nil {
		return c.Redirect(http.StatusFound, RedirectURLAfterAuth)
	}
	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("auth/new.plush.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	o := &models.Organization{}
	u := &models.User{}

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		c.Flash().Add("danger", "ログインできませんでした")

		return c.Render(http.StatusUnauthorized, r.HTML("auth/new.plush.html"))
	}

	if err := c.Bind(u); err != nil {
		return bad()
	}

	tx := c.Value("tx").(*pop.Connection)

	// find an organization by display_id
	if err := tx.Where("display_id = ?", u.DisplayID).First(o); err != nil {
		// return errors.WithStack(err)
		return bad()
	}

	// find a user with the email and organization_id
	if err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).Where("organization_id = ?", o.ID).First(u); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find a user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password)); err != nil {
		return bad()
	}

	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "PlaySecurity にログインしました")

	redirectURL := RedirectURLAfterAuth
	if redir, ok := c.Session().Get("redirectURL").(string); ok && redir != "" {
		redirectURL = redir
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "PlaySecurity からログアウトしました")
	return c.Redirect(http.StatusFound, "/")
}
