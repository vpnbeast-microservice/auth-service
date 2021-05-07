package jwt

import (
	"auth-service/pkg/logging"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"time"
)

const (
	PrivKey                    = "app_rsa" // openssl genrsa -out app.rsa keysize
	Issuer                     = "info@thevpnbeast.com"
	AccessTokenValidInMinutes  = 60
	RefreshTokenValidInMinutes = 600
)

var (
	signKey   *rsa.PrivateKey
	logger *zap.Logger
)


// read the key files before starting http handlers
func init() {
	logger = logging.GetLogger()

	signBytes, err := ioutil.ReadFile(PrivKey)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateAccessToken(username string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * AccessTokenValidInMinutes).Unix(),
		Issuer:    Issuer,
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
	}

	tokenString, err := t.SignedString(signKey)
	if err != nil {
		logger.Warn("", zap.String("error", err.Error()))
		return "", err
	} else {
		return tokenString, nil
	}
}

func GenerateRefreshToken(username string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * RefreshTokenValidInMinutes).Unix(),
		Issuer:    Issuer,
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
	}

	tokenString, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}