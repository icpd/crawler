package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/whoisix/subscribe2clash/pkg/clash/subscribe"
	"github.com/whoisix/subscribe2clash/utils/mybase64"
)

func Clash(c *gin.Context) {
	query := c.Request.URL.String()
	idx := strings.Index(query, "sub_link=")

	if idx < 0 {
		c.String(http.StatusBadRequest, "sub_link=订阅链接.")
		c.Abort()
		return
	}

	val := query[idx+9:]
	if val == "" {
		c.String(http.StatusBadRequest, "sub_link 不能为空")
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

	c.String(http.StatusOK, mybase64.UnicodeEmojiDecode(string(config)))
}
