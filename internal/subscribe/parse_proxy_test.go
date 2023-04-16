package subscribe

import (
	"fmt"
	"testing"
)

func Test_hysteriaConf(t *testing.T) {
	s := `hysteria://host:123?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarkse`
	conf := hysteriaConf(s)
	fmt.Printf("%+v\n", conf)
}
