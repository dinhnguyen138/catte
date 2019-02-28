package main

import (
	"net/http"

	"./routers"
	"./settings"
	"github.com/codegangsta/negroni"
)

func main() {
	settings.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(":8080", n)
}
