package data

import (
	"bitbucket.org/juztin/dingo"

	"bitbucket.org/juztin/virginia/config"
	"bitbucket.org/juztin/virginia/models/user"
)

type Data struct {
	IsProd			  bool
	StaticURL, Domain string
	User			  user.User
}

var anonymousUser = user.Anonymous()

func New(ctx dingo.Context) Data {
	t := new(Data)
	t.IsProd = config.IsProd
	t.StaticURL = config.StaticURL
	t.Domain = config.Domain
	t.User = anonymousUser

	return *t
}