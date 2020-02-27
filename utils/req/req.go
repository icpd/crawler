package req

import (
	"github.com/parnurzeal/gorequest"
	"time"
)

func HttpGet(url string) (string, error) {
	_, body, errs := gorequest.New().Get(url).Timeout(time.Minute).End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	return body, nil
}
