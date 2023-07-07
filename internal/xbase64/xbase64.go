package xbase64

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"
)

var (
	// emoji表情的数据表达式
	re = regexp.MustCompile(`(?i)\\\\u[0-9a-zA-Z]+`)
	// 提取emoji数据表达式
	reg = regexp.MustCompile(`(?i)\\\\u`)
)

func UnicodeEmojiDecode(s string) string {
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

func Base64DecodeStripped(s string) ([]byte, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	s = strings.ReplaceAll(s, " ", "+")
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(s)
	}
	return decoded, err
}

func Base64EncodeStripped(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
