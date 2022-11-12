package actions

import (
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func QuizzesNew(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	quiz := &models.Quiz{}
	quiz.QuizOptions = make([]models.QuizOption, 1)

	quizzes := &models.Quizzes{}

	if err := tx.Eager("QuizOptions").All(quizzes); err != nil {
		return errors.WithStack(err)
	}

	c.Set("quiz", quiz)
	c.Set("quizzes", quizzes)
	return c.Render(http.StatusOK, r.HTML("admin/quizzes/new.plush.html"))
}

func QuizzesCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	quiz := &models.Quiz{}

	if err := c.Bind(quiz); err != nil {
		return errors.WithStack(err)
	}

	quizzes := &models.Quizzes{}
	if err := tx.Eager("QuizOptions").All(quizzes); err != nil {
		return errors.WithStack(err)
	}

	if quiz.Question == "" && quiz.FailureMessage.String == "" {
		return c.Redirect(http.StatusFound, "/admin/quizzes/new/")
	}

	tran, err := tx.NewTransaction()
	defer tx.TX.Rollback()

	if err != nil {
		return errors.WithStack(err)
	}

	verrs, err := quiz.Create(tran)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("quiz", quiz)
		c.Set("quizzes", quizzes)
		c.Flash().Add("error", verrs.String())
		return c.Render(http.StatusOK, r.HTML("admin/quizzes/new.plush.html"))
	}

	for _, v := range quiz.QuizOptions {
		if v.Answer == "" {
			continue
		}

		v.QuizID = quiz.ID
		verrs, err := v.Create(tx)

		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			c.Set("quiz", quiz)
			c.Set("quizzes", quizzes)
			c.Flash().Add("error", verrs.String())
			return c.Render(http.StatusOK, r.HTML("admin/quizzes/new.plush.html"))
		}
	}

	tran.TX.Commit()
	return c.Redirect(http.StatusFound, "/admin/quizzes/new/")
}
