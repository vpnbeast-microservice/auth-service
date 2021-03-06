package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func validationResponse(ctx *gin.Context, errSlice []string) {
	ctx.JSON(http.StatusBadRequest, validationErrorResponse{
		ErrorMessage: errSlice,
		Status:       false,
		HttpCode:     http.StatusBadRequest,
		Timestamp:    time.Now(),
	})
}

func errorResponse(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, authFailResponse{
		ErrorMessage: err,
		Status:       false,
		HttpCode:     code,
		Timestamp:    time.Now(),
	})
}
