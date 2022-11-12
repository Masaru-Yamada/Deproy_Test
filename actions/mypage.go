package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// MypageHandler default implementation.
func MypageHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("mypage/index.html"))
}
