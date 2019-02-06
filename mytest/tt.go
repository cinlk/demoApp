package main

import (
	"demoApp/server/model/dbModel"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
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
		item.Id = "dqwdqwd" + strconv.Itoa(i)
		item.Name = "当前为多"
		item.CompanyID = "companyID"
		item.Icon = "http://icons.iconarchive.com/icons/blackvariant/button-ui-requests-6/1024/VirtualBox-icon.png"
		err := testDB.Create(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func creatCarrerTalk() {
	for i := 0; i < 12; i++ {
		item := dbModel.CarrerTalk{}
		item.Id = "宣讲会当前为多" + strconv.Itoa(i)
		item.Name = "宣讲会当前为多"
		item.CompanyID = "ompanyID"
		item.Icon = "http://file06.16sucai.com/2016/0902/994c8ad1edafbc538542708e79fa1bc5.jpg"
		err := testDB.Create(&item).Error
		if err != nil {
			fmt.Println(err)
			return
		}
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
}
