package subscribe

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/icpd/subscribe2clash/internal/xbase64"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

const (
	ssrHeader      = "ssr://"
	vmessHeader    = "vmess://"
	ssHeader       = "ss://"
	trojanHeader   = "trojan://"
	hysteriaHeader = "hysteria://"
)

var (
	//ssReg      = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)\?plugin=([^;]+);\w+=(\w+)(?:;obfs-host=)?([^#]+)?#(.+)`)
	ssReg2 = regexp.MustCompile(`(?m)([\-0-9a-z]+):(.+)@(.+):(\d+)(.+)?#(.+)`)
	ssReg  = regexp.MustCompile(`(?m)([^@]+)(@.+)?#?(.+)?`)

	trojanReg  = regexp.MustCompile(`(?m)^trojan://(.+)@(.+):(\d+)\?allowInsecure=\d&peer=(.+)#(.+)`)
	trojanReg2 = regexp.MustCompile(`(?m)^trojan://(.+)@(.+):(\d+)#(.+)$`)
)

func ParseProxy(contentSlice []string) []any {
	var proxies []any
	for _, v := range contentSlice {
		// try unmarshal clash config
		var c Clash
		if err := yaml.Unmarshal([]byte(v), &c); err == nil {
			for _, pg := range c.Proxies {
				proxies = append(proxies, pg)
			}
			continue
		}

		// ssd
		if strings.Contains(v, "airport") {
			ssSlice := ssdConf(v)
			for _, ss := range ssSlice {
				if ss.Name != "" {
					proxies = append(proxies, ss)
				}
			}
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(v))
		for scanner.Scan() {
			proxy := parseProxy(scanner.Text())
			if proxy != nil {
				proxies = append(proxies, proxy)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("parse proxy failed, err: %v", err)
		}
	}

	return proxies
}

func subProtocolBody(proxy string, prefix string) string {
	return strings.TrimSpace(proxy[len(prefix):])
}

func parseProxy(proxy string) any {
	switch {
	case strings.HasPrefix(proxy, ssrHeader):
		return ssrConf(subProtocolBody(proxy, ssrHeader))
	case strings.HasPrefix(proxy, vmessHeader):
		return v2rConf(subProtocolBody(proxy, vmessHeader))
	case strings.HasPrefix(proxy, ssHeader):
		return ssConf(subProtocolBody(proxy, ssHeader))
	case strings.HasPrefix(proxy, trojanHeader):
		return trojanConf(subProtocolBody(proxy, trojanHeader))
	case strings.HasPrefix(proxy, hysteriaHeader):
		return hysteriaConf(proxy)
	}

	return nil
}

type ClashHysteria struct {
	Name                string   `yaml:"name"`
	Type                string   `yaml:"type"`
	Server              string   `yaml:"server"`
	Port                int      `yaml:"port"`
	AuthStr             string   `yaml:"auth-str"`
	Obfs                string   `yaml:"obfs"`
	ObfsParams          string   `yaml:"obfs-param"`
	Alpn                []string `yaml:"alpn"`
	Protocol            string   `yaml:"protocol"`
	Up                  string   `yaml:"up"`
	Down                string   `yaml:"down"`
	Sni                 string   `yaml:"sni"`
	SkipCertVerify      bool     `yaml:"skip-cert-verify"`
	RecvWindowConn      int      `yaml:"recv-window-conn"`
	RecvWindow          int      `yaml:"recv-window"`
	Ca                  string   `yaml:"ca"`
	CaStr               string   `yaml:"ca-str"`
	DisableMtuDiscovery bool     `yaml:"disable_mtu_discovery"`
	Fingerprint         string   `yaml:"fingerprint"`
	FastOpen            bool     `yaml:"fast-open"`
}

// https://hysteria.network/docs/uri-scheme/
// hysteria://host:port?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarks
func hysteriaConf(body string) any {
	u, err := url.Parse(body)
	if err != nil {
		log.Printf("parse hysteria failed, err: %v", err)
		return nil
	}

	query := u.Query()
	return &ClashHysteria{
		Name:                u.Fragment,
		Type:                "hysteria",
		Server:              u.Hostname(),
		Port:                cast.ToInt(u.Port()),
		AuthStr:             query.Get("auth"),
		Obfs:                query.Get("obfs"),
		Alpn:                []string{query.Get("alpn")},
		Protocol:            query.Get("protocol"),
		Up:                  query.Get("upmbps"),
		Down:                query.Get("downmbps"),
		Sni:                 query.Get("peer"),
		SkipCertVerify:      cast.ToBool(query.Get("insecure")),
		RecvWindowConn:      cast.ToInt(query.Get("recv-window-conn")),
		RecvWindow:          cast.ToInt(query.Get("recv-window")),
		Ca:                  query.Get("ca"),
		CaStr:               query.Get("ca-str"),
		DisableMtuDiscovery: cast.ToBool(query.Get("disable_mtu_discovery")),
		Fingerprint:         query.Get("fingerprint"),
		FastOpen:            cast.ToBool(query.Get("fast-open")),
	}
}

func v2rConf(s string) ClashVmess {
	vmconfig, err := xbase64.Base64DecodeStripped(s)
	if err != nil {
		return ClashVmess{}
	}
	vmess := Vmess{}
	err = json.Unmarshal(vmconfig, &vmess)
	if err != nil {
		log.Printf("v2ray config json unmarshal failed, err: %v", err)
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

	findStr := ssReg.FindStringSubmatch(s)
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
		if err != nil {
			return ClashSS{}
		}

		ss.Plugin = queryMap["plugin"][0]
		p := new(PluginOpts)
		switch {
		case strings.Contains(ss.Plugin, "obfs"):
			ss.Plugin = "obfs"
			p.Mode = queryMap["obfs"][0]
			if strings.Contains(query, "obfs-host=") {
				p.Host = queryMap["obfs-host"][0]
			}
		case ss.Plugin == "v2ray-plugin":
			p.Mode = queryMap["mode"][0]
			if strings.Contains(query, "host=") {
				p.Host = queryMap["host"][0]
			}
			if strings.Contains(query, "path=") {
				p.Path = queryMap["path"][0]
			}
			p.Mux = strings.Contains(query, "mux")
			p.Tls = strings.Contains(query, "tls")
			p.SkipCertVerify = true
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
