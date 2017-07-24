package findMovies

import (
	"strings"
	"github.com/ddliu/go-httpclient"
	"github.com/PuerkitoBio/goquery"
	j "github.com/ricardolonga/jsongo"
)

var infoMap = map[string]string{
	"地区":   "region",
	"上映日期": "releasedAt",
	"更新日期": "updatedAt",
	"片长":   "timeLength",
	"语言":   "language",
	"类型":   "type",
	"导演":   "director",
	"豆瓣评分": "score",
}


func Find(keyword string, limit int) string {
	/**
		抓取电影列表
	 */
	result, _ := httpclient.Post("http://www.80s.tw/search", map[string]string{
		"keyword": keyword,
	})
	res, _ := goquery.NewDocumentFromResponse(result.Response)

	results := j.Array()

	resultsDom := res.Find(".search_list li")

	if resultsDom.Length() == 0 {
		return results.String()
	}
	resultsDom.Each(func(i int, s *goquery.Selection) {

		if i < limit {
			a := s.Find("a")
			title := Cleaar(a.Text())

			href, _ := a.Attr("href")
			url := "http://www.80s.tw" + href

			results.Put(FetchMovieItem(url, title))
		}

	})

	return results.String()
}

func FetchMovieItem(url string, title string) j.O {
	/**
		初始化 movieItem
	 */
	movieItem := j.Object()
	movieItem.Put("url", url).Put("title", title)

	res, _ := goquery.NewDocument(url)

	/**
		获取电影属性
	 */
	movieItem.Put("description", Cleaar(res.Find("#movie_content").Text()))
	infoList := res.Find("div[class=clearfix] span.span_block")
	infoLength := infoList.Length()
	infoList.Each(func(i int, infoItem *goquery.Selection) {
		nodes := infoItem.Find("span.font_888")
		if i < infoLength && nodes.Length() > 0 {
			attrName := strings.Replace(nodes.Text(), "：", "", -1)
			movieItem.Put(infoMap[attrName], Cleaar(infoItem.Text()))
		}
	})

	/**
		获取下载链接
	 */
	downloads := j.Array()
	downloadsDom := res.Find("div#cpdl2list li.dlurlelement")
	downloadsLength := downloadsDom.Length()
	downloadsDom.Each(func(i int, downloadDom *goquery.Selection) {
		if i != 0 && i != downloadsLength-1 {
			span := downloadDom.Find("span")
			span.Eq(0)
			download := j.Object().Put("title", Cleaar(span.Eq(0).Text()))
			url, _ = span.Eq(3).Find("a").Eq(0).Attr("href")
			download.Put("url", url)
			downloads.Put(download)
		}
	})

	movieItem.Put("downloads", downloads)

	return movieItem
}

func Cleaar(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	return str
}
