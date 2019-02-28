package authentication

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token, err := GetToken(req)
	if err == nil && token.Valid {
		c := token.Claims.(jwt.MapClaims)
		fmt.Println(c["exp"])
		fmt.Println(c["sub"])
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func GetToken(req *http.Request) (token *jwt.Token, err error) {
	authBackend := InitJWTAuthenticationBackend()

	return request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
}
