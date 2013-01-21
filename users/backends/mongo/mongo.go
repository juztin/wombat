package mongo

import (
	//"errors"
	"log"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
	"bitbucket.org/juztin/wombat/users"
)

const COL_NAME = "users"

type Backend struct {
	session *mgo.Session
}

func init() {
	if session, err := mgo.Dial(config.MongoURL); err != nil {
		log.Fatal("Failed to retrieve Mongo session: ", err)
	} else {
		// set monotonic mode
		session.SetMode(mgo.Monotonic, true)
		// register backend
		b := Backend{session}
		backends.Register("mongo:user-reader", b)
		backends.Register("mongo:user-writer", b)
	}
}

func (b Backend) col() (*mgo.Session, *mgo.Collection) {
	s := b.session.New()
	return s, s.DB(config.MongoDB).C(COL_NAME)
}

// Reader
func (b Backend) ByUsername(username string) (users.User, error) {
	s, c := b.col()
	defer s.Close()

	u := users.Model{b, new(users.Data)}
	if err := c.Find(bson.M{"username": username}).One(u.Data); err != nil {
		return u, backends.NewError(backends.StatusDatastoreError, "user authentication", err)
	}

	return u, nil
}

func (b Backend) BySession(key string) (users.User, error) {
	s, c := b.col()
	defer s.Close()

	u := users.Model{b, new(users.Data)}
	if err := c.Find(bson.M{"session": key}).One(&u.Data); err != nil {
		return u, backends.NewError(backends.StatusDatastoreError, "user cache", err)
	}

	return u, nil
}

// Writer
//func (b Backend) UpdateSession(username, key string) backends.Error {
func (b Backend) UpdateSession(username, key string) error {
	s, c := b.col()
	defer s.Close()

	selector := bson.M{"username": username}
	change := bson.M{"$set": bson.M{"session": key}}
	if err := c.Update(selector, change); err != nil {
		log.Println("Failed to update 'session' for: ", username, " : ", err)
		return backends.NewError(backends.StatusDatastoreError, "user update", err)
	}
	return nil
}

func (b Backend) UpdateLastSignin(username string, on time.Time) error {
	s, c := b.col()
	defer s.Close()

	selector := bson.M{"username": username}
	change := bson.M{"$set": bson.M{"lastSignin": on}}
	if err := c.Update(selector, change); err != nil {
		log.Println("Failed to update 'lastSignin' for: ", username, " : ", err)
		return backends.NewError(backends.StatusDatastoreError, "user update", err)
	}
	return nil
}
