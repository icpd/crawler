package main

import (
	"flag"

	"github.com/whoisix/subscribe2clash/pkg/clash/acl"
)

var (
	gc         bool
	h          bool
	baseFile   string
	outputFile string
	origin     string
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.BoolVar(&gc, "gc", false, "生成clash配置文件")
	flag.StringVar(&origin, "origin", "github", "acl规则获取地址。cn：国内镜像，github：github获取")
	flag.StringVar(&baseFile, "b", "", "clash基础配置文件")
	flag.StringVar(&outputFile, "o", "", "输出clash文件名")
	flag.Parse()
}

func main() {
	if h {
		flag.Usage()
		return
	}

	if gc {
		var options []acl.GenOption
		if origin != "" {
			options = append(options, acl.WithOrigin(origin))
		}
		if baseFile != "" {
			options = append(options, acl.WithBaseFile(baseFile))
		}
		if outputFile != "" {
			options = append(options, acl.WithOutputFile(outputFile))
		}
		acl.GenerateConfig(options...)
	}

	acl.GenerateConfig()
}
