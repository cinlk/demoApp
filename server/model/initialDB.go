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

	_ = orm.DB.Exec("CREATE TYPE role AS ENUM ('seeker', 'recruiter')").Error

	// 通知类型
	_ = orm.DB.Exec("CREATE TYPE notify AS ENUM ('system', 'channel', 'special')").Error

	// create table
	//orm.DB.SetLogger(gLog.GetLogUtil())
	err := orm.DB.AutoMigrate(&dbModel.AppGuidanceItem{}, &dbModel.NewsModel{}, &dbModel.Account{}, &dbModel.Recruiter{},
		&dbModel.User{}, &dbModel.SocialAccount{}, &dbModel.RecruiterVisitorUser{},
		&dbModel.Banners{}, &dbModel.LatestNews{}, &dbModel.JobCategory{}, &dbModel.TopJobs{},
		&dbModel.CareerTalk{}, &dbModel.UserApplyCarrerTalk{}, &dbModel.Company{}, &dbModel.CompuseJobs{}, &dbModel.InternJobs{},
		&dbModel.UserApplyJobs{}, &dbModel.ApplyClassify{}, &dbModel.TopWords{}, &dbModel.OnlineApply{},
		&dbModel.UserOnlineApply{}, &dbModel.UserCompanyRelate{}, &dbModel.LeanCloudAccount{},
		&dbModel.SingleConversation{}, &dbModel.SystemMessage{}, &dbModel.UserCheckSystemMessage{},
		&dbModel.NotifyMessage{}, &dbModel.ForumReplyMyTime{},
		&dbModel.ForumThumbUpTime{}, &dbModel.ForumArticle{}, &dbModel.ForumHotestArticle{}).Error

	if err != nil {
		gLog.LOG_PANIC(err)
	}

	// create replationship

	err = orm.DB.Model(&dbModel.User{}).AddForeignKey("uuid", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.SocialAccount{}).AddForeignKey("uuid", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.Recruiter{}).AddForeignKey("uuid", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.RecruiterVisitorUser{}).AddForeignKey("user_id", "\"user\"(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.RecruiterVisitorUser{}).AddForeignKey("recruiter_id", "recruiter(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	//err = orm.DB.Model(&dbModel.Recruiter{}).AddForeignKey("company_id", "company(id)", "CASCADE", "CASCADE").Error
	//if err != nil {
	//	gLog.LOG_PANIC(err)
	//}

	err = orm.DB.Model(&dbModel.LeanCloudAccount{}).AddForeignKey("uuid", "account(uuid)", "CASCADE", "CASCADE").Error
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
	err = orm.DB.Model(&dbModel.UserOnlineApply{}).AddForeignKey("user_id", "\"user\"(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserOnlineApply{}).AddForeignKey("online_apply_id", "online_apply(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.CompuseJobs{}).AddForeignKey("recruiter_uuid", "recruiter(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.InternJobs{}).AddForeignKey("recruiter_uuid", "recruiter(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserApplyCarrerTalk{}).AddForeignKey("career_talk_id", "career_talk(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserCheckSystemMessage{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.ForumReplyMyTime{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.ForumThumbUpTime{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	// leancloud
	err = orm.DB.Model(&dbModel.SingleConversation{}).AddForeignKey("job_id", "compuse_jobs(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.SingleConversation{}).AddForeignKey("user_id", "lean_cloud_account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.SingleConversation{}).AddForeignKey("recruiter_id", "lean_cloud_account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.ForumArticle{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.ForumHotestArticle{}).AddForeignKey("uuid", "forum_article(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

}

func CloseDB() {

	if err := orm.DB.Close(); err != nil {
		gLog.LOG_PANIC(err)
	}

}
