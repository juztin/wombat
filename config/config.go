package config

import (
	_ "bitbucket.org/juztin/virginia/backends/mongo"
)

var (
	//IsProd    = true
	IsProd    = false
	CryptIter = 12
	DebugURL  = "127.0.0.1"
	DebugPort = 9999
	Domain    = "juzt.in"
	Pepper	  = "sYNRXbVSc0RwPJXmcKXJK.8IcoiBHMX7Lc0cNrCW3aImXQO9U/t26"
	SockFile  = "/tmp/juztin.sock"
	StaticURL = "media.juzt.in/"

	ImageRoot = "/var/web/nginx/juztin/projects"

	Backend   = "sqlite"
	MongoHost = "localhost"
	SqliteFile= "./db.sqlite"
)
