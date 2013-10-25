// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongo

import (
	//"errors"
	"log"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"bitbucket.org/juztin/config"
	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/users"
)

const COL_NAME = "users"

var (
	db = "main"
)

type Backend struct {
	session *mgo.Session
}

type queryFunc func(c *mgo.Collection)

func init() {
	if url, ok := config.GroupString("db", "mongoURL"); !ok {
		log.Fatal("users mongo: MongoURL missing from configuration")
	} else if session, err := mgo.Dial(url); err != nil {
		log.Fatal("Failed to retrieve Mongo session: ", err)
	} else {
		// set monotonic mode
		session.SetMode(mgo.Monotonic, true)
		// register backend
		b := Backend{session}
		backends.Register("mongo:user-reader", b)
		backends.Register("mongo:user-writer", b)
	}

	if d, ok := config.GroupString("db", "mongoDB"); ok {
		db = d
	}
}

func (b Backend) Col() (*mgo.Session, *mgo.Collection) {
	s := b.session.New()
	return s, s.DB(db).C(COL_NAME)
}
func (b Backend) Query(fn queryFunc) {
	s, c := b.Col()
	defer s.Close()
	fn(c)
}

// Reader
func (b Backend) ByUsername(username string) (users.User, error) {
	s, c := b.Col()
	defer s.Close()

	u := users.Model{b, new(users.Data)}
	if err := c.Find(bson.M{"username": username}).One(u.Data); err != nil {
		return u, backends.NewError(backends.StatusDatastoreError, "user authentication", err)
	}

	return u, nil
}

func (b Backend) BySession(key string) (users.User, error) {
	s, c := b.Col()
	defer s.Close()

	u := users.Model{b, new(users.Data)}
	if err := c.Find(bson.M{"session": key}).One(&u.Data); err != nil {
		return u, backends.NewError(backends.StatusDatastoreError, "user cache", err)
	}

	return u, nil
}

// Writer
func (b Backend) UpdateSession(username, key string) error {
	s, c := b.Col()
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
	s, c := b.Col()
	defer s.Close()

	selector := bson.M{"username": username}
	change := bson.M{"$set": bson.M{"lastSignin": on}}
	if err := c.Update(selector, change); err != nil {
		log.Println("Failed to update 'lastSignin' for: ", username, " : ", err)
		return backends.NewError(backends.StatusDatastoreError, "user update", err)
	}
	return nil
}
