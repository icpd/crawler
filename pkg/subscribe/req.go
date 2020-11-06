package subscribe

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/whoisix/subscribe2clash/pkg/mybase64"
	"github.com/whoisix/subscribe2clash/pkg/req"
)

func GetSubContent(query string) ([]string, error) {
	subLinks := strings.Split(query, ",")

	var contentSlice []string
	for _, link := range subLinks {
		content, err := req.HttpGet(link)
		if err != nil {
			return nil, err
		}
		content = strings.TrimSpace(content)

		if strings.HasPrefix(content, "ssd://") {
			content = content[6:]
			decodeBody, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(content)
			if err != nil {
				log.Println("ssd decode err:", err)
				continue
			}

			contentSlice = append(contentSlice, string(decodeBody))
			continue
		}

		decoded, err := mybase64.Base64DecodeStripped(content)
		if err != nil {
			log.Println("base64 decode err:", err)
			continue
		}
		contentSlice = append(contentSlice, string(decoded))
	}

	return contentSlice, nil
}
