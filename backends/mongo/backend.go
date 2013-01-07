package mongo

import (
	"log"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

type userBackend struct {
	session *mgo.Session
}

func init() {
	if session, err := mgo.Dial(config.MongoURL); err != nil {
		log.Fatal("Failed to retrieve Mongo session: ", err)
	} else {
		// set monotonic mode
		session.SetMode(mgo.Monotonic, true)
		// register backend
		backends.Register("mongo", backends.Backend{
			User: userBackend{session},
		})
	}
}

func (b userBackend) db() (*mgo.Session, *mgo.Database) {
	s := b.session.New()
	return s, s.DB(config.MongoDB)
}

func (b userBackend) col(name string) (*mgo.Session, *mgo.Collection) {
	s := b.session.New()
	return s, s.DB(config.MongoDB).C(name)
}

func (b userBackend) GetByUsername(username string) (backends.UserData, backends.Error) {
	s, c := b.col("users")
	defer s.Close()

	d := new(backends.UserData)
	if err := c.Find(bson.M{"username": username}).One(&d); err != nil {
		return *d, backends.NewError(backends.StatusDatastoreError, "user authentication", err)
	}

	return *d, nil
}

func (b userBackend) GetBySession(key string) (backends.UserData, backends.Error) {
	s, c := b.col("users")
	defer s.Close()

	d := new(backends.UserData)
	if err := c.Find(bson.M{"session": key}).One(&d); err != nil {
		return *d, backends.NewError(backends.StatusDatastoreError, "user cache", err)
	}

	return *d, nil
}

func (b userBackend) SetSession(username, key string) backends.Error {
	s, c := b.col("users")
	defer s.Close()

	selector := bson.M{"username": username}
	change := bson.M{"$set": bson.M{"session": key}}
	if err := c.Update(selector, change); err != nil {
		log.Println("Failed to update session for: ", username, " : ", err)
		return backends.NewError(backends.StatusDatastoreError, "user update", err)
	}
	return nil
}
