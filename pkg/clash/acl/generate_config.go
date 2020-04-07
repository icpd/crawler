package acl

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

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
	urls := GetUrls(option.origin, false)

	for _, g := range Sort {
		if u, ok := urls[g]; ok {
			resp, _ := req.HttpGet(u)
			s = append(s, AddProxyGroup(resp, Group[g]))
		}
	}

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
	dir := path.Dir(outputFile)
	if !Exists(dir) {
		mkDir(dir)
	}

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
			builder.WriteString(scanner.Text() + "\n")
		}
	}
	return builder.String()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func mkDir(path string) {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
