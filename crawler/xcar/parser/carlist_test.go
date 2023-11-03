package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCarList(t *testing.T) {
	contents, err := ioutil.ReadFile(
		"carlist_test_data.html")

	if err != nil {
		panic(err)
	}

	result := ParseCarList(contents, "")

	const resultSize = 30
	const carModelSize = 20
	expectedCarModelUrls := []string{
		"http://newcar.xcar.com.cn/4007/",
		"http://newcar.xcar.com.cn/52/",
		"http://newcar.xcar.com.cn/3875/",
	}
	expectedCarListUrls := []string{
		"http://newcar.xcar.com.cn/car/0-0-0-0-0-0-0-0-0-0-0-1/",
		"http://newcar.xcar.com.cn/car/0-0-0-0-0-0-0-0-0-0-2-1/",
		"http://newcar.xcar.com.cn/car/0-0-0-0-0-0-0-0-0-0-3-1/",
	}

	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d "+
			"requests; but had %d",
			resultSize, len(result.Requests))
	}
	for i, url := range expectedCarModelUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but "+
				"was %s",
				i, url, result.Requests[i].Url)
		}
	}
	for i, url := range expectedCarListUrls {
		if result.Requests[carModelSize+i].Url != url {
			t.Errorf("expected url #%d: %s; but "+
				"was %s",
				carModelSize+i, url, result.Requests[i].Url)
		}
	}
}
