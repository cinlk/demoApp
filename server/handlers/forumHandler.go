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

// 帖子的回复
type forumSubReplyReq struct {
	PostId string `json:"post_id" binding:"required"`
	Offset int
	Limit  int
}

// 回复的回复
type forumSecondReplyReq struct {
	ReplyId string `json:"reply_id" binding:"required"`
	Offset  int
	Limit   int
}

type forumArticleLike struct {
	PostId string `json:"post_id" binding:"required"`
	Flag   bool   `json:"flag"`
}
type forumReplyLike struct {
	ReplyId string `json:"reply_id" binding:"required"`
	Flag    bool   `json:"flag"`
}
type forumSubReplyLike struct {
	SubReplyId string `json:"sub_reply_id" binding:"required"`
	Flag       bool   `json:"flag"`
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

// 二级回复请求
type forumPostSubReplyReq struct {
	ReplyId      string `json:"reply_id" binding:"required"`
	TalkedUserId string `json:"talked_user_id" binding:"required"`
	Content      string `json:"content" binding:"required"`
}

// 举报帖子
type forumAlertPostReq struct {
	PostId  string `json:"post_id"`
	Content string `json:"content"`
}

// 举报回复
type forumAlertReplyReq struct {
	ReplyId string `json:"reply_id"`
	Content string `json:"content"`
}

// 举报二级回复
type forumAlertSubReplyReq struct {
	SubReplyId string `json:"sub_reply_id"`
	Content    string `json:"content"`
}

// 搜索帖子请求
type forumSearchReq struct {
	Word   string `json:"word" binding:"required"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type forumHandler struct {
	baseHandler
	urlPrefix  string
	validate   handlerValider
	dbOperator *dbOperater.ForumDboperator
}


// 帖子分组
type  forumGroupNameReq struct {

	PostId string `json:"post_id" binding:"required"`
	GroupName []string `json:"group_name"`

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

func (f *forumHandler) oneArticle(w http.ResponseWriter, r *http.Request, param httprouter.Params)  {
	var userId = r.Header.Get(utils.USER_ID)
	var postId = param.ByName("postId")

	res, err := f.dbOperator.FindArticleBy(userId, postId)
	if err != nil{
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

// 删除一级回复
func (f *forumHandler) RemoveReply(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var replyId = para.ByName("replyId")
	var userId = r.Header.Get(utils.USER_ID)

	err := f.dbOperator.DeleteReply(replyId, userId)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		replyId,
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
	var userId = r.Header.Get(utils.USER_ID)
	res, err := f.dbOperator.PostContentInfo(req.PostId, userId, req.Offset, req.Limit)
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

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		"",
	}, http.StatusAccepted)
}

func (f *forumHandler) AlertPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req forumAlertPostReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}

	_ = f.dbOperator.AlertPost(req.PostId, userId, req.Content)

	w.WriteHeader(http.StatusAccepted)

}

func (f *forumHandler) AlertReply(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req forumAlertReplyReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	_ = f.dbOperator.AlertReply(req.ReplyId, userId, req.Content)

	w.WriteHeader(http.StatusAccepted)

}

func (f *forumHandler) AlertSubReply(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req forumAlertSubReplyReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	_ = f.dbOperator.AlertSubReply(req.SubReplyId, userId, req.Content)
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

// 帖子加入分组
func (f *forumHandler) CollectedPostInGroup(w http.ResponseWriter, r *http.Request, para httprouter.Params){

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

// 获取子回复的 二级回复内容
func (f *forumHandler) UserSubReplys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req forumSecondReplyReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res, err := f.dbOperator.SecondReplys(req.ReplyId, userId, req.Offset, req.Limit)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, res, http.StatusOK)

}

// 发布二级回复
func (f *forumHandler) NewSubReply(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req forumPostSubReplyReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}

	id, err := f.dbOperator.RecordUserSubReply(userId, req.TalkedUserId, req.ReplyId, req.Content)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		id,
	}, http.StatusCreated)

}

// 删除自己的二级回复
func (f *forumHandler) RemoveMySubReply(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var subReplyId = param.ByName("subReplyId")
	var userId = r.Header.Get(utils.USER_ID)

	err := f.dbOperator.DeleteSubReply(subReplyId, userId)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			Result: "success",
		},
		subReplyId,
	}, http.StatusAccepted)

}

// 子回复点赞
func (f *forumHandler) UserLikeReply(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var req forumReplyLike
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var userId = r.Header.Get(utils.USER_ID)

	err = f.dbOperator.LikeReply(userId, req.ReplyId, req.Flag)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			"success",
		},
		"",
	}, http.StatusAccepted)

}

// 二级回复点赞
func (f *forumHandler) UserLikeSubReply(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var req forumSubReplyLike
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	var userId = r.Header.Get(utils.USER_ID)
	err = f.dbOperator.LikeSubReply(userId, req.SubReplyId, req.Flag)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpForumResponse{
		httpModel.HttpResultModel{
			"success",
		},
		"",
	}, http.StatusAccepted)

}

// 搜索帖子 test

func (f *forumHandler) SearchForumPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req forumSearchReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}
	res, err := f.dbOperator.SearchPostBy(req.Word, userId, req.Offset, req.Limit)
	if err != nil {
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, res, http.StatusOK)

}


func (f *forumHandler) postNewGroupName(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var req forumGroupNameReq
	var userId = r.Header.Get(utils.USER_ID)
	err := f.validate.Validate(r, &req)
	if err != nil {
		f.ERROR(w, err, http.StatusBadRequest)
		return
	}

	err = f.dbOperator.RelatePostGroupName(userId, req.PostId, req.GroupName)
	if err != nil{
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, httpModel.HttpResultModel{
		Result:"success",
	}, http.StatusCreated)

}

func (f *forumHandler) postGroupName(w http.ResponseWriter, r *http.Request, param httprouter.Params){
	var userId = r.Header.Get(utils.USER_ID)
	var postId = param.ByName("postId")

	res, err := f.dbOperator.PostGroupNames(userId, postId)
	if err != nil{
		f.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	f.JSON(w, res, http.StatusOK)

}
