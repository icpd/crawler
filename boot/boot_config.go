package boot

import (
	"os"
	"time"

	"github.com/whoisix/subscribe2clash/pkg/acl"
	"github.com/whoisix/subscribe2clash/pkg/global"
)

func generateConfig() {
	// 配置文件相关设置
	options := Options()

	if global.GenerateConfig {
		acl.GenerateConfig(options...)
		os.Exit(0)
	}

	go func() {
		acl.GenerateConfig(options...)
		tick := time.Tick(time.Duration(global.Tick) * time.Hour)
		for {
			<-tick
			acl.GenerateConfig(options...)
		}
	}()
}
