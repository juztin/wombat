// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func Update(name string, backend Backend) {
	if _, exists := backends[name]; !exists {
		log.Fatal("backends: No driver found to update ", name)
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
