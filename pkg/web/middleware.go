package web

import (
	"github.com/gin-gonic/gin"
)

func authRequestValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authReq authRequest
		_, errSlice := isValidRequest(c, &authReq)
		if len(errSlice) != 0 {
			validationResponse(c, errSlice)
			c.Abort()
			return
		}

		c.Set("data", authReq)
		c.Next()
	}
}

func validateRequestValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var validateReq validateRequest
		_, errSlice := isValidRequest(c, &validateReq)
		if len(errSlice) != 0 {
			validationResponse(c, errSlice)
			c.Abort()
			return
		}

		c.Set("data", validateReq)
		c.Next()
	}
}
