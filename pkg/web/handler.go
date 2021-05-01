package web

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
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
	var request authenticationRequest
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
}