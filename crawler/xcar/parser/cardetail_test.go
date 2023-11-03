package parser

import (
	"io/ioutil"
	"testing"

	"imooc.com/ccmouse/learngo/crawler/engine"
	"imooc.com/ccmouse/learngo/crawler/model"
)

func TestParseCarDetail(t *testing.T) {
	contents, err := ioutil.ReadFile(
		"cardetail_test_data.html")

	if err != nil {
		panic(err)
	}

	expectedItem := engine.Item{
		Url:  "http://newcar.xcar.com.cn/m35001/",
		Type: "xcar",
		Id:   "m35001",
		Payload: model.Car{
			Name:         "奥迪TT双门2017款45 TFSI",
			Price:        47.18,
			ImageURL:     "http://img1.xcarimg.com/b63/s8386/m_20170616000036181753843373443.jpg-280x210.jpg",
			Size:         "4191×1832×1353mm",
			Fuel:         16.7,
			Transmission: "6挡双离合",
			Engine:       "169kW(2.0L涡轮增压)",
			Displacement: 2,
			MaxSpeed:     250,
			Acceleration: 5.9,
		},
	}

	result := ParseCarDetail(contents, "http://newcar.xcar.com.cn/m35001/")

	if len(result.Items) != 1 {
		t.Errorf("Must only return one item, but %d returned",
			len(result.Items))
	}

	actualItem := result.Items[0]
	if actualItem != expectedItem {
		t.Errorf("expected item: %+v, but was %+v",
			expectedItem, actualItem)
	}

	const resultSize = 8
	expectedUrls := []string{
		"http://newcar.xcar.com.cn/m45776/",
		"http://newcar.xcar.com.cn/m45776/",
		"http://newcar.xcar.com.cn/m32946/",
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
