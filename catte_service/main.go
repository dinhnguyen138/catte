package main

import (
	"net/http"

	"./db"
	"./routers"
	"./settings"
	"github.com/codegangsta/negroni"
)

func main() {
	settings.Init()
	db.InitDB()
	defer db.CloseDB()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(":8080", n)
}
