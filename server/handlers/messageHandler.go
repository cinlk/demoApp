package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type newConversationReq struct {
	ConversationId string `json:"conversation_id" binding:"required"`
	//UserId         string `json:"sender_id" binding:"required"`
	RecruiterId string `json:"recruiter_id" binding:"required"`
	JobId       string `json:"job_id" binding:"required"`
}

type visitorHistoryReq struct {
	UserId string `json:"user_id" binding:"required"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type visitorStatusReq struct {
	UserId      string `json:"user_id" binding:"required"`
	RecruiterId string `json:"recruiter_id" binding:"required"`
}

type messageHandler struct {
	baseHandler
	urlPrefix  string
	validate   handlerValider
	dbOperator *dbOperater.MessageDbOperater
}

// 获取看过我的hr 记录
func (m *messageHandler) myVisitor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req visitorHistoryReq
	err := m.validate.Validate(r, &req)
	if err != nil {
		m.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res, err := m.dbOperator.MyVisitors(req.UserId, req.Offset, req.Limit)
	if err != nil {
		m.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	m.JSON(w, res, http.StatusOK)

}

// 更新已经访问者的 查看状态
func (m *messageHandler) visitorChecked(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var req visitorStatusReq
	err := m.validate.Validate(r, &req)
	if err != nil {
		m.ERROR(w, err, http.StatusBadRequest)
		return
	}
	m.dbOperator.CheckVisitor(req.UserId, req.RecruiterId)
	w.WriteHeader(http.StatusAccepted)
}

//
// 访问者
func (m *messageHandler) checkVisitorTime(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var userId = para.ByName("userId")
	if strings.TrimSpace(userId) == "" {
		m.ERROR(w, errors.New("empty userid"), http.StatusBadRequest)
		return
	}
	m.dbOperator.UpdateTimeVisitor(userId)
	w.WriteHeader(http.StatusAccepted)
}

// 访问者
func (m *messageHandler) CheckNewVisitor(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	var userId = para.ByName("userId")
	if strings.TrimSpace(userId) == "" {
		m.ERROR(w, errors.New("empty userid"), http.StatusBadRequest)
		return
	}
	// 检查是否有最新的访问者
	exist := m.dbOperator.ExistNewVisitor(userId)
	m.JSON(w, map[string]interface{}{"exist": exist}, http.StatusOK)
}

func (m *messageHandler) recruiterInfo(w http.ResponseWriter, r *http.Request, para httprouter.Params) {
	userId := para.ByName("userId")
	if userId == "" {
		http.Error(w, errors.New("empty user").Error(), http.StatusBadRequest)
		return
	}

	res, err := m.dbOperator.RecruiterInfo(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	m.JSON(w, res, http.StatusOK)
}

// single conversation with recruiter
func (m *messageHandler) conversation(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

	// 创建conversation
	if r.Method == http.MethodPost {
		var req newConversationReq
		err := m.validate.Validate(r, &req)
		if err != nil {
			m.ERROR(w, err, http.StatusBadRequest)
			return
		}
		var userId = r.Header.Get(utils.USER_ID)

		// 不能和自己交谈
		if userId == req.RecruiterId {
			m.ERROR(w, errors.Wrap(err, "con't talk to mysel"), http.StatusConflict)
			return
		}
		err = m.dbOperator.CreateConversation(userId, req.RecruiterId, req.JobId, req.ConversationId)
		if err != nil {
			m.ERROR(w, err, http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

	// 获取conversationid
	if r.Method == http.MethodGet {

		userId := para.ByName("userId")
		JobId := para.ByName("jobId")

		id := m.dbOperator.ConversationBy(userId, JobId)
		if id == "" {
			m.ERROR(w, errors.New("not found conversation"), http.StatusNotFound)
			return
		}

		m.JSON(w, map[string]string{"conversation_id": id}, http.StatusOK)
	}

}
