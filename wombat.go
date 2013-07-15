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
	"bitbucket.org/juztin/wombat/users"
)

/*-----------------------------------Fields------------------------------------*/
const (
	VERSION  string = "0.2.2"
	ERR_TMPL string = "/errors/"
)

var (
	Wrap  Wrapper = wrap
	Users users.Users
)

//var httpCodes = []int64{401, 404, 500, 501}

type Server struct {
	*dingo.Server
	Wrapper Wrapper
}

type Context struct {
	dingo.Context
	User users.User
}

type Handler func(ctx Context)

type Wrapper func(fn Handler) func(ctx dingo.Context)

func SetWrapper(fn Wrapper) {
	Wrap = fn
}

func dingoHandler() (net.Listener, error) {
	if config.UnixSock {
		return dingo.SOCKListener(config.UnixSockFile, os.ModePerm)
	} else if config.TLS {
		return dingo.TLSListener(config.ServerHost, config.ServerPort, config.TLSCert, config.TLSKey)
	}
	return dingo.HttpListener(config.ServerHost, config.ServerPort)
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

func getUser(ctx dingo.Context) users.User {
	u, err := Users.BySession(GetCookieSession(ctx.Request))
	if err != nil {
		return users.NewAnonymous()
	}
	return u
}

func getExpireUser(ctx dingo.Context) users.User {
	if cookie, key, ok := UpdatedExpireCookie(ctx.Request); ok {
		http.SetCookie(ctx.Response, cookie)
		if u, err := Users.BySession(key); err == nil {
			return u
		}
	} else if cookie != nil {
		http.SetCookie(ctx.Response, cookie)
	}
	return users.NewAnonymous()
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
			http.SetCookie(ctx.Response, cookie)
		} else {
			c.User = user.Anonymous()
		}*/

		fn(*c)
	}
}

func Error(ctx dingo.Context, status int) bool {
	// default to 500 if in invalid status is given
	if status < 100 || status > 505 {
		status = 500
		// TODO Log an error
	}
	r := ctx.Response
	r.WriteHeader(status)

	// write matching error template
	n := fmt.Sprintf("%s%d.html", ERR_TMPL, status)
	if v := views.Get(n); v != nil {
		v.Execute(ctx, nil)
	} else {
		m := []byte(http.StatusText(status))
		r.Write(m)
	}
	return true
}

func Signin(ctx *Context, username, password string) (err error) {
	var u users.User
	if u, err = Users.Signin(username, password); err != nil {
		log.Println("User signin failed:", err)
	} else {
		// set Context user
		ctx.User = u
		// new session-key/cookie
		k, c := NewExpireCookie()
		// save the session-key
		u.SetSession(k)
		// add the cookie to the response
		http.SetCookie(ctx.Response, c)
	}
	return
}

func Signout(ctx *Context) {
	if !ctx.User.IsAnonymous() {
		ctx.User.SetSession("")
		ctx.User = users.NewAnonymous()
	}
	http.SetCookie(ctx.Response, expiredCookie())
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
/*func (s *Server) Route(rt wombat.Route, methods ...string) {
}*/
func (s *Server) SRoute(path string, handler Handler, methods ...string) {
	s.Server.SRoute(path, Wrap(handler), methods...)
}
func (s Server) ReRoute(path string, handler Handler, methods ...string) {
	s.Server.ReRoute(path, Wrap(handler), methods...)
}
func (s *Server) RRoute(path string, handler interface{}, methods ...string) {
	rt := NewRRoute(path, handler)
	s.Server.Route(rt, methods...)
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
	} else if config.TLS {
		log.Printf("Wombat - listening on TLS %s:%d\n", config.ServerHost, config.ServerPort)
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

	// set `users` factory
	Users = users.New()

	// return the server so that routes can be added
	return Server{&s, Wrap}
}
