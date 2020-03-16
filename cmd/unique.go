package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/whoisix/subscribe2clash/utils/req"
)

var url string

func init() {
	flag.StringVar(&url, "url", "", "路由文件地址")
	flag.Parse()
}

func main() {

	if url == "" {
		flag.Usage()
		return
	}

	req.Proxy = "http://127.0.0.1:7890"
	content, err := req.HttpGet(url)
	if err != nil {
		log.Fatal(err)
	}

	var filterMap = make(map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(content))

	f, err := os.Create("./x.list")
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for scanner.Scan() {
		if scanner.Text() == "" {
			_, _ = fmt.Fprintln(w)
			continue
		}
		if _, ok := filterMap[scanner.Text()]; !ok {
			filterMap[scanner.Text()] = struct{}{}
			_, _ = fmt.Fprintln(w, scanner.Text())
		} else {
			fmt.Println(scanner.Text())
		}
	}

	_ = w.Flush()
}
