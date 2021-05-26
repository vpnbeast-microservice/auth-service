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
func GenerateToken(username string, expiresAtInMinutes int32) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Duration(expiresAtInMinutes) * time.Minute).Unix(),
		Issuer:    opts.Issuer,
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
	}

	tokenString, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
