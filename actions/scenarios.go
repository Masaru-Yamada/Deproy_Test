package actions

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/tags/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func ScenariosList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := c.Data()["current_user"].(*models.User)
	scoreHistoryMap := map[uuid.UUID]int{}
	countHistoryMap := map[uuid.UUID]int{}

	scenarios := &models.Scenarios{}
	tx.Where("organization_id = ?", u.OrganizationID).Order("story_id")
	if err := tx.Eager("Story").Eager("Role").All(scenarios); err != nil {
		return errors.WithStack(err)
	}

	scenarioHistory := &models.UserScenarioHistories{}
	if err := scenarioHistory.ListByUserID(tx, u.ID); err != nil {
		return errors.WithStack(err)
	}

	for _, scenario := range *scenarios {
		scoreHistoryMap[scenario.ID] = 0
		countHistoryMap[scenario.ID] = 0
	}

	for _, scenario := range *scenarioHistory {
		if scoreHistoryMap[scenario.ScenarioID] == 0 {
			scoreHistoryMap[scenario.ScenarioID] = scenario.TotalScore
		}
		countHistoryMap[scenario.ScenarioID]++
	}

	c.Set("scenarios", scenarios)
	c.Set("scoreHistory", scoreHistoryMap)
	c.Set("countHistory", countHistoryMap)
	return c.Render(http.StatusOK, r.HTML("scenarios/index.plush.html"))
}

// ScenariosShow default implementation.
func ScenariosShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := c.Data()["current_user"].(*models.User)
	scenarioID, err := uuid.FromString(c.Param("scenario_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	s := &models.Scenario{}
	if err := tx.Eager("Story").Eager("Role").Find(s, scenarioID); err != nil {
		return errors.WithStack(err)
	}

	sq := &models.ScenarioQuizzes{}
	if err := tx.Where("scenario_id = ? and first = true", scenarioID).All(sq); err != nil {
		return errors.WithStack(err)
	}

	firstID := ""
	for _, sq := range *sq {
		firstID = sq.ID.String()
	}

	sh := &models.UserScenarioHistories{}

	count, err := tx.Where("scenario_id = ? and user_id = ? and played_at is not null", scenarioID, u.ID).Order("played_at desc").Count(sh)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Where("scenario_id = ? and user_id = ? and played_at is not null", scenarioID, u.ID).Order("played_at desc").All(sh); err != nil {
		return errors.WithStack(err)
	}

	c.Set("scenario", s)
	c.Set("scenarioQuizID", firstID)
	c.Set("history", sh)
	c.Set("historyCount", count)

	return c.Render(http.StatusOK, r.HTML("scenarios/show.plush.html"))
}

type outputScenariosResult struct {
	ScenarioID     uuid.UUID
	ScenarioQuizID uuid.UUID
	Title          string
	Role           string
	Message        string
	HighestScore   int
	TotalScore     int
	Tags           []models.StoryTagging
	ResultAnswers  []outputScenariosResultAnswer
}

type outputScenariosResultAnswer struct {
	Question string
	Answer   string
}

// ScenariosResult default implementation.
func ScenariosResult(c buffalo.Context) error {
	scenarioID, err := uuid.FromString(c.Param("scenario_id"))
	// redirect to the scenario list page if scenario_id is not set in parameter.
	if err != nil {
		return c.Redirect(http.StatusFound, "/scenarios")
	}

	tx := c.Value("tx").(*pop.Connection)

	s := &models.Scenario{}

	// redirect to the scenario list page if scenario_id does not exist in the database.
	if err := tx.Eager("Story.StoryTaggings.Tag").Eager("Role").Find(s, scenarioID); err != nil {
		return c.Redirect(http.StatusFound, "/scenarios")
	}

	// redirect to the first quiz page if user_scenario_history_id is not set in session.
	if c.Session().Get("user_scenario_history_id") == nil {
		return resumeQuiz(c, scenarioID, uuid.Nil, errors.New("user_scenario_history_id not found in session."))
	}

	historyID := c.Session().Get("user_scenario_history_id").(uuid.UUID)

	// session delete
	c.Session().Delete("user_scenario_history_id")

	qh := &models.UserQuizHistories{}

	// get answer histories
	if err := tx.Eager("ScenarioQuiz.Quiz").Eager("ScenarioQuizOption.QuizOption").Where("user_scenario_history_id = ?", historyID).All(qh); err != nil {
		return resumeQuiz(c, scenarioID, historyID, err)
	}

	// create answer histories
	totalScore := 0
	answers := make([]outputScenariosResultAnswer, len(*qh))
	for _, history := range *qh {
		totalScore += history.Score

		question := history.ScenarioQuiz.Question.String
		if question == "" {
			question = history.ScenarioQuiz.Quiz.Question
		}

		answer := history.ScenarioQuizOption.Answer.String
		if answer == "" {
			answer = history.ScenarioQuizOption.QuizOption.Answer
		}

		answers = append(answers, outputScenariosResultAnswer{
			Question: question,
			Answer:   answer,
		})
	}

	// get a message to display
	message := s.ResultMessage
	if c.Session().Get("failure_scenario_quiz_id") != nil {
		fid := c.Session().Get("failure_scenario_quiz_id").(uuid.UUID)
		// session delete
		c.Session().Delete("failure_scenario_quiz_id")

		sq := &models.ScenarioQuiz{}
		if err := tx.Find(sq, fid); err != nil {
			return resumeQuiz(c, scenarioID, historyID, err)
		}
		message = sq.FailureMessage.String
	}

	// update TotalScore and PlayedAt of UserScenarioHistory
	sh := &models.UserScenarioHistory{
		ID:         historyID,
		UserID:     c.Session().Get("current_user_id").(uuid.UUID),
		ScenarioID: scenarioID,
		TotalScore: totalScore,
		PlayedAt:   nulls.NewTime(time.Now()),
	}

	verrs, err := sh.Update(tx)
	if err != nil || verrs.HasAny() {
		return resumeQuiz(c, scenarioID, historyID, err)
	}

	output := &outputScenariosResult{
		ScenarioID:    scenarioID,
		Title:         s.Story.Title,
		Role:          s.Role.Name,
		Message:       message,
		HighestScore:  s.HighestScore,
		TotalScore:    totalScore,
		Tags:          s.Story.StoryTaggings,
		ResultAnswers: answers[1:],
	}

	c.Set("output", output)

	return c.Render(http.StatusOK, r.HTML("scenarios/result.html"))
}

// ScenariosQuizzesShow default implementation.
func ScenariosQuizzesShow(c buffalo.Context) error {
	_, err := setScenariosQuizzesShow(c)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, r.HTML("scenarios/quizzes/show.html"))
}

