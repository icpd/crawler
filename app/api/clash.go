package api

import (
	"errors"
	"fmt"
	"github.com/icpd/subscribe2clash/internal/xbase64"
	"github.com/wumansgy/goEncrypt/aes"
	"log"
	"net/http"
	"net/url"
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
	data, exists := c.GetQuery(v)
	if !exists || len(data) == 0 {
		c.String(http.StatusBadRequest, v+" 不能为空")
		c.Abort()
		return
	}
	if strings.HasPrefix(data, password) {
		data = subProtocolBody(data, password)
		dataByte, err := aes.AesCbcDecryptByHex(data, []byte(EasKey), nil)
		if err != nil {
			log.Printf("decode fail,%v\n", err.Error())
			c.String(http.StatusBadRequest, "decode failed.")
			c.Abort()
			return
		}
		data = string(dataByte)
	}

	switch {
	case strings.HasPrefix(data, ssrHeader):
		tmp, err := ssrGenerate(subProtocolBody(data, ssrHeader))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		data = tmp

	case strings.HasPrefix(data, ssHeader):
		tmp, err := ssGenerate(subProtocolBody(data, ssHeader))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		data = tmp

	default:
		c.String(http.StatusBadRequest, v+"前缀格式不支持")
		c.Abort()
		return
	}

	nodeOnly, _ := c.GetQuery("nodeonly")
	config, err := clash.Config(clash.Txt, data, cast.ToBool(nodeOnly))
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

	data, linkStr := "", ""
	switch {
	case linkType == "ss":
		data = fmt.Sprintf("%v:%v:%v:%v:%v", ip, port, method, pwd, name)
		linkStr = fmt.Sprintf("%v%v",
			ssHeader, xbase64.Base64EncodeStripped(fmt.Sprintf("%v:%v@%v:%v", method, pwd, ip, port)))
		data = fmt.Sprintf("%v%v", ssHeader, data)

	case linkType == "ssr":
		data = fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v", ip, port, protocol, method, blending, pwd, name)
		linkStr = fmt.Sprintf("%v%v", ssrHeader,
			xbase64.Base64EncodeStripped(fmt.Sprintf("%v:%v:%v:%v:%v:%v/?remarks=%v",
				ip, port, protocol, method, blending, pwd, xbase64.Base64EncodeStripped(name))))
		data = fmt.Sprintf("%v%v", ssrHeader, data)

	default:
		c.String(http.StatusOK, "不支持")
		return
	}
	data, err := aes.AesCbcEncryptHex([]byte(data), []byte(EasKey), nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	proto := c.Request.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	log.Println(c.Request.Host)
	log.Println(c.Request.URL.Port())

	decodeUrl := fmt.Sprintf("%v://%v/txt?v=%v",
		proto,
		c.Request.Host,
		url.PathEscape("pwd://"+data),
	)

	// 加密
	c.HTML(http.StatusOK, "link.html", gin.H{
		"url":  decodeUrl,
		"link": linkStr,
	})
}
