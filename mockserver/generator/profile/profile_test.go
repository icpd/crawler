package profile

import (
	"bytes"
	"html/template"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"imooc.com/ccmouse/learngo/crawler/engine"
	"imooc.com/ccmouse/learngo/crawler/model"
	"imooc.com/ccmouse/learngo/crawler/zhenai/parser"
	"imooc.com/ccmouse/learngo/mockserver/config"
	"imooc.com/ccmouse/learngo/mockserver/recommendation"
)

func TestGenerate(t *testing.T) {
	config.ServerAddress = "localhost:8080"
	g := Generator{
		Tmpl:           template.Must(template.ParseFiles("profile_tmpl.html")),
		Recommendation: recommendation.Client{},
	}

	rand.Seed(34534)
	var b bytes.Buffer
	err := g.generate(12345, &b)

	if err != nil {
		t.Fatalf("Cannot generate content: %v.", err)
	}

	want := engine.Item{
		Url:  "http://localhost:8080/mock/album.zhenai.com/u/12345",
		Type: "zhenai",
		Id:   "12345",
		Payload: model.Profile{
			Name:       "逍遥べ无痕病娇",
			Gender:     "男",
			Age:        36,
			Height:     41,
			Weight:     275,
			Income:     "8001-10000元",
			Marriage:   "未婚",
			Education:  "初中",
			Occupation: "人事/行政",
			Hokou:      "成都市",
			Xinzuo:     "魔羯座",
			House:      "有房",
			Car:        "有豪车",
		},
	}
	r := parser.NewProfileParser("逍遥べ无痕病娇").Parse(b.Bytes(), "http://localhost:8080/mock/album.zhenai.com/u/12345")
	if len(r.Items) != 1 {
		t.Errorf("want exactly 1 element, got %d items: %v", len(r.Items), r.Items)
	} else {
		if diff := cmp.Diff(want, r.Items[0]); diff != "" {
			t.Errorf("generated data is incorrect: diff: -want +got\n%s", diff)
		}
	}

	type reqData struct {
		URL    string
		Parser string
		Arg    interface{}
	}
	wantReq := []reqData{
		{
			URL:    "http://localhost:8080/mock/album.zhenai.com/u/12335",
			Parser: "ParseProfile",
			Arg:    "酒虑肉嘟嘟",
		},
		{
			URL:    "http://localhost:8080/mock/album.zhenai.com/u/12340",
			Parser: "ParseProfile",
			Arg:    "称霸全服胆小鬼",
		},
	}
	var gotReq []reqData
	for _, req := range r.Requests {
		parser, arg := req.Parser.Serialize()
		gotReq = append(gotReq, reqData{
			URL:    req.Url,
			Parser: parser,
			Arg:    arg,
		})
	}

	if diff := cmp.Diff(wantReq, gotReq); diff != "" {
		t.Errorf("generated request is incorrect: diff: -want +got\n%s", diff)
	}
}
