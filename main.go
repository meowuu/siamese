package main

import (
	"github.com/meowuu/siamese/lib"
)

func main() {
	name := "七大罪"
	url := "http://qidazui.feiwan.net/manhua"

	arr := lib.GetPage(url)

	sections := lib.Stretch(arr)

	lib.SaveBook(name, url, sections)
}