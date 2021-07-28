package acl

import (
	"bufio"
	_ "embed"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"unsafe"

	"github.com/whoisix/subscribe2clash/internal/req"
	"github.com/whoisix/subscribe2clash/internal/subscribe"
)

//go:embed config/default_base_config.yaml
var defaultBaseConfig []byte

type genOption struct {
	baseFile   string
	outputFile string
}

type GenOption func(option *genOption)

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
		outputFile: "./config/acl.yaml",
	}

	for _, fn := range genOptions {
		fn(&option)
	}

	subscribe.RwMtx.Lock()
	subscribe.OutputFile = option.outputFile
	subscribe.RwMtx.Unlock()

	var s []string
	rules := GetRules()
	for _, r := range rules {
		log.Println(r.url, r.rule)
		resp, _ := req.HttpGet(r.url)
		s = append(s, AddProxyGroup(resp, r.rule))
	}

	r := MergeRule(s...)
	r = unique(r)

	var (
		configContent []byte
		err           error
	)
	if option.baseFile != "" {
		configContent, err = ioutil.ReadFile(option.baseFile)
		if err != nil {
			log.Fatal("读取基础配置文件失败", err)
		}
	} else {
		configContent = defaultBaseConfig
	}

	writeNewFile(configContent, option.outputFile, r)
}

func writeNewFile(configContent []byte, outputFile, filler string) {
	ctt := *(*string)(unsafe.Pointer(&configContent))
	tpl, err := template.New("config").Parse(ctt)
	if err != nil {
		log.Fatal("解析配置模版失败", err)
	}

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

	err = tpl.Execute(file, filler)
	if err != nil {
		log.Fatal(err)
	}
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
		return os.IsExist(err)
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
