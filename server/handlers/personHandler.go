package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/model/httpModel"
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type briefInfoReq struct {
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	College string `json:"college"`
}

type deliveryHistory struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type personHandler struct {
	baseHandler
	UrlPrefix string
	validate  handlerValider
	db        *dbOperater.PersonDbOperator
}

func (p *personHandler) updateAvatar(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

	var userId = r.Header.Get(utils.USER_ID)
	// image 数据和名称

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, errors.Wrap(err, "can't get image file").Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	_, err = ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, errors.Wrap(err, "can't read file").Error(), http.StatusBadRequest)
		return
	}

	// 上传到七牛云 获取url TODO
	// 图片名称存入数据库 TODO

	p.db.Avatar(userId, header.Filename)

	// 返回数据 test
	var newIcon = "http://pic-public.yihu.bingfengtech.com/demo.jpeg"
	time.Sleep(time.Second * 3)
	p.JSON(w, httpModel.HttpPersonAvatarModel{
		IconUrl: newIcon,
	}, http.StatusAccepted)

}

// 跟新用户简要信息

func (p *personHandler) BriefInfos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req briefInfoReq
	var userId = r.Header.Get(utils.USER_ID)
	err := p.validate.Validate(r, &req)
	if err != nil {
		p.ERROR(w, err, http.StatusBadRequest)
		return
	}

	err = p.db.BriefInfos(userId, req.Name, req.Gender, req.College)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	p.JSON(w, "", http.StatusAccepted)

}

// 查询投递记录(网申的职位，校招和实习职位)
func (p *personHandler) DeliveryJobList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//var req deliveryHistory
	//err := p.validate.Validate(r, &req)
	//if err != nil {
	//	p.ERROR(w, err, http.StatusBadRequest)
	//	return
	//}
	//time.Sleep(time.Second * 3)
	var userId = r.Header.Get(utils.USER_ID)
	res := p.db.FindDeliveryInfos(userId)
	p.JSON(w, res, http.StatusOK)

}

// 职位投递的历史状态记录
func (p *personHandler) JobDeliveryHistoryStatus(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var userId = r.Header.Get(utils.USER_ID)
	var t = param.ByName("type")
	var jobId = param.ByName("jobId")
	res, err := p.db.JobDeliveryHistory(userId, jobId, t)
	if err != nil {
		p.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	p.JSON(w, res, http.StatusOK)
}

// 根据position id 查询online apply id

func (p *personHandler) FindOnlineApplyId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var positionId = param.ByName("positionId")
	oid := p.db.FindOnlineApplyIdBy(positionId)
	if oid == "" {
		p.ERROR(w, errors.New("not found online apply id"), http.StatusNotFound)
		return
	}

	p.JSON(w, map[string]interface{}{
		"online_apply_id": oid,
	}, http.StatusOK)
}
