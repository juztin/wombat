package user

import (
	"log"

	"code.google.com/p/go.crypto/bcrypt"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

var b backends.User

type User interface {
	Username() string
	Firstname() string
	Lastname() string
	Email() string
	IsAdmin() bool
	IsAnonymous() bool
	Role() int
	Status() int
}

type user struct {
	username    string
	firstname   string
	lastname    string
	email       string
	password    string
	isAdmin     bool
	isAnonymous bool
	role        int
	status      int
}

func init() {
	if backend, err := backends.Open(config.Backend); err != nil {
		log.Fatal("Failed to get backend")
	} else {
		b = backend.User
	}
}

/*---------------- user ----------------*/
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
func (u user) IsAnonymous() bool {
	return u.isAnonymous
}
func (u user) Role() int {
	return u.role
}
func (u user) Status() int {
	return u.status
}

/*--------------------------------------*/

func hashit(v string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(v), config.CryptIter)
}

func userFromData(m backends.UserData) user {
	u := new(user)
	u.username = m.Username
	u.firstname = m.Firstname
	u.lastname = m.Lastname
	u.email = m.Email
	u.password = m.Password
	u.isAdmin = (m.Role & 1) == 1
	u.isAnonymous = false
	u.role = m.Role
	u.status = m.Status
	return *u
}

func Anonymous() User {
	u := new(user)
	u.username = "anonymous"
	u.isAdmin = false
	u.isAnonymous = true
	return u
}

func Authenticate(username, password string) (User, bool) {
	m, err := b.GetByUsername(username)
	if err != nil {
		log.Printf("Failed to get user: ", username, " : ", err)
		return Anonymous(), false
	}

	p1, p2 := []byte(m.Password), []byte(password)
	if err := bcrypt.CompareHashAndPassword(p1, p2); err != nil {
		log.Printf("Failed to compare password hashes: ", err)
		return Anonymous(), false
	}

	return userFromData(m), true
}

func FromSession(key string) User {
	if key == "" {
		return Anonymous()
	}

	m, err := b.GetBySession(key)
	if err != nil {
		log.Printf("Failed to get user from cache:%v\n", err)
		return Anonymous()
	}

	return userFromData(m)
}

func SetSession(u User, s string) bool {
	if err := b.SetSession(u.Username(), s); err != nil {
		log.Printf("Failed to set session for: ", u.Username(), " : ", err)
		return false
	}
	return true
}
