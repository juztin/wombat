package data

import (
	//"bitbucket.org/juztin/dingo"

	"bitbucket.org/juztin/virginia"
	"bitbucket.org/juztin/virginia/config"
	"bitbucket.org/juztin/virginia/models/user"
)

type Data struct {
	cookieKey        string
	IsProd			 bool
	MediaURL, Domain string
	User			 user.User
}

var anonymousUser = user.Anonymous()

//func New(ctx dingo.Context) Data {
func New(ctx virginia.Context) Data {
	t := new(Data)
	t.IsProd = config.IsProd
	t.MediaURL = config.MediaURL
	t.Domain = config.Domain
	// TODO -> load the user from cache here?
	//t.User = anonymousUser
	//t.User = user.FromSession("jr", "password")
	//t.User = user.Authenticate("jr", "password")
	t.User = ctx.User

	return *t
}

/*func (d *Data) LoadUser() *Data {
	d.User = user.Authenticate("jr", "password")

	return d
}*/