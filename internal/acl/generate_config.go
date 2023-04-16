package acl

import (
	"bufio"
	_ "embed"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/icpd/subscribe2clash/internal/req"
)

//go:embed config/default_base_config.yaml
var defaultBaseConfig []byte

var GlobalGen *Gen

type Gen struct {
	OutputFile string
	baseFile   string

	configContent []byte
	rule          string
}

func (g *Gen) GenerateConfig() {
	var s []string
	rules := GetRules()
	for _, r := range rules {
		log.Println(r.url, r.rule)
		resp, err := req.HttpGet(r.url)
		if err != nil {
			log.Printf("获取规则失败: %s, err:%v", r.url, err)
		}

		s = append(s, AddProxyGroup(resp, r.rule))
	}

	r := MergeRule(s...)
	r = unique(r)
	g.rule = r

	if g.baseFile != "" {
		var err error
		g.configContent, err = os.ReadFile(g.baseFile)
		if err != nil {
			log.Fatal("读取基础配置文件失败", err)
		}
	} else {
		g.configContent = defaultBaseConfig
	}

	g.writeNewFile()
}

func (g *Gen) writeNewFile() {
	tpl, err := template.New("config").Parse(string(g.configContent))
	if err != nil {
		log.Fatal("解析配置模版失败", err)
	}

	dir := path.Dir(g.OutputFile)
	if !Exists(dir) {
		mkDir(dir)
	}

	file, err := os.OpenFile(
		g.OutputFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	err = tpl.Execute(file, g.rule)
	if err != nil {
		log.Fatal(err)
	}
}

type GenOption func(option *Gen)

func WithBaseFile(filepath string) GenOption {
	return func(option *Gen) {
		option.baseFile = filepath
	}
}

func WithOutputFile(filepath string) GenOption {
	return func(option *Gen) {
		option.OutputFile = filepath
	}
}

func unique(rules string) string {
	var filterMap = make(map[string]any)
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

func New(ops ...GenOption) *Gen {
	gen := Gen{
		OutputFile: "./config/acl.yaml",
	}

	for _, fn := range ops {
		fn(&gen)
	}
	GlobalGen = &gen

	return GlobalGen
}
