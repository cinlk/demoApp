package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/model/httpModel"
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type forumSectionItemReq struct {
	Type   string `json:"type" binding:"required"`
	Offset int
	Limit  int
}

type forumSubReplyReq struct {
	PostId string `json:"post_id" binding:"required"`
	Offset int
	Limit  int
}

type forumArticleLike struct {
	PostId string `json:"post_id"`
	Flag   bool   `json:"flag"`
}

type postArticleReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type forumPostReplyReq struct {
	PostId       string `json:"post_id" binding:"required"`
	ReplyContent string `json:"reply_content" binding:"required"`
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

	// 请求用户的id
	userId := r.Header.Get(utils.USER_ID)

	res, err := f.dbOperator.Articles(req.Type, req.Offset, req.Limit, userId)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, res, http.StatusOK)
}

// 发布帖子
func (f *forumHandler) NewArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var req postArticleReq
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	uuid, err := f.dbOperator.NewArticle(req.Title, req.Content, req.Type, userId)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		uuid,
	}, http.StatusCreated)

}

// 删除帖子
func (f *forumHandler) RemovePost(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var postId = para.ByName("postId")
	var userId = r.Header.Get(utils.USER_ID)

	err := f.dbOperator.DeletePostBy(postId, userId)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		postId,
	}, http.StatusAccepted)
}

// 帖子的一级回复数据
func (f *forumHandler) PostReply(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

	var req forumSubReplyReq
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res, err := f.dbOperator.PostContentInfo(req.PostId, req.Offset, req.Limit)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, res, http.StatusOK)
}

// 查看帖子数, 不区分用户
func (f *forumHandler) ReadPostCount(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var postId = para.ByName("postId")

	_ = f.dbOperator.PostReadCount(postId)

	w.WriteHeader(http.StatusAccepted)
}

// 帖子点赞
func (f *forumHandler) LikePost(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var req forumArticleLike
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	//var postId = para.ByName("postId")
	var userId = r.Header.Get(utils.USER_ID)
	err = f.dbOperator.UserLikePost(userId, req.PostId, req.Flag)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	//w.WriteHeader(http.StatusAccepted)
	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			"success",
		},
		"",
	}, http.StatusAccepted)

}

// 帖子收藏
func (f *forumHandler) CollectedPost(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var req forumArticleLike
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	//var postId = para.ByName("postId")
	var userId = r.Header.Get(utils.USER_ID)
	err = f.dbOperator.UserCollectedPost(userId, req.PostId, req.Flag)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	//w.WriteHeader(http.StatusAccepted)
	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			"success",
		},
		"",
	}, http.StatusAccepted)
}

// 回复帖子（一级回复）
func (f *forumHandler) UserReplyPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req forumPostReplyReq
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var userId = r.Header.Get(utils.USER_ID)

	rid, err := f.dbOperator.RecordUserReplyPost(userId, req.PostId, req.ReplyContent)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		rid,
	}, http.StatusCreated)
}
