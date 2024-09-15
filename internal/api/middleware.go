package api

import (
	"github.com/gin-gonic/gin"
)

func TokenInjector() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := "my-test-token"

		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		ctx.Next()
	}
}
