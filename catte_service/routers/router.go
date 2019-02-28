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
		"/join-room/{id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.JoinRoom),
		)).Methods("GET")

	router.Handle(
		"/leave-room/{id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.LeaveRoom),
		)).Methods("GET")

	router.Handle(
		"/ws/{id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.PlayerJoin),
		)).Methods("GET")
	return router
}
