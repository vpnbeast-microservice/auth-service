package web

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func pingHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	}
}

func authenticateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var request authRequest
		if err := context.ShouldBindJSON(&request); err != nil {
			errors, _ := err.(validator.ValidationErrors)
			e := make(map[string]string)
			for _, err := range errors {
				e[err.Field()] = err.Tag()
				logger.Info("", zap.Any("err", err))
			}
			context.JSON(400, e)
			return
		}

		// TODO: request validation
		// TODO: check if password is correct, return success response else return fail response
		// TODO: generate accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt, update database and return response
		// TODO: update last_login column on db and return response(aspect oriented programing? check https://github.com/gogap/aop)

		result := authSuccessResponse{}
		sqlStatement := fmt.Sprintf("SELECT uuid, id, created_at, updated_at, version, user_name, email, " +
			"last_login, enabled, email_verified, access_token, access_token_expires_at, refresh_token, " +
			"refresh_token_expires_at FROM users WHERE user_name='%s'", request.Username)
		row := db.QueryRow(sqlStatement)
		switch err := row.Scan(&result.Uuid, &result.Id, &result.CreatedAt, &result.UpdatedAt, &result.Version,
			&result.Username, &result.Email, &result.LastLogin, &result.Enabled, &result.EmailVerified, &result.AccessToken,
			&result.AccessTokenExpiresAt, &result.RefreshToken, &result.RefreshTokenExpiresAt); err {
		case sql.ErrNoRows:
			logger.Warn("no rows were returned!", zap.String("user", request.Username))
			failResponse := authFailResponse{
				Tag:          "getUser",
				ErrorMessage: "User not found!",
				Status:       false,
				HttpCode:     404,
				Timestamp:    time.Now(),
			}

			context.JSON(http.StatusBadRequest, failResponse)
			return
		case nil:
			logger.Info("user fetched")
			result.Tag = "getToken"
			context.JSON(http.StatusOK, result)
			return
		default:
			logger.Error("unknown error", zap.String("error", err.Error()))
			context.JSON(http.StatusInternalServerError, authFailResponse{
				Tag:          "getUser",
				ErrorMessage: "unknown error occured at the backend",
				Status:       false,
				HttpCode:     500,
				Timestamp:    time.Now(),
			})
			return
		}
	}
}