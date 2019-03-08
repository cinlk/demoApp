package handlers

import (
	"demoApp/server/model/dbOperater"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)

type pageQuery struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type jobHandler struct {
	baseHandler
	validate   *baseValidate
	UrlPrefix  string
	dbOperator *dbOperater.JobDbOperator
}

func (j *jobHandler) FindJobKind(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

	// TODO verify in path
	kind := para.ByName("kind")
	if kind == "" {
		j.ERROR(w, errors.New("job kind is empty"), http.StatusBadRequest)
		return
	}
	var query pageQuery
	err := j.validate.Validate(r, &query)
	if err != nil {
		j.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := j.dbOperator.JobsKind(kind, query.Offset, query.Limit)
	j.JSON(w, res, http.StatusAccepted)

	//res := j.dbOperator.JobsKind(kind, query.Offset, query.Limit)

	//var httpres []httpModel.HttpGraduateModel
	//for _, item := range res {
	//	m := httpModel.HttpGraduateModel{}
	//	m.Id = item.Id
	//	m.Recruiter = nil
	//	m.Validate = item.IsValidate
	//	m.Applied = item.UserApplyJob.Applied
	//	m.ApplyEndTime = item.ApplyEndTime.Unix()
	//	for _, b := range item.BussinesField {
	//		m.BussinesField = append(m.BussinesField, b)
	//	}
	//	//m.BussinesField = item.BussinesField
	//	m.Benefits = item.Benefits
	//	m.Collected = item.UserApplyJob.Collected
	//	m.Comapany = nil
	//	m.CreatedTime = item.CreatedTime.Unix()
	//	m.Education = item.Education
	//	m.Icon = item.Icon
	//	m.Link = item.Link
	//	for _, b := range item.LocationCity {
	//		m.LocationCity = append(m.LocationCity, b)
	//	}
	//	for _, b := range item.Major {
	//		m.Major = append(m.Major, b)
	//	}
	//
	//	m.Name = item.Name
	//	m.NeedSkills = item.NeedSkills
	//	m.WorkContent = item.WorkContent
	//	m.ReviewCounts = item.ReviewCounts
	//	m.Salary = item.Salary
	//	m.Tags = []string{}
	//
	//	m.Talked = item.UserApplyJob.Talked
	//	m.Type = string(item.Type)
	//
	//	httpres = append(httpres, m)
	//}

	//j.JSON(w, &httpres, http.StatusAccepted)
}
