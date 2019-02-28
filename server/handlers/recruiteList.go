package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)


type onlineQueryPage struct {
	queryPage
	Citys  []string `json:"citys"`
	BusinessField string `json:"business_field"`
}

type careerTalkQueryPage struct {
	queryPage
	College []string `json:"college"`
	Time  string `json:"time"`
	BusinessField string `json:"business_field"`
	// 如果没有选择学校，默认该城市所有学校或全国的学校
	City string `json:"city"`
}

type companyQueryPage struct {
	queryPage
	Citys []string `json:"citys"`
	BusinessField string `json:"business_field"`
	CompanyType string `json:"company_type"`
}


type graduateQueryPage struct {
	queryPage
	Citys []string `json:"citys"`
	SubBusinessField string `json:"sub_business_field"`
	Degree string `json:"degree"`
}

type internQueryPage struct {
	queryPage
	Citys  []string `json:"citys"`
	SubBusinessField string `json:"sub_business_field"`
	InternCondition map[string]interface{} `json:"intern_condition"`

}


type recruiteListHandle struct {

	baseHandler
	UrlPrefix string
	validate handlerValider
	dbOperator *dbOperater.RecruiteDboperator

}


func (re *recruiteListHandle) onlineApplys(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	var req onlineQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil{
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := re.dbOperator.ListApplyOnlines(req.Citys, req.BusinessField, req.Offset, req.Limit)
	re.JSON(w, res, http.StatusOK)
}


func (re *recruiteListHandle) findOnlineApply(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	 id := param.ByName("id")
	 if id == ""{
	 	re.ERROR(w, errors.New("empty id"), http.StatusBadRequest)
		 return
	 }

	 userId :=  r.Header.Get(utils.USER_ID)
	 res := re.dbOperator.OnlineApplyInfo(id, userId)

	 re.JSON(w, res, http.StatusOK)
}



func (re *recruiteListHandle) careerTalks(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var req careerTalkQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil{
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.ListCareerTalk(req.College, req.BusinessField,
		req.City, req.Time, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) findCareerTalk(w http.ResponseWriter, r *http.Request, params httprouter.Params){

}

func (re *recruiteListHandle) companys(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	var req companyQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil{
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := re.dbOperator.ListCompany(req.Citys, req.BusinessField, req.CompanyType, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)

}


func (re *recruiteListHandle) findCompany(w http.ResponseWriter, r *http.Request, param httprouter.Params){

}

func (re *recruiteListHandle) graduatejobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var req graduateQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil{
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.ListGraduateJobs(req.Citys, req.SubBusinessField, req.Degree, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)
}

func (re *recruiteListHandle) findGraduate(w http.ResponseWriter, r *http.Request, param httprouter.Params){

}


func (re *recruiteListHandle) internjobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	var req internQueryPage
	err := re.validate.Validate(r, &req)
	if err != nil{
		re.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := re.dbOperator.ListInternJobs(req.InternCondition, req.Citys, req.SubBusinessField, req.Offset, req.Limit)

	re.JSON(w, res, http.StatusOK)
}


func (re *recruiteListHandle) findInternJob(w http.ResponseWriter, r *http.Request, param httprouter.Params){

}
