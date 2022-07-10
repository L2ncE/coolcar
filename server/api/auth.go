package api

import (
	auth "coolcar/auth/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

func authRPC(ctx *gin.Context) {
	code := ctx.PostForm("code")
	res := auth.Test("127.0.0.1:8081", code)
	ctx.JSON(http.StatusOK, gin.H{
		"info": "成功",
		"data": res.AccessToken,
	})
	return
}
