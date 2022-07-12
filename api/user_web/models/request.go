package models

import (
	"github.com/dgrijalva/jwt-go"
)

//jwt 生成 token 用
type CustomClaims struct {
	ID uint
	NickName string
	AuthorityID uint
	jwt.StandardClaims
}