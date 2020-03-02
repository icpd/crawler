package req

import (
	"github.com/parnurzeal/gorequest"
	"time"
)

var Proxy string

func HttpGet(url string) (string, error) {
	reqIns := gorequest.New().Get(url).Timeout(time.Minute)
	if Proxy != "" {
		reqIns = reqIns.Proxy(Proxy)
	}
	_, body, errs := reqIns.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	return body, nil
}
