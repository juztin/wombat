// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import "bitbucket.org/juztin/config"

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
	if config.Load() == nil {
		setFromCfg()
	}
}

func setFromCfg() {
	// server
	config.SetCfgGroupBool("server", "isProd", &IsProd)
	config.SetCfgGroupString("server", "host", &ServerHost)
	config.SetCfgGroupInt("server", "port", &ServerPort)
	config.SetCfgGroupString("server", "domain", &ServerDomain)
	config.SetCfgGroupString("server", "signinUrl", &SigninURL)
	config.SetCfgGroupBool("server", "unixSock", &UnixSock)
	config.SetCfgGroupString("server", "unixSockFile", &UnixSockFile)
	config.SetCfgGroupBool("server", "tls", &TLS)
	config.SetCfgGroupString("server", "tlsCert", &TLSCert)
	config.SetCfgGroupString("server", "tlsKey", &TLSKey)

	// cookie
	config.SetCfgGroupString("cookie", "name", &Cookie)
	config.SetCfgGroupBool("cookie", "expires", &CookieExpires)
	config.SetCfgGroupString("cookie", "expireHash", &CookieExpireHash)
	config.SetCfgGroupInt("cookie", "expireTime", &CookieExpireTime)
	config.SetCfgGroupString("cookie", "path", &CookiePath)
	config.SetCfgGroupBool("cookie", "httpOnly", &CookieHttpOnly)
	config.SetCfgGroupBool("cookie", "secure", &CookieSecure)

	// media
	config.SetCfgGroupString("media", "url", &MediaURL)

	// templates
	config.SetCfgGroupString("templates", "editRoute", &TmplEditRt)
	config.SetCfgGroupBool("templates", "editErr", &TmplEditErr)
	config.SetCfgGroupString("templates", "path", &TmplPath)

	// database
	config.SetCfgGroupString("user", "reader", &UserReader)
}
