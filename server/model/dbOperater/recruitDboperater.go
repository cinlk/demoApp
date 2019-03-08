package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"demoApp/server/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"goframework/orm"
	"strings"
	"time"
)

type RecruiteDboperator struct {
	orm *gorm.DB
}

func (r *RecruiteDboperator) ListApplyOnlines(citys []string, field, t string, offset, limit int64) []httpModel.HttpOnlineApplyListModel {

	var res = []httpModel.HttpOnlineApplyListModel{}
	// TODO 查找指定数据
	var whereCitySql = ""
	if len(citys) != 0 {
		whereCitySql = fmt.Sprintf("array[%s] && online_apply.location_city", r.joinListSql(citys))
	}
	var whereFieldSql = ""
	if field != "" {
		whereFieldSql += fmt.Sprintf(" '%s' = any(online_apply.business_field)", field)
	}
	// 搜索名称？？
	var whereTypeSql = ""
	if t != "" {
		t = "%" + t + "%"
		whereTypeSql += fmt.Sprintf(" online_apply.name like  '%s' ", t)
	}

	_ = r.orm.Model(&dbModel.OnlineApply{}).
		Where(whereCitySql).
		Where(whereFieldSql).
		Where(whereTypeSql).
		Select("online_apply.id as online_apply_id, " +
			"online_apply.location_city as citys," +
			"online_apply.business_field, online_apply.end_time, online_apply.name, online_apply.outside as out_side, " +
			"online_apply.link, company.name as company_name, company.icon_url as company_icon_url ").
		Joins("left join company on company.id = online_apply.company_id").
		Order("online_apply.created_time desc").Offset(offset).Limit(limit).Scan(&res)

	return res
}

func (r *RecruiteDboperator) OnlineApplyInfo(onlineApplyId, userId string) httpModel.HttpOnlineApplyModel {

	var online httpModel.HttpOnlineApplyModel

	_ = r.orm.Model(&dbModel.UserOnlineApply{}).
		Where("user_id = ? and online_apply_id = ?", userId, onlineApplyId).Select("is_collected").Scan(&online)
	_ = r.orm.Model(&dbModel.OnlineApply{}).Where("id = ?", onlineApplyId).
		Select("link, name, icon_url, id, created_time, end_time, location_city as citys, " +
			"positions, major, content_type, content, company_id, outside as outer_side").Scan(&online)
	// company
	_ = r.orm.Model(&dbModel.Company{}).Where("company.id = ?", online.CompanyID).
		Select("id, link, name, icon_url, created_time, describe, simple_describe, feature_tags as tags, " +
			"review_counts, business_field, citys, is_validate, total_staff as staff ").Scan(&online.Company)

	return online
}

func (r *RecruiteDboperator) ListCareerTalk(college []string,
	field, city, timeStr string, offset, limit int64) []httpModel.HttpCareerTalkListModel {
	var res = []httpModel.HttpCareerTalkListModel{}
	// 过滤条件检测 TODO
	var collegeWhereSql = ""
	if len(college) > 0 {
		collegeWhereSql += fmt.Sprintf("career_talk.college = any(array[%s])", r.joinListSql(college))
	} else if city != "" {
		collegeWhereSql += fmt.Sprintf("career_talk.city = '%s' ", city)
	}

	var timeWhereSql = ""
	switch timeStr {
	case "将来":
		timeWhereSql += fmt.Sprintf("career_talk.start_time > '%s'", time.Now().Format(utils.BASETIME_FORMAT))
	case "过去":
		timeWhereSql += fmt.Sprintf("career_talk.end_time < '%s'", time.Now().Format(utils.BASETIME_FORMAT))
	default:
		if t, err := time.Parse(utils.TIME_YEAR_DAY, timeStr); err == nil {
			timeWhereSql += fmt.Sprintf(" career_talk.start_time < '%s' and  career_talk.start_time > '%s'",
				t.Add(time.Hour*24).Format(utils.BASETIME_FORMAT), t.Format(utils.BASETIME_FORMAT))
		}

	}
	var fieldWhereSql = ""

	if field != "" {
		fieldWhereSql += fmt.Sprintf(" '%s' = any(career_talk.business_field)", field)
	}

	_ = r.orm.Model(&dbModel.CareerTalk{}).Joins("left join company on company.id = career_talk.company_id").
		Select("career_talk.id as meeting_id, career_talk.college, career_talk.start_time, " +
			"career_talk.simplify_address, career_talk.icon_url as college_icon_url, " +
			"career_talk.business_field, career_talk.city, company.name as company_name").
		Where(collegeWhereSql).
		Where(timeWhereSql).
		Where(fieldWhereSql).Offset(offset).Limit(limit).Scan(&res)

	return res

}

