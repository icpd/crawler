package acl

import (
	_ "embed"
	"log"

	"gopkg.in/ini.v1"
)

//go:embed config/default_base_urls.ini
var defaultRulesConfig string

type Rules struct {
	url   string
	rule string
}

func GetRules() []Rules {
	var rs []Rules

	cfg, err := ini.Load("pkg/acl/config/default_base_urls.ini")
	if err != nil {
		log.Fatal(err)
	}

	host := cfg.Section("").Key("host").String()
	for _, cfgK := range cfg.Section("rules").Keys() {
		rs = append(rs, Rules{
			url:   host + "/" + cfgK.Name(),
			rule: cfgK.Value(),
		})
	}

	return rs
}
