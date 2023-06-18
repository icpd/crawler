package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpd/subscribe2clash/internal/clash"
	"github.com/spf13/cast"
)

const key = "link"

type ClashController struct{}

func (cc *ClashController) Clash(c *gin.Context) {
	links, exists := c.GetQuery(key)
	if !exists {
		links, _ = c.GetQuery("sub_link") // 兼容旧key
	}

	if links == "" {
		c.String(http.StatusBadRequest, key+"不能为空")
		c.Abort()
		return
	}

	nodeOnly, _ := c.GetQuery("nodeonly")
	config, err := clash.Config(clash.Url, links, cast.ToBool(nodeOnly))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, config)
}
