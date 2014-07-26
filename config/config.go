// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import "code.minty.io/config"

var (
	IsProd           = false
	Cookie           = "oatmeal"
	CookieExpires    = true
	CookieExpireHash = "$uP{r - s@lTy _ st|_|fF"
	CookieExpireTime = 15
	CookiePath       = "/"
	CookieHttpOnly   = true
	CookieSecure     = true
	ServerHost       = "127.0.0.1"
	ServerPort       = 9999
	ServerDomain     = "juzt.in"
	TLS              = false
	TLSCert          = ""
	TLSKey           = ""
	UnixSock         = false
	UnixSockFile     = "/tmp/juztin.sock"
	MediaURL         = "media.juzt.in/"
	SigninURL        = "/account/"
	TmplEditRt       = "^/_dt/$"
	TmplEditErr      = false
	TmplPath         = "./templates"
	UserReader       = "mongo:user-reader"
)

func init() {
	//if config.Load() == nil {
	setFromCfg()
	//}
}

func setBool(group, key string, b *bool) {
	if v, ok := config.GroupBool(group, key); ok {
		*b = v
	}
}
func setInt(group, key string, i *int) {
	if v, ok := config.GroupInt(group, key); ok {
		*i = v
	}
}
func setString(group, key string, s *string) {
	if v, ok := config.GroupString(group, key); ok {
		*s = v
	}
}

func setFromCfg() {
	// server
	setBool("server", "isProd", &IsProd)
	setString("server", "host", &ServerHost)
	setInt("server", "port", &ServerPort)
	setString("server", "domain", &ServerDomain)
	setString("server", "signinURL", &SigninURL)
	setBool("server", "unixSock", &UnixSock)
	setString("server", "unixSockFile", &UnixSockFile)
	setBool("server", "tls", &TLS)
	setString("server", "tlsCert", &TLSCert)
	setString("server", "tlsKey", &TLSKey)

	// cookie
	setString("cookie", "name", &Cookie)
	setBool("cookie", "expires", &CookieExpires)
	setString("cookie", "expireHash", &CookieExpireHash)
	setInt("cookie", "expireTime", &CookieExpireTime)
	setString("cookie", "path", &CookiePath)
	setBool("cookie", "httpOnly", &CookieHttpOnly)
	setBool("cookie", "secure", &CookieSecure)

	// media
	setString("media", "url", &MediaURL)

	// templates
	setString("templates", "editRoute", &TmplEditRt)
	setBool("templates", "TmplEditErr", &TmplEditErr)
	setString("templates", "TmplPath", &TmplPath)

	// database
	setString("user", "reader", &UserReader)
}
