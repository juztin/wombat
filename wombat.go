package wombat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	//"sort"

	"bitbucket.org/juztin/dingo"
	"bitbucket.org/juztin/dingo/views"

	"bitbucket.org/juztin/wombat/config"
	"bitbucket.org/juztin/wombat/models/user"
)

/*-----------------------------------Fields------------------------------------*/
const (
	VERSION  string = "0.0.1"
	ERR_TMPL string = "/errors/"
)

var (
	Wrap Wrapper = wrap
)

//var httpCodes = []int64{401, 404, 500, 501}

type Server struct {
	*dingo.Server
	Wrapper Wrapper
}

type Context struct {
	dingo.Context
	User user.User
}

type Handler func(ctx Context)

type Wrapper func(fn Handler) func(ctx dingo.Context)

func SetWrapper(fn Wrapper) {
	Wrap = fn
}

func dingoHandler() (net.Listener, error) {
	if !config.UnixSock {
		return dingo.HttpHandler(config.ServerHost, config.ServerPort)
	}
	return dingo.SOCKHandler(config.UnixSockFile, os.ModePerm)
}

func canEdit(ctx dingo.Context) bool {
	//return user.FromCookie(ctx.Request).IsAdmin()
	return getUser(ctx).IsAdmin()
}

func canEditExpire(ctx dingo.Context) bool {
	//return user.FromCookie(ctx.Request).IsAdmin()
	return getExpireUser(ctx).IsAdmin()
}

func addErrViews(path string) {
	files, err := ioutil.ReadDir(filepath.Join(path, ERR_TMPL))
	if err != nil {
		log.Printf("Failed to load error templates: %v\n", err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		bName := name[:len(name)-len(filepath.Ext(name))]
		// Add all error templates where the base filename exists within `httpCodes`
		if x, err := strconv.ParseInt(bName, 10, 16); err == nil {
			/*i := sort.Search(len(httpCodes), func(i int) bool { return httpCodes[i] >= x })
			if i < len(httpCodes) && httpCodes[i] == x {
				// x is present within httpCodes[i]
				views.New(fmt.Sprintf("%s%s", ERR_TMPL, name))
			} else {
				// x is NOT present within httpCodes[i]
				log.Println("Invalid error template found: ", name)
			}*/
			if x > 99 && x < 600 {
				t := fmt.Sprintf("%s%s", ERR_TMPL, name)
				views.New(t)
				if config.TmplEditErr {
					views.AddEditableView(t)
				}
			} else {
				log.Println("Invalid error template found: ", name)
			}
		}
	}
}

func getUser(ctx dingo.Context) user.User {
	return user.FromSession(GetCookieSession(ctx.Request))
}

func getExpireUser(ctx dingo.Context) user.User {
	if cookie, key, ok := UpdatedExpireCookie(ctx.Request); ok {
		http.SetCookie(ctx.Writer, cookie)
		return user.FromSession(key)
	} else if cookie != nil {
		http.SetCookie(ctx.Writer, cookie)
	}
	return user.Anonymous()
}

func wrap(fn Handler) func(ctx dingo.Context) {
	return func(ctx dingo.Context) {
		c := new(Context)
		c.Context = ctx
		//c.User = user.FromSession(GetCookieSession(ctx.Request))
		c.User = getUser(ctx)

		fn(*c)
	}
}

func wrapExpires(fn Handler) func(ctx dingo.Context) {
	return func(ctx dingo.Context) {
		c := new(Context)
		c.Context = ctx
		c.User = getExpireUser(ctx)
		/*if cookie, key, ok := UpdatedExpireCookie(ctx.Request); ok {
			c.User = user.FromSession(key)
			http.SetCookie(ctx.Writer, cookie)
		} else {
			c.User = user.Anonymous()
		}*/

		fn(*c)
	}
}

func Error(ctx dingo.Context, status int) bool {
	n := fmt.Sprintf("%s%d.html", ERR_TMPL, status)
	if v := views.Get(n); v != nil {
		v.Execute(ctx, nil)
		return true
	}
	return false
}

func Signin(ctx *Context, username, password string) {
	if u, ok := user.Authenticate(username, password); ok {
		// set Context user
		ctx.User = u
		// new session-key/cookie
		k, c := NewExpireCookie()
		// save the session-key
		user.SetSession(u, k)
		// add the cookie to the response
		http.SetCookie(ctx.Writer, c)
	}
}

func Signout(ctx *Context) {
	if !ctx.User.IsAnonymous() {
		user.SetSession(ctx.User, "")
		ctx.User = user.Anonymous()
	}
	http.SetCookie(ctx.Writer, expiredCookie())
}

/*-----------------------------------Routes-----------------------------------*/
type newRoute func(path string, h Handler) dingo.Route

func iroute(route newRoute) dingo.NewIRoute {
	var fn dingo.NewIRoute
	fn = func(path string, h interface{}) dingo.Route {
		switch h.(type) {
		case Handler:
			return route(path, h.(Handler))
		case func(Context):
			return route(path, Handler(h.(func(Context))))
		}

		log.Fatal("Handler is invalid: ", h)
		return nil
	}
	return fn
}

func NewSRoute(path string, handler Handler) dingo.Route {
	return dingo.NewSRoute(path, Wrap(handler))
}
func NewReRoute(re string, handler Handler) dingo.Route {
	return dingo.NewReRoute(re, Wrap(handler))
}

/*-----------------------------------Server-----------------------------------*/
func (s *Server) SRoute(path string, handler Handler, methods ...string) {
	s.Server.SRoute(path, Wrap(handler), methods...)
}
func (s Server) ReRoute(path string, handler Handler, methods ...string) {
	s.Server.ReRoute(path, Wrap(handler), methods...)
}
func (s *Server) RRoute(path string, handler interface{}, methods ...string) {
	rt := NewRRoute(path, handler)
	s.Route(rt, methods...)
}

func (s *Server) SRouter(p string) dingo.Router {
	return dingo.NewRouter(s.Server, p, iroute(NewSRoute))
}
func (s *Server) ReRouter(p string) dingo.Router {
	return dingo.NewRouter(s.Server, p, iroute(NewReRoute))
}
func (s *Server) RRouter(p string) dingo.Router {
	return dingo.NewRouter(s.Server, p, NewRRoute)
}

func (s *Server) Serve() {
	if config.UnixSock {
		log.Println("Wombat - listening on", config.UnixSockFile)
	} else {
		log.Printf("Wombat - listening on %s:%d\n", config.ServerHost, config.ServerPort)
	}

	s.Server.Serve()
}

/*----------------------------------------------------------------------------*/

func New() Server {
	// load config
	//config.Load()
	if config.CookieExpires {
		Wrap = wrapExpires
		views.CanEdit = canEditExpire
	} else {
		// secure editable views
		views.CanEdit = canEdit
	}

	// update empty template
	//views.EmptyTmpl = template.Empty

	// editable view media
	views.CodeMirrorJS = "//" + config.MediaURL + "js/vendor/codemirror.js"
	views.CodeMirrorCSS = "//" + config.MediaURL + "css/codemirror.css"

	// add all XXX.html, http status-code, templates
	addErrViews(config.TmplPath)

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
	return Server{&s, Wrap}
}
