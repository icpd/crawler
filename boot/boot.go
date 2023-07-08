package boot

import (
	"github.com/icpd/subscribe2clash/app/api"
	"github.com/icpd/subscribe2clash/internal/randkey"
	"log"
	"math"
	"os"
)

var (
	keyLen int = 16
)

func Run() {
	appKey := os.Getenv("APP_KEY")
	appKeyLen := len(appKey)
	if appKeyLen == 0 {
		appKey = randkey.GenerateRandomString(keyLen)
	} else {
		if appKeyLen > keyLen {
			appKey = appKey[0:keyLen]
		}
		if appKeyLen < keyLen {
			appKey = appKey + randkey.GenerateRandomString(keyLen-appKeyLen)
		}
	}
	api.EasKey = appKey

	log.Printf("appKey -> %v...\n", api.EasKey[0:int(math.Min(float64(keyLen/2), float64(len(appKey))))])

	initFlag()
	log.Println("规则获取中...")
	generateConfig()
	log.Println("服务启动中...")
	initHttpServer()
}
