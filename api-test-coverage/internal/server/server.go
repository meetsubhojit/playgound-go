package server

import (
	"api-test-coverage/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	r.GET("/time", func(ctx *gin.Context) {
		resp := api.GetTime(ctx)
		ctx.JSON(http.StatusOK, resp)
	})
	r.Run()
}
