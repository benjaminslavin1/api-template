package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.logRequest, app.recoverPanic)
	r.Get("/", app.dummyHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) { app.notFound(w) })
	return r
}
