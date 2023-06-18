package subscribe

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"

	"github.com/icpd/subscribe2clash/internal/acl"
)

func (c *Clash) LoadTemplate() []byte {
	if !c.NodeOnly() {
		_, err := os.Stat(c.path)
		if err != nil && os.IsNotExist(err) {
			log.Printf("[%s] template doesn't exist. err: %v", c.path, err)
			return nil
		}
		buf, err := os.ReadFile(c.path)
		if err != nil {
			log.Printf("[%s] template open the failure. err: %v", c.path, err)
			return nil
		}
		err = yaml.Unmarshal(buf, &c)
		if err != nil {
			log.Printf("[%s] Template format error. err: %v", c.path, err)
		}
	}

	var proxiesName []string
	names := map[string]int{}

	for _, proto := range c.rawProxies {
		proxy := make(map[string]any)
		p, ok := proto.(map[string]any)
		if ok {
			proxy = p
		} else {
			j, err := json.Marshal(proto)
			if err != nil {
				log.Printf("json marshal error: %v", err)
				continue
			}
			_ = json.Unmarshal(j, &proxy)
		}

		var name string
		switch reflect.TypeOf(proto).Kind() {
		case reflect.Struct:
			name = reflect.ValueOf(proto).FieldByName("Name").String()
		case reflect.Map:
			if v, ok := proto.(map[string]any); ok {
				name = v["name"].(string)
			}
		}

		if index, ok := names[name]; ok {
			names[name] = index + 1
			name = fmt.Sprintf("%s%d", name, index+1)
		} else {
			names[name] = 0
		}

		proxy["name"] = name
		c.Proxies = append(c.Proxies, proxy)
		proxiesName = append(proxiesName, name)
	}

	for _, group := range c.ProxyGroups {
		groupProxies, ok := group["proxies"].([]any)
		if !ok {
			continue
		}
		for i, p := range groupProxies {
			if p == "1" {
				groupProxies = groupProxies[:i]

				var groupProxiesName []string
				for _, s := range groupProxies {
					groupProxiesName = append(groupProxiesName, s.(string))
				}
				groupProxiesName = append(groupProxiesName, proxiesName...)

				group["proxies"] = groupProxiesName
				break
			}
		}

	}

	var yamlOut any = c
	if c.NodeOnly() {
		type nodeOnly struct {
			Proxies []map[string]any `yaml:"proxies"`
		}

		yamlOut = nodeOnly{
			Proxies: c.Proxies,
		}
	}

	d, err := yaml.Marshal(yamlOut)
	if err != nil {
		return nil
	}

	return d
}

func (c *Clash) NodeOnly() bool {
	return c.nodeOnly
}

func GenerateClashConfig(proxies []any, nodeOnly bool) ([]byte, error) {
	clash := Clash{
		path:       acl.GlobalGen.OutputFile,
		rawProxies: proxies,
		nodeOnly:   nodeOnly,
	}

	r := clash.LoadTemplate()
	if r == nil {
		return nil, fmt.Errorf("sublink 返回数据格式不对")
	}
	return r, nil
}
