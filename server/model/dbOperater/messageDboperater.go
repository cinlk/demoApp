package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"goframework/orm"
	"time"
)

type MessageDbOperater struct {
	orm *gorm.DB
}

func (m *MessageDbOperater) CreateConversation(userId, recruiterId, jobId, conId string) error {

	session := orm.DB.Begin()

	err := session.Set("gorm:insert_option", "ON CONFLICT(conversation_id) do nothing").Create(&dbModel.SingleConversation{
		UserID:         userId,
		RecruiterID:    recruiterId,
		JobID:          jobId,
		ConversationID: conId,
		IsValidate:     true,
	}).Error
	if err != nil {
		session.Rollback()
		return errors.Wrap(err, "create conversation error")
	}
	// 跟新已经和该职位交谈
	var apply = dbModel.UserApplyJobs{}
	err = session.Where("job_id = ?", jobId).Assign(dbModel.UserApplyJobs{
		UserId: userId,
		IsTalk: true,
		JobId:  jobId,
	}).FirstOrCreate(&apply).Error

	if err != nil {
		session.Rollback()
		return errors.Wrap(err, "update user_apply_jobs error")
	}

	session.Commit()
	return nil
}

func (m *MessageDbOperater) ConversationBy(userId, jobId string) string {

	var conId struct {
		ConversationId string `json:"conversation_id"`
	}
	_ = m.orm.Model(&dbModel.SingleConversation{}).
		Where("user_id = ? and job_id = ? and is_validate = ?", userId, jobId, true).
		Select("conversation_id").Scan(&conId)

	return conId.ConversationId

}

func (m *MessageDbOperater) RecruiterInfo(userId string) (*httpModel.HttpRecruiterModel, error) {

	var res httpModel.HttpRecruiterModel

	err := m.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", userId).
		Select("uuid as user_id, name, user_icon, " +
			"title, last_login as online_time, company_id, company").Scan(&res).Error
	if err != nil {
		return nil, err
	}

	// leancloud 账号
	_ = m.orm.Model(&dbModel.LeanCloudAccount{}).Where("uuid = ?", userId).
		Select("user_id as lean_cloud_account").Scan(&res).Error

	return &res, nil
}

func (m *MessageDbOperater) MyVisitors(userId string, offset, limit int) ([]httpModel.HttpRecruiterVisitorModel, error) {

	var res []httpModel.HttpRecruiterVisitorModel

	// 找到关注我的recruiter 的信息
	err := m.orm.Model(&dbModel.RecruiterVisitorUser{}).
		Joins("left join recruiter on recruiter.uuid =  recruiter_visitor_user.recruiter_id").
		Where("recruiter_visitor_user.user_id = ?", userId).
		Select("recruiter.uuid as recruiter_id, recruiter.name, " +
			"recruiter_visitor_user.created_at as visit_time, recruiter_visitor_user.checked, " +
			"recruiter.company, recruiter.title, recruiter.user_icon").
		Offset(offset).
		Limit(limit).Order("recruiter_visitor_user.created_at desc").
		Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *MessageDbOperater) CheckVisitor(userId, recruiterId string) {

	// 更新 已经查看了recruiter的访问
	_ = m.orm.Model(&dbModel.RecruiterVisitorUser{}).Where("user_id = ? and recruiter_id = ?",
		userId, recruiterId).Update("checked", true)
}

func (m *MessageDbOperater) ExistNewVisitor(userId string) bool {

	// 最新的访问者时间 和 check_visitor_time 比较
	var model dbModel.RecruiterVisitorUser
	err := m.orm.Model(&dbModel.RecruiterVisitorUser{}).Where("user_id = ?", userId).Last(&model).Error
	// 可能还没数据
	if err != nil {
		return false
	}
	var user dbModel.User
	err = m.orm.Model(&dbModel.User{}).Where("uuid = ?", userId).First(&user).Error
	if err != nil {
		return false
	}
	if user.CheckVisitorTime == nil {
		return true
	}
	return model.CreatedAt.After(*user.CheckVisitorTime)

}

func (m *MessageDbOperater) UpdateTimeVisitor(userId string) {

	_ = m.orm.Model(&dbModel.User{}).Where("uuid = ?", userId).
		Update("check_visitor_time", time.Now()).Error

}

func NewMessageDbOperater() *MessageDbOperater {

	return &MessageDbOperater{
		orm: orm.DB,
	}
}
