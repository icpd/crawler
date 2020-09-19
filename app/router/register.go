package router

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	clashRouter(r)
}