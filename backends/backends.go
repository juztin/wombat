package backends

import (
	"fmt"
	"log"
)

var backends = make(map[string]Backend)

func Register(name string, backend Backend) {
	if _, dup := backends[name]; dup {
		log.Fatal("backends: Register called twice for driver ", name)
	}
	backends[name] = backend
}

func Open(name string) (Backend, error) {
	b, ok := backends[name]
	if !ok {
		return *new(Backend), fmt.Errorf("backends: unknonwn backend %q (forgotten import?)", name)
	}
	return b, nil
}
