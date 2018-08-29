package main

import (
	"github.com/kataras/iris"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct {
	app *App
}

// NewAPI returns a new API for the "app".
func NewAPI(app *App) *API {
	return &API{
		app: app,
	}
}

// NOTE by @kataras
//
// Don't ask me, I found that code from the original repository,
// look my examples at
// github.com/kataras/iris/tree/master/_examples/structuring for folder structuring and organization instead.
//
// END NOTE

// Bind attaches api routes
func (api *API) Bind(party iris.Party) {
	party.Get("/v1/conf", api.ConfHandler)
}

// ConfHandler handle the app config, for example
func (api *API) ConfHandler(ctx iris.Context) {
	ctx.JSON(api.app.Conf.Root)
}
