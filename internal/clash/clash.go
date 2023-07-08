package clash

import (
	"errors"
	"os"
	"strings"
	"unsafe"

	"github.com/icpd/subscribe2clash/internal/subscribe"
	"github.com/icpd/subscribe2clash/internal/xbase64"
)

type SourceType int

const (
	Url SourceType = iota
	File
	Txt
)

func Config(sourceType SourceType, source string, nodeOnly bool) (string, error) {
	var contents []string
	var err error

	switch sourceType {
	case Url:
		contents, err = subscribe.GetRawProxiesFromLinks(source)
		if err != nil {
			return "", err
		}

	case File:
		file, err := os.ReadFile(source)
		if err != nil {
			return "", err
		}
		contents = append(contents, subscribe.ParseRawProxies(unsafe.String(&file[0], len(file))))

	case Txt:
		contents = append(contents, strings.Split(source, ",")...)

	default:
		return "", errors.New("unknown source type")
	}

	proxies := subscribe.ParseProxy(contents)
	config, err := subscribe.GenerateClashConfig(proxies, nodeOnly)
	if err != nil {
		return "", err
	}

	return xbase64.UnicodeEmojiDecode(string(config)), nil
}
