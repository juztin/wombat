package backends

import (
	"fmt"

	"bitbucket.org/juztin/wombat/config"
)

type User interface {
	ByUsername(username string) (UserData, Error)
	BySession(session string) (UserData, Error)
	SetSession(username, key string) Error
}

type UserData struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Session   string
	Role      int
	Status    int
}

func UserBackend() (User, error) {
	if b, err := Open(config.Backend + ":user"); err != nil {
		return nil, err
	} else if u, ok := b.(User); !ok {
		return nil, fmt.Errorf("backends: invalid user backend type %v", b)
	} else {
		return u, nil
	}
	return nil, nil
}
