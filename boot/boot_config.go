package boot

import (
	"os"
	"time"

	"github.com/whoisix/subscribe2clash/internal/acl"
	"github.com/whoisix/subscribe2clash/internal/global"
)

func generateConfig() {
	// 配置文件相关设置
	options := Options()

	if global.GenerateConfig {
		acl.GenerateConfig(options...)
		os.Exit(0)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(global.Tick) * time.Hour)
		for ; true; <-ticker.C {
			acl.GenerateConfig(options...)
		}
	}()
}
