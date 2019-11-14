package user

import "github.com/dgrijalva/jwt-go"

/*
JWT claims struct
*/
type Token struct {
	UserName string
	Email string
	ID       string
	jwt.StandardClaims
}
type ContextData struct {
	UserName   string
	SessionKey string
}
type ContextKey struct{
	Key int
}
