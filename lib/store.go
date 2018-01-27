package lib

import (
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
	return
}