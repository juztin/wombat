package data

import (
	"bitbucket.org/juztin/wombat"
	"bitbucket.org/juztin/wombat/config"
	"bitbucket.org/juztin/wombat/users"
)

type Data struct {
	cookieKey        string
	IsProd           bool
	MediaURL, Domain string
	User             users.User
}

func New(ctx wombat.Context) Data {
	t := new(Data)
	t.IsProd = config.IsProd
	t.MediaURL = config.MediaURL
	t.Domain = config.ServerDomain
	t.User = ctx.User

	return *t
}
