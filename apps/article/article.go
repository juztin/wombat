package article

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"bitbucket.org/juztin/wombat"
	"bitbucket.org/juztin/wombat/config"
	"bitbucket.org/juztin/wombat/imgconv"
)

var (
	mongo *mgo.Session
	imgRoot string
)

type Img struct {
	Src, Alt string
	W, H int
}

type Article struct {
	TitlePath string
	Title string
	Content string
	IsActive bool
	Created, Modified time.Time
	Img Img
	Imgs []Img
}

func init() {
	if session, err := mgo.Dial(config.MongoURL); err != nil {
		log.Fatal("Failed to retrieve Mongo session: ", err)
	} else {
		// set monotonic mode
		session.SetMode(mgo.Monotonic, true)
		mongo = session
	}

	imgRoot, _ = config.GetString("ArticleImgRoot")
}

func col() (*mgo.Session, *mgo.Collection) {
	s := mongo.New()
	return s, s.DB(config.MongoDB).C("articles")
}

func formFileImage(ctx wombat.Context, titlePath string) (string, *os.File, error) {
	// grab the image from the request
	f, h, err := ctx.Request.FormFile("image")
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	// create the project path
	imgPath := filepath.Join(imgRoot, titlePath)
	if err := os.MkdirAll(imgPath, 0755); err != nil {
		return "", nil, err
	}

	// replace spaces in image name
	imgName := strings.Replace(h.Filename, " ", "-", -1)

	// save the image, as a temp file
	//t, err := os.OpenFile(filepath.Join(imgPath, imgName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	t, err := ioutil.TempFile(imgPath, "."+imgName)
	if err != nil {
		return imgName, nil, err
	}
	defer t.Close()
	// copy the image to the temp file
	if _, err := io.Copy(t, f); err != nil {
		return imgName, nil, err
	}

	return imgName, t, nil
}

func convertToJpg(imgName string, f *os.File, isThumb bool) (p string, i image.Image, err error) {
	x := filepath.Ext(imgName)
	//n = imgName[:len(imgName)-len(x)]+".jpg"
	s := ""
	n := imgName[:len(imgName)-len(x)]
	if isThumb {
		s = ".thumb"
	}
	p = fmt.Sprintf("%s%s.jpg", n, s)

	if isThumb {
		i, err = imgconv.ResizeWidthToJPG(f.Name(), p, true, 200)
	} else {
		i, err = imgconv.ConvertToJPG(f.Name(), p, true)
	}
	return
}

func Load(ctx wombat.Context, titlepath string) (Article, bool) {
	s, c := col()
	defer s.Close()

	// load the article
	d := new(Article)
	if err := c.Find(bson.M{"titlepath": titlepath}).One(&d); err != nil {
		log.Println(err)
		return *d, false
	}
	return *d, true
}

func Create(ctx wombat.Context, title string) (string, bool) {
	s, c := col()
	defer s.Close()

	// create a new article, based on the current time
	t := time.Now()
	titlePath := fmt.Sprintf("/%d/%02d/%02d/%s/",
		t.Year(),
		t.Month(),
		t.Day(),
		strings.Replace(title, " ", "-", -1))

	// set the article's TitlePath
	a := &Article{TitlePath: titlePath, Title: title, Created: t}
	if err := c.Insert(a); err != nil {
		log.Println("Failed to create article: ", titlePath, " : ", err)
		return titlePath, false
	}

	return titlePath, true
}

//func UpdateContent(ctx wombat.Context, titlePath string, article *interface{}) error {
func UpdateContent(ctx wombat.Context, titlePath string, article interface{}) error {
	// get the content
	content := ctx.FormValue("content")

	s, c := col()
	defer s.Close()

	// update the article's content
	selector := bson.M{"titlepath": titlePath}
	change := bson.M{"$set": bson.M{"content": &content}}
	if err := c.Update(selector, change); err != nil {
		//ctx.HttpError(500)
		return err
	}

	// render either JSON|HTML
	if d := ctx.FormValue("d"); d == "json" {
		// return the new article's content (json)
		ctx.Writer.Header().Set("Content-Type", "application/json")
		if j, err := json.Marshal(map[string] string { "content": content }); err != nil {
			//ctx.HttpError(500)
			return err
		} else {
			ctx.Writer.Write(j)
		}
	} else {
		// render the article (html)
		//a := new(Article)
		//c.Find(bson.M{"titlepath": titlePath}).One(&a)
		//renderArticle(ctx, *a)
		c.Find(bson.M{"titlepath": titlePath}).One(article)
		//renderArticle(ctx, *a)
	}
	return nil
}

