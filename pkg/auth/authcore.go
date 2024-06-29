package auth

import (
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

const Secret = "secretSUKA" // Replace <jwt-secret> with your secret key that is private to you.

func NewAuth() (au *jwtauth.JWTAuth) {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
	return tokenAuth
}

func MakeToken(ID string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"userID": ID})
	return tokenString
}
