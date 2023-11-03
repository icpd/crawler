package parser

import (
	"fmt"
	"regexp"
	"strconv"

	"imooc.com/ccmouse/learngo/crawler/engine"
	"imooc.com/ccmouse/learngo/crawler/model"
)

var priceReTmpl = `<a href="/%s/baojia/".*>(\d+\.\d+)</a>`

var nameRe = regexp.MustCompile(`<title>【(.*)报价_图片_参数】.*</title>`)
var carImageRe = regexp.MustCompile(`<img class="color_car_img_new" src="([^"]+)"`)
var sizeRe = regexp.MustCompile(`<li.*车身尺寸.*<em>(\d+[^\d]\d+[^\d]\d+mm)`)
var fuelRe = regexp.MustCompile(`<li.*工信部油耗.*<em>(\d+\.\d+)L/100km`)
var transmissionRe = regexp.MustCompile(`<li.*变\s*速\s*箱.*<em>(.+)</em>`)
var engineRe = regexp.MustCompile(`发\s*动\s*机.*\s*.*<.*>(\d+kW[^<]*)<`)
var displacementRe = regexp.MustCompile(`<li.*排.*量.*(\d+\.\d+)L`)
var maxSpeedRe = regexp.MustCompile(`<td.*最高车速\(km/h\).*\s*<td[^>]*>(\d+)</td>`)
var accelRe = regexp.MustCompile(`<td.*0-100加速时间\(s\).*\s*<td[^>]*>([\d\.]+)</td>`)
var urlRe = regexp.MustCompile(`http://newcar.xcar.com.cn/(m\d+)/`)

func ParseCarDetail(contents []byte, url string) engine.ParseResult {
	id := extractString([]byte(url), urlRe)

	car := model.Car{
		Name:         extractString(contents, nameRe),
		ImageURL:     "http:" + extractString(contents, carImageRe),
		Size:         extractString(contents, sizeRe),
		Fuel:         extractFloat(contents, fuelRe),
		Transmission: extractString(contents, transmissionRe),
		Engine:       extractString(contents, engineRe),
		Displacement: extractFloat(contents, displacementRe),
		MaxSpeed:     extractFloat(contents, maxSpeedRe),
		Acceleration: extractFloat(contents, accelRe),
	}
	priceRe, err := regexp.Compile(
		fmt.Sprintf(priceReTmpl, regexp.QuoteMeta(id)))
	if err == nil {
		car.Price = extractFloat(contents, priceRe)
	}

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Id:      id,
				Url:     url,
				Type:    "xcar",
				Payload: car,
			},
		},
	}

	carModelResult := ParseCarModel(contents, url)
	result.Requests = carModelResult.Requests

	return result
}

func extractString(
	contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

func extractFloat(contents []byte, re *regexp.Regexp) float64 {
	f, err := strconv.ParseFloat(extractString(contents, re), 64)
	if err != nil {
		return 0
	}
	return f
}
