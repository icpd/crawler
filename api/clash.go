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
	val := query[idx+9:]

	content, err := subscribe.GetSubContent(val)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
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