func (r *RecruiteDboperator) ListCompany(city []string, field, companyType string, offset, limit int64) []httpModel.HttpCompanyListModel {

	var res = []httpModel.HttpCompanyListModel{}

	var cityWhereSql = ""
	if len(city) > 0 {
		cityWhereSql += fmt.Sprintf("array[%s] && company.citys", r.joinListSql(city))
	}

	var fieldWhereSql = ""
	if field != "" {
		fieldWhereSql += fmt.Sprintf(" '%s' = any(company.business_field)", field)
	}
	var companyTypeWhereSql = ""
	if companyType != "" {
		companyTypeWhereSql += fmt.Sprintf(" company.type = '%s'", companyType)
	}

	_ = r.orm.Model(&dbModel.Company{}).Select("id as company_id, icon_url as company_icon_url," +
		"name as company_name, review_counts, business_field, citys").
		Where(cityWhereSql).
		Where(fieldWhereSql).
		Where(companyTypeWhereSql).
		Offset(offset).
		Limit(limit).
		Order("company.created_time desc").
		Scan(&res)

	return res
}

func (r *RecruiteDboperator) ListGraduateJobs(citys []string, field, degree string, offset, limit int64) []httpModel.HttpJobListModel {

	var res = []httpModel.HttpJobListModel{}

	var whereSql = ""
	if len(citys) > 0 {
		whereSql += fmt.Sprintf("array[%s] && compuse_jobs.location_city", r.joinListSql(citys))
	}
	if field != "" {
		if len(citys) > 0 {
			whereSql += " and "
		}
		whereSql += fmt.Sprintf(" '%s' =  any(compuse_jobs.business_field)", field)
	}
	if degree != "" {
		if len(citys) > 0 || field != "" {
			whereSql += " and "
		}
		whereSql += fmt.Sprintf("compuse_jobs.education = '%s'", degree)
	}

	_ = r.orm.Model(&dbModel.CompuseJobs{}).Joins("left join  company on company.id = compuse_jobs.company_id").
		Where(whereSql).
		Select("compuse_jobs.id as job_id, compuse_jobs.icon_url, compuse_jobs.type, compuse_jobs.name as job_name," +
			"compuse_jobs.location_city as address, compuse_jobs.education as degree, " +
			"compuse_jobs.review_counts as review_count, compuse_jobs.created_time, " +
			"compuse_jobs.business_field, company.type as company_type , company.name as company_name ").
		Order(" compuse_jobs.created_time desc").
		Offset(offset).Limit(limit).Scan(&res)

	return res
}

func (r *RecruiteDboperator) GraduatJobInfo(jobId, userId string) httpModel.HttpGraduateModel {

	var res httpModel.HttpGraduateModel
	// 用户关联的状态
	// 公司
	// job 内容
	// 关联的recruiter
	_ = r.orm.Model(&dbModel.CompuseJobs{}).Where("id = ?", jobId).
		Select("id, name, is_validate, created_time, icon_url, type, benefits, company_id," +
			"recruiter_uuid, address, location_city as city, need_skills, work_content, business_field, major, " +
			"tags as job_tags, review_counts, apply_end_time, salary, education").Scan(&res)
	_ = r.orm.Model(&dbModel.Company{}).Where("id = ?", res.CompanyId).
		Select("id as company_id, icon_url , name as company_name, " +
			"citys, business_field, total_staff as staff").Scan(&res.Company)
	_ = r.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", res.RecruiterUUID).
		Select("uuid as user_id, name, user_icon, title, last_login as online_time, company").Scan(&res.Recruiter)

	_ = r.orm.Model(&dbModel.UserApplyJobs{}).Where("job_id = ? and user_id = ?", jobId, userId).
		Select("is_collected, is_apply, is_talk").Scan(&res)
	return res

}

