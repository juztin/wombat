package backends

type User interface {
	Authenticate(username, password string) (UserData, Error)
	FromCache(key string) (UserData, Error)
}

type UserData struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Role      int
	Status    int
}