package backends

import (
	"fmt"
)

var backends = make(map[string]Backend)

func Register(name string, backend Backend) {
	/*if backend == nil {
		panic("backends: Register backend is nil")
	}*/
	if _, dup := backends[name]; dup {
		panic("backends: Register called twice for driver " + name)
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
