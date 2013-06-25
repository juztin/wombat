package users

import (
	"errors"
	"log"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

type Data struct {
	Writer     `-`
	Username   string    `username`
	Firstname  string    `firstname`
	Lastname   string    `lastname`
	Email      string    `email`
	Password   string    `password`
	Session    string    `session`
	Locale     string    `locale`
	Role       int       `role`
	Status     int       `status`
	LastSignin time.Time `lastSignin`
	CreatedOn  time.Time `createdOn`
}

type Model struct {
	Writer
	Data *Data
}

type Reader interface {
	ByUsername(username string) (User, error)
	BySession(key string) (User, error)
}

type Writer interface {
	UpdateSession(username, key string) error
	UpdateLastSignin(username string, on time.Time) error
}

type Users interface {
	Reader
	Signin(username, password string) (User, error)
}

type User interface {
	Writer
	CheckPassword(p string) bool
	Username() string
	Firstname() string
	Lastname() string
	Email() string
	Password() string
	Session() string
	Locale() string
	IsAdmin() bool
	IsAnonymous() bool
	IsActive() bool
	LastSignin() time.Time
	CreatedOn() time.Time

	SetSession(key string) error
}

type users struct {
	Reader
}

type usersFn func(r Reader) Users

var NewUsersFn usersFn = newUsers

func New() Users {
	var r Reader
	if p, err := backends.Open(config.UserReader); err != nil {
		log.Fatal("No 'user' reader available")
	} else {
		if o, ok := p.(Reader); !ok {
			log.Fatal("Invalid 'user' reader")
		} else {
			r = o
		}
	}
	return NewUsersFn(r)
}

func newUsers(r Reader) Users {
	return &users{r}
}

func NewAnonymous() User {
	o := new(Data)
	o.Username = "anonymous"
	o.Role = 0
	o.Status = 0
	return &Model{nil, o}
}

// Model
func (m Model) CheckPassword(p string) bool {
	p1, p2 := []byte(m.Password()), []byte(p)
	if err := bcrypt.CompareHashAndPassword(p1, p2); err != nil {
		log.Printf("Failed to compare password hashes: ", err)
		return false
	} else {
		if err = m.UpdateLastSignin(m.Username(), time.Now()); err != nil {
			log.Println("Failed to update 'lastSignin' time ", err)
		}
	}
	return true
}
func (m Model) Username() string {
	return m.Data.Username
}
func (m Model) Firstname() string {
	return m.Data.Firstname
}
func (m Model) Lastname() string {
	return m.Data.Lastname
}
func (m Model) Email() string {
	return m.Data.Email
}
func (m Model) Password() string {
	return m.Data.Password
}
func (m Model) Session() string {
	return m.Data.Session
}
func (m Model) Locale() string {
	return m.Data.Locale
}
func (m Model) IsAdmin() bool {
	// admin flag is 1 AND status is active
	return (m.Data.Role&1) == 1 && m.Data.Status == 1
}
func (m Model) IsAnonymous() bool {
	return m.Data.Role == 0
}
func (m Model) IsActive() bool {
	return m.Data.Status == 1
}
func (m Model) LastSignin() time.Time {
	return m.Data.LastSignin
}
func (m Model) CreatedOn() time.Time {
	return m.Data.CreatedOn
}
func (m Model) SetSession(key string) error {
	return m.Writer.UpdateSession(m.Data.Username, key)
}

// Users
func (u *users) Signin(username, password string) (User, error) {
	o, err := u.Reader.ByUsername(username)
	if err != nil {
		return NewAnonymous(), err
	}

	user, ok := o.(User)
	if !ok {
		return NewAnonymous(), errors.New("Invalid 'users' object")
	} else if ok = user.CheckPassword(password); !ok {
		return NewAnonymous(), errors.New("Authentication")
	}

	return user, nil
}
