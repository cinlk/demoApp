package handlers

import (
	"demoApp/server/model/dbOperater"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type forumSectionItemReq struct {
	Type   string `json:"type" binding:"required"`
	Offset int
	Limit  int
}

type forumHandler struct {
	baseHandler
	urlPrefix  string
	validate   handlerValider
	dbOperator *dbOperater.ForumDboperator
}

// 获取帖子数据
func (f *forumHandler) SectionArticles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req forumSectionItemReq
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res, err := f.dbOperator.Articles(req.Type, req.Offset, req.Limit)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, res, http.StatusOK)
}

// 发布帖子
func (f *forumHandler) NewArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
