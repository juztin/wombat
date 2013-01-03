package sqlite

import (
	"bitbucket.org/juztin/wombat/backends"
)

type userBackend struct {
}

func init() {
	/*
	// connection
	conn := newConn()

	// user
	bUser := userBackend{newBackend(conn, userStmts)}

	// address
	bAddr := addrBackend{newBackend(conn, addrStmts), updateAddrStmt}

	// addressMailing
	b := addrBackend{newBackend(conn, addrMailingStmts), updateAddrMailingStmt}
	bAddrMailing := addrMailingBackend{b}

	// register
	*/
	backends.Register("sqlite", backends.Backend{
		User: userBackend{},
		/*User:           bUser,
		Address:        bAddr,
		AddressMailing: bAddrMailing,*/
	})
}

func (b userBackend) Authenticate(username, password string) (backends.UserData, backends.Error) {
	u := new(backends.UserData)
	u.Username = username
	
	return *u, nil
}
