package api

import (
	"errors"
	"fmt"
	"github.com/icpd/subscribe2clash/internal/crypto"
	"github.com/icpd/subscribe2clash/internal/xbase64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icpd/subscribe2clash/internal/clash"
	"github.com/spf13/cast"
)

const (
	ssrHeader      = "ssr://"
	password       = "pwd://"
	vmessHeader    = "vmess://"
	ssHeader       = "ss://"
	trojanHeader   = "trojan://"
	hysteriaHeader = "hysteria://"
	key            = "link"
	v              = "v"
)

var (
	EasKey = "JIofjdsiowhj&*(hf"
)

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

func (cc *ClashController) Txt(c *gin.Context) {
	value, exists := c.GetQuery(v)
	if !exists || len(value) == 0 {
		c.String(http.StatusBadRequest, v+" 不能为空")
		c.Abort()
		return
	}
	if strings.HasPrefix(value, password) {
		decodeStr, err := xbase64.Base64DecodeStripped(subProtocolBody(value, password))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}

		value, err = crypto.Decrypt(EasKey, string(decodeStr))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
	}

	switch {
	case strings.HasPrefix(value, ssrHeader):
		tmp, err := ssrGenerate(subProtocolBody(value, ssrHeader))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		value = tmp

	case strings.HasPrefix(value, ssHeader):
		tmp, err := ssGenerate(subProtocolBody(value, ssHeader))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		value = tmp

	default:
		c.String(http.StatusBadRequest, v+"前缀格式不支持")
		c.Abort()
		return
	}

	_, exists = c.GetQuery("ssTxt")
	if exists {

	}

	nodeOnly, _ := c.GetQuery("nodeonly")
	config, err := clash.Config(clash.Txt, value, cast.ToBool(nodeOnly))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, config)
}

func subProtocolBody(proxy string, prefix string) string {
	return strings.TrimSpace(proxy[len(prefix):])
}

func ssGenerate(value string) (string, error) {
	params := strings.Split(value, `:`)

	if len(params) != 5 {
		return "", errors.New("数据格式错误，示例：ip:port:method:password:nodeName")
	}

	value = "ss://" + fmt.Sprintf(
		"%v@%v:%v#%v",
		xbase64.Base64EncodeStripped(fmt.Sprintf("%v:%v", params[2], params[3])),
		params[0],
		params[1],
		params[4],
	)

	return value, nil
}
func ssrGenerate(value string) (string, error) {
	params := strings.Split(value, `:`)

	if len(params) != 7 {
		return "", errors.New("数据格式错误，示例：ip:port:protocol:method:blending:password:nodeName")
	}

	params[5] = xbase64.Base64EncodeStripped(params[5])
	if len(params) > 6 {
		params[6] = xbase64.Base64EncodeStripped(params[6])
	} else {
		params[6] = xbase64.Base64EncodeStripped(params[0])
	}

	// group
	if len(params) < 8 {
		params = append(params, "")
	}
	params[7] = xbase64.Base64EncodeStripped("x.com")

	// 覆盖
	value = "ssr://" + xbase64.Base64EncodeStripped(
		fmt.Sprintf("%v:%v:%v:%v:%v:%v/?remarks=%v&group=%v",
			params[0],
			params[1],
			params[2],
			params[3],
			params[4],
			params[5],
			params[6],
			params[7],
		))

	return value, nil
}

func (cc *ClashController) Base64(c *gin.Context) {
	data, exists := c.GetQuery(key)
	if !exists {
		c.String(http.StatusBadRequest, key+"不能为空，本接口对ss连接进行base64加密\n"+
			"参考文档：https://github.com/hoochanlon/fq-book/blob/master/docs/append/srvurl.md\n"+
			"ss://method:password@server:port\n"+
			"ssr://ip:port:protocol:method:blending:password")
		c.Abort()
		return
	}
	ssLink := "ss://" + xbase64.Base64EncodeStripped(data)

	c.String(http.StatusOK, ssLink)
}

func (cc *ClashController) GenerateUrl(c *gin.Context) {
	// 获取表单参数
	linkType := strings.ToLower(c.PostForm("linkType"))
	ip := c.PostForm("ip")
	port := c.PostForm("port")
	method := c.PostForm("method")
	pwd := c.PostForm("password")
	name := c.PostForm("name")
	protocol := c.PostForm("protocol")
	blending := c.PostForm("blending")

	if len(linkType) == 0 || len(ip) == 0 || len(port) == 0 || len(method) == 0 || len(pwd) == 0 || len(name) == 0 {
		c.String(http.StatusOK, "参数不全 "+linkType)
		return
	}
	if linkType == "ssr" && (len(protocol) == 0 || len(blending) == 0) {
		c.String(http.StatusOK, "参数不全 "+linkType)
		return
	}

	data := ""
	switch {
	case linkType == "ss":
		data = fmt.Sprintf("ss://%v:%v:%v:%v:%v", ip, port, method, pwd, name)

	case linkType == "ssr":
		data = fmt.Sprintf("ssr://%v:%v:%v:%v:%v:%v:%v", ip, port, protocol, method, blending, pwd, name)

	default:
		c.String(http.StatusOK, "不支持")
		return
	}
	data, err := crypto.Encrypt(EasKey, data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	data = xbase64.Base64EncodeStripped(data)

	proto := c.Request.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	decodeUrl := fmt.Sprintf("%v://%v:%v/txt?v=pwd://%v",
		proto,
		c.Request.Host,
		c.Request.URL.Port(),
		data,
	)

	// 加密
	c.HTML(http.StatusOK, "link.html", gin.H{
		"url": decodeUrl,
	})
}
