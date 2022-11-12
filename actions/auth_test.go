package actions

import (
	"net/http"
)

func (as *ActionSuite) Test_Auth_New() {
	res := as.HTML("/signin").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "ログイン")
}

func (as *ActionSuite) Test_Auth_Create() {
	c, err := as.createCompany()
	as.NoError(err)
	o, err := as.createOrganization(c)
	as.NoError(err)
	u, err := as.createUser(o)
	as.NoError(err)

	res := as.HTML("/signin").Post(u)
	as.Equal(http.StatusFound, res.Code)
	as.Equal("/", res.Location())
}

func (as *ActionSuite) Test_Auth_Create_UnknownUser() {
	c, err := as.createCompany()
	as.NoError(err)
	o, err := as.createOrganization(c)
	as.NoError(err)
	u, err := as.createUser(o)
	as.NoError(err)

	u.Email = "unknown@esample.com"
	u.PasswordConfirmation = ""
	res := as.HTML("/signin").Post(u)
	as.Equal(http.StatusUnauthorized, res.Code)
	as.Contains(res.Body.String(), "ログインできませんでした")
}

func (as *ActionSuite) Test_Auth_Create_InvalidPassword() {
	c, err := as.createCompany()
	as.NoError(err)
	o, err := as.createOrganization(c)
	as.NoError(err)
	u, err := as.createUser(o)
	as.NoError(err)

	u.Password = "invalid"
	res := as.HTML("/signin").Post(u)
	as.Equal(http.StatusUnauthorized, res.Code)
	as.Contains(res.Body.String(), "ログインできませんでした")
}
