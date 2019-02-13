package main

import (
	"demoApp/server/model/dbModel"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

var testDB *gorm.DB

func loadDB() {

	db, err := gorm.Open("postgres",
		"host=localhost user=wxadmin dbname=app sslmode=disable password=wxadmin")

	if err != nil {
		fmt.Println(err)
	}

	testDB = db
	testDB.LogMode(true)
	//AddDepartment()
	testDB.SingularTable(true)
	testDB.Set("gorm:table_options", "charset=utf8")
	testDB.DB().SetMaxIdleConns(3)
	testDB.DB().SetMaxIdleConns(1)

}

func createJobs() {

	for i := 0; i < 5; i++ {
		item := dbModel.CompuseJobs{}
		item.Id = "ffewfwef" + strconv.Itoa(i)
		item.Name = "实习大豆纤维多"
		item.CompanyID = "companyID"
		item.IconURL = "http://icons.iconarchive.com/icons/blackvariant/button-ui-system-folders-drives/1024/Developer-icon.png"
		item.LocationCity = []string{"city1", "city2", "city3"}
		item.ReviewCounts = 100
		item.Education = "大专"
		item.Type = "intern"

		err := testDB.FirstOrCreate(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func creatCarrerTalk() {
	for i := 0; i < 12; i++ {
		item := dbModel.CareerTalk{}
		item.Id = "宣讲会当前为多" + strconv.Itoa(i)
		item.Name = "宣讲会当前为多"
		item.CompanyID = "companyID"
		item.IconURL = "https://cdn1.iconfinder.com/data/icons/education-icons-3/155/vector_313_21-512.png"
		item.College = "我的大学"
		item.SimplifyAddress = "某个楼"
		t := time.Now()
		item.StartTime = &t
		item.ContentType = dbModel.Text

		err := testDB.Create(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createNews() {
	for i := 0; i < 6; i++ {
		item := dbModel.NewsModel{}
		var t = time.Now()
		item.Time = &t
		item.Icon = "http://pic34.photophoto.cn/20150128/0007020160374237_b.jpg"
		item.Link = "http://gorm.book.jasperxu.com/database.html"
		item.Title = "gorm 中文文档"
		item.Author = "github"

		testDB.Create(&item)
	}
}

func main() {

	loadDB()

	//var company dbModel.Company
	// find talks
	//c := &dbModel.CarrerTalk{}
	//c.Id = "dwqdqwd"
	//c.CompanyID = "companyID"
	//err := testDB.Create(&c).Error
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//c := dbModel.Company{}
	//c.Id = "12345"
	//c.Name = "年你啊你按"
	//err := testDB.Model(&dbModel.Company{}).Create(&c).Error
	//if err != nil {
	//	fmt.Println(err)
	//}
	//c := []dbModel.Banners{
	//	dbModel.Banners{
	//		ImageURL: "http://www.papertraildesign.com/wp-content/uploads/2017/08/Buffalo-Plaid-Banner-Letters-B.png",
	//		Link:     "https://news.sina.com.cn/c/2019-02-06/doc-ihrfqzka3892512.shtml",
	//	},
	//	dbModel.Banners{
	//		ImageURL: "https://tse3.mm.bing.net/th?id=OIP.Dt_kEuhIXzMgVaoiiwJbrQHaFj&pid=Api",
	//		Link:     "https://news.sina.com.cn/c/2019-02-06/doc-ihrfqzka3892512.shtml",
	//	},
	//	dbModel.Banners{
	//		ImageURL: "https://tse3.mm.bing.net/th?id=OIP.mmQ2swrg-DbkW4TT2nTIjgHaEK&pid=Api",
	//		Link:     "https://news.sina.com.cn/c/2019-02-06/doc-ihrfqzka3892512.shtml",
	//	},
	//	dbModel.Banners{
	//		ImageURL: "http://www.obfuscata.com/wp-content/uploads/2017/11/Youtube-banner-template-23.jpg",
	//		Link:     "https://news.sina.com.cn/c/2019-02-06/doc-ihrfqzka3892512.shtml",
	//	},
	//}

	//createJobs()
	creatCarrerTalk()
	//createNews()

}
