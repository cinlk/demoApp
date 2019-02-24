package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/model/httpModel"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type searchTargetReq struct {
	Word string `json:"word" bind:"required"`
}

type searchQuery struct {
	Type string `json:"type"`
	Word string `json:"word"`
}

type searchHandler struct {
	baseHandler
	validate  handlerValider
	db        *dbOperater.SearchDboperator
	UrlPrefix string
}

func (s *searchHandler) TopWords(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	t := param.ByName("type")
	if t == "" {
		s.ERROR(w, errors.New("invalidate type"), http.StatusBadRequest)
		return
	}

	res := s.db.TopWorks(t)

	s.JSON(w, res, http.StatusOK)
}

// 实时搜索
func (s *searchHandler) searchKeyword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req searchQuery
	err := s.validate.Validate(r, &req)
	if err != nil {
		s.ERROR(w, err, http.StatusBadRequest)
		return
	}
	// 模糊匹配搜索 TODO
	res := httpModel.HttpSearchWord{
		Type: req.Type,
	}

	for i := 0; i < 10; i++ {
		res.Words = append(res.Words, req.Word+strconv.Itoa(i))
	}
	s.JSON(w, res, http.StatusOK)

}

//
func (s *searchHandler) searchOnlineApply(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req searchTargetReq
	err := s.validate.Validate(r, &req)
	if err != nil {
		s.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := s.db.SearchOnlineItem(req.Word)

	s.JSON(w, res, http.StatusOK)

}

func (s *searchHandler) searchCompany(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req searchTargetReq
	err := s.validate.Validate(r, &req)
	if err != nil {
		s.ERROR(w, err, http.StatusBadRequest)
		return
	}
	//
	res := s.db.SearchCompany(req.Word)

	s.JSON(w, res, http.StatusOK)

}

func (s *searchHandler) searchCarrerTalk(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req searchTargetReq
	err := s.validate.Validate(r, &req)
	if err != nil {
		s.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := s.db.SearchCarrerTalk(req.Word)
	s.JSON(w, res, http.StatusOK)

}

func (s *searchHandler) searchGraduateJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req searchTargetReq
	err := s.validate.Validate(r, &req)
	if err != nil {
		s.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := s.db.SearchGraduateJob(req.Word)

	s.JSON(w, res, http.StatusOK)

}

func (s *searchHandler) searchInternJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req searchTargetReq
	err := s.validate.Validate(r, &req)
	if err != nil {
		s.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res := s.db.SearchInternJobs(req.Word)

	s.JSON(w, res, http.StatusOK)
}
