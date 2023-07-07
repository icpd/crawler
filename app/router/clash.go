package router

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/icpd/subscribe2clash/app/api"
)

func clashRouter(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")

	clash := api.ClashController{}

	r.GET("/clash", clash.Clash)
	r.GET("/txt", clash.Txt)
	r.GET("/base64", clash.Base64)

	r.GET("/build", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/generate-url", clash.GenerateUrl)
}
