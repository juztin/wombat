package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bitbucket.org/juztin/dingo/views"
	"bitbucket.org/juztin/wombat"
	"bitbucket.org/juztin/wombat/apps/article"
	"bitbucket.org/juztin/wombat/config"
	"bitbucket.org/juztin/wombat/template/data"
)

type articleData struct {
	data.Data
	Article article.Article
}

var (
	imgRoot     string
	articlePath string
	listView    string
	articleView string
	createView  string
	updateView  string
)

func Init(s wombat.Server, basePath, list, article, create, update string) {
	imgRoot, _ = config.GetString("ArticleImgRoot")

	articlePath = basePath
	listView = list
	articleView = article
	createView = create
	updateView = update

	// routes
	s.ReRouter(fmt.Sprintf("^%s/$", articlePath)).
		Get(listArticles).
		Post(newArticle)

	s.RRouter(fmt.Sprintf("^%s(/\\d{4}/\\d{2}/\\d{2}/[a-zA-Z0-9-]+/)$", articlePath)).
		Get(getArticle).
		Post(postArticle)
}

/* -------------------------------- Handlers -------------------------------- */
func listArticles(ctx wombat.Context) {
	v := listView
	if ctx.User.IsAdmin() {
		if action := ctx.FormValue("action"); action == "create" {
			v = createView
		}
	}

	views.Execute(ctx.Context, v, data.New(ctx))
}

func newArticle(ctx wombat.Context) {
	if ctx.User.IsAdmin() {
		if t, ok := Create(ctx); ok {
			ctx.Redirect(articlePath + t)
		}
	}
	views.Execute(ctx.Context, listView, data.New(ctx))
}

func renderArticle(ctx wombat.Context, a article.Article) {
	d := &articleData{data.New(ctx), a}
	views.Execute(ctx.Context, updateView, d)
}

func getArticle(ctx wombat.Context, titlePath string) {
	v := articleView
	if ctx.User.IsAdmin() {
		if action := ctx.FormValue("action"); action == "update" {
			v = updateView
		}
	}

	a, _ := article.ByTitlePath(titlePath)
	d := &articleData{data.New(ctx), a}
	views.Execute(ctx.Context, v, d)
}

func postArticle(ctx wombat.Context, titlePath string) {
	if ctx.User.IsAdmin() {
		if action := ctx.FormValue("action"); action != "" {
			switch action {
			default:
				getArticle(ctx, titlePath)
			case "update":
				UpdateContent(ctx, titlePath)
			case "delete":
				Delete(ctx, titlePath)
			case "addImage":
				AddImage(ctx, titlePath)
			case "addThumb":
				AddThumb(ctx, titlePath)
			case "delImage":
				DelImage(ctx, titlePath)
			}
		} else {
			getArticle(ctx, titlePath)
		}
	} else {
		getArticle(ctx, titlePath)
	}
}

/* ------------------------------------  ------------------------------------ */

func Create(ctx wombat.Context) (string, bool) {
	title := ctx.FormValue("title")
	if title != "" {
		return title, false
	}

	return article.Create(title)
}

//func UpdateContent(ctx wombat.Context, titlePath string, article *interface{}) error {
func UpdateContent(ctx wombat.Context, titlePath string) {
	a, ok := article.ByTitlePath(titlePath)
	if !ok {
		ctx.HttpError(404)
		return
	}

	// get the content
	content := ctx.FormValue("content")
	a.UpdateContent(content)

	// render either JSON|HTML
	if d := ctx.FormValue("d"); d == "json" {
		// return the new article's content (json)
		ctx.Writer.Header().Set("Content-Type", "application/json")
		if j, err := json.Marshal(map[string]string{"content": content}); err != nil {
			log.Println("Failed to marshal article's content to JSON : ", err)
			ctx.HttpError(500)
		} else {
			ctx.Writer.Write(j)
		}
	} else {
		//renderArticle(ctx, *a)
	}
}

func Delete(ctx wombat.Context, titlePath string) {
	//getArticle(ctx, titlePath)
}

func ImgHandler(ctx wombat.Context, titlePath string, isThumb bool) {
	// get the article
	a, ok := article.ByTitlePath(titlePath)
	if !ok {
		ctx.HttpError(404)
	}

	// create the image, from the POST
	imgName, f, err := formFileImage(ctx, titlePath)
	if err != nil {
		log.Println("Failed to create temporary image from form-file: ", titlePath, " : ", err)
		ctx.HttpError(500)
		return
	}

	// convert image to jpeg
	n, i, err := convertToJpg(imgName, f, isThumb)
	if err != nil {
		log.Println("Failed to convert image to jpeg for article: ", titlePath, " : ", err)
		ctx.HttpError(500)
		return
	}

	// create the image object
	var imgs []article.Img
	exists := false
	s := i.Bounds().Size()
	if isThumb {
		// TODO -> maybe remove the image upon successful addition of the new one
		// remove current thumb
		imgPath := filepath.Join(imgRoot, titlePath)
		os.Remove(filepath.Join(imgPath, a.Img.Src))
	} else {
		l := len(a.Imgs)
		imgs = make([]article.Img, l, l+1)
		copy(imgs, a.Imgs)
		for _, v := range imgs {
			if v.Src == n {
				v.W, v.H = s.X, s.Y
				exists = true
			}
		}
		if !exists {
			imgs = append(imgs, article.Img{n, "", s.X, s.Y})
		}
	}

	// update article images
	if isThumb {
		ok = a.SetImg(article.Img{n, imgName, s.X, s.Y})
	} else {
		ok = a.SetImgs(imgs)
	}

	if !ok {
		log.Println("Failed to persit new image: ", imgName, " for article: ", titlePath)
		ctx.HttpError(500)
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
	} else {
		//renderArticle(ctx, a)
	}
}

func AddThumb(ctx wombat.Context, titlePath string) {
	ImgHandler(ctx, titlePath, true)
}

func AddImage(ctx wombat.Context, titlePath string) {
	ImgHandler(ctx, titlePath, false)
}

func DelImage(ctx wombat.Context, titlePath string) {
	// get the article
	a, ok := article.ByTitlePath(titlePath)
	if !ok {
		//getArticle(ctx, titlePath)
		return
	}

	// get the image to be deleted
	src := ctx.FormValue("image")

	// update the articles images
	n := []article.Img{}
	for _, i := range a.Imgs {
		if i.Src != src {
			n = append(n, i)
		} else {
			imgPath := filepath.Join(imgRoot, titlePath)
			os.Remove(filepath.Join(imgPath, i.Src))
		}
		//p.Imgs = append(p.Imgs[:i], p.Imgs[i+1:]...)
	}

	// if a matching image was found, remove it
	if len(n) != len(a.Imgs) {
		if ok = a.SetImgs(n); !ok {
			log.Println("Failed to persit image deletion: ", src, " for article: ", titlePath)
			ctx.HttpError(500)
			return
		}
	}

	// redirect back to the update page, when a regular POST
	if d := ctx.FormValue("d"); d != "json" {
		ctx.Writer.WriteHeader(303)
	}
}
