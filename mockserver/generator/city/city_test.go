package city

import (
	"bytes"
	"html/template"
	"testing"

	"imooc.com/ccmouse/learngo/crawler/zhenai/parser"
	"imooc.com/ccmouse/learngo/mockserver/config"
	"imooc.com/ccmouse/learngo/mockserver/generator/profile"
)

func TestGenerate(t *testing.T) {
	config.ServerAddress = "localhost:8080"
	pg := &profile.Generator{Tmpl: template.Must(template.ParseFiles("../profile/profile_tmpl.html"))}
	g := Generator{
		Tmpl:       template.Must(template.ParseFiles("city_tmpl.html")),
		ProfileGen: pg,
	}

	var b bytes.Buffer
	err := g.generate(params{
		City: "fuxin",
		Page: 34,
	}, &b)

	if err != nil {
		t.Fatalf("Cannot generate content: %v.", err)
	}

	r := parser.ParseCity(b.Bytes(), "")

	wantItems, wantRequests := 0, 24
	if len(r.Items) != wantItems {
		t.Errorf("generate() want %d items, got %d: %v", wantItems, len(r.Items), r.Items)
	}

	if len(r.Requests) != wantRequests {
		t.Errorf("generate() want %d requests, got %d: %v", wantRequests, len(r.Requests), r.Requests)
	}

	verify := []struct {
		i          int
		wantURL    string
		wantParser string
		wantArg    interface{}
	}{
		{
			i:          0,
			wantURL:    "http://localhost:8080/mock/album.zhenai.com/u/484971159322053275",
			wantParser: "ParseProfile",
			wantArg:    "与你度余生迁就",
		},
		{
			i:          23,
			wantURL:    "http://localhost:8080/mock/www.zhenai.com/zhenghun/fuxin/37",
			wantParser: "ParseCity",
		},
	}

	for _, v := range verify {
		gotURL := r.Requests[v.i].Url
		gotParser, gotArg := r.Requests[v.i].Parser.Serialize()
		if v.wantURL != gotURL || v.wantParser != gotParser || v.wantArg != gotArg {
			t.Errorf("generate() want %d-th request (%q, %q, %v), got (%q, %q, %v)",
				v.i, v.wantURL, v.wantParser, v.wantArg, gotURL, gotParser, gotArg)
		}
	}
}
