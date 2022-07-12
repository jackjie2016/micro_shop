package models

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID uint
	NickName string
	AuthorityID uint
	jwt.StandardClaims
}