func Delete(ctx wombat.Context, titlePath string) {
	//getArticle(ctx, titlePath)
}

func ImgHandler(ctx wombat.Context, titlePath string, isThumb bool) (error, bool) {
	// get the article
	a, ok := Load(ctx, titlePath)
	if !ok {
		//ctx.HttpError(404)
		return errors.New("Article not found"), false
	}

	// create the image, from the POST
	imgName, f, err := formFileImage(ctx, titlePath)
	if err != nil {
		//ctx.HttpError(500)
		//return
		return err, false
	}

	// convert image to jpeg
	n, i, err := convertToJpg(imgName, f, isThumb)
	if err != nil {
		//ctx.HttpError(500)
		//return
		return err, false
	}

	// create the image object
	var imgs []Img
	exists := false
	s := i.Bounds().Size()
	if isThumb {
		// TODO -> maybe remove the image upon successful addition of the new one
		// remove current thumb
		imgPath := filepath.Join(imgRoot, titlePath)
		os.Remove(filepath.Join(imgPath, a.Img.Src))
		// set the new thumb
		a.Img = Img{n, imgName, s.X, s.Y}
	} else {
		l := len(a.Imgs)
		imgs = make([]Img, l, l+1)
		copy(imgs, a.Imgs)
		for _, v := range imgs {
			if v.Src == n {
				v.W, v.H = s.X, s.Y
				exists = true
			}
		}
		if !exists {
			imgs = append(imgs, Img{n, "", s.X, s.Y})
		}
	}

	// update article images
	session, col := col()
	defer session.Close()

	selector := bson.M{"titlepath": titlePath}
	var change bson.M
	if isThumb {
		change = bson.M{"$set": bson.M{"img": &a.Img}}
	} else {
		change = bson.M{"$set": bson.M{"imgs": imgs}}//&a.Imgs}}
	}
	if err := col.Update(selector, change); err != nil {
		//ctx.HttpError(500)
		//return
		return err, false
	}

	// append the image to the article
	if !isThumb {
		a.Imgs = imgs
	}

	// return either a JSON/HTML response
	if d := ctx.FormValue("d"); d == "json" {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		k := "image"
		if isThumb {
			k = "thumb"
		}
		j := fmt.Sprintf(`{"%s":"%s","w":%d,"h":%d}`, k, n, s.X, s.Y)
		ctx.Writer.Write([]byte(j))
		return nil, true
	} else {
		return nil, false
		//renderArticle(ctx, a)
	}
	return nil, true
}

func AddThumb(ctx wombat.Context, titlePath string) {
	ImgHandler(ctx, titlePath, true)
}

func AddImage(ctx wombat.Context, titlePath string) {
	ImgHandler(ctx, titlePath, false)
}

func DelImage(ctx wombat.Context, titlePath string) error {
	// get the article
	a, ok := Load(ctx, titlePath)
	if !ok {
		//getArticle(ctx, titlePath)
		//return
		return errors.New("Article not found")
	}

	// get the image to be deleted
	src := ctx.FormValue("image")

	// update the articles images
	n := []Img{}
	for _, i := range a.Imgs {
		if i.Src != src {
			n = append(n, i)
		} else {
			imgPath := filepath.Join(imgRoot, titlePath)
			os.Remove(filepath.Join(imgPath, i.Src))
		}
		//p.Imgs = append(p.Imgs[:i], p.Imgs[i+1:]...)
	}

	// if a matchine image was found, remove it
	if len(n) != len(a.Imgs) {
		a.Imgs = n

		s, c := col()
		defer s.Close()

		selector := bson.M{"titlepath": titlePath}
		change := bson.M{"$set": bson.M{"imgs": &a.Imgs}}
		if err := c.Update(selector, change); err != nil {
			//ctx.HttpError(500)
			return err
		}
	}

	// redirect back to the update page, when a regular POST
	if d := ctx.FormValue("d"); d != "json" {
		ctx.Writer.WriteHeader(303)
	}
	return nil
}
