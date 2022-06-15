package router

import (
	"github.com/gin-gonic/gin"

	"github.com/icpd/subscribe2clash/app/api"
)

func clashRouter(r *gin.Engine) {
	clash := api.ClashController{}

	r.GET("/", clash.Clash)
}
