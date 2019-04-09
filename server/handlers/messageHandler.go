package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)

type newConversationReq struct {
	ConversationId string `json:"conversation_id" binding:"required"`
	//UserId         string `json:"sender_id" binding:"required"`
	RecruiterId string `json:"recruiter_id" binding:"required"`
	JobId       string `json:"job_id" binding:"required"`
}

type messageHandler struct {
	baseHandler
	urlPrefix  string
	validate   handlerValider
	dbOperator *dbOperater.MessageDbOperater
}

// who is the recuiter  has pay attension to me
func (m *messageHandler) attentionsToMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (m *messageHandler) checkAttentions(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

}

//

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
