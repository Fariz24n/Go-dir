package handlers

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("ashdjqy9283409bsdklkg8hda01")

type JWTClaim struct {
	Name string
	jwt.RegisteredClaims
}
