package secure_tools

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func CSRF(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			Path:     "/",
			HttpOnly: true,
			Secure:   utils.GetBoolValue(utils.Config.CSRF_SECURE), // TODO
		})

		next.ServeHTTP(w, r)
	}
}
