package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/dinhnguyen138/catte/catte_service/db"
	"github.com/dinhnguyen138/catte/catte_service/routers"
	"github.com/dinhnguyen138/catte/catte_service/settings"
	"github.com/kataras/golog"
	"github.com/natefinch/lumberjack"
)

func main() {
	golog.SetOutput(&lumberjack.Logger{
		Filename:   "log/service/daily.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})

	golog.Info(os.Getenv("ENV"))
	settings.Init()
	db.InitDB()
	defer db.CloseDB()
	router := routers.InitRoutes()
	negroniLog := &negroni.Logger{ALogger: golog.Default}
	negroniLog.SetFormat(negroni.LoggerDefaultFormat)
	n := negroni.New(negroni.NewRecovery(), negroniLog, negroni.NewStatic(http.Dir("public")))
	n.UseHandler(router)

	if os.Getenv("ENV") == "prod" {
		err := http.ListenAndServeTLS(":443", settings.Get().ServerCertPath, settings.Get().ServerKeyPath, n)
		if err != nil {
			golog.Fatal("ListenAndServe: ", err)
		}
	} else {
		err := http.ListenAndServe(":8080", n)
		if err != nil {
			golog.Fatal("ListenAndServe: ", err)
		}
	}
}
