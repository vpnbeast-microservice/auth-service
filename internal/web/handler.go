package web

import (
	"auth-service/internal/database"
	"auth-service/internal/jwt"
	"auth-service/internal/model"
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

func whoamiHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("Authorization")[7:]
		subject, _, err, code := jwt.ValidateToken(token)
		if err != nil {
			validateRes := validateResponse{
				Status:       false,
				ErrorMessage: err.Error(),
				HttpCode:     code,
				Timestamp:    time.Now().Format(time.RFC3339),
			}
			context.JSON(code, validateRes)
			context.Abort()
			return
		}
		var user model.User
		db := database.GetDatabase()
		switch err := db.Preload("Roles").Where(queryUsername, subject).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn(errNoRowsReturned, zap.String("user", subject))
			errorResponse(context, http.StatusNotFound, errUserNotFound)
			context.Abort()
			return
		case nil:
			authRes := authSuccessResponse{
				Uuid:                       user.Uuid,
				Id:                         user.Id,
				CreatedAt:                  user.CreatedAt,
				UpdatedAt:                  user.UpdatedAt,
				Version:                    user.Version,
				Username:                   user.UserName,
				Email:                      user.Email,
				LastLogin:                  user.LastLogin,
				Enabled:                    user.Enabled,
				EmailVerified:              user.EmailVerified,
				AccessToken:                user.AccessToken,
				AccessTokenExpiresAt:       user.AccessTokenExpiresAt,
				RefreshToken:               user.RefreshToken,
				RefreshTokenExpiresAt:      user.RefreshTokenExpiresAt,
				VerificationCodeCreatedAt:  user.VerificationCodeCreatedAt,
				VerificationCodeVerifiedAt: user.VerificationCodeVerifiedAt,
				Roles:                      user.Roles,
			}
			context.JSON(http.StatusOK, authRes)
			context.Abort()
			return
		}
	}
}

func refreshHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("Authorization")[7:]
		subject, roles, err, code := jwt.ValidateToken(token)
		if err != nil {
			validateRes := validateResponse{
				Status:       false,
				ErrorMessage: err.Error(),
				HttpCode:     code,
				Timestamp:    time.Now().Format(time.RFC3339),
			}
			context.JSON(code, validateRes)
			context.Abort()
			return
		}

		var user model.User
		db := database.GetDatabase()
		switch err := db.Preload("Roles").Where(queryUsername, subject).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn(errNoRowsReturned, zap.String("user", subject))
			errorResponse(context, http.StatusNotFound, errUserNotFound)
			context.Abort()
			return
		case nil:
			accessToken, err := jwt.GenerateToken(subject, roles, int32(opts.AccessTokenValidInMinutes))
			if err != nil {
				logger.Error("an error occurred generating access token",
					zap.String("error", err.Error()))
				errorResponse(context, http.StatusInternalServerError, errUnknown)
				context.Abort()
				return
			}

			refreshToken, err := jwt.GenerateToken(subject, roles, int32(opts.RefreshTokenValidInMinutes))
			if err != nil {
				logger.Error("an error occurred generating refresh token",
					zap.String("error", err.Error()))
				errorResponse(context, http.StatusInternalServerError, errUnknown)
				context.Abort()
				return
			}

			now := time.Now().Format(time.RFC3339)
			user.LastLogin = now
			user.UpdatedAt = now
			user.AccessToken = accessToken
			user.AccessTokenExpiresAt = time.Now().Add(time.Duration(opts.AccessTokenValidInMinutes) * time.Minute).Format(time.RFC3339)
			user.RefreshToken = refreshToken
			user.RefreshTokenExpiresAt = time.Now().Add(time.Duration(opts.RefreshTokenValidInMinutes) * time.Minute).Format(time.RFC3339)
			user.Version = user.Version + 1

			switch err := db.Save(&user).Error; err {
			case nil:
				authRes := authSuccessResponse{
					Uuid:                       user.Uuid,
					Id:                         user.Id,
					CreatedAt:                  user.CreatedAt,
					UpdatedAt:                  user.UpdatedAt,
					Version:                    user.Version,
					Username:                   user.UserName,
					Email:                      user.Email,
					LastLogin:                  user.LastLogin,
					Enabled:                    user.Enabled,
					EmailVerified:              user.EmailVerified,
					AccessToken:                user.AccessToken,
					AccessTokenExpiresAt:       user.AccessTokenExpiresAt,
					RefreshToken:               user.RefreshToken,
					RefreshTokenExpiresAt:      user.RefreshTokenExpiresAt,
					VerificationCodeCreatedAt:  user.VerificationCodeCreatedAt,
					VerificationCodeVerifiedAt: user.VerificationCodeVerifiedAt,
					Roles:                      user.Roles,
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
		}
	}
}

func validateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		req, _ := context.Get("data")
		validateReq := req.(validateRequest)
		subject, roles, err, code := jwt.ValidateToken(validateReq.Token)
		if err != nil {
			validateRes := validateResponse{
				Status:       false,
				ErrorMessage: err.Error(),
				HttpCode:     code,
				Timestamp:    time.Now().Format(time.RFC3339),
			}
			context.JSON(code, validateRes)
			context.Abort()
			return
		}

		var user model.User
		db := database.GetDatabase()
		switch err := db.Where(queryUsername, subject).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn(errNoRowsReturned, zap.String("user", subject))
			validateRes := validateResponse{
				Status:       false,
				ErrorMessage: "no such user",
				HttpCode:     404,
				Timestamp:    time.Now().Format(time.RFC3339),
			}
			context.JSON(http.StatusNotFound, validateRes)
			context.Abort()
			return
		case nil:
			validateRes := validateResponse{
				Status:    true,
				Username:  subject,
				Roles:     roles,
				HttpCode:  200,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			context.JSON(http.StatusOK, validateRes)
			context.Abort()
			return
		}
	}
}

func authenticateHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		req, _ := context.Get("data")
		logger.Info("", zap.Any("req", req))
		authReq := req.(authRequest)
		logger.Info("", zap.Any("authReq", authReq.Username))
		var user model.User
		db := database.GetDatabase()
		switch err := db.Preload("Roles").Where(queryUsername, authReq.Username).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn(errNoRowsReturned, zap.String("user", authReq.Username))
			errorResponse(context, http.StatusNotFound, errUserNotFound)
			context.Abort()
			return
		case nil:
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

			var roles []string
			for _, v := range user.Roles {
				logger.Info("appending role to roles", zap.String("role", v.Name))
				roles = append(roles, v.Name)
			}

			if encryptRes.Status {
				accessToken, err := jwt.GenerateToken(authReq.Username, roles, int32(opts.AccessTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating access token",
						zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}

				refreshToken, err := jwt.GenerateToken(authReq.Username, roles, int32(opts.RefreshTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating refresh token",
						zap.String("error", err.Error()))
					errorResponse(context, http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}

				now := time.Now().Format(time.RFC3339)
				user.LastLogin = now
				user.UpdatedAt = now
				user.AccessToken = accessToken
				user.AccessTokenExpiresAt = time.Now().Add(time.Duration(opts.AccessTokenValidInMinutes) * time.Minute).Format(time.RFC3339)
				user.RefreshToken = refreshToken
				user.RefreshTokenExpiresAt = time.Now().Add(time.Duration(opts.RefreshTokenValidInMinutes) * time.Minute).Format(time.RFC3339)
				user.Version = user.Version + 1

				switch err := db.Save(&user).Error; err {
				case nil:
					authRes := authSuccessResponse{
						Uuid:                       user.Uuid,
						Id:                         user.Id,
						CreatedAt:                  user.CreatedAt,
						UpdatedAt:                  user.UpdatedAt,
						Version:                    user.Version,
						Username:                   user.UserName,
						Email:                      user.Email,
						LastLogin:                  user.LastLogin,
						Enabled:                    user.Enabled,
						EmailVerified:              user.EmailVerified,
						AccessToken:                user.AccessToken,
						AccessTokenExpiresAt:       user.AccessTokenExpiresAt,
						RefreshToken:               user.RefreshToken,
						RefreshTokenExpiresAt:      user.RefreshTokenExpiresAt,
						VerificationCodeCreatedAt:  user.VerificationCodeCreatedAt,
						VerificationCodeVerifiedAt: user.VerificationCodeVerifiedAt,
						Roles:                      user.Roles,
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
