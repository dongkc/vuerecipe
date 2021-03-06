package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/packr/v2"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/vuerecipe/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_vue_session",
			SessionStore: CustomCookieStore([]byte("fzdiotplatform")),
		})

		// Automatically redirect to SSL
		app.Use(forcessl.Middleware(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.New("../locales", "../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.ServeFiles("/assets", assetsBox)

		api := app.Group("/api")
		band := api.Resource("/bands", BandsResource{&buffalo.BaseResource{}})
		band.Resource("/members", MembersResource{&buffalo.BaseResource{}})

		app.GET("/", HomeHandler)
		//AuthMiddlewares
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		// Routes for Auth
		auth := app.Group("/auth")
		auth.GET("/", AuthLanding)
		auth.GET("/new", AuthNew)
		auth.POST("/new", AuthCreate)
		auth.DELETE("/", AuthDestroy)
		auth.Middleware.Skip(Authorize, AuthLanding, AuthNew, AuthCreate)

		// Routes for User registration
		users := app.Group("/users")
		users.GET("/new", UsersNew)
		users.POST("/", UsersCreate)
		users.Middleware.Remove(Authorize)

		// 匹配其他路径
		app.GET("/{path:.+}", HomeHandler)
	}

	return app
}
