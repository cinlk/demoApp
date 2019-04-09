package handlers

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"demoApp/server/utils/cache"
	"demoApp/server/utils/errorStatus"
	"demoApp/server/utils/sms"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"time"
)

const phoneRex = `^1([38][0-9]|14[57]|5[^4])\d{8}$`

type loginPwd struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type loginCode struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type loginByRelatedAccount struct {
	RelatedID string `json:"related_id"`
	Type      string `json:"type"`
}

type logout struct {
	Token string `json:"token"`
}

type bindRelatedAccount struct {
	RelatedID string `json:"related_id"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
	Type      string `json:"type"`
}

type registryAccount struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type resetPassword struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type tokenRes = struct {
	Token       string `json:"token"`
	LeanCloudId string `json:"lean_cloud_id"`
}
type codeRes = struct {
	Code string `json:"code"`
}

type accountHandle struct {
	baseHandler
	validate  handlerValider
	UrlPrefix string
	db        *dbOperater.AccountDbOperator
}

type recruiterInfoReq struct {
	Name      string `json:"name"`
	Company   string `json:"company"`
	CompanyId string `json:"company_id"`
	UserIcon  string `json:"user_icon"`
	Title     string `json:"title"`
}

type phoneFormater string

func (p phoneFormater) isValide() bool {

	// 正则表达式判断phone
	reg := regexp.MustCompile(phoneRex)
	return reg.MatchString(string(p))

}

// 接口ip 和 phone 频率限制 TODO
// redis 不可用 查数据库 ？
// 数据库数据和redis 同步 TODO
func (a *accountHandle) SecurityCode(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	phone := param.ByName("phone")
	if !phoneFormater(phone).isValide() {
		a.ERROR(w, errors.New("invalide phone"), http.StatusBadRequest)
		return
	}
	time.Sleep(time.Second * 2)
	// 第三方短信服务 发送发送验证码
	a.sendVeifyCode(w, phone)
}

// 手机号已经存在 在发送验证码
func (a *accountHandle) ExistAccountCode(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	phone := param.ByName("phone")
	if !phoneFormater(phone).isValide() {
		a.ERROR(w, errors.New("not phone type"), http.StatusBadRequest)
		return
	}
	// 判断账号 是否存在
	err := a.db.ExistAccount(phone)
	if err != nil {
		a.ERROR(w, err, http.StatusNotFound)
		return
	}

	a.sendVeifyCode(w, phone)

}

// 匿名访问
func (a *accountHandle) LoginWithAnonymous(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	token, err := a.db.LoginAnonymous()
	if err != nil {
		a.ERROR(w, err, http.StatusInternalServerError)
		return
	}

	a.JSON(w, tokenRes{
		Token: token,
	}, http.StatusOK)
}

func (a *accountHandle) LogoutUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// 查找该用户
	role := r.Header.Get("USER_ROLE")
	uuid := r.Header.Get("USER_UUID")

	err := a.clearCache(uuid)
	if err != nil {
		a.ERROR(w, err, http.StatusInternalServerError)
		return
	}
	a.db.SetLogOut(role, uuid)

	a.JSON(w, logout{
		Token: "",
	}, http.StatusAccepted)
}

func (a *accountHandle) LoginWithPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req loginPwd
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
	}

	token, lid, err := a.db.LoginByPwd(req.Phone, req.Password)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	a.JSON(w, tokenRes{
		Token:       token,
		LeanCloudId: lid,
	}, http.StatusOK)

}

func (a *accountHandle) LoginWithCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginCode
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
	}

	code, err := a.verifyCode(w, req.Phone, req.Code)
	if err != nil {
		a.ERROR(w, err, code)
		return
	}

	token, lid, err := a.db.LoginByCode(req.Phone, req.Code)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	err = a.clearCache(req.Phone)
	if err != nil {
		a.ERROR(w, errors.Wrap(err, "clear cache key failed"), http.StatusInternalServerError)
		return
	}

	a.JSON(w, tokenRes{
		Token:       token,
		LeanCloudId: lid,
	}, http.StatusAccepted)

}

// 注册新账号
func (a *accountHandle) RegistryAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req registryAccount
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
		return
	}

	code, err := a.verifyCode(w, req.Account, req.Code)
	if err != nil {
		a.ERROR(w, err, code)
		return
	}

	token, err := a.db.RegistryAccount(req.Account, req.Password)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	err = a.clearCache(req.Account)
	if err != nil {
		a.ERROR(w, errors.Wrap(err, "clear cache key failed"), http.StatusInternalServerError)
		return
	}
	a.JSON(w, tokenRes{
		Token: token,
	}, http.StatusAccepted)

}

//
func (a *accountHandle) ResetAccountPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req resetPassword
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
		return
	}

	code, err := a.verifyCode(w, req.Account, req.Code)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	token, err := a.db.ResetPassword(req.Account, req.Password)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	err = a.clearCache(req.Account)
	if err != nil {
		a.ERROR(w, errors.Wrap(err, "clear cache key failed"), http.StatusInternalServerError)
		return
	}

	a.JSON(w, tokenRes{
		Token: token,
	}, http.StatusAccepted)

}

// 第三方账号关联 处理
func (a *accountHandle) LoginWithRelatedAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginByRelatedAccount
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
		return
	}

	token, err := a.db.LoginByRelatedAccount(req.RelatedID, req.Type)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	a.JSON(w, tokenRes{
		Token: token,
	}, http.StatusAccepted)

}

func (a *accountHandle) BindRelatedAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req bindRelatedAccount
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
		return
	}

	code, err := a.verifyCode(w, req.Phone, req.Code)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	token, err := a.db.BindRelatedAccount(req.Phone, req.RelatedID, req.Type)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	err = a.clearCache(req.Phone)
	if err != nil {
		a.ERROR(w, errors.Wrap(err, "clear cache key failed"), http.StatusInternalServerError)
		return
	}

	a.JSON(w, tokenRes{
		Token: token,
	}, http.StatusAccepted)
}

// userinfo

// 用户信息
func (a *accountHandle) userInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var userId = r.Header.Get(utils.USER_ID)
	var role = r.Header.Get(utils.USER_ROLE)

	res := a.db.GetUserInfo(role, userId)
	switch t := res.(type) {
	case error:
		a.ERROR(w, t, http.StatusUnprocessableEntity)
	case dbOperater.SeekerUser:
		a.JSON(w, t, http.StatusOK)
	case dbModel.Recruiter:
		a.JSON(w, t, http.StatusOK)
	}
	//a.JSON(w, res, http.stat)

}

// recruiter

func (a *accountHandle) RecuiterInfo(w http.ResponseWriter, r *http.Request, para httprouter.Params) {

	recruiterId := para.ByName("recruiterId")
	if recruiterId == "" {
		a.ERROR(w, errors.New("empty recruiter"), http.StatusBadRequest)
		return
	}

	var req recruiterInfoReq
	err := a.validate.Validate(r, &req)
	if err != nil {
		a.ERROR(w, err, http.StatusBadRequest)
		return
	}

	orm := a.db.Dborm()

	// 更新用户信息
	err = orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", recruiterId).
		Update(*utils.Struct2Map(req)).Error
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

//  function tools
func (a *accountHandle) verifyCode(w http.ResponseWriter, key, code string) (int, error) {

	c, err := cache.GetKeyFromCache(key)
	//
	if e, ok := err.(errorStatus.AppErrorRt); ok {
		return e.State(), e
	}

	if string(c) != code {
		return http.StatusForbidden, errors.New("code not match")
	}
	return 0, nil
}

func (a *accountHandle) clearCache(key string) error {

	return cache.ClearTokenBy(key)
}

func (a *accountHandle) sendVeifyCode(w http.ResponseWriter, phone string) {

	sender := sms.GetSmsService("test")
	code, err := sender.Send(phone)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}
	// 验证码存入 缓存中 5分钟有效期，
	err = cache.SetKeyInCache(phone, []byte(code), 5*60)
	if err != nil {
		a.ERROR(w, err, http.StatusUnprocessableEntity)
		return
	}

	a.JSON(w, codeRes{
		Code: code,
	}, http.StatusCreated)
}
