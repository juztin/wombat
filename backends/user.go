package backends

type User interface {
	Authenticate(username, password string) (UserData, Error)
}

type UserData struct {
	Username  string
	Firstname string
	Lastname  string
	Hash      string
	Status    int
}