package model

import (
	"demoApp/server/model/dbModel"
	"goframework/gLog"
	"goframework/orm"
)

func CreateTables() {

	orm.InitalDB()
	// create job enum type
	_ = orm.DB.Exec("CREATE TYPE mood AS ENUM ('intern', 'graduate', 'all', 'onlineApply')").Error

	// create recruite meeting content enum type
	_ = orm.DB.Exec("CREATE TYPE contentType AS ENUM ('text', 'html')").Error

	// create table
	//orm.DB.SetLogger(gLog.GetLogUtil())
	err := orm.DB.AutoMigrate(&dbModel.AppGuidanceItem{}, &dbModel.NewsModel{}, &dbModel.Account{}, &dbModel.Recruiter{},
		&dbModel.User{}, &dbModel.SocialAccount{},
		&dbModel.Banners{}, &dbModel.LatestNews{}, &dbModel.JobCategory{}, &dbModel.TopJobs{},
		&dbModel.CareerTalk{}, &dbModel.UserApplyCarrerTalk{}, &dbModel.Company{}, &dbModel.CompuseJobs{}, &dbModel.InternJobs{},
		&dbModel.UserApplyJobs{}, &dbModel.ApplyClassify{}, &dbModel.TopWords{}, &dbModel.OnlineApply{},
		&dbModel.UserOnlineApply{}).Error

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
	err = orm.DB.Model(&dbModel.CompuseJobs{}).AddForeignKey("company_id", "company(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.InternJobs{}).AddForeignKey("company_id", "company(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
}

func CloseDB() {

	if err := orm.DB.Close(); err != nil {
		gLog.LOG_PANIC(err)
	}

}