func (r *RecruiteDboperator) InternJobInfo(jobId, userId string) httpModel.HttpInternJobModel {

	var res httpModel.HttpInternJobModel
	_ = r.orm.Model(&dbModel.InternJobs{}).Where("id = ?", jobId).
		Select("id, name, is_validate, created_time, icon_url, type, benefits, company_id," +
			"recruiter_uuid, address, location_city as city, need_skills, work_content, business_field, major, " +
			"tags as job_tags, review_counts, apply_end_time, salary, education, pay_day, days, months, can_transfer").Scan(&res)
	_ = r.orm.Model(&dbModel.Company{}).Where("id = ?", res.CompanyId).
		Select("id as company_id, icon_url , name as company_name, " +
			"citys, business_field, total_staff as staff").Scan(&res.Company)
	_ = r.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", res.RecruiterUUID).
		Select("uuid as user_id, name, user_icon, title, last_login as online_time, company").Scan(&res.Recruiter)

	_ = r.orm.Model(&dbModel.UserApplyJobs{}).Where("job_id = ? and user_id = ?", jobId, userId).
		Select("is_collected, is_apply, is_talk").Scan(&res)

	return res
}

func (r *RecruiteDboperator) ListInternJobs(condition map[string]interface{}, citys []string,
	field string, offset, limit int64) []httpModel.HttpInternListModel {

	// 细化判断条件规则 TODO
	var res = []httpModel.HttpInternListModel{}

	var daysWhereSql = ""
	if day, ok := condition["days"]; ok {
		daysWhereSql = fmt.Sprintf("intern_jobs.days = '%v'", day)
	}
	var monthsWhereSql = ""
	if months, ok := condition["month"]; ok {
		monthsWhereSql = fmt.Sprintf("intern_jobs.months <= '%v'", months)
	}

	var payDayWhereSql = ""
	if pays, ok := condition["pay"]; ok {
		payDayWhereSql = fmt.Sprintf("intern_jobs.pay_day >= '%v'", pays)
	}
	var isStaffWhereSql = ""
	if chance, ok := condition["transfer"]; ok {
		isStaffWhereSql = fmt.Sprintf("intern_jobs.can_transfer = '%v'", chance.(bool))
	}

	var cityWhereSql = ""
	if len(citys) > 0 {
		cityWhereSql += fmt.Sprintf("array[%s] && intern_jobs.location_city", r.joinListSql(citys))
	}

	var fieldWhereSql = ""
	if field != "" {
		fieldWhereSql += fmt.Sprintf(" '%s' =  any(intern_jobs.business_field)", field)
	}

	_ = orm.DB.Model(&dbModel.InternJobs{}).Joins("left join  company on company.id = intern_jobs.company_id").
		Select("intern_jobs.id as job_id, intern_jobs.icon_url, intern_jobs.type, intern_jobs.name as job_name,"+
			"intern_jobs.location_city as address, intern_jobs.education as degree, "+
			"intern_jobs.review_counts as review_count, intern_jobs.created_time, "+
			"intern_jobs.days, intern_jobs.months, intern_jobs.pay_day, "+
			"intern_jobs.can_transfer as is_transfer, intern_jobs.business_field, company.name as company_name ").
		Where("intern_jobs.type = ?", "intern").
		Where(cityWhereSql).
		Where(fieldWhereSql).
		Where(daysWhereSql).
		Where(monthsWhereSql).
		Where(payDayWhereSql).
		Where(isStaffWhereSql).
		Order("intern_jobs.created_time desc").
		Limit(limit).
		Offset(offset).
		Scan(&res)

	return res

}

// 获取招聘者相关信息
func (r *RecruiteDboperator) FindRecruiterInfo(id string) httpModel.HttpRecruiterMainModel {
	// 自身信息
	// 所属公司信息
	// 自己发布的职位
	var res httpModel.HttpRecruiterMainModel

	_ = r.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", id).
		Select("uuid as user_id, name, user_icon, " +
			"title, last_login as online_time, company_id, company").Scan(&res.Recruiter)
	_ = r.orm.Model(&dbModel.Company{}).Where("id = ?", res.Recruiter.CompanyID).
		Select("id as company_id, icon_url, name as company_name, " +
			"citys, business_field, total_staff as staff ").Scan(&res.Company)
	_ = r.orm.Model(&dbModel.CompuseJobs{}).Where("recruiter_uuid = ?", res.Recruiter.UserID).
		Select("id, name, address, education, type, created_time").Scan(&res.Jobs)

	var interns []httpModel.HttpRecruiterJobsModel
	_ = r.orm.Model(&dbModel.InternJobs{}).Where("recruiter_uuid = ?", res.Recruiter.UserID).
		Select("id, name, address, education, type, created_time").Scan(&interns)
	for _, item := range interns {
		res.Jobs = append(res.Jobs, item)
	}

	return res

}

