package boot

import (
	"flag"
	"fmt"
	"os"

	"github.com/whoisix/subscribe2clash/library/global"
	"github.com/whoisix/subscribe2clash/library/req"
)

func init() {
	flag.BoolVar(&global.H, "h", false, "help")
	flag.BoolVar(&global.Gc, "gc", false, "生成clash配置文件")
	flag.StringVar(&global.Origin, "origin", "github", "acl规则获取地址。cn：国内镜像，github：github获取")
	flag.StringVar(&global.BaseFile, "b", "", "clash基础配置文件")
	flag.StringVar(&global.OutputFile, "o", "", "输出clash文件名")
	flag.StringVar(&global.ListenAddr, "l", "0.0.0.0", "listen address")
	flag.StringVar(&global.ListenPort, "p", "8162", "listen port")
	flag.StringVar(&req.Proxy, "proxy", "", "http proxy")
	flag.IntVar(&global.T, "t", 6, "规则更新频率（小时）")
	flag.Parse()
	fmt.Println("init flag ......")
}

func initFlag() {
	if global.H {
		flag.Usage()
		os.Exit(0)
	}
}
