package jwt

import (
	"auth-service/pkg/config"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var (
	signKey *rsa.PrivateKey
	err     error

	privateKey, issuer string
)

func init() {
	issuer = config.GetStringEnv("ISSUER", "info@thevpnbeast.com")
	privateKey = strings.Replace(config.GetStringEnv("PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----\\nMIICdwIBADANBgkqhkiG"+
		"9w0BAQEFAASCAmEwggJdAgEAAoGBANnmLifeLBsiXe/J\\n8O3ophHHaCfJ+EdAUYn7vArJTUtankCD3I8O3n+QM0KNsXzXd+eN6VmNm3bjLp"+
		"Hq\\nVjI/jCr2m1EqXgvRQP74/wOU1sHN3zSRQbcPR0dfJiDfTRmfh/LVrKgcU0kQ4yrG\\nlc0KGB2uslzrKLJCmQ4G0WeM3tKNAgMBAAECgY"+
		"EApsep+FXzSGmLoOfegxqZUe5g\\n6GOMp2yxfH2ztkXR5aVcj2DeRplI8DZ9Jamyei2p1xAl1aevoNXOZV0J0LgXHbm+\\nP6MGU7d+IYD2hI"+
		"CWPfD4pqJafkYc7Q94eQaIiShlYEOoEiLDt09m2V3J/VWxEWw0\\nGTzT1T6zDuwD5epXY/kCQQDv+Xeq+SU5+avfysvm/8bITu/WBRXKxQ7V"+
		"2dg9rJIF\\nrAZSTPUIqdKm2F+o8DIX4sSMouFMgo81Ad4S8D13iQCPAkEA6HNRwmcfQCHzsuuT\\n1407mEPFAcgIckU6e9ubXRRepWPjE6MJ"+
		"IyrDeIkCJfgFPiK8OcNvFLUCD8NaySD1\\nQuHRIwJAJtIqo8QOW6SiQ1/hQItcMwdiETNdZSIf1kSZkNCcBsLfeuzsLuyaIVeb\nkg7Za7fJ"+
		"qB6pZ+EvHZohvNqUdwP4zQJBALf6piiG9C4PcVIYsOA3cYa3hNM/HqhK\\n8NodW9+VAsBGyfC95rqF20aosiGZJ5UhavcRHvc1uNb/GPj99A"+
		"EmuB8CQCy9M89/\\nWGs0V60TrOWn2cmNlvexvxJgtWIjzdtp5rBj/E7Dmfx9nE6sG+uJqob389HYb0fF\\nj8MrN6RCirNhupc=\\n-----END"+
		" PRIVATE KEY-----"), "\\n", "\n", -1)

	signBytes := []byte(privateKey)
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
}

func GenerateToken(username string, expiresAtInMinutes int32) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Duration(expiresAtInMinutes) * time.Minute).Unix(),
		Issuer:    issuer,
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
