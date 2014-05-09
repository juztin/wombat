// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wombat

import (
	"reflect"
	"regexp"

	"code.minty.io/dingo"
)

type rRoute struct {
	path    string
	expr    *regexp.Regexp
	handler reflect.Value
}

func NewRRoute(re string, handler interface{}) dingo.Route {
	r := new(rRoute)
	r.path = re
	r.expr = regexp.MustCompile(re)

	if fn, ok := handler.(reflect.Value); ok {
		r.handler = fn
	} else {
		r.handler = reflect.ValueOf(handler)
	}

	return *r
}

func (r rRoute) Path() string {
	return r.path
}

func (r rRoute) IsCanonical() bool {
	return true
}

func (r rRoute) Matches(url string) bool {
	return r.expr.MatchString(url)
}

func (r rRoute) Execute(ctx dingo.Context) {
	Wrap(func(c Context) {
		// TODO it would be nice if we could detect numbers and cast them as such prior to invoking the func
		args := []reflect.Value{reflect.ValueOf(c)}
		matches := r.expr.FindStringSubmatch(c.URL.Path)
		for _, a := range matches[1:] {
			args = append(args, reflect.ValueOf(a))
		}
		r.handler.Call(args)
	})(ctx)
}
