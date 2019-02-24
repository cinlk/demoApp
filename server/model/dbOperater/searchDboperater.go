package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"goframework/orm"
	"time"
)

type SearchDboperator struct {
	orm *gorm.DB
}

func (s *SearchDboperator) TopWorks(t string) []httpModel.HttpSearchTopWords {

	var res []httpModel.HttpSearchTopWords
	// 如何直接 获取列表数据？？？
	s.orm.Model(&dbModel.TopWords{}).Where("type = ?", t).Select("name").Scan(&res)

	return res
}

// 搜索关键字
func (s *SearchDboperator) SearchOnlineItem(word string) []httpModel.HttpOnlineApplyListModel {

	var res = []httpModel.HttpOnlineApplyListModel{}
	time.Sleep(time.Second * 2)
	// TODO 模糊查询
	_ = s.orm.Model(&dbModel.OnlineApply{}).Select("online_apply.id as online_apply_id, " +
		"online_apply.location_city as citys," +
		"online_apply.business_field, online_apply.end_time, online_apply.name, online_apply.outside, " +
		"online_apply.link, company.name as company_name, company.icon_url as company_icon_url ").
		Joins("left join company on company.id = online_apply.company_id").Scan(&res)
	//var online []httpModel.httpon

	return res
}

func (s *SearchDboperator) SearchCompany(word string) []httpModel.HttpCompanyListModel {

	var res = []httpModel.HttpCompanyListModel{}
	time.Sleep(time.Second * 2)
	// TODO 查询
	_ = s.orm.Model(&dbModel.Company{}).Select("id as company_id, icon_url as company_icon_url," +
		"name as company_name, review_counts, business_field, citys").Scan(&res)

	return res

}

func (s *SearchDboperator) SearchCarrerTalk(word string) []httpModel.HttpCareerTalkListModel {

	var res = []httpModel.HttpCareerTalkListModel{}
	// 根据关键字搜索
	time.Sleep(time.Second * 2)
	_ = orm.DB.Model(&dbModel.CareerTalk{}).
		Joins("left join company on company.id = career_talk.company_id").
		Select("career_talk.id as meeting_id, career_talk.college, career_talk.start_time, " +
			"career_talk.simplify_address, career_talk.icon_url as college_icon_url, " +
			"career_talk.business_field, career_talk.city, company.name as company_name").Scan(&res)

	return res

}

func (s *SearchDboperator) SearchGraduateJob(word string) []httpModel.HttpJobListModel {
	// graduate
	var res = []httpModel.HttpJobListModel{}
	time.Sleep(time.Second * 2)
	_ = orm.DB.Model(&dbModel.CompuseJobs{}).Joins("left join  company on company.id = compuse_jobs.company_id").
		Select("compuse_jobs.id as job_id, compuse_jobs.icon_url, compuse_jobs.type, compuse_jobs.name as job_name,"+
			"compuse_jobs.location_city as address, compuse_jobs.education as degree, "+
			"compuse_jobs.review_counts as review_count, compuse_jobs.created_time, "+
			"compuse_jobs.business_field, company.type as company_type , company.name as company_name ").
		Where("compuse_jobs.type = ?", "graduate").Scan(&res)

	return res

}

func (s *SearchDboperator) SearchInternJobs(word string) []httpModel.HttpInternListModel {
	// intern
	var res = []httpModel.HttpInternListModel{}
	time.Sleep(time.Second * 2)
	_ = orm.DB.Model(&dbModel.InternJobs{}).Joins("left join  company on company.id = intern_jobs.company_id").
		Select("intern_jobs.id as job_id, intern_jobs.icon_url, intern_jobs.type, intern_jobs.name as job_name,"+
			"intern_jobs.location_city as address, intern_jobs.education as degree, "+
			"intern_jobs.review_counts as review_count, intern_jobs.created_time, "+
			"intern_jobs.days, intern_jobs.months, intern_jobs.pay_day, "+
			"intern_jobs.can_transfer as is_transfer, intern_jobs.business_field, company.name as company_name ").
		Where("intern_jobs.type = ?", "intern").Scan(&res)

	return res

}

func NewSearchDboperator() *SearchDboperator {
	return &SearchDboperator{
		orm: orm.DB,
	}
}
