package routers

import (
	"../controllers"
	"../core/authentication"
	"github.com/codegangsta/negroni"
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

	router.Handle(
		"/refresh-token-auth",
		negroni.New(
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("GET")

	router.Handle(
		"/logout",
		negroni.New(
			negroni.HandlerFunc(
				authentication.RequireTokenAuthentication,
			),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("GET")

	router.Handle(
		"/get-rooms",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetRooms),
		)).Methods("GET")

	router.Handle(
		"/get-info",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.GetInfo),
		)).Methods("GET")

	return router
}
