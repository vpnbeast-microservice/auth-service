package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		logger.Error("an error occured while writing response", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	var request authRequest
	err := decodeJSONBody(w, r, &request)
	if err != nil {
		logger.Error("an error occured while decoding json body", zap.String("error", err.Error()))
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// TODO: request validation

	if (len(request.Username) < 3 || len(request.Username) > 40) || (len(request.Password) < 5 || len(request.Password) > 200) {
		logger.Warn("validation failed", zap.String("user", request.Username),
			zap.String("password", request.Password))
		return
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
		failResponse := authFailResponse{
			Tag:          "getUser",
			ErrorMessage: "User not found!",
			Status:       false,
			HttpCode:     404,
			Timestamp:    time.Now(),
		}

		responseBytes, err := json.Marshal(failResponse)
		if err != nil {
			logger.Error("an error occured while marshaling response", zap.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBytes)
		if err != nil {
			logger.Error("an error occured while writing response", zap.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case nil:
		// TODO: check if password is correct, return success response else return fail response
		logger.Info("user fetched")
		result.Tag = "getToken"
		responseBytes, err := json.Marshal(result)
		if err != nil {
			logger.Error("an error occured while marshaling response", zap.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBytes)
		if err != nil {
			logger.Error("an error occured while writing response", zap.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		panic(err)
	}


}