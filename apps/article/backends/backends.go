package backends

import (
	"fmt"
	"time"

	b "bitbucket.org/juztin/wombat/backends"
	"bitbucket.org/juztin/wombat/config"
)

type Article interface {
	ByTitlePath(titlePath string) (ArticleData, b.Error)
	Create(a ArticleData) b.Error
	UpdateContent(titlePath, content string, modified time.Time) b.Error
	Delete(titlepath string) b.Error
	SetImg(titlePath string, img ImgData) b.Error
	SetImgs(titlePath string, imgs []ImgData) b.Error
}

type ImgData struct {
	Src, Alt string
	W, H     int
}

type ArticleData struct {
	TitlePath         string
	Title             string
	Content           string
	IsActive          bool
	Created, Modified time.Time
	Img               ImgData
	Imgs              []ImgData
}

func ArticleBackend() (Article, error) {
	if b, err := b.Open(config.Backend + ":apps:article"); err != nil {
		return nil, err
	} else if u, ok := b.(Article); !ok {
		return nil, fmt.Errorf("backends: invalid apps:article backend type %v", b)
	} else {
		return u, nil
	}
	return nil, nil
}
