package main

import (
	"fmt"
	"net/http"
)

func (app *application) dummyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Mochi!")
}
