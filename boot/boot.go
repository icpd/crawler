package boot

import (
	"github.com/icpd/subscribe2clash/app/api"
	"github.com/icpd/subscribe2clash/internal/randkey"
	"log"
	"math"
	"os"
)

func Run() {
	appKey := os.Getenv("APP_KEY")
	if len(appKey) == 0 {
		appKey = randkey.RandomKey()
	}
	api.EasKey = appKey

	log.Printf("appKey -> %v...\n", api.EasKey[0:int(math.Min(float64(16), float64(len(appKey))))])

	initFlag()
	log.Println("规则获取中...")
	generateConfig()
	log.Println("服务启动中...")
	initHttpServer()
}
