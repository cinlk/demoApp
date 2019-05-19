package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/model/httpModel"
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type onlineQueryPage struct {
	queryPage
	Citys         []string `json:"citys"`
	BusinessField string   `json:"business_field"`
	// 搜索某个类型
	TypeField string `json:"type_field"`
}

type careerTalkQueryPage struct {
	queryPage
	College       []string `json:"college"`
	Time          string   `json:"time"`
	BusinessField string   `json:"business_field"`
	// 如果没有选择学校，默认该城市所有学校或全国的学校
	City string `json:"city"`
}

type companyQueryPage struct {
	queryPage
	Citys         []string `json:"citys"`
	BusinessField string   `json:"business_field"`
	CompanyType   string   `json:"company_type"`
}

type graduateQueryPage struct {
	queryPage
	Citys            []string `json:"citys"`
	SubBusinessField string   `json:"sub_business_field"`
	Degree           string   `json:"degree"`
}

type internQueryPage struct {
	queryPage
	Citys            []string               `json:"citys"`
	SubBusinessField string                 `json:"sub_business_field"`
	InternCondition  map[string]interface{} `json:"intern_condition"`
}

type companyTagJobQuery struct {
	queryPage
	Tag       string `json:"tag"`
	CompanyID string `json:"company_id" binding:"required"`
}

type companyRecruitMeeting struct {
	queryPage
	CompanyID string `json:"company_id" binding:"required"`
}

type recruiteListHandle struct {
	baseHandler
	UrlPrefix  string
	validate   handlerValider
	dbOperator *dbOperater.RecruiteDboperator
}

func (re *recruiteListHandle) onlineApplys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req onlineQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := re.dbOperator.ListApplyOnlines(req.Citys, req.BusinessField, req.TypeField, req.Offset, req.Limit)
	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) findOnlineApply(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	id := param.ByName("id")
	if id == "" {
		re.ERROR(w, errors.New("empty id"), http.StatusBadRequest)
		return
	}
	userId := r.Header.Get(utils.USER_ID)
	res := re.dbOperator.OnlineApplyInfo(id, userId)

	re.JSON(w, res, http.StatusOK)
}

// 申请网申职位
func (re *recruiteListHandle) applyOnlineJob(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var onlineApplyId = param.ByName("onlineId")
	var positiobId = param.ByName("positionId")
	pid, err := strconv.Atoi(positiobId)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	b, err := re.dbOperator.ApplyOnlieJob(userId, onlineApplyId, pid)
	if err != nil {
		re.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	re.JSON(w, map[string]interface{}{
		"exist": b,
	}, http.StatusCreated)
}

func (re *recruiteListHandle) careerTalks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req careerTalkQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.ListCareerTalk(req.College, req.BusinessField,
		req.City, req.Time, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) findCareerTalk(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var id = param.ByName("id")
	if id == "" {
		re.ERROR(w, errors.New("empty job id"), http.StatusBadRequest)
		return
	}

	userId := r.Header.Get(utils.USER_ID)
	res := re.dbOperator.RecruitMeetingInfo(id, userId)

	re.JSON(w, res, http.StatusOK)

}

func (re *recruiteListHandle) companys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req companyQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := re.dbOperator.ListCompany(req.Citys, req.BusinessField, req.CompanyType, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)

}

func (re *recruiteListHandle) findCompany(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var id = param.ByName("id")
	if id == "" {
		re.ERROR(w, errors.New("empty company id"), http.StatusBadRequest)
		return
	}

	userId := r.Header.Get(utils.USER_ID)
	res := re.dbOperator.CompanyInfo(id, userId)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) graduatejobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req graduateQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.ListGraduateJobs(req.Citys, req.SubBusinessField, req.Degree, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) findGraduate(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var id = param.ByName("id")
	if id == "" {
		re.ERROR(w, errors.New("empty job id"), http.StatusBadRequest)
		return
	}

	userId := r.Header.Get(utils.USER_ID)
	res := re.dbOperator.GraduatJobInfo(id, userId)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) internjobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req internQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.ListInternJobs(req.InternCondition, req.Citys, req.SubBusinessField, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) findInternJob(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var id = param.ByName("id")
	if id == "" {

		re.ERROR(w, errors.New("empty job id"), http.StatusBadRequest)
		return
	}
	user := r.Header.Get(utils.USER_ID)
	res := re.dbOperator.InternJobInfo(id, user)

	re.JSON(w, res, http.StatusOK)

}

// 申请校招和实习职位
func (re *recruiteListHandle) applyJob(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var jobId = param.ByName("jobId")
	var t = param.ByName("type")

	err := re.dbOperator.ApplyJob(userId, jobId, t)
	if err != nil {
		re.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	re.JSON(w, httpModel.HttpResultModel{
		Result: "success",
	}, http.StatusCreated)
}

// 根据招聘者信息 查找他所在的公司和发布的职位
func (re *recruiteListHandle) recruiterCompanyAndJobs(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var id = param.ByName("id")
	if id == "" {
		re.ERROR(w, errors.New("empty user id"), http.StatusBadRequest)
		return
	}

	res := re.dbOperator.FindRecruiterInfo(id)

	re.JSON(w, res, http.StatusOK)

}

func (re *recruiteListHandle) companyTagJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req companyTagJobQuery
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.CompanyRelatedJobs(req.CompanyID, req.Tag, req.Offset, req.Limit)
	re.JSON(w, res, http.StatusOK)

}

func (re *recruiteListHandle) companyRecruitMeeting(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req companyRecruitMeeting
	err := re.validate.Validate(r, &req)
	if err != nil {
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := re.dbOperator.CompanyRelatedCareerTalk(req.CompanyID, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)

}
