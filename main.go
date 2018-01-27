package main

import (
	"strconv"
	"regexp"
	"github.com/meowuu/siamese/lib"
)

func main() {
	name := "七大罪"
	url := "http://qidazui.feiwan.net/manhua"

	arr := lib.GetPage(url)
	var sections []lib.Sections
	for _, item := range arr {
		valid := regexp.MustCompile("/manhua/(\\d+).html")
		regstr := valid.FindString(item.Url)
	
		id, _ := strconv.Atoi(valid.ReplaceAllString(regstr, "$1"))

		section := lib.GetPictureToSection(item.Url, item.Title, id)
		sections = append(sections, section)
	}
	lib.SaveBook(name, url, sections)
}