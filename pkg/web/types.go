package web

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type authRequest struct {
	Username string `json:"userName" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=3,max=16"`
}

type selectResult struct {
	EncryptedPassword string
	UserName string
}

type authSuccessResponse struct {
	Uuid string `json:"uuid"`
	Id int64 `json:"id"`
	EncryptedPassword string `json:"encryptedPassword,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Version int `json:"version"`
	Username string `json:"username"`
	Email string `json:"email"`
	LastLogin time.Time `json:"lastLogin"`
	Enabled bool `json:"enabled"`
	EmailVerified bool `json:"emailVerified"`
	Tag string `json:"tag"`
	AccessToken string `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
	RefreshToken string `json:"refreshToken"`
	RefreshTokenExpiresAt string `json:"refreshTokenExpiresAt"`
}

// when user not found or auth failed
type authFailResponse struct {
	Tag string `json:"tag"`
	ErrorMessage string `json:"errorMessage"`
	Status bool `json:"status"`
	HttpCode int `json:"httpCode"`
	Timestamp time.Time `json:"timestamp"`
}

type validationErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	HttpCode int `json:"httpCode"`
	Tag string `json:"tag"`
	Status bool `json:"status"`
	ErrorMessage []string `json:"errorMessage"`
}

type encryptRequest struct {
	PlainText string `json:"plainText"`
	EncryptedText string `json:"encryptedText"`
}

func (req encryptRequest) encrypt(plainText, encrypted string) (encryptResponse, error) {
	postBody, err := json.Marshal(encryptRequest{
		PlainText:     plainText,
		EncryptedText: encrypted,
	})
	if err != nil {
		logger.Error("an error occured while marshalling request", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(encryptionServiceUrl, "application/json", responseBody)
	if err != nil {
		logger.Error("an error occured while making remote request", zap.String("error", err.Error()))
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
		logger.Error("an error occured while reading response", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	var response encryptResponse
	responseString := string(body)
	err = json.Unmarshal([]byte(responseString), &response)
	if err != nil {
		logger.Error("an error occured while unmarshaling response", zap.String("error", err.Error()))
		return encryptResponse{}, err
	}

	return response, nil
}

type encryptResponse struct {
	Tag string `json:"tag"`
	Status bool `json:"status"`
}
