package acl

import "strings"

var (
	// https://github.com/ACL4SSR/ACL4SSR
	github = map[string]string{
		"Apple":            "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/Apple.list",
		"BanAD":            "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/BanAD.list",
		"BanProgramAD":     "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/BanProgramAD.list",
		"ChinaCompanyIp":   "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/ChinaCompanyIp.list",
		"ChinaDomain":      "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/ChinaDomain.list",
		"ChinaIp":          "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/ChinaIp.list",
		"LocalAreaNetwork": "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/LocalAreaNetwork.list",
		"Microsoft":        "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/Microsoft.list",
		"ProxyGFWlist":     "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/ProxyGFWlist.list",
		"ProxyLite":        "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/ProxyLite.list",
		"BanEasyListChina": "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/BanEasyListChina.list",
	}

	Group = map[string]string{
		"Apple":            "DIRECT",
		"BanAD":            "广告拦截",
		"BanProgramAD":     "广告拦截",
		"ChinaCompanyIp":   "DIRECT",
		"ChinaDomain":      "DIRECT",
		"ChinaIp":          "DIRECT",
		"LocalAreaNetwork": "DIRECT",
		"Microsoft":        "DIRECT",
		"ProxyGFWlist":     "节点选择",
		"ProxyLite":        "节点选择",
		"BanEasyListChina": "广告拦截",
	}
	Sort = []string{
		"BanAD", "BanProgramAD", "BanEasyListChina",
		"ProxyGFWlist", "ProxyLite",
		"ChinaDomain", "Microsoft", "LocalAreaNetwork", "ChinaCompanyIp", "ChinaIp",
	}
)

func GetUrls(origin string, ProxyLite bool) map[string]string {
	var list map[string]string
	switch origin {
	case "cn":
		list = getCn()
	case "github":
		fallthrough
	default:
		list = github
	}

	if ProxyLite {
		delete(list, "ProxyGFWlist")
	} else {
		delete(list, "ProxyLite")
	}

	return list
}

func getCn() map[string]string {
	cn := make(map[string]string, len(github))
	for r, u := range github {
		cn[r] = strings.ReplaceAll(u, "raw.githubusercontent.com", "raw.sevencdn.com")
	}

	return cn
}
