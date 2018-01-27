package lib

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
	"fmt"
)

type Section struct {
	Title string
	Url string
}

type Sections struct {
	Pics []string
	Section Section
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

	doc.Find("#list_com.articletxt3 li").Each(func (_ int, s *goquery.Selection) {
		band := s.Find("a")
		url, _ := band.Attr("href")
		if band.Length() > 0 {
			section = append(section, Section{
				Title: band.Text(),
				Url: url,
			})
		}
	})

	return
}

func GetPictureToSection(url string, title string, id int) (sections Sections) {
	fmt.Printf("🐹开始转换 -> %s     ", title)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	resultReader := transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder())
	
	document, _ := goquery.NewDocumentFromReader(resultReader)

	document.Find("[name=listNarImg]").First().Find("option").Each(func (_ int, tag *goquery.Selection) {
		nextPage, _ := tag.Attr("value")
		sections.Pics = append(sections.Pics, fmt.Sprintf("http://img.feiwan.net/qidazui/manhua/%d/%s.jpg\n", id, nextPage))
	})
	sections.Section = Section{
		Title: title,
		Url: url,
	}
	fmt.Println("转换完成✨")	
	return
}