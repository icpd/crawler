package subscribe

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/whoisix/subscribe2clash/internal/xbase64"
)

var (
	//ssReg      = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)\?plugin=([^;]+);\w+=(\w+)(?:;obfs-host=)?([^#]+)?#(.+)`)
	ssReg2 = regexp.MustCompile(`(?m)ss://([\-0-9a-z]+):(.+)@(.+):(\d+)(.+)?#(.+)`)

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
	vmconfig, err := xbase64.Base64DecodeStripped(s)
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
	if vmess.Net == "ws" {
		clashVmess.Network = vmess.Net
		clashVmess.WSOpts.Path = vmess.Path
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
		ss.PluginOpts = &PluginOpts{
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
	rawSSRConfig, err := xbase64.Base64DecodeStripped(s)
	if err != nil {
		return ClashRSSR{}
	}
	params := strings.Split(string(rawSSRConfig), `:`)
	if len(params) != 6 {
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
	if ssr.Protocol == "origin" && ssr.OBFS == "plain" {
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
	if len(suffix) != 2 {
		return ClashRSSR{}
	}
	passwordBase64 := suffix[0]
	password, err := xbase64.Base64DecodeStripped(passwordBase64)
	if err != nil {
		return ClashRSSR{}
	}
	ssr.Password = string(password)

	m, err := url.ParseQuery(suffix[1])
	if err != nil {
		return ClashRSSR{}
	}

	for k, v := range m {
		de, err := xbase64.Base64DecodeStripped(v[0])
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

	findStr := regexp.MustCompile(`(?m)ss://([\/\+=\w]+)(@.+)?#(.+)`).FindStringSubmatch(s)
	if len(findStr) < 4 {
		return ClashSS{}
	}

	rawSSRConfig, err := xbase64.Base64DecodeStripped(findStr[1])
	if err != nil {
		return ClashSS{}
	}

	s = strings.ReplaceAll(s, findStr[1], string(rawSSRConfig))
	findStr = ssReg2.FindStringSubmatch(s)

	ss := ClashSS{}
	ss.Type = "ss"
	ss.Cipher = findStr[1]
	ss.Password = findStr[2]
	ss.Server = findStr[3]
	ss.Port = findStr[4]
	ss.Name = findStr[6]

	if findStr[5] != "" && strings.Contains(findStr[5], "plugin") {
		query := findStr[5][strings.Index(findStr[5], "?")+1:]
		queryMap, err := url.ParseQuery(query)
		if err != nil{
			return ClashSS{}
		}

		ss.Plugin = queryMap["plugin"][0]
		p := new(PluginOpts)
		switch {
		case strings.Contains(ss.Plugin, "obfs"):
			ss.Plugin = "obfs"
			p.Mode = queryMap["obfs"][0]
			if strings.Contains(query, "obfs-host="){
				p.Host = queryMap["obfs-host"][0]
			}
			break
		case ss.Plugin == "v2ray-plugin":
			p.Mode = queryMap["mode"][0]
			if strings.Contains(query, "host="){
				p.Host = queryMap["host"][0]
			}
			if strings.Contains(query, "path="){
				p.Path = queryMap["path"][0]
			}
			p.Mux = strings.Contains(query, "mux");
			p.Tls = strings.Contains(query, "tls");
			p.SkipCertVerify = true
			break
		}
		ss.PluginOpts = p
	}

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
