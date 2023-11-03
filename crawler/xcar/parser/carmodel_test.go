package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCarModel(t *testing.T) {
	contents, err := ioutil.ReadFile(
		"carmodel_test_data.html")

	if err != nil {
		panic(err)
	}

	result := ParseCarModel(contents, "")

	const resultSize = 3
	expectedUrls := []string{
		"http://newcar.xcar.com.cn/m37326/",
		"http://newcar.xcar.com.cn/m35001/",
		"http://newcar.xcar.com.cn/m35002/",
	}

	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d "+
			"requests; but had %d",
			resultSize, len(result.Requests))
	}
	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but "+
				"was %s",
				i, url, result.Requests[i].Url)
		}
	}
}
