package acl

// https://github.com/ACL4SSR/ACL4SSR
var github = map[string]string{
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
}

// 国内镜像 https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@latest/
var cn = map[string]string{
	"Apple":            "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/Apple.list",
	"BanAD":            "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/BanAD.list",
	"BanProgramAD":     "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/BanProgramAD.list",
	"ChinaCompanyIp":   "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/ChinaCompanyIp.list",
	"ChinaDomain":      "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/ChinaDomain.list",
	"ChinaIp":          "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/ChinaIp.list",
	"LocalAreaNetwork": "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/LocalAreaNetwork.list",
	"Microsoft":        "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/Microsoft.list",
	"ProxyGFWlist":     "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/ProxyGFWlist.list",
	"ProxyLite":        "https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master/Clash/ProxyLite.list",
}

func GetUrls(origin string, ProxyLite bool) map[string]string {
	var list map[string]string
	switch origin {
	case "cn":
		list = cn
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
