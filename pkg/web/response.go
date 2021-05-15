package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func validationResponse(ctx *gin.Context, errSlice []string) {
	ctx.JSON(http.StatusBadRequest, validationErrorResponse{
		Tag:          "authUser",
		ErrorMessage: errSlice,
		Status:       false,
		HttpCode:     http.StatusBadRequest,
		Timestamp:    time.Now(),
	})
}

func errorResponse(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, authFailResponse{
		Tag:          "authUser",
		ErrorMessage: err,
		Status:       false,
		HttpCode:     code,
		Timestamp:    time.Now(),
	})
}