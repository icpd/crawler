package acl

import (
	"fmt"
	"regexp"
	"strings"
)

var itemReg = regexp.MustCompile(`(?m)([\w-]+,[\w./:-]*[^\x00-\xff]*)(,no-resolve)?`)
var comReg = regexp.MustCompile(`(?m)#.+`)

func AddProxyGroup(debris, groupName string) string {
	return comReg.ReplaceAllString(itemReg.ReplaceAllString(debris, fmt.Sprintf("  - $1,%s$2", groupName)), "  $0")
}

func MergeRule(debris ...string) string {
	return strings.Join(debris, strings.Repeat("\r\n", 2))
}
