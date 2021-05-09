package web

import (
	"auth-service/pkg/jwt"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
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
		var request authRequest

		_, errSlice := isValidRequest(context, &request)
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

		result := selectResult{}
		sqlStatement := fmt.Sprintf("SELECT encrypted_password, user_name FROM users WHERE user_name='%s'",
			request.Username)
		row := db.QueryRow(sqlStatement)
		switch err := row.Scan(&result.EncryptedPassword, &result.UserName); err {
		case sql.ErrNoRows:
			logger.Warn("no rows were returned!", zap.String("user", request.Username))
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
			postBody, err := json.Marshal(encryptRequest{
				PlainText:     request.Password,
				EncryptedText: result.EncryptedPassword,
			})
			if err != nil {
				logger.Error("an error occured while marshalling request body", zap.String("error", err.Error()))
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

			responseBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(encryptionServiceUrl, "application/json", responseBody)
			if err != nil {
				logger.Error("an error occured while making request to encryption-service",
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

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					panic(err)
				}
			}()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error("an error occured while reading response body", zap.String("error", err.Error()))
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

			var encryptResponse encryptResponse
			responseString := string(body)
			err = json.Unmarshal([]byte(responseString), &encryptResponse)
			if err != nil {
				logger.Error("an error occured while unmarshalling response to struct",
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

			if encryptResponse.Status {
				accessToken, err := jwt.GenerateToken(request.Username, int32(accessTokenValidInMinutes))
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

				refreshToken, err := jwt.GenerateToken(request.Username, int32(refreshTokenValidInMinutes))
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
					request.Username)
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

				result := authSuccessResponse{}
				sqlStatement := fmt.Sprintf("SELECT uuid, id, encrypted_password, created_at, updated_at, version, " +
					"user_name, email, last_login, enabled, email_verified, access_token, access_token_expires_at, " +
					"refresh_token, refresh_token_expires_at FROM users WHERE user_name='%s'", request.Username)
				row := db.QueryRow(sqlStatement)
				switch err := row.Scan(&result.Uuid, &result.Id, &result.EncryptedPassword, &result.CreatedAt, &result.UpdatedAt,
					&result.Version, &result.Username, &result.Email, &result.LastLogin, &result.Enabled, &result.EmailVerified,
					&result.AccessToken, &result.AccessTokenExpiresAt, &result.RefreshToken, &result.RefreshTokenExpiresAt); err {
				case nil:
					result.Tag = "authUser"
					result.EncryptedPassword = ""
					context.JSON(http.StatusOK, result)
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