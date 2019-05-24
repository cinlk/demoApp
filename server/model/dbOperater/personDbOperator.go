package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	utils2 "demoApp/server/utils"
	"demoApp/server/utils/errorStatus"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"goframework/orm"
	"goframework/utils"
	"net/http"
	"sort"
	"strconv"
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

func (p *PersonDbOperator) GetMyResumes(userId string) ([]httpModel.ResumeListModel, error) {
	var res []httpModel.ResumeListModel
	err := p.orm.Model(&dbModel.MyResume{}).Where("user_id = ?", userId).
		Select("uuid as resume_id, type, is_primary, name, created_at as created_time").Scan(&res).Error

	return res, err
}

func (p *PersonDbOperator) NewResume(userId string) (*httpModel.ResumeListModel, error) {

	var rtype = dbModel.ResumeText

	// 检查数量 不能超过5个
	var count int
	err := p.orm.Model(&dbModel.MyResume{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, &errorStatus.AppError{
			Code: 40111,
			Err:  errors.New("maxsize resumes"),
		}
	}

	var uuid = utils.GetUUID()
	session := p.orm.Begin()

	// 名字怎么自动取 TODO
	var last int
	session.Model(&dbModel.MyResume{}).Where("user_id = ?", userId).Count(&last)
	var name = "我的简历" + strconv.Itoa(int(last+1))

	err = session.Create(&dbModel.MyResume{
		UserId:    userId,
		Uuid:      uuid,
		Name:      name,
		Type:      string(rtype),
		IsPrimary: count == 0,
	}).Error
	if err != nil {
		session.Rollback()
		return nil, err
	}

	err = session.Create(&dbModel.TextResume{
		ResumeId: uuid,
	}).Error
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Create(&dbModel.TextResumeBaseInfo{
		ResumeId: uuid,
	}).Error
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Create(&dbModel.TextResumeEstimate{
		ResumeId: uuid,
	}).Error
	if err != nil{
		session.Rollback()
		return nil, err
	}
	session.Commit()

	return &httpModel.ResumeListModel{
		Name:        name,
		IsPrimary:   count == 0,
		CreatedTime: httpModel.TStringFormat(time.Now()),
		Type:        string(rtype),
		ResumeId:    uuid,
	}, nil

}

func (p *PersonDbOperator) CreateAttachFile(data []byte, fileName, userId string) error {

	var rtype = dbModel.ResumeAttach
	session := p.orm.Begin()

	// 存入数据库
	var count int
	err := session.Model(&dbModel.MyResume{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		session.Rollback()
		return err
	}
	if count >= 5 {
		return &errorStatus.AppError{
			Code: http.StatusNotAcceptable,
			Err:  errors.New("maxsize resumes"),
		}
	}

	var uuid = utils.GetUUID()

	// 名字怎么自动取 TODO
	var last int
	session.Model(&dbModel.MyResume{}).Where("user_id = ?", userId).Count(&last)
	var name = "我的简历" + strconv.Itoa(int(last+1))

	err = session.Create(&dbModel.MyResume{
		UserId:    userId,
		Uuid:      uuid,
		Name:      name,
		Type:      string(rtype),
		IsPrimary: count == 0,
	}).Error
	if err != nil {
		session.Rollback()
		return err
	}
	// test
	err = session.Create(&dbModel.AttachFileResume{
		ResumeId: uuid,
		FileUrl:  "https://xueqiu.com/",
	}).Error
	if err != nil {
		session.Rollback()
		return err
	}

	// 上传到七牛云 TODO

	session.Commit()
	return nil
}

func (p *PersonDbOperator) SetPrimaryResume(userId, resumeId string) error {

	var target dbModel.MyResume
	err := p.orm.Model(&dbModel.MyResume{}).Where("user_id = ? and uuid = ?", userId, resumeId).First(&target).Error
	if err != nil {
		return err
	}
	target.IsPrimary = true
	session := p.orm.Begin()
	err = p.orm.Model(&dbModel.MyResume{}).Where("user_id = ? and uuid <> ?", userId, resumeId).Update("is_primary", false).Error
	if err != nil {
		session.Rollback()
		return err
	}
	err = p.orm.Model(target).Update("is_primary", true).Error
	if err != nil {
		session.Rollback()
		return err
	}

	return nil

}

func (p *PersonDbOperator) ChangeResumeName(userId, resumeId, name string) error {

	return p.orm.Model(&dbModel.MyResume{}).Where("user_id = ? and uuid = ?", userId, resumeId).
		Update("name", name).Error
}

func (p *PersonDbOperator) DeleteResume(userId, resumeId, t string) error {
	var rtype = dbModel.ResumeType(t)
	if rtype.Validate() == false {
		return &errorStatus.AppError{
			Code: http.StatusBadRequest,
			Err:  errors.New("bad resume type"),
		}
	}
	session := p.orm.Begin()

	err := session.Unscoped().Delete(dbModel.MyResume{}, "user_id = ? and uuid = ?", userId, resumeId).Error
	if err != nil {
		session.Rollback()
		return err
	}
	switch rtype {
	case dbModel.ResumeText:
		err = session.Unscoped().Delete(dbModel.TextResume{}, "resume_id = ?", resumeId).Error
		// TODO 删除关联的数据
		err = session.Unscoped().Delete(dbModel.TextResumeBaseInfo{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeEducation{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeWorkExperience{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeProject{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeCollegeActivity{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeSocialPractice{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeSkills{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeOther{}, "resume_id = ?", resumeId).Error
		err = session.Unscoped().Delete(dbModel.TextResumeEstimate{}, "resume_id = ?", resumeId).Error


	case dbModel.ResumeAttach:
		err = session.Unscoped().Delete(dbModel.AttachFileResume{}, "resume_id = ?", resumeId).Error
	default:
		break
	}
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()
	return nil

}

func (p *PersonDbOperator)  TextResumeInfo(userId, resumeId string) (*httpModel.TextResumeContentModel, error)  {
	var res httpModel.TextResumeContentModel
	err := p.orm.Model(&dbModel.TextResume{}).Where("resume_id = ?", resumeId).
		Select("resume_id, complete_level as level").Scan(&res).Error
	if err != nil{
		return nil, err
	}
	// 查询相关信息
	err = p.orm.Model(&dbModel.TextResumeBaseInfo{}).Where("resume_id = ?", res.ResumeId).
		Select("id, icon, name, college, gender, city, degree, birthday, phone, email").Scan(&res.BaseInfo).Error

	err = p.orm.Model(&dbModel.TextResumeEducation{}).Where("resume_id = ?", res.ResumeId).
		Select("id, college, major, rank, degree, describe, start_time, end_time").Scan(&res.Educations).Error

	err = p.orm.Model(&dbModel.TextResumeWorkExperience{}).Where("resume_id = ?", res.ResumeId).
		Select("id, company_name, work_type, city, position, describe, start_time, end_time").Scan(&res.Works).Error

	err = p.orm.Model(&dbModel.TextResumeProject{}).Where("resume_id = ?", res.ResumeId).
		Select("id, project_name, project_level, position, describe, start_time, end_time").Scan(&res.Projects).Error
	err = p.orm.Model(&dbModel.TextResumeCollegeActivity{}).Where("resume_id = ?", res.ResumeId).
		Select("id, college, orgnization, position, describe, start_time, end_time").Scan(&res.Activities).Error
	err = p.orm.Model(&dbModel.TextResumeSocialPractice{}).Where("resume_id = ?", res.ResumeId).
		Select("id, name as practice_name, describe, start_time, end_time ").Scan(&res.Practices).Error
	err = p.orm.Model(&dbModel.TextResumeSkills{}).Where("resume_id = ?", res.ResumeId).
		Select("id, skill_name, describe").Scan(&res.Skills).Error
	err = p.orm.Model(&dbModel.TextResumeOther{}).Where("resume_id = ?", res.ResumeId).
		Select("id, title, describe").Scan(&res.Others).Error
	err = p.orm.Model(&dbModel.TextResumeEstimate{}).Where("resume_id = ?", res.ResumeId).
		Select("id, content").Scan(&res.SelfEstimate).Error


	return &res, nil
}

func (p *PersonDbOperator) BaseInfoAvatar(resumeId, name string, data []byte) error {

	// 图片数据上传到七牛云 TODO
	var icon = "http://pic34.photophoto.cn/20150120/0020033095117762_b.jpg"

	session := p.orm.Begin()
	// 存入数据库
	err := session.Where("resume_id = ?", resumeId).Assign(&dbModel.TextResumeBaseInfo{
		Icon: icon,
	}).FirstOrCreate(&dbModel.TextResumeBaseInfo{
		ResumeId: resumeId,
	}).Error
	if err != nil{
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func (p *PersonDbOperator) UpdateBaseInfo(resumeId string, id int,data map[string]interface{}) error{
	return  p.orm.Model(&dbModel.TextResumeBaseInfo{}).Where("resume_id = ? and id  = ?", resumeId, id).
		Updates(data).Error

}

func (p *PersonDbOperator) NewEducationInfo(resumeId string, data map[string]interface{}) (string,error){

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		endTime = t
	}

	var m = dbModel.TextResumeEducation{
		ResumeId: resumeId,
		College: data["college"].(string),
		Major: data["major"].(string),
		Rank: data["rank"].(string),
		Degree: data["degree"].(string),
		Describe: data["describe"].(string),
		StartTime: &startTime,
		EndTime: &endTime,
	}
	err :=  p.orm.Create(&m).Error
	if err != nil{
		return "", err
	}

	return strconv.Itoa(int(m.ID)), nil
}



func (p *PersonDbOperator) UpdateEducationInfo(resumeId string, id int, data map[string]interface{}) error {

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		endTime = t
	}

	data["start_time"] = &startTime
	data["end_time"] = &endTime

	return  p.orm.Model(&dbModel.TextResumeEducation{}).Where("resume_id = ? and id = ?", resumeId, id).
		Updates(data).Error
}

func (p *PersonDbOperator) DeleteEducation(resumeId, id string) error {
	return p.orm.Unscoped().Delete(&dbModel.TextResumeEducation{},"resume_id = ? and id = ?", resumeId, id).Error
}

func (p *PersonDbOperator) CreateWorkExperience(data map[string]interface{}) (string, error) {

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		endTime = t
	}
	var target  = dbModel.TextResumeWorkExperience{
		ResumeId: data["resume_id"].(string),
		CompanyName: data["company_name"].(string),
		WorkType: data["work_type"].(string),
		City: data["city"].(string),
		Position: data["position"].(string),
		Describe: data["describe"].(string),
		StartTime: &startTime,
		EndTime: &endTime,
	}
	err := p.orm.Create(&target).Error
	if err != nil{
		return  "", err
	}

	return strconv.Itoa(int(target.ID)), nil
}

func (p *PersonDbOperator) UpdateWorkExperience(id, resumeId string, data map[string]interface{}) error{
	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		endTime = t
	}

	data["start_time"] = &startTime
	data["end_time"] = &endTime


	return 	p.orm.Model(&dbModel.TextResumeWorkExperience{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}


func (p *PersonDbOperator) DeleteWorkExperience(id, resumeId string) error{
	return p.orm.Unscoped().
		Delete(&dbModel.TextResumeWorkExperience{}, "id = ? and resume_id = ?", id, resumeId).Error
}



func (p *PersonDbOperator) CreateProjectExperience(data map[string]interface{}) (string, error) {

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		endTime = t
	}
	var target  = dbModel.TextResumeProject{
		ResumeId: data["resume_id"].(string),
		ProjectName: data["project_name"].(string),
		ProjectLevel: data["project_level"].(string),
		Position: data["position"].(string),
		Describe: data["describe"].(string),
		StartTime: &startTime,
		EndTime: &endTime,
	}
	err := p.orm.Create(&target).Error
	if err != nil{
		return  "", err
	}

	return strconv.Itoa(int(target.ID)), nil
}

func (p *PersonDbOperator) UpdateProjectExperience(id, resumeId string, data map[string]interface{}) error{
	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		endTime = t
	}

	data["start_time"] = &startTime
	data["end_time"] = &endTime


	return 	p.orm.Model(&dbModel.TextResumeProject{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}


func (p *PersonDbOperator) DeleteProjectExperience(id, resumeId string) error{
	return p.orm.Unscoped().
		Delete(&dbModel.TextResumeProject{}, "id = ? and resume_id = ?", id, resumeId).Error
}





func (p *PersonDbOperator) CreateCollegeActive(data map[string]interface{}) (string, error) {

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		endTime = t
	}
	var target  = dbModel.TextResumeCollegeActivity{
		ResumeId: data["resume_id"].(string),
		College: data["college"].(string),
		Orgnization: data["orgnization"].(string),
		Position: data["position"].(string),
		Describe: data["describe"].(string),
		StartTime: &startTime,
		EndTime: &endTime,
	}
	err := p.orm.Create(&target).Error
	if err != nil{
		return  "", err
	}

	return strconv.Itoa(int(target.ID)), nil
}

func (p *PersonDbOperator) UpdateCollegeActive(id, resumeId string, data map[string]interface{}) error{
	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		endTime = t
	}

	data["start_time"] = &startTime
	data["end_time"] = &endTime

	return 	p.orm.Model(&dbModel.TextResumeCollegeActivity{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}


func (p *PersonDbOperator) DeleteCollegeActive(id, resumeId string) error{
	return p.orm.Unscoped().
		Delete(&dbModel.TextResumeCollegeActivity{}, "id = ? and resume_id = ?", id, resumeId).Error
}




func (p *PersonDbOperator) CreateSocialPractice(data map[string]interface{}) (string, error) {

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return "", err
		}
		endTime = t
	}
	var target  = dbModel.TextResumeSocialPractice{
		ResumeId: data["resume_id"].(string),
		Name: data["practice_name"].(string),
		Describe: data["describe"].(string),
		StartTime: &startTime,
		EndTime: &endTime,
	}
	err := p.orm.Create(&target).Error
	if err != nil{
		return  "", err
	}

	return strconv.Itoa(int(target.ID)), nil
}

func (p *PersonDbOperator) UpdateSocialPractice(id, resumeId string, data map[string]interface{}) error{

	var startTime, endTime  time.Time

	if s, ok := data["start_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		startTime = t
	}
	if s, ok := data["end_time"].(string); ok{
		t, err := time.Parse(utils2.RESUME_TIME_FORMAT,s)
		if err != nil{
			return  err
		}
		endTime = t
	}

	data["start_time"] = &startTime
	data["end_time"] = &endTime


	return 	p.orm.Model(&dbModel.TextResumeSocialPractice{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}


func (p *PersonDbOperator) DeleteSocialPractice(id, resumeId string) error{
	return p.orm.Unscoped().
		Delete(&dbModel.TextResumeSocialPractice{}, "id = ? and resume_id = ?", id, resumeId).Error
}



func (p *PersonDbOperator) CreateResumeSkill(data map[string]interface{}) (string, error) {


	var target  = dbModel.TextResumeSkills{
		ResumeId: data["resume_id"].(string),
		SkillName: data["skill_name"].(string),
		Describe: data["describe"].(string),

	}
	err := p.orm.Create(&target).Error
	if err != nil{
		return  "", err
	}

	return strconv.Itoa(int(target.ID)), nil
}

func (p *PersonDbOperator) UpdateResumeSkill(id, resumeId string, data map[string]interface{}) error{

	return 	p.orm.Model(&dbModel.TextResumeSkills{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}


func (p *PersonDbOperator) DeleteResumeSkill(id, resumeId string) error{
	return p.orm.Unscoped().
		Delete(&dbModel.TextResumeSkills{}, "id = ? and resume_id = ?", id, resumeId).Error
}


func (p *PersonDbOperator) CreateResumeOther(data map[string]interface{}) (string, error) {


	var target  = dbModel.TextResumeOther{
		ResumeId: data["resume_id"].(string),
		Title: data["title"].(string),
		Describe: data["describe"].(string),

	}
	err := p.orm.Create(&target).Error
	if err != nil{
		return  "", err
	}

	return strconv.Itoa(int(target.ID)), nil
}

func (p *PersonDbOperator) UpdateResumeOther(id, resumeId string, data map[string]interface{}) error{
	return 	p.orm.Model(&dbModel.TextResumeOther{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}


func (p *PersonDbOperator) DeleteResumeOther(id, resumeId string) error{
	return p.orm.Unscoped().
		Delete(&dbModel.TextResumeOther{}, "id = ? and resume_id = ?", id, resumeId).Error
}



//func (p *PersonDbOperator) CreateResumeEstimate(data map[string]interface{}) (string, error) {
//
//
//	var target  = dbModel.TextResumeEstimate{
//		ResumeId: data["resume_id"].(string),
//		Content: data["content"].(string),
//	}
//	err := p.orm.Create(&target).Error
//	if err != nil{
//		return  "", err
//	}
//
//	return strconv.Itoa(int(target.ID)), nil
//}

func (p *PersonDbOperator) UpdateResumeEstimate(id, resumeId string, data map[string]interface{}) error{
	return 	p.orm.Model(&dbModel.TextResumeEstimate{}).Where("id = ? and resume_id = ?", id, resumeId).
		Updates(data).Error
}

func (p *PersonDbOperator) AttachResume(resumeId string) (string, error) {
	var url struct{
		FileUrl string
	}
	err := p.orm.Model(&dbModel.AttachFileResume{}).Where("resume_id = ?", resumeId).
		Select("file_url").Scan(&url).Error

	if err != nil{
		return "", err
	}
	return url.FileUrl, nil
}


//func (p *PersonDbOperator) DeleteResumeEstimate(id, resumeId string) error{
//	return p.orm.Unscoped().
//		Delete(&dbModel.TextResumeEstimate{}, "id = ? and resume_id = ?", id, resumeId).Error
//}

func NewPersonDbOperator() *PersonDbOperator {
	return &PersonDbOperator{
		orm: orm.DB,
	}
}
