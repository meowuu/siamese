package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Section struct {
	Title string
	Url   string
}

type Sections struct {
	Pics    []string
	Section Section
	IdNum   int
	Index   int
}

type Datas []Sections

func (d Datas) Len() int           { return len(d) }
func (d Datas) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Datas) Less(i, j int) bool { return d[i].Index < d[j].Index }

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

func GetPictureToSection(url string, title string, id int, index int, c chan Sections) {
	var sections Sections

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	resultReader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())

	valid := regexp.MustCompile("http://img.feiwan.net/qidazui/manhua/\\w+/\\d+.(jpg|png)")
	data, _ := ioutil.ReadAll(resultReader)
	for _, match := range valid.FindAllString(string(data), -1) {
		sections.Pics = append(sections.Pics, match)
	}

	sections.Section = Section{
		Title: title,
		Url:   url,
	}

	sections.IdNum = id
	sections.Index = index

	c <- sections
}

func Stretch(arr []Section) (sections Datas) {
	c := make(chan Sections)
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
			sort.Sort(sections)

			fmt.Println("获取完成 ✨")
			return
		}
	}
	return
}
