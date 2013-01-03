package user

import (
	"net/http"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

var b backends.User

type User struct {
	Username string
	IsAdmin  bool
}

func init() {
	if backend, err := backends.Open(config.Backend); err != nil {
		panic("Failed to get backend")
	} else {
		b = backend.User
	}
}

func Anonymous() User {
	u := new(User)
	u.Username = "anonymous"
	return *u
}

func Authenticate(username, password string) User {
	m, _ := b.Authenticate(username, password)
	return User{ m.Username, true }
}

func FromCookie(r *http.Request) User {
	c, err := r.Cookie(config.Cookie)
	if err != nil {
		return Anonymous()
	}

	return User{ c.Value, true }
}