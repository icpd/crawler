package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/icpd/subscribe2clash/internal/subscribe"
	"github.com/icpd/subscribe2clash/internal/xbase64"
)

const substr = "sub_link"

type ClashController struct{}

func (cc *ClashController) Clash(c *gin.Context) {
	query := c.Request.URL.String()
	idx := strings.Index(query, substr)

	if idx < 0 {
		c.String(http.StatusBadRequest, substr+"=订阅链接.")
		c.Abort()
		return
	}

	val := query[idx+len(substr)+1:]
	if val == "" {
		c.String(http.StatusBadRequest, substr+"不能为空")
		c.Abort()
		return
	}

	content, err := subscribe.GetSubContent(val)
	if err != nil {
		c.String(http.StatusBadRequest, "请求失败:"+err.Error())
		c.Abort()
		return
	}

	proxies := subscribe.ParseProxy(content)
	config, err := subscribe.GenerateClashConfig(proxies)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, xbase64.UnicodeEmojiDecode(string(config)))
}
