package main

import (
	"log"

	"github.com/whoisix/subscribe2clash/boot"
)

func main() {
	if err := boot.Boot(); err != nil {
		log.Fatalln("启动失败：", err)
	}
}
