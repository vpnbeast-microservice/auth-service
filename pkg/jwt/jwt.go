package jwt

import (
	"auth-service/pkg/options"
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey
	privateKeyBytes []byte
	publicKeyBytes []byte
	err     error
	opts    *options.AuthServiceOptions
)

func init() {
	opts = options.GetAuthServiceOptions()
	privateKeyBytes = []byte(opts.PrivateKey)
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	publicKeyBytes = []byte(opts.PublicKey)
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
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

	tokenString, err := t.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates JWT token by checking if issuer is registered user, expiration time not passed etc
func ValidateToken(signedToken string) (string, error, int) {
	// TODO: refactor
	token, err := jwt.ParseWithClaims(
		signedToken,
		&VpnbeastClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

	if err != nil {
		log.Println(err.Error())
		return "", err, 500
	}

	claims, ok := token.Claims.(*VpnbeastClaim)
	if !ok {
		err = errors.New("could not parse claims")
		return "", err, 500
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("jwt is already expired")
		return "", err, 401
	}

	return claims.Subject, nil, 200
}
