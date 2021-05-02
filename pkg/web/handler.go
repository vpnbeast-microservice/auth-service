package web

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func pingHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
		return
	}
}

func authenticateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var request authRequest
		err := decodeJSONBody(context.Writer, context.Request, &request)
		if err != nil {
			logger.Error("an error occured while decoding json body", zap.String("error", err.Error()))
			var mr *malformedRequest
			if errors.As(err, &mr) {
				context.JSON(mr.status, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			} else {
				context.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			}

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
		default:
			panic(err)
		}
	}
}