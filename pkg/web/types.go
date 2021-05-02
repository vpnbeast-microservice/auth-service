package web

import (
	"database/sql"
	"time"
)

type authRequest struct {
	Username string `json:"userName" binding:"required,min=3,max=16"`
	Password string `json:"password" binding:"required,min=3,max=16"`
}

type authSuccessResponse struct {
	Uuid string `json:"uuid"`
	Id int64 `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Version int `json:"version"`
	Username string `json:"username"`
	Email string `json:"email"`
	LastLogin sql.NullTime `json:"lastLogin"`
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

/*type validationErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	HttpCode int `json:"httpCode"`
	Tag string `json:"tag"`
	Status bool `json:"status"`
	ErrorMessage []string `json:"errorMessage"`
}*/

/*type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}*/