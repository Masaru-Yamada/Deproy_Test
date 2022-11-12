package actions

import (
	"net/http"

	"github.com/aiit2022-pbl-okuhara/play-security/locales"
	"github.com/aiit2022-pbl-okuhara/play-security/models"
	"github.com/aiit2022-pbl-okuhara/play-security/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n/v2"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app *buffalo.App
	T   *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_play_security_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)

		// AuthMiddlewares
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		app.GET("/signin", AuthNew)
		app.POST("/signin", AuthCreate)
		app.GET("/signout", AuthDestroy)
		app.Middleware.Skip(Authorize, HomeHandler, UsersCreate, AuthNew, AuthCreate, AdminHandler)

		app.GET("/mypage", MypageHandler)
		app.GET("/scenarios", ScenariosList)
		app.GET("/scenarios/{scenario_id}", ScenariosShow).Name("scenariosShow")
		app.GET("/scenarios/{scenario_id}/quizzes/{scenario_quiz_id}", ScenariosQuizzesShow)
		app.POST("/scenarios/{scenario_id}/quizzes/{scenario_quiz_id}", ScenariosQuizzesAnswer).Name("answer")
		app.GET("/scenarios/{scenario_id}/result", ScenariosResult)

		g := app.Group("/admin")
		g.Use(AdminAuthorize)
		g.GET("/", AdminHandler)

		g.POST("/tags", TagsCreate)
		g.GET("/tags", TagsList)

		g.GET("/quizzes/new", QuizzesNew)
		g.POST("/quizzes", QuizzesCreate)

		g.GET("/companies/new", CompaniesNew)
		g.POST("/companies", CompaniesCreate)
		g.GET("/organizations/new", OrganizationsNew)
		g.POST("/organizations", OrganizationsCreate)

		g.POST("/organizations/{organization_id}/users", UsersCreate).Name("adminUsers")
		g.GET("/organizations/{organization_id}/users", UsersList).Name("adminUsers")

		g.GET("/organizations/{organization_id}/roles/new", RolesNew).Name("newAdminRoles")
		g.POST("/organizations/{organization_id}/roles", RolesCreate).Name("adminRoles")
		g.GET("/organizations/{organization_id}/roles", RolesList).Name("adminRoles")

		g.GET("/organizations/{organization_id}/stories/new", StoriesNew).Name("newAdminStories")
		g.POST("/organizations/{organization_id}/stories", StoriesCreate).Name("adminStories")
		g.GET("/organizations/{organization_id}/stories", StoriesList).Name("adminStories")

		// TODO: scenario and scenarios_quizzes and scenarios_quizzes_options

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		if err := app.Stop(err); err != nil {
			return nil
		}
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
