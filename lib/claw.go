package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Section struct {
	Pics  []string
	Title string
	Url   string
	IdNum int
	Index int
}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns: 10,
	},
}

// GetPage is get book info from url
func GetPage(url string) (section []Section) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	resultReader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())

	doc, err := goquery.NewDocumentFromReader(resultReader)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#list_com.articletxt3 li").Each(func(_ int, s *goquery.Selection) {
		band := s.Find("a")
		url, _ := band.Attr("href")
		if band.Length() > 0 {
			section = append(section, Section{
				Title: band.Text(),
				Url:   url,
			})
		}
	})

	return
}

func GetPictureToSection(url string, title string, id int, index int, c chan Section) {
	var section Section

	res, err := client.Get(url)

	if err != nil {
		fmt.Printf("获取-> %s <-数据出现错误", title)
		log.Fatal(err)
	}
	defer res.Body.Close()
	resultReader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())

	valid := regexp.MustCompile("http://img.feiwan.net/qidazui/manhua/\\w+/\\d+.(jpg|png)")
	data, _ := ioutil.ReadAll(resultReader)
	for _, match := range valid.FindAllString(string(data), -1) {
		section.Pics = append(section.Pics, match)
	}

	section.Title = title
	section.Url = url
	section.IdNum = id
	section.Index = index

	c <- section
}

func Stretch(arr []Section) (sections []Section) {
	c := make(chan Section)
	fmt.Println("🐣 开始获取章节内容")

	for i, item := range arr {
		valid := regexp.MustCompile("/manhua/(\\d+).html")
		regstr := valid.FindString(item.Url)

		id, _ := strconv.Atoi(valid.ReplaceAllString(regstr, "$1"))

		go GetPictureToSection(item.Url, item.Title, id, i, c)
	}

	index := 0
	for i := range c {
		sections = append(
			sections,
			i,
		)

		index++

		if index == len(arr) {
			fmt.Println("获取完成 ✨")
			return
		}
	}
	return
}
