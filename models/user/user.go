package user

import (
	"log"
	"net/http"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

var b backends.User

type User interface {
	Username() string
	IsAdmin()  bool
}

type user struct {
	username string
	isAdmin bool
}
func (u user) Username() string {
	return u.username
}
func (u user) IsAdmin() bool {
	return u.isAdmin
}

func init() {
	if backend, err := backends.Open(config.Backend); err != nil {
		log.Fatal("Failed to get backend")
	} else {
		b = backend.User
	}
}

func Anonymous() User {
	return &user{"anonymous", false}
}

func Authenticate(username, password string) User {
	m, _ := b.Authenticate(username, password)
	return user{ m.Username, true }
}

func FromCookie(r *http.Request) User {
	c, err := r.Cookie(config.Cookie)
	if err != nil {
		return Anonymous()
	}

	return user{ c.Value, true }
}