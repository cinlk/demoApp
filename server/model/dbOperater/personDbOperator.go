package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"goframework/orm"
	"sort"
	"time"
)

type PersonDbOperator struct {
	orm *gorm.DB
}

// TODO
func (p *PersonDbOperator) Avatar(userId string, iconUrl string) {

}

func (p *PersonDbOperator) BriefInfos(userId, name, gender, college string) error {
	var s = -1
	if gender == "male" {
		s = 0
	} else if gender == "female" {
		s = 1
	}
	return p.orm.Model(&dbModel.User{}).Where("uuid = ?", userId).Updates(map[string]interface{}{
		"name":    name,
		"gender":  s,
		"college": college,
	}).Error

}

type deliveryHistory []httpModel.DeliveryJob

func (d deliveryHistory) Len() int {
	return len(d)
}

func (d deliveryHistory) Less(i, j int) bool {
	// 时间先后排序

	return time.Time((d[i].CreatedTime)).After(time.Time(d[j].CreatedTime))
}

func (d deliveryHistory) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// 查出所有的投递记录
func (p *PersonDbOperator) FindDeliveryInfos(userId string) []httpModel.DeliveryJob {

	var res deliveryHistory
	// 查询网申
	_ = p.orm.Model(&dbModel.UserOnlineApplyPosition{}).Where("user_online_apply_position.user_id = ?", userId).
		Joins("inner join online_apply on user_online_apply_position.online_apply_id = online_apply.id").
		Joins("inner join online_apply_position on  user_online_apply_position.position_id = online_apply_position.id").
		Select("user_online_apply_position.position_id as  job_id, user_online_apply_position.status, user_online_apply_position.feed_back," +
			" user_online_apply_position.created_at as created_time, online_apply.company_id, online_apply.location_city as address, online_apply_position.name as job_name").
		Scan(&res).Error

	// 获取公司名称
	for i := 0; i < len(res); i++ {
		//_ = p.orm.Model(&dbModel.Company{}).Where("id = ?", res[i].CompanyId).
		//	Select("name as company_name ").Scan(&res[i])
		//res.OnlineApplys[i].Type = "online_apply"
		res[i].Type = "onlineApply"

	}

	var compuse []httpModel.DeliveryJob
	// 查询校招
	_ = p.orm.Model(&dbModel.UserApplyJobs{}).Where("user_apply_jobs.user_id = ? and user_apply_jobs.is_apply = ? "+
		"and user_apply_jobs.job_type =?", userId, true, "graduate").
		Joins("inner join  compuse_jobs on user_apply_jobs.job_id = compuse_jobs.id").
		Select("user_apply_jobs.job_id, user_apply_jobs.status,user_apply_jobs.feed_back, user_apply_jobs.created_at as created_time," +
			"user_apply_jobs.job_type as type, compuse_jobs.name as job_name,compuse_jobs.company_id, " +
			"compuse_jobs.location_city as address").Scan(&compuse).Error

	for _, item := range compuse {
		res = append(res, item)
	}
	// 查询实习
	var interns []httpModel.DeliveryJob
	_ = p.orm.Model(&dbModel.UserApplyJobs{}).Where("user_apply_jobs.user_id = ? and user_apply_jobs.is_apply = ? "+
		"and user_apply_jobs.job_type =?", userId, true, "intern").
		Joins("inner join  intern_jobs on user_apply_jobs.job_id = intern_jobs.id").
		Select("user_apply_jobs.job_id, user_apply_jobs.status, user_apply_jobs.feed_back, user_apply_jobs.created_at as created_time," +
			"user_apply_jobs.job_type as type, intern_jobs.name as job_name,intern_jobs.company_id, intern_jobs.location_city as address").Scan(&interns).Error

	for _, item := range interns {
		res = append(res, item)
	}

	// 获取公司名称
	for i := 0; i < len(res); i++ {
		_ = p.orm.Model(&dbModel.Company{}).Where("id = ?", res[i].CompanyId).
			Select("name as company_name, icon_url as company_icon").Scan(&res[i])

	}
	sort.Sort(res)

	return res
}

func (p *PersonDbOperator) JobDeliveryHistory(userId, jobId string, t string) ([]httpModel.DeliveryJobStatusHistory, error) {
	if dbModel.JobType(t).Validate() == false {
		return nil, errors.New("not validate job type")
	}
	var res []httpModel.DeliveryJobStatusHistory
	time.Sleep(time.Second * 3)
	err := p.orm.Model(&dbModel.UserDeliveryStatusHistory{}).Where("user_id = ? and job_id = ? and type = ?",
		userId, jobId, t).Select("time, status, describe").Order("time desc").Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (p *PersonDbOperator) FindOnlineApplyIdBy(positionbId string) string {
	var onlineId struct {
		Id string
	}
	_ = p.orm.Model(&dbModel.UserOnlineApplyPosition{}).Where("position_id = ?", positionbId).
		Select("online_apply_id as id").Scan(&onlineId).Error
	return onlineId.Id
}

func NewPersonDbOperator() *PersonDbOperator {
	return &PersonDbOperator{
		orm: orm.DB,
	}
}
