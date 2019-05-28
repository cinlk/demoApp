package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/model/httpModel"
	"demoApp/server/utils"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type briefInfoReq struct {
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	College string `json:"college"`
}

type deliveryHistory struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type resumeNewNameReq struct {
	Name     string `json:"name"`
	ResumeId string `json:"resume_id"`
}

type resumeBaseInfoReq struct {
	Id string `json:"id" binding:"required"`
	ResumeId string `json:"resume_id" binding:"required"`
	Name string `json:"name,omitempty"`
	College string `json:"college,omitempty"`
	Gender string `json:"gender,omitempty"`
	City   string `json:"city,omitempty"`
	Degree  string `json:"degree,omitempty"`
	Birthday string `json:"birthday,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
}

type resumeEducationInfoReq struct {


	ResumeId string `json:"resume_id" binding:"required"`
	College string `json:"college,omitempty"`
	Major string `json:"major,omitempty"`
	Rank string `json:"rank,omitempty"`
	// 学历
	Degree string `json:"degree,omitempty"`
	// 描述
	Describe string `json:"describe,omitempty"`
	// YYYY-MM 格式
	StartTime string `json:"start_time,omitempty"`
	EndTime string  `json:"end_time,omitempty"`
}

type resumeWorkExperienceReq struct {

	ResumeId string `json:"resume_id,omitempty"  binding:"required"`
	CompanyName string `json:"company_name,omitempty"`
	WorkType string `json:"work_type,omitempty"`
	City string `json:"city,omitempty"`
	Position string `json:"position,omitempty"`
	// 描述
	Describe string `json:"describe,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime string  `json:"end_time,omitempty"`
}

type resumeProjectReq struct {

	ResumeId string `json:"resume_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	// 人数规模 ??
	ProjectLevel  string `json:"project_level,omitempty"`
	Position string `json:"position,omitempty"`
	Describe string `json:"describe,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime string  `json:"end_time,omitempty"`
}

type resumeCollegeActiveReq struct {

	ResumeId string `json:"resume_id,omitempty"`
	College string `json:"college,omitempty"`
	Orgnization string `json:"orgnization,omitempty"`
	Position string `json:"position,omitempty"`
	Describe string `json:"describe,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime string  `json:"end_time,omitempty"`

}


type resumeSocialPracticeReq struct {
	ResumeId string `json:"resume_id,omitempty"`
	PracticeName string  `json:"practice_name,omitempty"`
	Describe string `json:"describe,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime string  `json:"end_time,omitempty"`
}


type resumeSkillReq struct {

	ResumeId string `json:"resume_id,omitempty"`
	SkillName string `json:"skill_name,omitempty"`
	Describe string `json:"describe,omitempty"`
}

type resumeOtherReq struct {

	ResumeId string `json:"resume_id,omitempty"`
	Title string `json:"title,omitempty"`
	Describe string `json:"describe,omitempty"`
}

type resumeEsitmateReq struct {
	ResumeId string `json:"resume_id,omitempty"`
	Content string `json:"content,omitempty"`
}

type collectedJobReq struct {
	Type string `json:"type"`
	Offset int `json:"offset"`
	Limit int  `json:"limit"`
}

type collectedReq struct {
	Offset int `json:"offset"`
	Limit int  `json:"limit"`
}


type unCollectedJobsReq struct {
	Type string `json:"type"`
	JobIds []string `json:"job_ids" binding:"required"`

}

type unCollectedTargetReq struct {
	Ids []string `json:"ids"`
}


type renamePostGroupNameReq struct {
	GroupId string `json:"group_id" binding:"required"`
	NewName string `json:"new_name" binding:"required"`
}

type jobSubscribeReq struct {
	Type string `json:"type" binding:"required"`
	Fields  string `json:"fields" binding:"required"`
	Citys   []string `json:"citys" binding:"required"`
	Degree  string `json:"degree" binding:"required"`
	InternDay     string `json:"intern_day"`
	InternMonth   string `json:"intern_month"`
	InternSalary  string `json:"intern_salary"`
	Salary  string `json:"salary"`
}

type personHandler struct {
	baseHandler
	UrlPrefix string
	validate  handlerValider
	db        *dbOperater.PersonDbOperator
}

