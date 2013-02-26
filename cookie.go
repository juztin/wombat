package wombat

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/juztin/wombat/config"
)

// taken from https://groups.google.com/forum/?fromgroups=#!msg/golang-nuts/owCogizIuZs/ZzmwkQGrlnEJ
// by Russ Cox
func uuid() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
		//log.Println("Failed to create uuid: ", err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func expireCookieHash(s string) string {
	hash := hmac.New(sha256.New, []byte(config.CookieExpireHash))
	hash.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func cookie(val string) *http.Cookie {
	c := new(http.Cookie)
	c.Name = config.Cookie
	c.Value = val
	c.Domain = config.ServerDomain
	c.Path = config.CookiePath
	c.HttpOnly = config.CookieHttpOnly
	return c
}

func expiredCookie() *http.Cookie {
	c := cookie("")
	c.MaxAge = -1
	return c
}

func signedExpireCookie(val string, ut int64) *http.Cookie {
	v := fmt.Sprintf("%d|%s", ut, val)
	v = fmt.Sprintf("%s|%s", v, expireCookieHash(v))
	return cookie(v)
}

func expireCookieParts(c *http.Cookie) (t int64, k string, h string) {
	// split cookit into [timestamp, key, signature]
	if p := strings.Split(c.Value, "|"); len(p) != 3 {
		log.Println("Invalid expire cookie format: ", c.Value)
	} else {
		if ct, err := strconv.ParseInt(p[0], 10, 64); err != nil {
			log.Println("Invalid cookie expiration: ", p[0], " : ", err)
		} else {
			t = ct
			k = p[1]
			h = p[2]
		}
	}

	return
}

func NewCookie() (string, *http.Cookie) {
	key := uuid()
	return key, cookie(key)
}

func NewExpireCookie() (string, *http.Cookie) {
	key := uuid()
	cookie := signedExpireCookie(key, time.Now().UTC().Unix())
	return key, cookie
}

func UpdatedExpireCookie(r *http.Request) (*http.Cookie, string, bool) {
	c, err := r.Cookie(config.Cookie)
	if err != nil {
		return nil, "", false
	}

	t, k, h := expireCookieParts(c)
	now := time.Now().UTC().Unix()
	// verify time
	if (now - t) > int64(config.CookieExpireTime*60) {
		// cookie has expired
		return expiredCookie(), "", false
	}

	// verify signature
	if s := expireCookieHash(fmt.Sprintf("%d|%s", t, k)); s != h {
		log.Println("Invalid cookie signature: ", c)
		return expiredCookie(), "", false
	}

	// return update cookie
	return signedExpireCookie(k, now), k, true
}

func GetCookieSession(r *http.Request) string {
	if c, err := r.Cookie(config.Cookie); err == nil {
		return c.Value
	}
	return ""
}
