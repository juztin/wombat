package user

import (
	"log"
	"net/http"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

var b backends.User

type User interface {
	Username()  string
	Firstname() string
	Lastname()  string
	Email()     string
	IsAdmin()   bool
	Role()      int
	Status()    int
}

type user struct {
	username  string
	firstname string
	lastname  string
	email     string
	password  string
	isAdmin   bool
	role      int
	status    int
}
func (u user) Username() string {
	return u.username
}
func (u user) Firstname() string {
	return u.firstname
}
func (u user) Lastname() string {
	return u.lastname
}
func (u user) Email() string {
	return u.email
}
func (u user) IsAdmin() bool {
	return u.isAdmin
}
func (u user) Role() int {
	return u.role
}
func (u user) Status() int {
	return u.status
}

func init() {
	if backend, err := backends.Open(config.Backend); err != nil {
		log.Fatal("Failed to get backend")
	} else {
		b = backend.User
	}
}

func userFromData(m backends.UserData) user {
	u := new(user)
	u.username = m.Username
	u.firstname = m.Firstname
	u.lastname = m.Lastname
	u.email = m.Email
	u.password = m.Password
	u.isAdmin = m.Role == 666
	u.role = m.Role
	u.status = m.Status
	return *u
}

func Anonymous() User {
	u := new(user)
	u.username = "anonymous"
	u.isAdmin = false
	return u
}

func Authenticate(username, password string) (User, bool) {
	if m, err := b.Authenticate(username, password); err == nil {
		return userFromData(m), true
	}
	return nil, false
}

func FromCookie(r *http.Request) User {
	c, err := r.Cookie(config.Cookie)
	if err != nil {
		return Anonymous()
	}

	m, err := b.FromCache(c.Value)
	if err != nil {
		log.Println("Failed to get user from cache: ", err)
		return nil
	}

	return userFromData(m)
}
