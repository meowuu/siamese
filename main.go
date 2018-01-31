package main

import (
	"fmt"

	"github.com/meowuu/siamese/lib/leancloud"
)

func main() {
	// name := "七大罪"
	// url := "http://qidazui.feiwan.net/manhua"

	// arr := lib.GetPage(url)

	// sections := lib.Stretch(arr)

	// lib.SaveBook(name, url, sections)
	fmt.Println(leancloud.GetClient())
}