type inputScenarioQuizOption struct {
	ID uuid.UUID `form:"quiz_option"`
}

// ScenariosQuizzesAnswer default implementation.
func ScenariosQuizzesAnswer(c buffalo.Context) error {
	// redirect to the scenario list page if scenario quiz data cannot get.
	output, err := setScenariosQuizzesShow(c)
	if err != nil {
		return c.Redirect(http.StatusFound, "/scenarios")
	}

	tx := c.Value("tx").(*pop.Connection)

	historyID := uuid.Nil
	// redirect to the first quiz page if the quiz has already been answered.
	if !output.First {
		// redirect to the first quiz page if user_scenario_history_id is not set in session.
		if c.Session().Get("user_scenario_history_id") == nil {
			return resumeQuiz(c, output.ScenarioID, uuid.Nil, errors.New("user_scenario_history_id not found in session."))
		}

		historyID = c.Session().Get("user_scenario_history_id").(uuid.UUID)

		qh := &models.UserQuizHistories{}
		count, err := tx.Where("user_scenario_history_id = ? and scenario_id = ? and scenario_quiz_id = ?", historyID, output.ScenarioID, output.ScenarioQuizID).Count(qh)
		if err != nil {
			log.Println(errors.WithStack(err))
		}

		if count > 0 {
			return resumeQuiz(c, output.ScenarioID, historyID, errors.New("the quiz has already been answered."))
		}
	}

	// helper function to handle bad attempts
	bad := func(err error) error {
		c.Flash().Add("danger", "回答を選択してください")
		log.Println(errors.WithStack(err))
		return c.Render(http.StatusBadRequest, r.HTML("scenarios/quizzes/show.html"))
	}

	input := &inputScenarioQuizOption{}
	if err := c.Bind(input); err != nil {
		return bad(err)
	}

	sqo := &models.ScenarioQuizOption{}

	// err if no data
	if err := tx.Eager("QuizOption").Find(sqo, input.ID); err != nil {
		return bad(err)
	}

	// err if parameter and database value are different
	if sqo.ScenarioQuizID != output.ScenarioQuizID {
		return bad(errors.New("parameter and database value are different."))
	}

	// save UserScenarioHistory when first answered
	uid := c.Session().Get("current_user_id").(uuid.UUID)
	if output.First {
		ush := &models.UserScenarioHistory{
			UserID:     uid,
			ScenarioID: output.ScenarioID,
		}
		verrs, err := ush.Create(tx)
		if err != nil || verrs.HasAny() {
			return bad(err)
		}
		historyID = ush.ID
		// save UserScenarioHistory.ID to session
		c.Session().Set("user_scenario_history_id", ush.ID)
	}

	// save UserQuizHistory
	uqh := &models.UserQuizHistory{
		UserScenarioHistoryID: historyID,
		UserID:                uid,
		ScenarioID:            output.ScenarioID,
		ScenarioQuizID:        output.ScenarioQuizID,
		ScenarioQuizOptionID:  input.ID,
		Score:                 sqo.Score,
	}
	verrs, err := uqh.Create(tx)
	if err != nil || verrs.HasAny() {
		return bad(err)
	}

	switch sqo.Status {
	case models.StatusSuccess:
		return c.Redirect(http.StatusFound, fmt.Sprintf("/scenarios/%v/result/", output.ScenarioID))
	case models.StatusFailure:
		c.Session().Set("failure_scenario_quiz_id", output.ScenarioQuizID)
		return c.Redirect(http.StatusFound, fmt.Sprintf("/scenarios/%v/result/", output.ScenarioID))
	case models.StatusNextQuiz:
		return c.Redirect(http.StatusFound, fmt.Sprintf("/scenarios/%v/quizzes/%v/", output.ScenarioID, sqo.NextScenarioQuizID.UUID.String()))
	default:
		return bad(errors.New("invalid quiz status."))
	}
}

