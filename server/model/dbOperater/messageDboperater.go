package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"goframework/orm"
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

func NewMessageDbOperater() *MessageDbOperater {

	return &MessageDbOperater{
		orm: orm.DB,
	}
}
