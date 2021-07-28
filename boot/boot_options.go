package boot

import (
	"github.com/whoisix/subscribe2clash/pkg/acl"
	"github.com/whoisix/subscribe2clash/pkg/global"
)

func Options() []acl.GenOption {
	var options []acl.GenOption
	if global.BaseFile != "" {
		options = append(options, acl.WithBaseFile(global.BaseFile))
	}
	if global.OutputFile != "" {
		options = append(options, acl.WithOutputFile(global.OutputFile))
	}
	return options
}
