package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"github.com/nu7hatch/gouuid"
	"github.com/olebedev/config"
)

// App struct.
// There is no singleton anti-pattern,
// all variables defined locally inside
// this struct.
type App struct {
	Server *iris.Application
	Conf   *config.Config
	React  *React
	API    *API
}

// NewApp returns initialized struct
// of main server application.
func NewApp(opts ...AppOptions) *App {
	options := AppOptions{}
	for _, i := range opts {
		options = i
		break
	}

	options.init()

	// Parse config yaml string from ./conf.go
	conf, err := config.ParseYaml(confString)
	Must(err)

	// Set config variables delivered from main.go:11
	// Variables defined as ./conf.go:3
	conf.Set("debug", debug)
	conf.Set("commitHash", commitHash)

	// Parse environ variables for defined
	// in config constants
	conf.Env()

	// Make an engine
	srv := iris.New()

	// Use precompiled embedded templates
	srv.RegisterView(iris.HTML("./data/templates", ".html").Binary(Asset, AssetNames))

	// Set up debug level for iris logger
	if conf.UBool("debug") {
		srv.Logger().SetLevel("debug")
	}

	// Regular middlewares
	srv.UseGlobal(recover.New())

	// Map app and uuid for every requests
	srv.UseGlobal(func(ctx iris.Context) {
		id, err := uuid.NewV4()
		if err == nil {
			ctx.Values().Set("uuid", id)
		}
		// ctx.Values().Set("app", app)
		ctx.Next()
	})

	// Favicon
	srv.Favicon("./data/static/images/favicon.ico")

	// Initialize the application
	app := &App{
		Conf:   conf,
		Server: srv,
		React: NewReact(
			conf.UString("duktape.path"),
			conf.UBool("debug"),
			srv,
		),
	}

	// Serve static via bindata
	srv.StaticEmbedded("/static", "./data/static", Asset, AssetNames)

	// Request Logger with columns
	srv.Use(logger.New(logger.Config{
		Columns: true,
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
	}))

	api := NewAPI(app)

	// Bind api hadling for URL api.prefix
	api.Bind(
		app.Server.Party(
			app.Conf.UString("api.prefix"),
		),
	)

	// Handle via react app
	//
	// Registers / and /anything/here on GET and HEAD http methods.
	// srv.HandleMany("GET HEAD", "/ /{p:path}", app.React.Handle)
	// Or:
	//
	// handle root with react
	srv.Get("/", app.React.Handle)
	// handle anything expect /static/ and /api/v1/conf with react as well
	srv.Get("/{p:path}", app.React.Handle)

	return app
}

// Run runs the app
func (app *App) Run() {
	Must(app.Server.Run(
		iris.Addr(":"+app.Conf.UString("port")),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed)))
}

// AppOptions is options struct
type AppOptions struct{}

func (ao *AppOptions) init() { /* write your own*/ }
