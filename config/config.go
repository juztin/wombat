package config

//import (
	//_ "bitbucket.org/juztin/virginia/backends/mongo"
	//_ "bitbucket.org/juztin/virginia/backends/sqlite"
//)

var (
	IsProd    = false
	CryptIter = 12
	Cookie    = "oatmeal"
	DebugURL  = "127.0.0.1"
	DebugPort = 9999
	Domain    = "juzt.in"
	ImageRoot = "/var/web/nginx/juztin/projects"
	Pepper	  = "sYNRXbVSc0RwPJXmcKXJK.8IcoiBHMX7Lc0cNrCW3aImXQO9U/t26"
	SockFile  = "/tmp/juztin.sock"
	MediaURL  = "media.juzt.in/"
	TmplEditRt= "^/_dt/$"

	Backend   = "sqlite"
	MongoHost = "localhost"
	SqliteFile= "./db.sqlite"
)
