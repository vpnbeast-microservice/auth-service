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
			"code" : http.StatusOK,
			"message": "pong",
			"timestamp": time.Now(),
		})
	}
}

func authenticateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var authReq authRequest

		_, errSlice := isValidRequest(context, &authReq)
		if len(errSlice) != 0 {
			context.JSON(http.StatusBadRequest, validationErrorResponse{
				Tag:          "authUser",
				ErrorMessage: errSlice,
				Status:       false,
				HttpCode:     400,
				Timestamp:    time.Now(),
			})
			context.Abort()
			return
		}

		selectRes := selectResult{}
		sqlStatement := fmt.Sprintf("SELECT encrypted_password, user_name FROM users WHERE user_name='%s'",
			authReq.Username)
		row := db.QueryRow(sqlStatement)
		switch err := row.Scan(&selectRes.EncryptedPassword, &selectRes.UserName); err {
		case sql.ErrNoRows:
			logger.Warn("no rows were returned!", zap.String("user", authReq.Username))
			context.JSON(http.StatusBadRequest, authFailResponse{
				Tag:          "authUser",
				ErrorMessage: "User not found!",
				Status:       false,
				HttpCode:     404,
				Timestamp:    time.Now(),
			})
			context.Abort()
			return
		case nil:
			encryptReq := encryptRequest{
				PlainText:     authReq.Password,
				EncryptedText: selectRes.EncryptedPassword,
			}

			encryptRes, err := encryptReq.encrypt(authReq.Password, selectRes.EncryptedPassword)
			if err != nil {
				logger.Error("an error occured while making encryption request", zap.String("error", err.Error()))
				context.JSON(http.StatusInternalServerError, authFailResponse{
					Tag:          "authUser",
					ErrorMessage: "unknown error occured at the backend",
					Status:       false,
					HttpCode:     500,
					Timestamp:    time.Now(),
				})
				context.Abort()
				return
			}

			if encryptRes.Status {
				accessToken, err := jwt.GenerateToken(authReq.Username, int32(accessTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occured generating access token",
						zap.String("error", err.Error()))
					context.JSON(http.StatusInternalServerError, authFailResponse{
						Tag:          "authUser",
						ErrorMessage: "Unknown error occured at the backend!",
						Status:       false,
						HttpCode:     500,
						Timestamp:    time.Now(),
					})
					context.Abort()
					return
				}

				refreshToken, err := jwt.GenerateToken(authReq.Username, int32(refreshTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occured generating refresh token",
						zap.String("error", err.Error()))
					context.JSON(http.StatusInternalServerError, authFailResponse{
						Tag:          "authUser",
						ErrorMessage: "Unknown error occured at the backend!",
						Status:       false,
						HttpCode:     500,
						Timestamp:    time.Now(),
					})
					context.Abort()
					return
				}

				lastLogin := time.Now().Format(time.RFC3339)
				accessTokenExpiresAt := time.Now().Add(time.Duration(accessTokenValidInMinutes) * time.Minute).
					Format(time.RFC3339)
				refreshTokenExpiresAt := time.Now().Add(time.Duration(refreshTokenValidInMinutes) * time.Minute).
					Format(time.RFC3339)
				updateStatement := fmt.Sprintf("UPDATE users SET version = version + 1, last_login='%v', " +
					"access_token='%v', access_token_expires_at='%v', refresh_token='%s', refresh_token_expires_at='%v' " +
					"WHERE user_name='%s'", lastLogin, accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt,
					authReq.Username)
				_, err = db.Exec(updateStatement)
				if err != nil {
					logger.Warn("an error  occured while updating db", zap.String("error", err.Error()))
					context.JSON(http.StatusBadRequest, authFailResponse{
						Tag:          "authUser",
						ErrorMessage: "Unknown error occured at the backend!",
						Status:       false,
						HttpCode:     500,
						Timestamp:    time.Now(),
					})
					context.Abort()
					return
				}

				authRes := authSuccessResponse{}
				sqlStatement := fmt.Sprintf("SELECT uuid, id, encrypted_password, created_at, updated_at, version, " +
					"user_name, email, last_login, enabled, email_verified, access_token, access_token_expires_at, " +
					"refresh_token, refresh_token_expires_at FROM users WHERE user_name='%s'", authReq.Username)
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
				context.JSON(http.StatusBadRequest, authFailResponse{
					Tag:          "authUser",
					ErrorMessage: "Invalid password!",
					Status:       false,
					HttpCode:     400,
					Timestamp:    time.Now(),
				})
				context.Abort()
				return
			}

		default:
			logger.Error("unknown error", zap.String("error", err.Error()))
			context.JSON(http.StatusInternalServerError, authFailResponse{
				Tag:          "authUser",
				ErrorMessage: "unknown error occured at the backend",
				Status:       false,
				HttpCode:     500,
				Timestamp:    time.Now(),
			})
			context.Abort()
			return
		}
	}
}