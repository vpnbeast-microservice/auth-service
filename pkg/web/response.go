package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func validationResponse(ctx *gin.Context, tag string, errSlice []string) {
	ctx.JSON(http.StatusBadRequest, validationErrorResponse{
		Tag:          tag,
		ErrorMessage: errSlice,
		Status:       false,
		HttpCode:     http.StatusBadRequest,
		Timestamp:    time.Now(),
	})
}

func errorResponse(ctx *gin.Context, tag string, code int, err string) {
	ctx.JSON(code, authFailResponse{
		Tag:          tag,
		ErrorMessage: err,
		Status:       false,
		HttpCode:     code,
		Timestamp:    time.Now(),
	})
}
