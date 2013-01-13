package article

import (
	"fmt"
	"log"
	"strings"
	"time"

	"bitbucket.org/juztin/wombat/apps/article/backends"
)

var b backends.Article

type Img struct {
	Src, Alt string
	W, H     int
}

type Article struct {
	TitlePath         string
	Title             string
	Content           string
	IsActive          bool
	Created, Modified time.Time
	Img               Img
	Imgs              []Img
}

func init() {
	if backend, err := backends.ArticleBackend(); err != nil {
		log.Fatal("Failed to get apps:article backend")
	} else {
		b = backend
	}
}

func imgFromData(d backends.ImgData) Img {
	return Img{d.Src, d.Alt, d.W, d.H}
}

func imgsFromData(d []backends.ImgData) []Img {
	imgs := make([]Img, len(d))
	for i, c := range d {
		imgs[i] = imgFromData(c)
	}
	return imgs
}

func dataFromImg(i Img) backends.ImgData {
	return backends.ImgData{i.Src, i.Alt, i.W, i.H}
}

func dataFromImgs(s []Img) []backends.ImgData {
	d := make([]backends.ImgData, len(s))
	for i, c := range s {
		d[i] = dataFromImg(c)
	}
	return d
}

func articleFromData(m backends.ArticleData) Article {
	a := new(Article)
	a.TitlePath = m.TitlePath
	a.Title = m.Title
	a.Content = m.Content
	a.IsActive = m.IsActive
	a.Created = m.Created
	a.Modified = m.Modified
	a.Img = imgFromData(m.Img)
	a.Imgs = imgsFromData(m.Imgs)
	return *a
}

func ByTitlePath(titlePath string) (a Article, ok bool) {
	if d, err := b.ByTitlePath(titlePath); err != nil {
		log.Println(err)
	} else {
		a, ok = articleFromData(d), true
	}
	return
}

func Create(title string) (string, bool) {
	// create a new article, based on the current time
	t := time.Now()
	titlePath := fmt.Sprintf("/%d/%02d/%02d/%s/",
		t.Year(),
		t.Month(),
		t.Day(),
		strings.Replace(title, " ", "-", -1))

	// create a new article, based on the current time
	d := backends.ArticleData{TitlePath: titlePath, Title: title, Created: t}
	if err := b.Create(d); err != nil {
		log.Println(err)
		return titlePath, false
	}
	return titlePath, true
}

func (a *Article) UpdateContent(content string) bool {
	t := time.Now()
	if err := b.UpdateContent(a.TitlePath, content, t); err != nil {
		log.Println(err)
		return false
	}
	a.Content = content
	a.Modified = t
	return true
}

func (a *Article) Delete() bool {
	return false
}

func (a *Article) SetImg(img Img) bool {
	d := dataFromImg(img)
	if err := b.SetImg(a.TitlePath, d); err != nil {
		log.Println(err)
		return false
	}
	a.Img = img
	return true
}

func (a *Article) SetImgs(imgs []Img) bool {
	d := dataFromImgs(imgs)
	if err := b.SetImgs(a.TitlePath, d); err != nil {
		log.Println(err)
		return false
	}
	a.Imgs = imgs
	return true
}
