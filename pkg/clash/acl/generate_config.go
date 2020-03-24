package acl

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/whoisix/subscribe2clash/pkg/clash/subscribe"
	"github.com/whoisix/subscribe2clash/utils/req"
)

type genOption struct {
	origin     string
	baseFile   string
	outputFile string
}

type GenOption func(option *genOption)

func WithOrigin(s string) GenOption {
	return func(option *genOption) {
		option.origin = s
	}
}

func WithBaseFile(filepath string) GenOption {
	return func(option *genOption) {
		option.baseFile = filepath
	}
}

func WithOutputFile(filepath string) GenOption {
	return func(option *genOption) {
		option.outputFile = filepath
	}
}

func GenerateConfig(genOptions ...GenOption) {
	option := genOption{
		origin:     "github",
		baseFile:   "./config/clash/base_clash.yaml",
		outputFile: "./config/clash/acl.yaml",
	}

	for _, fn := range genOptions {
		fn(&option)
	}

	subscribe.OutputFile = option.outputFile

	var s []string
	var wg sync.WaitGroup
	urls := GetUrls(option.origin, false)
	urlCh := make(chan map[string]string)

	workerCount := len(urls)
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for ch := range urlCh {
				resp, _ := req.HttpGet(ch["url"])
				s = append(s, AddProxyGroup(resp, Group[ch["group"]]))
			}
		}()
	}

	for k, url := range urls {
		urlCh <- map[string]string{
			"url":   url,
			"group": k,
		}
	}
	close(urlCh)
	wg.Wait()
	r := MergeRule(s...)
	r = unique(r)

	writeNewFile(option.baseFile, option.outputFile, r)
}

func writeNewFile(baseFile, outputFile, filler string) {
	fileBytes, err := ioutil.ReadFile(baseFile)
	if err != nil {
		log.Fatal(err)
	}

	configStr := fmt.Sprintf(string(fileBytes), filler)

	writeFile(outputFile, configStr)
}

func writeFile(outputFile string, content string) {
	file, err := os.OpenFile(
		outputFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice := []byte(content)
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}

func unique(rules string) string {
	var filterMap = make(map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(rules))

	var builder strings.Builder

	for scanner.Scan() {
		if scanner.Text() == "" {
			builder.WriteString("\n")
			continue
		}
		if _, ok := filterMap[scanner.Text()]; !ok {
			filterMap[scanner.Text()] = struct{}{}
			builder.WriteString(scanner.Text()+"\n")
		}
	}
	return builder.String()
}
