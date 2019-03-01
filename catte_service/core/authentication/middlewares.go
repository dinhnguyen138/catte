package authentication

import (
	"net/http"
)

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	//
	authBackend := InitJWTAuthenticationBackend()
	authBackend.middleWare.HandlerWithNext(rw, req, next)
}
