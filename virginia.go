package virginia

import (
	"fmt"
	"log"
	"net"
	"os"

	"bitbucket.org/juztin/dingo"
	"bitbucket.org/juztin/dingo/views"

	"bitbucket.org/juztin/virginia/config"
	"bitbucket.org/juztin/virginia/models/user"
	"bitbucket.org/juztin/virginia/template"
)

/*-----------------------------------Fields------------------------------------*/
const (
	VERSION string = "0.0.1"
)

type Context struct {
	dingo.Context
	User user.User
}

type Handler func(ctx Context)

func New() dingo.Server {
	// update empty template
	views.EmptyTmpl = template.Empty

	// secure editable views
	views.CanEdit = canEdit

	// editable view media
	views.CodeMirrorJS = "//" + config.MediaURL + "js/vendor/codemirror.js"
	views.CodeMirrorCSS = "//" + config.MediaURL + "css/codemirror.css"

	// create the handler
	h, e := dingoHandler()
	if e != nil {
		log.Fatalf("Failed to create socket: %v", e)
	}
	s := dingo.New(h)

	// add the editable route
	s.ReRoute(config.TmplEditRt, views.EditHandler, "GET", "POST")

	// set a default error handler
	dingo.ErrorHandler = Error

	// return the server so that routes can be added
	return s
}

func dingoHandler() (net.Listener, error) {
	if !config.IsProd {
		return dingo.HttpHandler(config.DebugURL, config.DebugPort)
	}
	return dingo.SOCKHandler(config.SockFile, os.ModePerm)
}

func canEdit(ctx dingo.Context) bool {
	return user.FromCookie(ctx.Request).IsAdmin
}

func V(fn Handler) func(ctx dingo.Context) {
	return func(ctx dingo.Context) {
		c := new(Context)
		c.Context = ctx
		c.User = user.FromCookie(ctx.Request)

		fn(*c)
	}
}

func Error(ctx dingo.Context, status int) bool {
	n := fmt.Sprintf("%d.html", status)
	if v := views.Get(n); v != nil {
		v.Execute(ctx, nil)
		return true
	}
	return false
}

/* virginia context, handler funcs */
