package main

import (
	"demoApp/server/model/dbModel"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var testDB *gorm.DB

func initialDB() {

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

func main() {

	initialDB()

	testDB.Create(&dbModel.AppGuidanceItem{
		ImageURL: "www.dwdw.com/jpg",
		Title:    "测试",
		Detail:   "描述",
	})

}
