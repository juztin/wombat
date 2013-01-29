package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const CONFIG_FILE = "config.json"

var (
	IsProd           = false
	Cookie           = "oatmeal"
	CookieExpires    = true
	CookieExpireHash = "$uP{r - s@lTy _ st|_|fF"
	CookieExpireTime = 15
	ServerHost       = "127.0.0.1"
	ServerPort       = 9999
	ServerDomain     = "juzt.in"
	UnixSock         = false
	UnixSockFile     = "/tmp/juztin.sock"
	MediaURL         = "media.juzt.in/"
	TmplEditRt       = "^/_dt/$"
	TmplEditErr      = false
	TmplPath         = "./templates"
	UserReader       = "mongo:user-reader"
	//UserWriter       = "mongo:user-writer"

	//MongoURL   = "localhost"
	//MongoDB    = "wombat"
	//SqliteFile = "./db.sqlite"

	cfg map[string]interface{}
)

//func Load() {
func init() {
	// get|read configuration from file
	p, c, err := getConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration file: %s, from: %s\n%v", CONFIG_FILE, p, err)
	}

	var j interface{}
	if err := json.Unmarshal(c, &j); err != nil {
		log.Fatalf("Failed to read configuration file: %s, from: %s\n%v", CONFIG_FILE, p, err)
	}

	cfg = j.(map[string]interface{})
	setFromCfg()
}

func setFromCfg() {
	// server
	setCfgGroupBool("server", "isProd", &IsProd)
	setCfgGroupString("server", "host", &ServerHost)
	setCfgGroupInt("server", "port", &ServerPort)
	setCfgGroupString("server", "domain", &ServerDomain)
	setCfgGroupBool("server", "unixSock", &UnixSock)
	setCfgGroupString("server", "unixSockFile", &UnixSockFile)

	// cookie
	setCfgGroupString("cookie", "name", &Cookie)
	setCfgGroupBool("cookie", "expires", &CookieExpires)
	setCfgGroupString("cookie", "expireHash", &CookieExpireHash)
	setCfgGroupInt("cookie", "expireTime", &CookieExpireTime)

	// media
	setCfgGroupString("media", "url", &MediaURL)

	// templates
	setCfgGroupString("templates", "editRoute", &TmplEditRt)
	setCfgGroupBool("templates", "editErr", &TmplEditErr)
	setCfgGroupString("templates", "path", &TmplPath)

	// database
	setCfgGroupString("user", "reader", &UserReader)
	//setCfgGroupString("user", "writer", &UserWriter)
	//setCfgGroupString("db", "mongoURL", &MongoURL)
	//setCfgGroupString("db", "mongoDB", &MongoDB)
	//setCfgGroupString("db", "sqliteFile", &SqliteFile)
}

func getConfig() (p string, c []byte, e error) {
	p = filepath.Dir(os.Args[0])
	f := filepath.Join(p, CONFIG_FILE)

	// if a config file exists within the executables path
	if _, err := os.Stat(f); err == nil {
		c, e = ioutil.ReadFile(f)
		return
	}

	// if a config file exists within the current working dir
	if p, e = os.Getwd(); e == nil {
		f = filepath.Join(p, CONFIG_FILE)
		if _, e = os.Stat(f); e == nil {
			c, e = ioutil.ReadFile(f)
			return
		}
	}

	// no configuration was found
	p = ""
	e = errors.New(fmt.Sprintf("Failed to find a configuration file: %s", CONFIG_FILE))

	return
}

// accessors
func colBool(key string, col map[string]interface{}) (bool, bool) {
	if v, ok := col[key]; ok {
		b, ok := v.(bool)
		return b, ok
	}
	return false, false
}

func colString(key string, col map[string]interface{}) (string, bool) {
	if v, ok := col[key]; ok {
		s, ok := v.(string)
		return s, ok
	}
	return *new(string), false
}

func colInt(key string, col map[string]interface{}) (int, bool) {
	if v, ok := col[key]; ok {
		switch v.(type) {
		case int:
			return v.(int), true
		case float64:
			return int(v.(float64)), true
		}
	}
	return -1, false
}

func colVal(key string, col map[string]interface{}) (interface{}, bool) {
	if v, ok := col[key]; ok {
		return v, true
	}
	return nil, false
}

func Bool(key string) (bool, bool) {
	return colBool(key, cfg)
}

func String(key string) (string, bool) {
	return colString(key, cfg)
}

func Int(key string) (int, bool) {
	return colInt(key, cfg)
}

func Val(key string) (interface{}, bool) {
	return colVal(key, cfg)
}

func GroupBool(group, key string) (v bool, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colBool(key, col)
		}
	}
	return
}

func GroupString(group, key string) (v string, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colString(key, col)
		}
	}
	return
}

func GroupInt(group, key string) (v int, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colInt(key, col)
		}
	}
	return
}

func GroupVal(group, key string) (v interface{}, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colVal(key, col)
		}
	}
	return
}

// root
func setCfgInt(key string, val *int) {
	if v, ok := Int(key); ok {
		*val = v
	}
}
func setCfgBool(key string, val *bool) {
	if v, ok := Bool(key); ok {
		*val = v
	}
}
func setCfgString(key string, val *string) {
	if v, ok := String(key); ok {
		*val = v
	}
}

// group
func setCfgGroupInt(group, key string, val *int) {
	if v, ok := GroupInt(group, key); ok {
		*val = v
	}
}
func setCfgGroupBool(group, key string, val *bool) {
	if v, ok := GroupBool(group, key); ok {
		*val = v
	}
}
func setCfgGroupString(group, key string, val *string) {
	if v, ok := GroupString(group, key); ok {
		*val = v
	}
}
