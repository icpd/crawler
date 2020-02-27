package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/whoisix/subscribe2clash/pkg/clash/subscribe"
	"github.com/whoisix/subscribe2clash/utils/mybase64"
)

func Clash(c *gin.Context) {
	content, err := subscribe.GetSubContent(c.Query("sub_link"))
	if err != nil {
		c.String(http.StatusBadRequest, "sub_link=订阅链接,英文逗号分割")
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
