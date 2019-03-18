package main

import (
	"log"
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

	if os.Getenv("ENV") == "prod" {
		err := http.ListenAndServeTLS(":443", settings.Get().ServerCertPath, settings.Get().ServerKeyPath, n)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
	else {
		err := http.ListenAndServe(":8080", n)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
