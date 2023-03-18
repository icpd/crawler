package boot

import (
	"os"
	"time"

	"github.com/icpd/subscribe2clash/internal/acl"
	"github.com/icpd/subscribe2clash/internal/global"
)

func generateConfig() {
	// 配置文件相关设置
	options := Options()
	g := acl.New(options...)

	if global.GenerateConfig {
		g.GenerateConfig()
		os.Exit(0)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(global.Tick) * time.Hour)
		for ; true; <-ticker.C {
			g.GenerateConfig()
		}
	}()
}
