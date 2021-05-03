package web

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
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

		translator := en.New()
		uni := ut.New(translator, translator)

		trans, found := uni.GetTranslator("en")
		if !found {
			// TODO: return json error
			logger.Error("translator not found")
			context.JSON(http.StatusBadRequest, errorResponse{
				Tag:          "authUser",
				ErrorMessage: "unknown error occured at the backend",
				Status:       false,
				HttpCode:     500,
				Timestamp:    time.Now(),
			})
			context.Abort()
			return
		}

		v := validator.New()

		if err := entranslations.RegisterDefaultTranslations(v, trans); err != nil {
			logger.Error("can not register translation", zap.String("error", err.Error()))
			context.JSON(http.StatusBadRequest, errorResponse{
				Tag:          "authUser",
				ErrorMessage: "unknown error occured at the backend",
				Status:       false,
				HttpCode:     500,
				Timestamp:    time.Now(),
			})
			context.Abort()
			return
		}

		if err := context.ShouldBindJSON(&request); err == nil {
			if err := v.Struct(&request); err != nil {
				var errSlice []string
				for _, e := range err.(validator.ValidationErrors) {
					errSlice = append(errSlice, e.Translate(trans))
				}
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
		}

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
			// TODO: check if password is correct, return success response else return fail response
			// TODO: generate accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt, update database and return response
			// TODO: update last_login column on db and return response(aspect oriented programing? check https://github.com/gogap/aop)
			logger.Info("user fetched")
			result.Tag = "authUser"
			context.JSON(http.StatusOK, result)
			context.Abort()
			return
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