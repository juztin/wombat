package wraps

import (
	"fmt"
	"net/http"

	"bitbucket.org/juztin/dingo/views"
	"bitbucket.org/juztin/wombat"
	"bitbucket.org/juztin/wombat/config"
	"bitbucket.org/juztin/wombat/template/data"
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
