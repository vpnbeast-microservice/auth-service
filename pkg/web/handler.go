package web

import (
	"auth-service/pkg/jwt"
	"auth-service/pkg/model"
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

func validateHandler() gin.HandlerFunc {
	// TODO: refactor
	return func(context *gin.Context) {
		ok, req := validateJsonRequest(context)
		if !ok {
			return
		}

		validateReq := req.(validateRequest)
		subject, roles, err, code := jwt.ValidateToken(validateReq.Token)
		if err != nil {
			validateRes := validateResponse{
				Tag: "validateToken",
				Status: false,
				ErrorMessage: err.Error(),
				HttpCode: code,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			context.JSON(code, validateRes)
			context.Abort()
			return
		}

		var user model.User
		switch err := db.Where("user_name = ?", subject).First(&user).Error; err {
		// switch err := db.Where("user_name = ?", authReq.Username).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn("no rows were returned!", zap.String("user", subject))
			validateRes := validateResponse{
				Tag: "validateToken",
				Status: false,
				ErrorMessage: "no such user",
				HttpCode: 404,
				Timestamp: time.Now().Format(time.RFC3339),
			}
			context.JSON(http.StatusNotFound, validateRes)
			context.Abort()
			return
		case nil:
			validateRes := validateResponse{
				Tag: "validateToken",
				Status: true,
				Username: subject,
				Roles: roles,
				HttpCode: 200,
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
		var authReq authRequest
		_, errSlice := isValidRequest(context, &authReq)
		if len(errSlice) != 0 {
			validationResponse(context, "authUser", errSlice)
			context.Abort()
			return
		}

		var user model.User
		switch err := db.Preload("Roles").Where("user_name = ?", authReq.Username).First(&user).Error; err {
		// switch err := db.Where("user_name = ?", authReq.Username).First(&user).Error; err {
		case gorm.ErrRecordNotFound:
			logger.Warn("no rows were returned!", zap.String("user", authReq.Username))
			errorResponse(context, "authUser", http.StatusNotFound, errUserNotFound)
			context.Abort()
			return
		case nil:
			// logger.Info("", zap.Any("role of user", user.Roles))
			encryptReq := encryptRequest{
				PlainText:     authReq.Password,
				EncryptedText: user.EncryptedPassword,
			}

			encryptRes, err := encryptReq.encrypt(authReq.Password, user.EncryptedPassword)
			if err != nil {
				logger.Error("an error occurred while making encryption request", zap.String("error", err.Error()))
				errorResponse(context, "authUser", http.StatusInternalServerError, errUnknown)
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
					errorResponse(context, "authUser", http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}

				refreshToken, err := jwt.GenerateToken(authReq.Username, roles, int32(opts.RefreshTokenValidInMinutes))
				if err != nil {
					logger.Error("an error occurred generating refresh token",
						zap.String("error", err.Error()))
					errorResponse(context, "authUser", http.StatusInternalServerError, errUnknown)
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
						Tag:                        "authUser",
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
					errorResponse(context, "authUser", http.StatusInternalServerError, errUnknown)
					context.Abort()
					return
				}
			} else {
				logger.Error("password validation failed")
				errorResponse(context, "authUser", http.StatusBadRequest, errInvalidPass)
				context.Abort()
				return
			}
		}
	}
}
