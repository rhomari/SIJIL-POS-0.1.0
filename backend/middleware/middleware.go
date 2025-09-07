package middleware

import (
	utils "SIJIL-POS/utils"
	"log"
	"net/http"
)

type Middleware struct {
}

func (mw *Middleware) Serve(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(utils.Green + "[middleware] activated" + utils.Reset)

		next.ServeHTTP(w, r)
	})
}
