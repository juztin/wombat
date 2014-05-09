// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wraps

import (
	"fmt"
	"net/http"

	"code.minty.io/dingo/views"
	"code.minty.io/wombat"
	"code.minty.io/wombat/config"
	"code.minty.io/wombat/template/data"
)

func RequireAdmin(h wombat.Handler) wombat.Handler {
	return func(ctx wombat.Context) {
		if ctx.User.IsAnonymous() {
			ctx.Redirect(config.SigninURL)
		} else if !ctx.User.IsAdmin() {
			if v := views.Get(fmt.Sprintf("%s%s", wombat.ERR_TMPL, "401.html")); v != nil {
				v.Execute(ctx.Context, data.New(ctx))
			} else {
				ctx.HttpError(http.StatusUnauthorized)
			}
		} else {
			h(ctx)
		}
	}
}

func RequireAuth(h wombat.Handler) wombat.Handler {
	return func(ctx wombat.Context) {
		if ctx.User.IsAnonymous() {
			ctx.Redirect(config.SigninURL)
		} else {
			h(ctx)
		}
	}
}
