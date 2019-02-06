package model

import (
	"demoApp/server/model/dbModel"
	"goframework/gLog"
	"goframework/orm"
)

func CreateTables() {

	orm.InitalDB()
	// create table
	//orm.DB.SetLogger(gLog.GetLogUtil())
	err := orm.DB.AutoMigrate(&dbModel.AppGuidanceItem{}, &dbModel.Account{}, &dbModel.Recruiter{},
		&dbModel.User{}, &dbModel.SocialAccount{},
		&dbModel.Banners{}, &dbModel.LatestNews{}, &dbModel.JobCategory{}, &dbModel.TopJobs{},
		&dbModel.CarrerTalk{}, &dbModel.Company{}, &dbModel.CompuseJobs{}, &dbModel.ApplyClassify{}).Error

	if err != nil {
		gLog.LOG_PANIC(err)
	}

	// create replationship

	err = orm.DB.Model(&dbModel.User{}).AddForeignKey("phone", "account(phone)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.SocialAccount{}).AddForeignKey("phone", "account(phone)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
}

func CloseDB() {

	if err := orm.DB.Close(); err != nil {
		gLog.LOG_PANIC(err)
	}

}
