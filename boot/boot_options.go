package boot

import (
	"fmt"
	"os"

	"github.com/whoisix/subscribe2clash/library/acl"
	"github.com/whoisix/subscribe2clash/library/global"
)

func Options() []acl.GenOption {
	var options []acl.GenOption
	if global.Origin != "" {
		switch global.Origin {
		case "cn", "github":
		default:
			fmt.Println("the origin argument can only be github or cn")
			os.Exit(0)
		}
		options = append(options, acl.WithOrigin(global.Origin))
	}
	if global.BaseFile != "" {
		options = append(options, acl.WithBaseFile(global.BaseFile))
	}
	if global.OutputFile != "" {
		options = append(options, acl.WithOutputFile(global.OutputFile))
	}
	return options
}
