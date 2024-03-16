package model

import "github.com/golang-jwt/jwt"

// UserClaims is custom wrapper for jwt claims.
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
