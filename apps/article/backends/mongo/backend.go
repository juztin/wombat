package mongo

import (
	"log"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	ab "bitbucket.org/juztin/wombat/apps/article/backends"
	"bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

const COL_NAME = "articles"

type ArticleBackend struct {
	session *mgo.Session
}

func init() {
	if session, err := mgo.Dial(config.MongoURL); err != nil {
		log.Fatal("Failed to retrieve Mongo session: ", err)
	} else {
		// set monotonic mode
		session.SetMode(mgo.Monotonic, true)
		// register backend
		backends.Register("mongo:apps:article", ArticleBackend{session})
	}
}

func (b ArticleBackend) db() (*mgo.Session, *mgo.Database) {
	s := b.session.New()
	return s, s.DB(config.MongoDB)
}

func (b ArticleBackend) col() (*mgo.Session, *mgo.Collection) {
	s := b.session.New()
	return s, s.DB(config.MongoDB).C(COL_NAME)
}

func (b ArticleBackend) ByTitlePath(titlePath string) (ab.ArticleData, backends.Error) {
	s, c := b.col()
	defer s.Close()

	// load the article
	d := new(ab.ArticleData)
	if err := c.Find(bson.M{"titlepath": titlePath}).One(&d); err != nil {
		//return *d, backends.NewError(backends.StatusDatastoreError, "Find by article", err)
		return *d, backends.NewError(backends.StatusNotFound, "Article not found", err)
	}
	return *d, nil
}

func (b ArticleBackend) Create(a ab.ArticleData) backends.Error {
	s, c := b.col()
	defer s.Close()

	if err := c.Insert(a); err != nil {
		//log.Println("Failed to create article: ", titlePath, " : ", err)
		return backends.NewError(backends.StatusDatastoreError, "Failed to create article", err)
	}

	return nil
}

func (b ArticleBackend) UpdateContent(titlePath, content string, modified time.Time) backends.Error {
	s, c := b.col()
	defer s.Close()

	// update the article's content
	selector := bson.M{"titlepath": titlePath}
	change := bson.M{"$set": bson.M{"content": &content, "modified": modified}}
	if err := c.Update(selector, change); err != nil {
		//ctx.HttpError(500)
		return backends.NewError(backends.StatusDatastoreError, "Failed to update article's content", err)
	}
	return nil
}

func (b ArticleBackend) Delete(titlepath string) backends.Error {
	return nil
}

func (b ArticleBackend) SetImg(titlePath string, img ab.ImgData) backends.Error {
	// update article image/thumb
	session, col := b.col()
	defer session.Close()

	selector := bson.M{"titlepath": titlePath}
	change := bson.M{"$set": bson.M{"img": img}}
	if err := col.Update(selector, change); err != nil {
		return backends.NewError(backends.StatusDatastoreError, "Failed to update image/thumb", err)
	}
	return nil
}

func (b ArticleBackend) SetImgs(titlePath string, imgs []ab.ImgData) backends.Error {
	// update article images
	session, col := b.col()
	defer session.Close()

	selector := bson.M{"titlepath": titlePath}
	change := bson.M{"$set": bson.M{"imgs": imgs}} //&a.Imgs}}
	if err := col.Update(selector, change); err != nil {
		return backends.NewError(backends.StatusDatastoreError, "Failed to update images", err)
	}
	return nil
}
