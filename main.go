package main

import (
	"github.com/meowuu/siamese/lib"
	"github.com/meowuu/siamese/lib/leancloud"
)

func main() {
	name := "七大罪"
	url := "http://qidazui.feiwan.net/manhua"

	arr := lib.GetPage(url)

	sections := lib.Stretch(arr)
	var leanSections []leancloud.Section

	for _, section := range sections {
		leanSections = append(leanSections, leancloud.Section{
			Name:   section.Title,
			Url:    section.Url,
			Images: section.Pics,
			ID:     section.IdNum,
			Index:  section.Index,
		})
	}

	client, _ := leancloud.GetClient()
	client.Save(leancloud.Book{
		Name: name,
		Url:  url,
	}, leanSections)
}
