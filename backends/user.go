package backends

type User interface {
	GetByUsername(username string) (UserData, Error)
	GetBySession(session string) (UserData, Error)
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
