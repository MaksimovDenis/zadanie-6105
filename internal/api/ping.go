package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) CheckServer(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}
