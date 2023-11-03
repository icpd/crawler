// Package profile implements profile generator.
package profile

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"imooc.com/ccmouse/learngo/crawler/model"
	"imooc.com/ccmouse/learngo/mockserver/config"
)

// Recommendation defines the interface for recommendation subsystem.
type Recommendation interface {
	NextGuess(id int64) []int64
}

// Generator represents the profile generator.
type Generator struct {
	Tmpl           *template.Template
	Recommendation Recommendation
}

// HandleRequest is the gin request handler for profile generation.
func (g *Generator) HandleRequest(c *gin.Context) {
	params := struct {
		ID int64 `uri:"id" binding:"required"`
	}{}
	err := c.BindUri(&params)
	if err != nil {
		log.Printf("BindUri(): %v.", err)
		return
	}

	err = g.generate(params.ID, c.Writer)
	if err != nil {
		log.Printf("Cannot generate profile for user %d: %v.", params.ID, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// GuessListItem defines guess list item bound into html template.
type GuessListItem struct {
	URL      string
	Name     string
	PhotoURL string
}

// PhotoProfile defiens profile with photo bound into html template.
type PhotoProfile struct {
	// Embed a *model.Profile here to save typing in template.
	*model.Profile
	PhotoURL string
}

// Content defines the master content bound into html template.
type Content struct {
	*PhotoProfile
	GuessList []GuessListItem
}

func (g *Generator) generate(id int64, w io.Writer) error {
	p := g.GenerateProfile(id)
	var guessItems []GuessListItem
	if g.Recommendation != nil {
		guesses := g.Recommendation.NextGuess(id)
		guessItems = make([]GuessListItem, len(guesses))
		for i, v := range guesses {
			gp := g.GenerateProfile(v)
			guessItems[i].URL = fmt.Sprintf("http://%s/mock/album.zhenai.com/u/%d", config.ServerAddress, v)
			guessItems[i].Name = gp.Name
			guessItems[i].PhotoURL = gp.PhotoURL
		}
	}
	return g.Tmpl.Execute(w, Content{
		PhotoProfile: p,
		GuessList:    guessItems,
	})
}

// GenerateProfile generates a photo profile given a user id.
func (g *Generator) GenerateProfile(id int64) *PhotoProfile {
	r := rand.New(rand.NewSource(id))
	n1 := elementFromSlice(r, []string{
		"断念",
		"街痞",
		"浪瘾",
		"猖狂",
		"归戾",
		"洒脱",
		"故我",
		"简白",
		"酒虑",
		"野性入骨",
		"学霸的芯",
		"一笑奈何",
		"隐匿微笑",
		"限时拥抱",
		"何必怀念",
		"厌与深情",
		"寂寞成影",
		"桀骜不驯",
		"浪痞孤王",
		"高冷绅士",
		"心事痕迹",
		"逍遥浪子",
		"隐身守候",
		"执手不离",
		"独久厌闹",
		"逆天飞翔",
		"醉生忧愁",
		"全球焦點",
		"孤者何惧",
		"君临天下",
		"无戏配角",
		"一身傲气",
		"冷暖自知",
		"海阔天空",
		"今朝醉者",
		"不良骚年",
		"薄情寡义",
		"称霸全服",
		"久别无恙°",
		"與我無關°",
		"逍遥べ无痕",
		"混過也愛過",
		"与你度余生",
		"一见不钟情",
		"孤独比酒暖°",
		"原来无话可说",
		"陪你浪迹天涯",
	})
	n2 := elementFromSlice(r, []string{
		"萌宝",
		"丁当",
		"遇到",
		"病娇",
		"莓哒",
		"深碍",
		"面码",
		"娇喘",
		"爱你",
		"迁就",
		"丸子",
		"初心",
		"baby",
		"猫儿.",
		"傻蛋.",
		"呆萌i",
		"如你*",
		"肉嘟嘟",
		"小气鬼",
		"胆小鬼",
		"欣欣雨",
		"酷酷猫",
		"柠小萌",
		"小丸子",
		"小可爱",
		"小仙女",
		"考莎啦",
		"晴妹儿.",
		"记得笑i",
		"稚性女",
		"小短腿i",
		"万能萌妹",
		"蓝莓布朗",
		"天瓓少女",
		"草莓裙摆",
	})
	p := &model.Profile{
		Name:   n1 + n2,
		Gender: elementFromSlice(r, []string{"男", "女"}),
		Age:    r.Intn(100),
		Height: r.Intn(200),
		Weight: r.Intn(300),
		Income: elementFromSlice(r, []string{
			"1-2000元",
			"2001-3000元",
			"3001-5000元",
			"5001-8000元",
			"8001-10000元",
			"10001-20000元",
			"财务自由",
		}),
		Marriage: elementFromSlice(r, []string{"未婚", "离异"}),
		Education: elementFromSlice(r, []string{
			"小学",
			"初中",
			"高中",
			"大学",
			"硕士",
			"博士及以上",
		}),
		Occupation: elementFromSlice(r, []string{
			"人事/行政",
			"程序员",
			"产品经理",
			"测试工程师",
			"财务",
			"总经理",
			"金融",
			"销售",
			"其它",
		}),
		Hokou: elementFromSlice(r, []string{
			"北京市", "上海市", "广州市", "深圳市",
			"成都市", "杭州市", "武汉市", "重庆市", "南京市", "天津市", "苏州市", "西安市", "长沙市", "沈阳市", "青岛市", "郑州市", "大连市", "东莞市", "宁波市",
			"其它",
		}),
		Xinzuo: elementFromSlice(r, []string{
			"白羊座",
			"金牛座",
			"双子座",
			"巨蟹座",
			"狮子座",
			"处女座",
			"天秤座",
			"天蝎座",
			"射手座",
			"魔羯座",
			"水瓶座",
			"双鱼座",
		}),
		House: elementFromSlice(r, []string{
			"有房",
			"租房",
			"无房",
		}),
		Car: elementFromSlice(r, []string{
			"无车",
			"有车",
			"有豪车",
		}),
	}

	return &PhotoProfile{
		Profile:  p,
		PhotoURL: fmt.Sprintf("https://picsum.photos/seed/%d/300/300", r.Intn(100000)),
	}
}

func elementFromSlice(r *rand.Rand, s []string) string {
	return s[r.Intn(len(s))]
}
