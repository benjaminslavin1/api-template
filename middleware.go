package main

import (
	"fmt"
	"net/http"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		// prints out some metadata around each request
		app.log.Info().Msgf("%s %s://%s%s: %s - %s %s", r.Method, scheme, r.Host, r.RequestURI, r.RemoteAddr, r.Proto, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// when panic occures and bubbles up, deferred statements are run
		// recover() halts the panicing process and returns an error
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