func (p *personHandler) updateAvatar(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

	var userId = r.Header.Get(utils.USER_ID)
	// image 数据和名称

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, errors.Wrap(err, "can't get image file").Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	_, err = ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, errors.Wrap(err, "can't read file").Error(), http.StatusBadRequest)
		return
	}

	// 上传到七牛云 获取url TODO
	// 图片名称存入数据库 TODO

	p.db.Avatar(userId, header.Filename)

	// 返回数据 test
	var newIcon = "http://pic-public.yihu.bingfengtech.com/demo.jpeg"
	time.Sleep(time.Second * 3)
	p.JSON(w, httpModel.HttpPersonAvatarModel{
		IconUrl: newIcon,
	}, http.StatusAccepted)

}

// 跟新用户简要信息

func (p *personHandler) BriefInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req briefInfoReq
	var userId = r.Header.Get(utils.USER_ID)
	err := p.validate.Validate(r, &req)
	if err != nil {
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	err = p.db.BriefInfos(userId, req.Name, req.Gender, req.College)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, "", http.StatusAccepted)

}

// 查询投递记录(网申的职位，校招和实习职位)
func (p *personHandler) DeliveryJobList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//var req deliveryHistory
	//err := p.validate.Validate(r, &req)
	//if err != nil {
	//	p.ERROR(w, err, http.StatusBadRequest)
	//	return
	//}
	//time.Sleep(time.Second * 3)
	var userId = r.Header.Get(utils.USER_ID)
	res := p.db.FindDeliveryInfos(userId)
	p.JSON(w, res, http.StatusOK)

}

// 职位投递的历史状态记录
func (p *personHandler) JobDeliveryHistoryStatus(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var t = param.ByName("type")
	var jobId = param.ByName("jobId")
	res, err := p.db.JobDeliveryHistory(userId, jobId, t)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)
}

// 根据position id 查询online apply id

func (p *personHandler) FindOnlineApplyId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var positionId = param.ByName("positionId")
	oid := p.db.FindOnlineApplyIdBy(positionId)
	if oid == "" {
		p.ERROR(w, errors.New("not found online apply id"), http.StatusNotFound)
		return
	}

	p.JSON(w, map[string]interface{}{
		"online_apply_id": oid,
	}, http.StatusOK)
}

//  获取当前的简历列表数据
func (p *personHandler) MyResumes(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var userId = r.Header.Get(utils.USER_ID)

	res, err := p.db.GetMyResumes(userId)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)

}

// 创建新的文本简历
func (p *personHandler) createTextResume(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)

	res, err := p.db.NewResume(userId)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)

}

// 创建新的附件简历
func (p *personHandler) createNewAttachResume(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	// 附件数据
	file, header, err := r.FormFile("attach")
	if err != nil {
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		p.ERROR(w, errors.Wrap(err, "not found file data"), http.StatusBadRequest)
		return
	}

	fileName := header.Filename

	err = p.db.CreateAttachFile(data, fileName, userId)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)

}

// 设置默认投递简历
func (p *personHandler) primaryResume(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var userId = r.Header.Get(utils.USER_ID)
	var rid = param.ByName("resumeId")
	err := p.db.SetPrimaryResume(userId, rid)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)

}

func (p *personHandler) resumeName(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var req resumeNewNameReq
	err := p.validate.Validate(r, &req)
	if err != nil {
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	err = p.db.ChangeResumeName(userId, req.ResumeId, req.Name)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteResume(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var resumeId = param.ByName("resumeId")
	var t = param.ByName("type")

	err := p.db.DeleteResume(userId, resumeId, t)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)

}

// 复制文本简历 TODO
func (p *personHandler) CopyTextResume(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}



// 获取文本简历 相关数据
func (p *personHandler) TextResumeInfo(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var resumeId = param.ByName("resumeId")

	res ,err := p.db.TextResumeInfo(userId, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)
}


// 更新文本简历 头像
func(p *personHandler) BaseInfoAvatar(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var resumeId = param.ByName("resumeId")

	file, header, err := r.FormFile("attach")
	if err != nil {
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		p.ERROR(w, errors.Wrap(err, "not found file data"), http.StatusBadRequest)
		return
	}
	var name = header.Filename
	err = p.db.BaseInfoAvatar(resumeId, name, data)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpPersonAvatarModel{
		IconUrl: "http://pic34.photophoto.cn/20150120/0020033095117762_b.jpg",
	}, http.StatusAccepted)
}

