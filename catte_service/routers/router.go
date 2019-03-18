package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/dinhnguyen138/catte/catte_service/controllers"
	"github.com/dinhnguyen138/catte/catte_service/core/authentication"
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetAuthenticationRoutes(router)
	return router
}

func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc(
		"/login",
		controllers.Login,
	).Methods("POST")

	router.HandleFunc(
		"/register",
		controllers.Register,
	).Methods("POST")

	router.HandleFunc(
		"/login3rd",
		controllers.Login3rd,
	).Methods("POST")

	router.HandleFunc(
		"/register-host",
		controllers.RegisterHost,
	).Methods("POST")

	router.Handle(
		"/refresh-token",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("GET")

	router.Handle(
		"/checkin",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.CheckIn),
		)).Methods("GET")

	router.Handle(
		"/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("GET")

	router.Handle(
		"/get-rooms",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetRooms),
		)).Methods("GET")

	router.Handle(
		"/create-room",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.CreateRoom),
		)).Methods("GET")

	router.Handle(
		"/get-info",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetInfo),
		)).Methods("GET")

	return router
}
