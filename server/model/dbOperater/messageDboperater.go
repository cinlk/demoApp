package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"demoApp/server/utils"
	"demoApp/server/utils/leancloud"
	"fmt"
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

func (m *MessageDbOperater) HasNewSystemMessage(userId string) bool {
	// 检查消息表 和 用户查看系统消息表
	// 对比时间

	// 最后一条系统消息记录
	var last dbModel.SystemMessage
	err := m.orm.Model(&dbModel.SystemMessage{}).Last(&last).Error
	if err != nil {
		return false
	}

	var check dbModel.UserCheckSystemMessage
	err = m.orm.Model(&dbModel.UserCheckSystemMessage{}).Where("user_id = ?", userId).First(&check).Error
	// 可能还没数据

	if err == gorm.ErrRecordNotFound {
		return true
	}
	if err != nil {
		return false
	}

	return last.CreatedAt.After(*check.CheckTime)

}

func (m *MessageDbOperater) ReviewSystemMessage(userId string) error {

	var model dbModel.UserCheckSystemMessage
	var t = time.Now()
	return m.orm.Where(dbModel.UserCheckSystemMessage{
		UserId: userId,
	}).Assign(dbModel.UserCheckSystemMessage{
		CheckTime: &t,
	}).FirstOrCreate(&model).Error

}

// 检查是否存在新的点赞消息
func (m *MessageDbOperater) HasNewThumbUpMessage(userId string) bool {
	var model dbModel.ForumThumbUpTime
	err := m.orm.Model(&dbModel.ForumThumbUpTime{}).Where("user_id = ?", userId).First(&model).Error
	if err != nil {
		return false
	}

	if model.CheckTime == nil && model.LatestThumbTime != nil {
		return true
	}
	if model.CheckTime != nil && model.LatestThumbTime != nil {
		return model.LatestThumbTime.After(*model.CheckTime)
	}

	return false
}

func (m *MessageDbOperater) ReviewThumbUpMessage(userId string) error {
	// 只能更新
	var t = time.Now()
	return m.orm.Model(&dbModel.ForumThumbUpTime{}).Where("user_id = ?", userId).
		Update("check_time", &t).Error

}

// 论坛回复我的最新 消息提醒
func (m *MessageDbOperater) HasNewForumReply2Me(userId string) bool {

	var model dbModel.ForumReplyMyTime
	err := m.orm.Model(&dbModel.ForumReplyMyTime{}).Where("user_id = ?", userId).First(&model).Error
	if err != nil {
		return false
	}
	fmt.Println(model)

	if model.CheckTime == nil && model.LatestReplyTime != nil {
		return true
	}
	if model.CheckTime != nil && model.LatestReplyTime != nil {
		return model.LatestReplyTime.After(*model.CheckTime)
	}

	return false
}

func (m *MessageDbOperater) ReviewForumReply2Me(userId string) error {

	// 只能更新
	var t = time.Now()
	return m.orm.Model(&dbModel.ForumReplyMyTime{}).Where("user_id = ?", userId).
		Update("check_time", &t).Error
}

func (m *MessageDbOperater) SystemMessage(message string) error {

	session := m.orm.Begin()
	var notifiyId = utils.RandStrings(16)

	// 加入到系统消息表
	err := session.Create(&dbModel.SystemMessage{
		Content: message,
		Title:   "[系统]",
	}).Error
	if err != nil {
		session.Rollback()
		return err
	}

	// 加入到通知表
	err = session.Create(&dbModel.NotifyMessage{
		Content:        message,
		NotificationId: notifiyId,
		Type:           dbModel.System,
		Title:          "系统通知",
		Channels:       []string{"systemNotify"},
	}).Error
	if err != nil {
		session.Rollback()
		return err

	}
	// 发送推送（订阅系统消息设备接受推送）
	err = leancloud.LeanCloudSendSystemNotify("systemNotify", "系统通知", notifiyId, message)
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()
	return nil
}

func (m *MessageDbOperater) SpecialUserNotify(leanCloudUserId, message string) error {
	session := m.orm.Begin()
	var notifyId = utils.RandStrings(16)

	err := session.Create(&dbModel.NotifyMessage{
		Content:        message,
		NotificationId: notifyId,
		Type:           dbModel.Special,
		Title:          "新消息",
		Channels:       []string{leanCloudUserId},
	}).Error
	if err != nil {
		session.Rollback()
		return err
	}

	err = leancloud.LeanCloudSendUserNotify(leanCloudUserId, "新消息", notifyId, message)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func NewMessageDbOperater() *MessageDbOperater {

	return &MessageDbOperater{
		orm: orm.DB,
	}
}
