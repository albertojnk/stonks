package domains

import "github.com/dgrijalva/jwt-go"

type Auth struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
}

type JWTClaims struct {
	LoggedAuth Auth `json:"logged_auth"`
	IsAdmin    bool `json:"is_admin"`
	jwt.StandardClaims
}
