package boot

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/icpd/subscribe2clash/app/api"
	"log"
)

func Run() {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return
	}
	api.EasKey = hex.EncodeToString(key)
	log.Printf("appKey -> %v...\n", api.EasKey[0:16])

	initFlag()
	log.Println("规则获取中...")
	generateConfig()
	log.Println("服务启动中...")
	initHttpServer()
}
