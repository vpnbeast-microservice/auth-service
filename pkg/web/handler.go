package web

import (
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