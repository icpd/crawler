package router

import (
	"github.com/gin-gonic/gin"
	"github.com/icpd/subscribe2clash/app/api"
	"net/http"
)

func clashRouter(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")

	clash := api.ClashController{}

	// ssr/sss格式
	// https://github.com/hoochanlon/fq-book/blob/master/docs/append/srvurl.md

	r.GET("/clash", clash.Clash)
	r.GET("/txt", clash.Txt)
	r.POST("/generate-url", clash.GenerateUrl)

	r.GET("/build", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