// 获取公司信息
func (r *RecruiteDboperator) CompanyInfo(id, userID string) httpModel.HttpCompanyModel {
	var res httpModel.HttpCompanyModel
	// 公司信息
	// 职位信息
	// 宣讲信息
	// 用户关联信息
	_ = r.orm.Model(&dbModel.Company{}).Where("id = ?", id).Select("id, name, link, " +
		"icon_url, created_time, describe, simple_describe as simple_des, citys, total_staff as staff, business_field," +
		" feature_tags as tags, review_counts,job_tags").Scan(&res)
	//_ = r.orm.Model(&dbModel.InternJobs{}).Where("company_id = ?", id).Select("").Scan(&res.Jobs)

	//_ = r.orm.Model(&dbModel.CompuseJobs{}).Where("company_id = ?", id).Select().Scan()
	//_ = r.orm.Model(&dbModel.CareerTalk{}).Where("company_id = ?", id).Select().Scan()
	_ = r.orm.Model(&dbModel.UserCompanyRelate{}).Where("user_id = ? and company_id = ?",
		userID, id).Select("is_collected").Scan(&res)

	return res

}

// 获取宣讲会信息
func (r *RecruiteDboperator) RecruitMeetingInfo(id, userID string) httpModel.HttpCareerTalkModel {
	var res httpModel.HttpCareerTalkModel
	// 获取公司
	// 宣讲会信息
	// 关联的用户投递信息
	_ = r.orm.Model(&dbModel.CareerTalk{}).Where("id = ?", id).
		Select("id, name, is_validate, created_time, icon_url, college, address, simplify_address," +
			"start_time, end_time, content, content_type, link, majors, business_field, reference, company_id").Scan(&res)
	_ = r.orm.Model(&dbModel.Company{}).Where("id = ?", res.CompanyID).
		Select("id as company_id, icon_url, name as company_name, " +
			"citys, business_field, total_staff as staff").Scan(&res.Company)
	_ = r.orm.Model(&dbModel.UserApplyCarrerTalk{}).
		Where("user_id = ? and career_talk_id = ?", userID, id).
		Select("is_collected").Scan(&res)

	return res
}

// 公司关联的jobs
func (r *RecruiteDboperator) CompanyRelatedJobs(id, tag string, offset, limit int64) []httpModel.HttpCompanyTagJobsModel {

	// 职位
	var res []httpModel.HttpCompanyTagJobsModel
	//var company dbModel.Company
	//err := r.orm.Model(&dbModel.Company{}).Where("id = ?", id).First(&company).Error
	//if err != nil {
	//	return res
	//}
	var whereSql = ""
	if tag == "" {
		whereSql = fmt.Sprintf("company_id = '%s'", id)
	} else {
		whereSql = fmt.Sprintf("company_id = '%s' and '%s'=any(tags) ", id, tag)
	}
	// 不考虑 类型 顺序
	//_ = r.orm.Model(&company).Select().Offset().Limit().Association("CompuseJob").
	_ = r.orm.Model(&dbModel.CompuseJobs{}).Where(whereSql).
		Select("id,  type, name , address, education, created_time, tags").
		Offset(offset).Limit(limit).Order("created_time desc").
		Scan(&res)

	var jbs []httpModel.HttpCompanyTagJobsModel

	_ = r.orm.Model(&dbModel.InternJobs{}).Where(whereSql).
		Select("id, type, name , address, education, created_time, tags").
		Offset(offset).Limit(limit).Order("created_time desc").
		Scan(&jbs)

	for _, item := range jbs {
		res = append(res, item)
	}

	return res
}

func (r *RecruiteDboperator) CompanyRelatedCareerTalk(id string, offset, limit int64) []httpModel.HttpCareerTalkListModel {

	var res []httpModel.HttpCareerTalkListModel
	// 关联的宣讲会
	//var company dbModel.Company
	//err := r.orm.Model(&dbModel.Company{}).Where("id = ?", id).First(&company).Error
	//if err != nil {
	//	return res
	//}
	_ = r.orm.Model(&dbModel.CareerTalk{}).Where("company_id = ?", id).Offset(offset).Limit(limit).
		Select("id as meeting_id, college, start_time, icon_url as college_icon_url, " +
			" simplify_address, business_field, city").Order("created_time desc").Scan(&res)
	return res

}

// 拼接数据
func (r *RecruiteDboperator) joinListSql(l []string) string {
	if len(l) < 1 {
		return ""
	}
	for i := 0; i < len(l); i++ {
		l[i] = fmt.Sprintf(" '%s' ", l[i])
	}
	return strings.Join(l, ",")
}

func NewRecruiteDboperator() *RecruiteDboperator {

	return &RecruiteDboperator{
		orm: orm.DB,
	}
}
