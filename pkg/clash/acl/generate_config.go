package acl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

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
				s = append(s, AddProxyGroup(resp, ch["group"]))
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

	fileBytes, err := ioutil.ReadFile(option.baseFile)
	if err != nil {
		log.Fatal(err)
	}

	configStr := fmt.Sprintf(string(fileBytes), r)

	file, err := os.OpenFile(
		option.outputFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice := []byte(configStr)
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}
