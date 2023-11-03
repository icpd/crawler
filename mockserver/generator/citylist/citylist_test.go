package citylist

import (
	"bytes"
	"html/template"
	"testing"

	"imooc.com/ccmouse/learngo/crawler/zhenai/parser"
	"imooc.com/ccmouse/learngo/mockserver/config"
)

func TestGenerate(t *testing.T) {
	config.ServerAddress = "localhost:8080"
	g := Generator{
		Tmpl: template.Must(template.ParseFiles("citylist_tmpl.html")),
	}

	var b bytes.Buffer
	err := g.generate(&b)

	if err != nil {
		t.Fatalf("Cannot generate content: %v.", err)
	}

	r := parser.ParseCityList(b.Bytes(), "")

	wantRequests := 470

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
			wantURL:    "http://localhost:8080/mock/www.zhenai.com/zhenghun/aba",
			wantParser: "ParseCity",
		},
		{
			i:          23,
			wantURL:    "http://localhost:8080/mock/www.zhenai.com/zhenghun/baotou",
			wantParser: "ParseCity",
		},
		{
			i:          469,
			wantURL:    "http://localhost:8080/mock/www.zhenai.com/zhenghun/zunyi",
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
