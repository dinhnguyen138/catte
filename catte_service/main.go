package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/dinhnguyen138/catte/catte_service/db"
	"github.com/dinhnguyen138/catte/catte_service/routers"
	"github.com/dinhnguyen138/catte/catte_service/settings"
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
