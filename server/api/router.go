package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()

	engine.POST("/v1/auth/login", authRPC)

	err := engine.Run(":8080")
	if err != nil {
		return
	}
}
