package jwt

import "github.com/dgrijalva/jwt-go"

type VpnbeastClaim struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}
