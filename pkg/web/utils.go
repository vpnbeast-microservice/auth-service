package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func validateJsonRequest(context *gin.Context) (bool, interface{}) {
	var validateReq validateRequest
	_, errSlice := isValidRequest(context, &validateReq)
	if len(errSlice) != 0 {
		validateRes := validateResponse{
			Status:       false,
			ErrorMessage: "not a valid json request",
			HttpCode:     400,
			Timestamp:    time.Now().Format(time.RFC3339),
		}
		context.JSON(http.StatusBadRequest, validateRes)
		context.Abort()
		return false, validateRequest{}
	}

	return true, validateReq
}
