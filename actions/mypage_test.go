package actions

import "net/http"

func (as *ActionSuite) Test_MypageHandler() {
	res := as.HTML("/mypage").Get()
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), "/signin")
}
