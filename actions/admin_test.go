package actions

import "net/http"

func (as *ActionSuite) Test_AdminHandler() {
	res := as.HTML("/admin/").Get()
	as.Equal(http.StatusOK, res.Code)
}
