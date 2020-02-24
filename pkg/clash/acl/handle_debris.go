package acl

import (
	"fmt"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`(?m)^[^#\s]+`)

func AddProxyGroup(debris, groupName string) string {
	return reg.ReplaceAllString(debris, fmt.Sprintf("- $0,%s", groupName))
}

func MergeRule(debris ...string) string {
	return strings.Join(debris, strings.Repeat("\r\n", 2))
}
