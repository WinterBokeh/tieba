package Model

import "github.com/dgrijalva/jwt-go"

type MyCustomClaims struct {
	Name string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}
