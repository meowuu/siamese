package lib

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
	"fmt"
)

func GetPage() {
	res, err := http.Get("http://qidazui.feiwan.net/manhua")
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
			fmt.Printf("标题: %s 地址:%s\n", band.Text(), url)
		}
		// if band != "" {
		// }
	})
}