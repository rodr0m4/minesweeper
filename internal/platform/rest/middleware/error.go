package middleware

import (
	"errors"
	"minesweeper/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		last := ctx.Errors.Last()

		if last == nil {
			return
		}

		err := last.Err

		var json JSONError

		json.Code = errorToStatusCode(err)
		json.Message = err.Error()

		ctx.AbortWithStatusJSON(json.Code, json)
	}
}

func errorToStatusCode(err error) int {
	if errors.Is(err, internal.ErrInvalidOperation) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
