package subscribe

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"unsafe"

	"github.com/icpd/subscribe2clash/internal/req"
	"github.com/icpd/subscribe2clash/internal/xbase64"
)

func GetRawProxiesFromLinks(links string) ([]string, error) {
	// url decode
	links, err := url.QueryUnescape(links)
	if err != nil {
		log.Printf("Warning: url decode %s 失败, err:%v", links, err)
		return nil, fmt.Errorf("url decode %s 失败, err:%w", links, err)
	}

	subLinks := strings.Split(links, ",")

	var rawProxiesSlice []string
	for _, link := range subLinks {
		content, err := req.HttpGet(link)
		if err != nil {
			log.Printf("Warning: 获取 %s 内容失败, err:%v", link, err)
			continue
		}

		rawProxiesSlice = append(rawProxiesSlice, ParseRawProxies(content))
	}

	return rawProxiesSlice, nil
}

func ParseRawProxies(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "ssd://") {
		content = content[6:]
		rawProxies, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(content)
		if err != nil {
			log.Println("ssd decode err:", err)
			return ""
		}

		return unsafe.String(&rawProxies[0], len(rawProxies))
	}

	// if is clash config content return directly
	if strings.Contains(content, "proxies") {
		return content
	}

	rawProxies, err := xbase64.Base64DecodeStripped(content)
	if err != nil {
		log.Println("base64 decode err:", err)
		return ""
	}

	return unsafe.String(&rawProxies[0], len(rawProxies))
}