type outputScenariosQuizzes struct {
	ScenarioID     uuid.UUID
	ScenarioQuizID uuid.UUID
	Title          string
	Role           string
	First          bool
	Tags           []models.StoryTagging
	Question       string
	Options        []tags.Options
}

func setScenariosQuizzesShow(c buffalo.Context) (*outputScenariosQuizzes, error) {
	scenarioID, err := uuid.FromString(c.Param("scenario_id"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	scenarioQuizID, err := uuid.FromString(c.Param("scenario_quiz_id"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s := &models.Scenario{}
	sq := &models.ScenarioQuiz{}

	tx := c.Value("tx").(*pop.Connection)

	if err := tx.Eager("Story.StoryTaggings.Tag").Eager("Role").Find(s, scenarioID); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := tx.Eager("Quiz").Eager("ScenarioQuizOptions.QuizOption").Find(sq, scenarioQuizID); err != nil {
		return nil, errors.WithStack(err)
	}

	question := sq.Question.String
	if question == "" {
		question = sq.Quiz.Question
	}

	options := make([]tags.Options, len(sq.ScenarioQuizOptions))

	for i, v := range sq.ScenarioQuizOptions {
		answer := v.Answer.String
		if answer == "" {
			answer = v.QuizOption.Answer
		}

		o := tags.Options{}
		o["id"] = v.ID.String()
		o["label"] = answer
		options[i] = o
	}

	output := &outputScenariosQuizzes{
		ScenarioID:     scenarioID,
		ScenarioQuizID: scenarioQuizID,
		Title:          s.Story.Title,
		Role:           s.Role.Name,
		First:          sq.First,
		Tags:           s.Story.StoryTaggings,
		Question:       question,
		Options:        options,
	}

	c.Set("output", output)

	return output, nil
}

func resumeQuiz(c buffalo.Context, scenarioID, historyID uuid.UUID, err error) error {
	c.Flash().Add("danger", "不正な操作を検知しました。もう一度やりなおしてください。")

	log.Println(errors.WithStack(err))

	sh := &models.UserScenarioHistory{}
	qh := &models.UserQuizHistories{}

	tx := c.Value("tx").(*pop.Connection)

	if historyID != uuid.Nil {
		// delete histories
		if err := qh.DeleteByUserScenarioHistoryID(tx, historyID); err != nil {
			return c.Redirect(http.StatusFound, "/scenarios")
		}
		sh.ID = historyID
		if err := sh.Delete(tx); err != nil {
			return c.Redirect(http.StatusFound, "/scenarios")
		}
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/scenarios/%v", scenarioID))
}
