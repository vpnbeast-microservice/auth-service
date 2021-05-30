package jwt

import (
	"auth-service/pkg/options"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	signKey *rsa.PrivateKey
	err     error
	opts    *options.AuthServiceOptions
)

func init() {
	opts = options.GetAuthServiceOptions()
	signBytes := []byte(opts.PrivateKey)
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
}

// GenerateToken generates JWT token with username and expiresAtInMinutes in RS256 signing method
func GenerateToken(username string, roles []string, expiresAtInMinutes int32) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &VpnbeastClaim{
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiresAtInMinutes) * time.Minute).Unix(),
			Issuer:    opts.Issuer,
			Subject:   username,
			IssuedAt:  time.Now().Unix(),
		},
	}

	tokenString, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates JWT token by checking if issuer is registered user, expiration time not passed etc
func ValidateToken(token string) (bool, error) {
	// TODO: implement
	// check https://betterprogramming.pub/hands-on-with-jwt-in-golang-8c986d1bb4c0
	return false, nil
}
