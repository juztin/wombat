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
	CryptIter   = 12
	Cookie      = "oatmeal"
	Host        = "127.0.0.1"
	Port        = 9999
	Domain      = "juzt.in"
	IsProd      = false
	//Media     = "/var/web/nginx/juztin/projects"
	ProjectPath = "/var/web/nginx/juztin/projects"
	Pepper	    = "sYNRXbVSc0RwPJXmcKXJK.8IcoiBHMX7Lc0cNrCW3aImXQO9U/t26"
	UnixSocket  = false
	SockFile    = "/tmp/juztin.sock"
	MediaURL    = "media.juzt.in/"
	TmplEditRt  = "^/_dt/$"
	TmplEditErr = false
	TmplPath    = "./templates"

	Backend     = "sqlite"
	MongoHost   = "localhost"
	SqliteFile  = "./db.sqlite"

	cfg map[string]interface{}
)

func setCfgInt(key string, val *int) {
	if v, ok := GetInt(key); ok {
		*val = v
	}
}
func setCfgBool(key string, val *bool) {
	if v, ok := GetBool(key); ok {
		*val = v
	}
}
func setCfgString(key string, val *string) {
	if v, ok := GetString(key); ok {
		*val = v
	}
}

func setFromCfg() {
	setCfgInt("cryptIter", &CryptIter)
	setCfgString("cookie", &Cookie)
	setCfgString("host", &Host)
	setCfgInt("port", &Port)
	setCfgString("domain", &Domain)
	setCfgBool("isProd", &IsProd)
	//setCfgString("media", &Media)
	setCfgString("projectPath", &ProjectPath)
	setCfgString("salt", &Pepper)
	setCfgString("sockFile", &SockFile)
	setCfgString("mediaUrl", &MediaURL)
	setCfgString("templateEditRoute", &TmplEditRt)
	setCfgBool("templateEditErr", &TmplEditErr)
	setCfgString("templatePath", &TmplPath)

	setCfgString("backend", &Backend)
	setCfgString("mongoHost", &MongoHost)
	setCfgString("sqliteFile", &SqliteFile)
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

func Load() {
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

func GetBool(key string) (bool, bool) {
	if v, ok := cfg[key]; ok {
		b, ok := v.(bool)
		return b, ok
	}
	return false, false
}

func GetString(key string) (string, bool) {
	if v, ok := cfg[key]; ok {
		s, ok := v.(string)
		return s, ok
	}
	return *new(string), false
}

func GetInt(key string) (int, bool) {
	if v, ok := cfg[key]; ok {
		switch v.(type) {
		case int:
			return v.(int), true
		case float64:
			return int(v.(float64)), true
		}
	}
	return -1, false
}

func GetVal(key string) (interface{}, bool) {
	if v, ok := cfg[key]; ok {
		return v, true
	}
	return nil, false
}