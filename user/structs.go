package user

import "github.com/dgrijalva/jwt-go"

/*
JWT claims struct
*/
type Token struct {
	UserName string
	ID       string
	jwt.StandardClaims
}
type ContextData struct {
	UserName string
}
