package actions

import "net/http"

func (as *ActionSuite) Test_ScenariosList() {
	res := as.HTML("/scenarios").Get()
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), "/signin")
}

func (as *ActionSuite) Test_ScenariosShow() {
	res := as.HTML("/scenarios/1").Get()
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), "/signin")
}

func (as *ActionSuite) Test_ScenariosResult() {
	res := as.HTML("/scenarios/1/result").Get()
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), "/signin")
}

func (as *ActionSuite) Test_ScenariosQuizzesShow() {
	res := as.HTML("/scenarios/1/quizzes/1").Get()
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), "/signin")
}

func (as *ActionSuite) Test_ScenariosQuizzesAnswer() {

}
