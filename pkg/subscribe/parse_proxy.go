package subscribe

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/whoisix/subscribe2clash/pkg/mybase64"
)

var (
	ssReg      = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)\?plugin=([^;]+);\w+=(\w+)(?:;obfs-host=)?([^#]+)?#(.+)`)
	trojanReg  = regexp.MustCompile(`(?m)^trojan://(.+)@(.+):(\d+)\?allowInsecure=\d&peer=(.+)#(.+)`)
	trojanReg2 = regexp.MustCompile(`(?m)^trojan://(.+)@(.+):(\d+)#(.+)$`)
)

func ParseProxy(contentSlice []string) []interface{} {
	var proxies []interface{}
	for _, v := range contentSlice {
		// ssd
		if strings.Contains(v, "airport") {
			ssSlice := ssdConf(v)
			for _, ss := range ssSlice {
				if !filterNode(ss.Name) {
					proxies = append(proxies, ss)
				}
			}
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(v))
		for scanner.Scan() {
			switch {
			case strings.HasPrefix(scanner.Text(), "ssr://"):
				s := scanner.Text()[6:]
				s = strings.TrimSpace(s)
				ssr := ssrConf(s)
				if ssr.Name != "" && !filterNode(ssr.Name) {
					proxies = append(proxies, ssr)
				}
			case strings.HasPrefix(scanner.Text(), "vmess://"):
				s := scanner.Text()[8:]
				s = strings.TrimSpace(s)
				clashVmess := v2rConf(s)
				if clashVmess.Name != "" && !filterNode(clashVmess.Name) {
					proxies = append(proxies, clashVmess)
				}
			case strings.HasPrefix(scanner.Text(), "ss://"):
				s := strings.TrimSpace(scanner.Text())
				ss := ssConf(s)
				if ss.Name != "" && !filterNode(ss.Name) {
					proxies = append(proxies, ss)
				}
			case strings.HasPrefix(scanner.Text(), "trojan://"):
				s := scanner.Text()
				s = strings.TrimSpace(s)
				trojan := trojanConf(s)
				if trojan.Name != "" && !filterNode(trojan.Name) {
					proxies = append(proxies, trojan)
				}
			}
		}
	}

	return proxies
}

func v2rConf(s string) ClashVmess {
	vmconfig, err := mybase64.Base64DecodeStripped(s)
	if err != nil {
		return ClashVmess{}
	}
	vmess := Vmess{}
	err = json.Unmarshal(vmconfig, &vmess)
	if err != nil {
		log.Println(err)
		return ClashVmess{}
	}
	clashVmess := ClashVmess{}
	clashVmess.Name = vmess.PS

	clashVmess.Type = "vmess"
	clashVmess.Server = vmess.Add
	switch vmess.Port.(type) {
	case string:
		clashVmess.Port, _ = vmess.Port.(string)
	case int:
		clashVmess.Port, _ = vmess.Port.(int)
	case float64:
		clashVmess.Port, _ = vmess.Port.(float64)
	default:

	}
	clashVmess.UUID = vmess.ID
	clashVmess.AlterID = vmess.Aid
	clashVmess.Cipher = vmess.Type
	if strings.EqualFold(vmess.TLS, "tls") {
		clashVmess.TLS = true
	} else {
		clashVmess.TLS = false
	}
	if "ws" == vmess.Net {
		clashVmess.Network = vmess.Net
		clashVmess.WSPATH = vmess.Path
	}

	return clashVmess
}

func ssdConf(ssdJson string) []ClashSS {
	var ssd SSD
	err := json.Unmarshal([]byte(ssdJson), &ssd)
	if err != nil {
		log.Println("ssd json unmarshal err:", err)
		return nil
	}

	var clashSSSlice []ClashSS
	for _, server := range ssd.Servers {
		options, err := url.ParseQuery(server.PluginOptions)
		if err != nil {
			continue
		}

		var ss ClashSS
		ss.Type = "ss"
		ss.Name = server.Remarks
		ss.Cipher = server.Encryption
		ss.Password = server.Password
		ss.Server = server.Server
		ss.Port = server.Port
		ss.Plugin = server.Plugin
		ss.PluginOpts = PluginOpts{
			Mode: options["obfs"][0],
			Host: options["obfs-host"][0],
		}

		switch {
		case strings.Contains(ss.Plugin, "obfs"):
			ss.Plugin = "obfs"
		}

		clashSSSlice = append(clashSSSlice, ss)
	}

	return clashSSSlice
}

func ssrConf(s string) ClashRSSR {
	rawSSRConfig, err := mybase64.Base64DecodeStripped(s)
	if err != nil {
		return ClashRSSR{}
	}
	params := strings.Split(string(rawSSRConfig), `:`)
	if 6 != len(params) {
		return ClashRSSR{}
	}
	ssr := ClashRSSR{}
	ssr.Type = "ssr"
	ssr.Server = params[SSRServer]
	ssr.Port = params[SSRPort]
	ssr.Protocol = params[SSRProtocol]
	ssr.Cipher = params[SSRCipher]
	ssr.OBFS = params[SSROBFS]

	// 如果兼容ss协议，就转换为clash的ss配置
	// https://github.com/Dreamacro/clash
	if "origin" == ssr.Protocol && "plain" == ssr.OBFS {
		switch ssr.Cipher {
		case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm",
			"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
			"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
			"rc4-md5", "chacha20", "chacha20-ietf", "xchacha20",
			"chacha20-ietf-poly1305", "xchacha20-ietf-poly1305":
			ssr.Type = "ss"
		}
	}
	suffix := strings.Split(params[SSRSuffix], "/?")
	if 2 != len(suffix) {
		return ClashRSSR{}
	}
	passwordBase64 := suffix[0]
	password, err := mybase64.Base64DecodeStripped(passwordBase64)
	if err != nil {
		return ClashRSSR{}
	}
	ssr.Password = string(password)

	m, err := url.ParseQuery(suffix[1])
	if err != nil {
		return ClashRSSR{}
	}

	for k, v := range m {
		de, err := mybase64.Base64DecodeStripped(v[0])
		if err != nil {
			return ClashRSSR{}
		}
		switch k {
		case "obfsparam":
			ssr.OBFSParam = string(de)
			continue
		case "protoparam":
			ssr.ProtocolParam = string(de)
			continue
		case "remarks":
			ssr.Name = string(de)
			continue
		case "group":
			continue
		}
	}

	return ssr
}

func ssConf(s string) ClashSS {
	s, err := url.PathUnescape(s)
	if err != nil {
		return ClashSS{}
	}

	findStr := ssReg.FindStringSubmatch(s)
	if len(findStr) < 6 {
		return ClashSS{}
	}
	rawSSRConfig, err := mybase64.Base64DecodeStripped(findStr[1])
	if err != nil {
		return ClashSS{}
	}
	params := strings.Split(string(rawSSRConfig), `:`)
	if 2 != len(params) {
		return ClashSS{}
	}

	ss := ClashSS{}
	ss.Type = "ss"
	ss.Cipher = params[0]
	ss.Password = params[1]
	ss.Server = findStr[2]
	ss.Port = findStr[3]

	ss.Plugin = findStr[4]
	switch {
	case strings.Contains(ss.Plugin, "obfs"):
		ss.Plugin = "obfs"
	}

	p := PluginOpts{
		Mode: findStr[5],
	}
	p.Host = findStr[6]

	ss.Name = findStr[7]
	ss.PluginOpts = p

	return ss
}

func trojanConf(s string) Trojan {
	s, err := url.PathUnescape(s)
	if err != nil {
		return Trojan{}
	}

	findStr := trojanReg.FindStringSubmatch(s)
	if len(findStr) == 6 {
		return Trojan{
			Name:     findStr[5],
			Type:     "trojan",
			Server:   findStr[2],
			Password: findStr[1],
			Sni:      findStr[4],
			Port:     findStr[3],
		}
	}

	findStr = trojanReg2.FindStringSubmatch(s)
	if len(findStr) < 5 {
		return Trojan{}
	}

	return Trojan{
		Name:     findStr[4],
		Type:     "trojan",
		Server:   findStr[2],
		Password: findStr[1],
		Port:     findStr[3],
	}
}

func filterNode(nodeName string) bool {
	// 过滤剩余流量
	if strings.Contains(nodeName, "剩余流量") {
		return true
	}

	if strings.Contains(nodeName, "过期时间") {
		return true
	}

	return false
}
