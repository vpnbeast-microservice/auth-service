package web

import (
	"auth-service/pkg/jwt"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func pingHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"code":      http.StatusOK,
			"message":   "pong",
			"timestamp": time.Now(),
		})
	}
}

func authenticateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var authReq authRequest
		_, errSlice := isValidRequest(context, &authReq)
		if len(errSlice) != 0 {
			validationResponse(context, errSlice)
			context.Abort()
			return
		}

		selectRes := selectResult{}
		sqlStatement := fmt.Sprintf(SqlSelectUsernamePass, authReq.Username)
		row := db.QueryRow(sqlStatement)
		switch err := row.Scan(&selectRes.EncryptedPassword, &selectRes.UserName); err {
		case sql.ErrNoRows:
			logger.Warn("no rows were returned!", zap.String("user", authReq.Username))
			errorResponse(context, http.StatusBadRequest, ErrUserNotFound)
			context.Abort()
			return
		case nil:
			encryptReq := encryptRequest{
				PlainText:     authReq.Password,
				EncryptedText: selectRes.EncryptedPassword,
			}

			encryptRes, err := encryptReq.encrypt(authReq.Password, selectRes.EncryptedPassword)
			if err != nil {
				logger.Error("an error occurred while making encryption request", zap.String("error", err.Error()))
				errorResponse(context, http.StatusInternalServerError, ErrUnknown)
				context.Abort()
				return
			}

			if encryptRes.Status {
				accessToken, err := jwt.GenerateToken(authReq.Username, int32(accessTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating access token",
						zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, ErrUnknown)
					context.Abort()
					return
				}

				refreshToken, err := jwt.GenerateToken(authReq.Username, int32(refreshTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating refresh token",
						zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, ErrUnknown)
					context.Abort()
					return
				}

				lastLogin := time.Now().Format(time.RFC3339)
				accessTokenExpiresAt := time.Now().Add(time.Duration(accessTokenValidInMinutes) * time.Minute).
					Format(time.RFC3339)
				refreshTokenExpiresAt := time.Now().Add(time.Duration(refreshTokenValidInMinutes) * time.Minute).
					Format(time.RFC3339)
				updateStatement := fmt.Sprintf(SqlUpdateUser, lastLogin, accessToken, accessTokenExpiresAt,
					refreshToken, refreshTokenExpiresAt, authReq.Username)
				_, err = db.Exec(updateStatement)
				if err != nil {
					logger.Warn("an error  occurred while updating db", zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, ErrUnknown)
					context.Abort()
					return
				}

				authRes := authSuccessResponse{}
				sqlStatement := fmt.Sprintf(SqlSelectUserAll, authReq.Username)
				row := db.QueryRow(sqlStatement)
				switch err := row.Scan(&authRes.Uuid, &authRes.Id, &authRes.EncryptedPassword, &authRes.CreatedAt, &authRes.UpdatedAt,
					&authRes.Version, &authRes.Username, &authRes.Email, &authRes.LastLogin, &authRes.Enabled, &authRes.EmailVerified,
					&authRes.AccessToken, &authRes.AccessTokenExpiresAt, &authRes.RefreshToken, &authRes.RefreshTokenExpiresAt); err {
				case nil:
					authRes.Tag = "authUser"
					authRes.EncryptedPassword = ""
					context.JSON(http.StatusOK, authRes)
					context.Abort()
					return
				default:

				}
			} else {
				logger.Error("password validation failed")
				errorResponse(context, http.StatusBadRequest, ErrInvalidPass)
				context.Abort()
				return
			}

		default:
			logger.Error("unknown error", zap.String("error", err.Error()))
			errorResponse(context, http.StatusInternalServerError, ErrUnknown)
			context.Abort()
			return
		}
	}
}
