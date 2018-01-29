package lib

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/meowuu/siamese/models"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:159357@/cartoon")
}

// SaveBook save book to database
func SaveBook(name string, url string, sections []Sections) (status int64, err error) {
	fmt.Println("💼 开始写入到数据库")

	o := orm.NewOrm()
	o.Using("default")

	book := new(models.Books)
	book.Name = name
	book.Url = url

	o.Insert(book)
	
	for _, sectiondata := range sections {
		section := &models.Section{
			Name: sectiondata.Section.Title,
			Bookid: book,
		}
		o.Insert(section)

		for _, picture := range sectiondata.Pics {
			picture := &models.Picture{
				Url: picture,
				Secid: section,
			}
			o.Insert(picture)
		}
	}

	fmt.Println("写入完成 🐾")
	return
}