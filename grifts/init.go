package grifts

import (
	"github.com/aiit2022-pbl-okuhara/play-security/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
