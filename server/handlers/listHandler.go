package handlers

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/dbOperater"
	"demoApp/server/model/httpModel"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)

const (
	meetings = "meetings"
	company  = "company"
)

type nearByReq struct {
	Distance  uint    `json:"distance"`
	Latitude  float64 `json:"latitude"`
	Lontitude float64 `json:"lontitude"`
	Type      string  `json:"type"`
}

type queryPage struct {
	Offset int64 `json:"offset" binding:"required"`
	Limit  int64 `json:"limit"  binding:"required"`
}

type recommands struct {
	News          []dbModel.LatestNews                `json:"news"`
	JobCategory   []dbModel.JobCategory               `json:"job_category"`
	TopJobs       []dbModel.TopJobs                   `json:"top_jobs"`
	CareerTalk    []httpModel.HttpCareerTalkListModel `json:"career_talk"`
	ApplyClassify []dbModel.ApplyClassify             `json:"apply_classify"`
}

type listHandler struct {
	baseHandler
	UrlPrefix string
	db        *dbOperater.ListDboperater
	validate  handlerValider
}

// 个性化 TODO
func (l *listHandler) bannerInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	banners := l.db.Banners(5)

	l.JSON(w, banners, http.StatusOK)

}

func (l *listHandler) latestNews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	news := l.db.LatestNews()
	l.JSON(w, news, http.StatusOK)
}

func (l *listHandler) jobCategories(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	categories := l.db.LatestJobCategory(5)
	l.JSON(w, categories, http.StatusOK)
}

func (l *listHandler) jobTops(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jobs := l.db.TopJobs()
	l.JSON(w, jobs, http.StatusOK)
}

// TODO 个性化推荐
func (l *listHandler) careerTalks(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	talks := l.db.CarrerTalks(12)
	l.JSON(w, talks, http.StatusOK)
}

// TODO 热门 或 推荐
func (l *listHandler) onlineApply(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	apply := l.db.OnlineApplyClass()
	l.JSON(w, apply, http.StatusOK)
}

// TODO 个性化推荐
func (l *listHandler) personalityJobs(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var page queryPage
	err := l.validate.Validate(r, &page)
	if err != nil {
		l.ERROR(w, errors.Wrap(err, "bad request"), http.StatusBadRequest)
		return
	}

	jobs := l.db.JobList(page.Offset, page.Limit)
	l.JSON(w, jobs, http.StatusOK)
}

// TODO 推荐 个性化
func (l *listHandler) personalRecommand(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var res = recommands{}

	res.News = l.db.LatestNews()
	res.ApplyClassify = l.db.OnlineApplyClass()
	res.TopJobs = l.db.TopJobs()
	res.JobCategory = l.db.LatestJobCategory(10)
	res.CareerTalk = l.db.CarrerTalks(16)

	l.JSON(w, res, http.StatusOK)

	//l.JSON(w, nil, http.StatusOK)
}

// TODO 附近的
func (l *listHandler) nearBy(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var req nearByReq
	err := l.validate.Validate(r, &req)
	if err != nil {
		l.ERROR(w, err, http.StatusBadRequest)
		return
	}

	switch req.Type {
	case meetings:
		res, err := l.db.NearByMeetings(req.Latitude, req.Lontitude, req.Distance)
		if err != nil {
			l.ERROR(w, err, http.StatusUnprocessableEntity)
			return
		}
		l.JSON(w, res, http.StatusOK)

	case company:

		res, err := l.db.NearyByCompany(req.Latitude, req.Lontitude, req.Distance)
		if err != nil {
			l.ERROR(w, err, http.StatusUnprocessableEntity)
			return
		}
		l.JSON(w, res, http.StatusOK)

	default:
		l.ERROR(w, errors.New(fmt.Sprintf("invalidate type %s", req.Type)), http.StatusBadRequest)
		return
	}

}
