package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
	"rest/consts"
	"rest/handlers"
)

func main() {

	mux := http.NewServeMux()
	n := negroni.Classic()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to REST api Demo.")
	})

	mux.HandleFunc("/users", handler.UserHandler)

	n.UseHandler(mux)
	n.Run(":" + consts.PORT)

}