// 更新文本简历的内容
func (p *personHandler) BaseInfoContent(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	var req resumeBaseInfoReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	id, _  := strconv.Atoi(req.Id)
	var m = utils.Struct2Map(req)
	err = p.db.UpdateBaseInfo(req.ResumeId, id, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}



func (p *personHandler) NewEducationInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var req resumeEducationInfoReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}


	var m = utils.Struct2Map(req)
	if req.College == "" || req.Major == "" || req.Degree == "" || req.Describe == "" || req.StartTime == "" ||
		req.EndTime == ""{
		p.ERROR(w, errors.New("invalidate data"), http.StatusBadRequest)
		return
	}
	id , err := p.db.NewEducationInfo(req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) UpdateEducationInfo(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var req resumeEducationInfoReq
	err := p.validate.Validate(r, &req)
	var id = param.ByName("id")
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	educationId, err := strconv.Atoi(id)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)

	err = p.db.UpdateEducationInfo(req.ResumeId, educationId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}



func (p *personHandler) DeleteEducationInfo(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var resumedId = param.ByName("resumeId")
	var id = param.ByName("id")

	err := p.db.DeleteEducation(resumedId, id)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)

}



func (p *personHandler) newWorkExperience(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var req resumeWorkExperienceReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)
	id, err := p.db.CreateWorkExperience(*m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) updateWorkExperience(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeWorkExperienceReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	err = p.db.UpdateWorkExperience(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteWorkExperience(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var id = para.ByName("id")
	var resumeId = para.ByName("resumeId")

	err := p.db.DeleteWorkExperience(id, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}





func (p *personHandler) newProjectExperience(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var req resumeProjectReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)
	id, err := p.db.CreateProjectExperience(*m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) updateProjectExperience(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeProjectReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	err = p.db.UpdateProjectExperience(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteProjectExperience(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var id = para.ByName("id")
	var resumeId = para.ByName("resumeId")

	err := p.db.DeleteProjectExperience(id, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}




func (p *personHandler) newCollegeActive(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var req resumeCollegeActiveReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)
	id, err := p.db.CreateCollegeActive(*m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) updateCollegeActive(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeCollegeActiveReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	err = p.db.UpdateCollegeActive(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteCollegeActive(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var id = para.ByName("id")
	var resumeId = para.ByName("resumeId")

	err := p.db.DeleteCollegeActive(id, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}





func (p *personHandler) newResumeSkill(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var req resumeSkillReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)
	id, err := p.db.CreateResumeSkill(*m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) updateResumeSkill(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeSkillReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	err = p.db.UpdateResumeSkill(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteResumeSkill(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var id = para.ByName("id")
	var resumeId = para.ByName("resumeId")

	err := p.db.DeleteResumeSkill(id, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}




func (p *personHandler) newSocialPractice(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var req resumeSocialPracticeReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)
	id, err := p.db.CreateSocialPractice(*m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) updateSocialPractice(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeSocialPracticeReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)
	fmt.Println(m)
	err = p.db.UpdateSocialPractice(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteSocialPractice(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var id = para.ByName("id")
	var resumeId = para.ByName("resumeId")

	err := p.db.DeleteSocialPractice(id, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}



func (p *personHandler) newResumeOther(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var req resumeOtherReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	var m = utils.Struct2Map(req)
	id, err := p.db.CreateResumeOther(*m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)
}

func (p *personHandler) updateResumeOther(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeOtherReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	err = p.db.UpdateResumeOther(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

func (p *personHandler) deleteResumeOther(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var id = para.ByName("id")
	var resumeId = para.ByName("resumeId")

	err := p.db.DeleteResumeOther(id, resumeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}



//func (p *personHandler) newResumeEstimate(w http.ResponseWriter, r *http.Request, para httprouter.Params){
//	var req resumeEsitmateReq
//	err := p.validate.Validate(r, &req)
//	if err != nil{
//		p.ERROR(w, err, http.StatusBadRequest)
//		return
//	}
//
//	var m = utils.Struct2Map(req)
//	id, err := p.db.CreateResumeEstimate(*m)
//	if err != nil{
//		p.ERROR(w, err, http.StatusUnprocessableEntity)
//		return
//	}
//
//	p.JSON(w, map[string]interface{}{
//		"id": id,
//	}, http.StatusCreated)
//}

func (p *personHandler) updateResumeEstimate(w http.ResponseWriter, r *http.Request, para httprouter.Params){
	var id = para.ByName("id")
	var req  resumeEsitmateReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	err = p.db.UpdateResumeEstimate(id, req.ResumeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}

//func (p *personHandler) deleteResumeEstimate(w http.ResponseWriter, r *http.Request, para httprouter.Params){
//
//	var id = para.ByName("id")
//	var resumeId = para.ByName("resumeId")
//
//	err := p.db.DeleteResumeEstimate(id, resumeId)
//	if err != nil{
//		p.ERROR(w, err, http.StatusUnprocessableEntity)
//		return
//	}
//	p.JSON(w, httpModel.HttpResultModel{
//		Result: "success",
//	}, http.StatusCreated)
//}

// 附件简历的url 地址
func (p *personHandler) attachResumeUrl(w http.ResponseWriter, r *http.Request, para httprouter.Params){

	var resumeId = para.ByName("resumeId")

	url, err := p.db.AttachResume(resumeId)
	if err != nil{
		p.ERROR(w, err,  http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, map[string]interface{}{
		"url": url,
	}, http.StatusOK)

}


// 收藏

func (p *personHandler) collectedJobs(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req collectedJobReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var res []httpModel.CollectedJobModel
	switch req.Type {
	case "intern":
		res, err = p.db.CollectedInternJobs(userId,req.Offset, req.Limit)
	case "graduate":
		res, err = p.db.CollecteCampusJobs(userId, req.Offset, req.Limit)
	default:
		break
	}

	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)
}


func (p *personHandler) collectedCompany(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req collectedReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res, err := p.db.CollectedCompany(userId, req.Offset, req.Limit)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, res, http.StatusOK)
}

func (p *personHandler) collectedOnlineApply(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req collectedReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res, err := p.db.CollectedOnlineApply(userId, req.Offset, req.Limit)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, res, http.StatusOK)

}

func (p *personHandler) collectedCareerTalk(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req collectedReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res, err := p.db.CollectedCareerTalk(userId, req.Offset, req.Limit)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)


}

func (p *personHandler) collectedPost(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)

	res, err := p.db.CollectedPost(userId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, res, http.StatusOK)
}


func (p *personHandler) unSubscribeCollectedJobs(w http.ResponseWriter, r *http.Request, param httprouter.Params){

	var userId = r.Header.Get(utils.USER_ID)
	var req unCollectedJobsReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	if len(req.JobIds) ==  0 {
		p.ERROR(w, errors.New("empty job list"), http.StatusBadRequest)
		return
	}

	err = p.db.UnCollectedJobs(userId, req.Type, req.JobIds)
	if err != nil{
		p.ERROR(w,err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result:"success",
	}, http.StatusAccepted)
}


func (p *personHandler) unSubScribeCollectedCompany(w http.ResponseWriter, r *http.Request, param httprouter.Params){

	var userId = r.Header.Get(utils.USER_ID)
	var req unCollectedTargetReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err ,http.StatusBadRequest)
		return
	}

	err = p.db.UnCollectedCompany(userId, req.Ids)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusOK)
}


func (p *personHandler) unSubScribeCollectedCareerTalk(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req unCollectedTargetReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err ,http.StatusBadRequest)
		return
	}

	err = p.db.UnCollectedCareerTalk(userId, req.Ids)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusOK)
}

func (p *personHandler) unSubScribeCollectedOnlineApply(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req unCollectedTargetReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err ,http.StatusBadRequest)
		return
	}

	err = p.db.UnCollectedOnlineApply(userId, req.Ids)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusOK)
}





func (p *personHandler) removePostGroup(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var name = param.ByName("name")

	err := p.db.RemovePostGroup(userId, name)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}



func (p *personHandler) renamePostGroup(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req renamePostGroupNameReq

	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err ,http.StatusBadRequest)
		return
	}

	err = p.db.RenamePostGroup(userId, req.GroupId, req.NewName)
	if err != nil{
		p.ERROR(w, err ,http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)


}

func (p *personHandler) myJobSubscribeCondition(w http.ResponseWriter, r *http.Request, param httprouter.Params)  {
	var userId  = r.Header.Get(utils.USER_ID)
	res, err := p.db.MySubscribeJobCondition(userId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, res, http.StatusOK)
}

// 职位订阅
func (p *personHandler) createJobSubscribe(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var req jobSubscribeReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)

	id, err := p.db.CreateJobSubscribe(userId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusCreated)

}

func (p *personHandler) updateJobSubscribe(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var subscribeId = param.ByName("subscribeId")
	var req jobSubscribeReq
	err := p.validate.Validate(r, &req)
	if err != nil{
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var m = utils.Struct2Map(req)
	err = p.db.UpdateJobSubscribe(userId, subscribeId, *m)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)

}


func (p *personHandler) deleteJobSubscribe(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var subscribeId = param.ByName("subscribeId")
	err := p.db.DeleteJobSubscribe(userId, subscribeId)
	if err != nil{
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusAccepted)
}