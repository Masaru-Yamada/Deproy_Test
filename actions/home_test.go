package actions

import (
	"net/http"
)

func (as *ActionSuite) Test_HomeHandler() {
	res := as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code)
}

func (as *ActionSuite) Test_HomeHandler_LoggedIn() {
	c, err := as.createCompany()
	as.NoError(err)
	o, err := as.createOrganization(c)
	as.NoError(err)
	u, err := as.createUser(o)
	as.NoError(err)

	as.Session.Set("current_user_id", u.ID)

	as.Session.Clear()
	res := as.HTML("/signin").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "ログイン")
}
