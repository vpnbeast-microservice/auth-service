package web

import (
	"auth-service/internal/model"
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

// authRequest represents the incoming auth request to the auth-service
type authRequest struct {
	Username string `json:"userName" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=3,max=16"`
}

// authRequest represents the response of the auth-service to the auth request
type authSuccessResponse struct {
	Uuid                       string        `json:"uuid"`
	Id                         uint          `json:"id"`
	CreatedAt                  string        `json:"createdAt"`
	UpdatedAt                  string        `json:"updatedAt"`
	Version                    uint          `json:"version"`
	Username                   string        `json:"username"`
	Email                      string        `json:"email"`
	LastLogin                  string        `json:"lastLogin"`
	Enabled                    bool          `json:"enabled"`
	EmailVerified              bool          `json:"emailVerified"`
	AccessToken                string        `json:"accessToken"`
	AccessTokenExpiresAt       string        `json:"accessTokenExpiresAt"`
	RefreshToken               string        `json:"refreshToken"`
	RefreshTokenExpiresAt      string        `json:"refreshTokenExpiresAt"`
	VerificationCodeCreatedAt  string        `json:"verificationCodeCreatedAt"`
	VerificationCodeVerifiedAt string        `json:"verificationCodeVerifiedAt"`
	Roles                      []*model.Role `json:"roles"`
}

type authFailResponse struct {
	ErrorMessage string    `json:"errorMessage"`
	Status       bool      `json:"status"`
	HttpCode     int       `json:"httpCode"`
	Timestamp    time.Time `json:"timestamp"`
}

type validationErrorResponse struct {
	Timestamp    time.Time `json:"timestamp"`
	HttpCode     int       `json:"httpCode"`
	Status       bool      `json:"status"`
	ErrorMessage []string  `json:"errorMessage"`
}

// validateRequest represents the request of the encryption-service validate request
type validateRequest struct {
	Token string `json:"token"`
}

// validateResponse represents the request of the encryption-service validate response
type validateResponse struct {
	Status       bool     `json:"status"`
	Username     string   `json:"username,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	ErrorMessage string   `json:"errorMessage"`
	HttpCode     int      `json:"httpCode"`
	Timestamp    string   `json:"timestamp"`
}

// encryptRequest represents the request of the encryption-service encrypt request
type encryptRequest struct {
	PlainText     string `json:"plainText"`
	EncryptedText string `json:"encryptedText"`
}

// encrypt method gets the plainText and makes request to encyrption-service to encrypt plainText
func (req encryptRequest) encrypt(plainText, encrypted string) (encryptResponse, error) {
	postBody, err := json.Marshal(encryptRequest{
		PlainText:     plainText,
		EncryptedText: encrypted,
	})
	if err != nil {
		logger.Error("an error occurred while marshalling request", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(opts.EncryptionServiceUrl, "application/json", responseBody)
	if err != nil {
		logger.Error("an error occurred while making remote request", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("an error occurred while reading response", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	var response encryptResponse
	responseString := string(body)
	err = json.Unmarshal([]byte(responseString), &response)
	if err != nil {
		logger.Error("an error occurred while unmarshaling response", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	return response, nil
}

// encryptResponse represents the response of the encryption-service request
type encryptResponse struct {
	Status bool `json:"status"`
}
