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

	// 简历类型
	_ = orm.DB.Exec("CREATE TYPE resume AS ENUM ('text', 'attache')").Error

	// create table
	//orm.DB.SetLogger(gLog.GetLogUtil())
	err := orm.DB.AutoMigrate(&dbModel.AppGuidanceItem{}, &dbModel.NewsModel{}, &dbModel.Account{}, &dbModel.Recruiter{},
		&dbModel.User{}, &dbModel.SocialAccount{}, &dbModel.RecruiterVisitorUser{},
		&dbModel.Banners{}, &dbModel.LatestNews{}, &dbModel.JobCategory{}, &dbModel.TopJobs{},
		&dbModel.CareerTalk{}, &dbModel.UserApplyCarrerTalk{}, &dbModel.Company{}, &dbModel.CompuseJobs{}, &dbModel.InternJobs{},
		&dbModel.UserApplyJobs{}, &dbModel.ApplyClassify{}, &dbModel.TopWords{}, &dbModel.OnlineApply{},
		&dbModel.UserOnlineApplyPosition{}, &dbModel.OnlineApplyPosition{}, &dbModel.UserCollectedOnlineApply{},
		&dbModel.UserCompanyRelate{}, &dbModel.LeanCloudAccount{},
		&dbModel.SingleConversation{}, &dbModel.SystemMessage{}, &dbModel.UserCheckSystemMessage{},
		&dbModel.NotifyMessage{}, &dbModel.ForumReplyMyTime{},
		&dbModel.ForumThumbUpTime{}, &dbModel.ForumArticle{}, &dbModel.ForumHotestArticle{},
		&dbModel.ReplyForumPost{}, &dbModel.SecondReplyPost{}, &dbModel.UserLikePost{},
		&dbModel.UserLikeReply{}, &dbModel.UserLikeSubReply{},
		&dbModel.UserCollectedPost{}, &dbModel.UserAlertPost{},
		&dbModel.UserAlertReply{}, &dbModel.UserAlertSubReply{},
		&dbModel.UserDeliveryStatusHistory{},
		// 简历
		&dbModel.MyResume{}, &dbModel.AttachFileResume{}, &dbModel.TextResume{}, &dbModel.TextResumeBaseInfo{},
		&dbModel.TextResumeEducation{}, &dbModel.TextResumeWorkExperience{}, &dbModel.TextResumeProject{},
		&dbModel.TextResumeCollegeActivity{}, &dbModel.TextResumeSocialPractice{}, &dbModel.TextResumeSkills{},
		&dbModel.TextResumeOther{}, &dbModel.TextResumeEstimate{},
	).Error

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
	err = orm.DB.Model(&dbModel.UserOnlineApplyPosition{}).AddForeignKey("user_id", "\"user\"(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserOnlineApplyPosition{}).AddForeignKey("online_apply_id", "online_apply(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserOnlineApplyPosition{}).AddForeignKey("position_id", "online_apply_position(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.OnlineApplyPosition{}).AddForeignKey("online_apply_id", "online_apply(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserCollectedOnlineApply{}).AddForeignKey("online_apply_id", "online_apply(id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserCollectedOnlineApply{}).AddForeignKey("user_id", "\"user\"(uuid)", "CASCADE", "CASCADE").Error
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

	// 论坛
	err = orm.DB.Model(&dbModel.ForumArticle{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.ForumHotestArticle{}).AddForeignKey("uuid", "forum_article(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.ReplyForumPost{}).AddForeignKey("post_uuid", "forum_article(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.ReplyForumPost{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.SecondReplyPost{}).AddForeignKey("reply_id", "reply_forum_post(reply_id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.SecondReplyPost{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.SecondReplyPost{}).AddForeignKey("talked_user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserLikePost{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserLikePost{}).AddForeignKey("post_uuid", "forum_article(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserLikeReply{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserLikeReply{}).AddForeignKey("reply_id", "reply_forum_post(reply_id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserLikeSubReply{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserLikeSubReply{}).AddForeignKey("second_reply_id", "second_reply_post(second_reply_id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserCollectedPost{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserCollectedPost{}).AddForeignKey("post_uuid", "forum_article(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserAlertPost{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserAlertPost{}).AddForeignKey("post_id", "forum_article(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserAlertReply{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserAlertReply{}).AddForeignKey("reply_id", "reply_forum_post(reply_id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.UserAlertSubReply{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.UserAlertSubReply{}).AddForeignKey("second_reply_id", "second_reply_post(second_reply_id)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	// person
	err = orm.DB.Model(&dbModel.UserDeliveryStatusHistory{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	// 1 简历
	err = orm.DB.Model(&dbModel.MyResume{}).AddForeignKey("user_id", "account(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.AttachFileResume{}).AddForeignKey("resume_id", "my_resume(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.TextResume{}).AddForeignKey("resume_id", "my_resume(uuid)", "CASCADE", "CASCADE").Error
	if err != nil {
		gLog.LOG_PANIC(err)
	}


	err = orm.DB.Model(&dbModel.TextResumeEducation{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}

	err = orm.DB.Model(&dbModel.TextResumeWorkExperience{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.TextResumeProject{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.TextResumeCollegeActivity{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.TextResumeSocialPractice{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.TextResumeSkills{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.TextResumeOther{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}
	err = orm.DB.Model(&dbModel.TextResumeEstimate{}).AddForeignKey("resume_id", "text_resume(resume_id)", "CASCADE", "CASCADE").Error
	if err != nil{
		gLog.LOG_PANIC(err)
	}


}

func CloseDB() {

	if err := orm.DB.Close(); err != nil {
		gLog.LOG_PANIC(err)
	}

}
