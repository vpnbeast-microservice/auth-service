package web

import (
	"auth-service/pkg/jwt"
	"auth-service/pkg/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

		var user types.User
		switch err := db.Where("user_name = ?", authReq.Username).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn("no rows were returned!", zap.String("user", authReq.Username))
			errorResponse(context, http.StatusBadRequest, errUserNotFound)
			context.Abort()
			return
		case nil:
			logger.Info("", zap.Any("user", user))
			encryptReq := encryptRequest{
				PlainText:     authReq.Password,
				EncryptedText: user.EncryptedPassword,
			}

			encryptRes, err := encryptReq.encrypt(authReq.Password, user.EncryptedPassword)
			if err != nil {
				logger.Error("an error occurred while making encryption request", zap.String("error", err.Error()))
				errorResponse(context, http.StatusInternalServerError, errUnknown)
				context.Abort()
				return
			}

			if encryptRes.Status {
				accessToken, err := jwt.GenerateToken(authReq.Username, int32(opts.AccessTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating access token",
						zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}

				refreshToken, err := jwt.GenerateToken(authReq.Username, int32(opts.RefreshTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating refresh token",
						zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}

				//lastLoginTime, _ := time.Parse(time.RFC3339, time.Now().String())
				//user.LastLogin = lastLoginTime
				//now := time.Now().Format(time.RFC3339)
				//user.LastLogin = now
				//user.UpdatedAt = now
				user.AccessToken = accessToken
				// accessTokenExpiresAtTime, _ := time.Parse(time.RFC3339, time.Now().Add(time.Duration(opts.AccessTokenValidInMinutes) * time.Minute).String())
				//user.AccessTokenExpiresAt = time.Now().Add(time.Duration(opts.AccessTokenValidInMinutes) * time.Minute).Format(time.RFC3339)
				user.RefreshToken = refreshToken
				//user.RefreshTokenExpiresAt = time.Now().Add(time.Duration(opts.RefreshTokenValidInMinutes) * time.Minute).Format(time.RFC3339)
				user.Version = user.Version + 1

				// TODO: save all fields while fixed timestamp problem
				// db.Save(&user)

				switch err := db.Model(&user).Updates(map[string]interface{}{"access_token": accessToken,
					"refresh_token": refreshToken, "version": user.Version}).Error; err {
				case nil:
					authRes := authSuccessResponse{
						Uuid:          user.Uuid,
						Id:            user.Id,
						CreatedAt:     user.CreatedAt,
						UpdatedAt:     user.UpdatedAt,
						Version:       user.Version,
						Username:      user.UserName,
						Email:         user.Email,
						LastLogin:     user.LastLogin,
						Enabled:       user.Enabled,
						EmailVerified: user.EmailVerified,
						Tag:           "authUser",
						AccessToken:   user.AccessToken,
						// AccessTokenExpiresAt:  user.AccessTokenExpiresAt,
						RefreshToken: user.RefreshToken,
						// RefreshTokenExpiresAt: user.RefreshTokenExpiresAt,
					}
					context.JSON(http.StatusOK, authRes)
					context.Abort()
					return
				default:
					logger.Warn("an error  occurred while updating db", zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}
			} else {
				logger.Error("password validation failed")
				errorResponse(context, http.StatusBadRequest, errInvalidPass)
				context.Abort()
				return
			}
		}
	}
}
