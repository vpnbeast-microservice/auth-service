package web

import "time"

type authenticationRequest struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type authenticationResponse struct {
	Uuid string `json:"uuid"`
	Id int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}