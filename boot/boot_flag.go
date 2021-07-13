package boot

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/whoisix/subscribe2clash/constant"
	"github.com/whoisix/subscribe2clash/pkg/global"
	"github.com/whoisix/subscribe2clash/pkg/req"
)

func init() {
	flag.BoolVar(&global.H, "h", false, "帮助信息")
	flag.BoolVar(&global.Gc, "gc", false, "生成clash配置文件")
	flag.StringVar(&global.BaseFile, "b", "", "clash基础配置文件")
	flag.StringVar(&global.OutputFile, "o", "", "输出clash文件名")
	flag.StringVar(&global.Listen, "l", "0.0.0.0:8162", "监听地址")
	flag.StringVar(&req.Proxy, "proxy", "", "http代理")
	flag.IntVar(&global.T, "t", 6, "规则更新频率（小时）")
	flag.BoolVar(&global.Version, "version", false, "查看版本信息")
	flag.Parse()
}

func initFlag() {
	if global.H {
		flag.Usage()
		os.Exit(0)
	}

	if global.Version {
		fmt.Printf("subscribe2clash %s %s %s %s\n", constant.Version, runtime.GOOS, runtime.GOARCH, constant.BuildTime)
		os.Exit(0)
	}
}
