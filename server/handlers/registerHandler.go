package handlers

import (
	"demoApp/server/handlers/middleware"
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"demoApp/server/utils/errorStatus"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"goframework/extension"
	"goframework/utils/appJson"
	"net/http"
)

//  协议处理基础功能
type handlerValider interface {
	// 验证请求body
	Validate(r *http.Request, obj interface{}) error
}

type response struct {
	ReturnCode string      `json:"return_code"`
	ReturnMsg  string      `json:"return_msg"`
	Code       int         `json:"code"`
	Body       interface{} `json:"body,omitempty"`
}

type baseValidate struct{}

func (b *baseValidate) Validate(r *http.Request, obj interface{}) error {

	if r.Body == nil {
		return &errorStatus.AppError{
			Code: http.StatusBadRequest,
			Err:  errors.New("empty body"),
		}
	}
	err := appJson.JsonDecode(r.Body, obj)
	if err != nil {
		return &errorStatus.AppError{
			Code: http.StatusBadRequest,
			Err:  errors.WithMessage(err, fmt.Sprintf("decode data %+v failed", obj)),
		}
	}
	defer r.Body.Close()

	return nil
}

type baseHandler struct {
	Response response
}

func (b baseHandler) JSON(w http.ResponseWriter, data interface{}, code int) {

	b.Response.ReturnMsg = ""
	b.Response.ReturnCode = "SUCCESS"
	b.Response.Code = code
	b.Response.Body = data

	utils.JsonResponse(w, b.Response, code)
}

func (b baseHandler) ERROR(w http.ResponseWriter, err error, code int) {

	b.Response.ReturnCode = "FAILED"

	switch e := err.(type) {
	case errorStatus.AppErrorRt:

		b.Response.ReturnMsg = e.Error()
		b.Response.Code = e.State()
	case error:
		b.Response.Code = code
		b.Response.ReturnMsg = e.Error()
	default:
		http.Error(w, "invalide error type", http.StatusInternalServerError)
	}

	utils.JsonResponse(w, b.Response, b.Response.Code)
}

var rg = extension.RouterGroup{
	BasePath: utils.API_NAMESPACE,
}

func RegisterRouter(router *httprouter.Router) {

	var apphandler = appHandler{

		UrlPrefix:  "/global",
		dbOperator: dbOperater.NewAppDBoperator(),
		validate:   &baseValidate{},
	}

	var accoutHandler = accountHandle{
		validate:  &baseValidate{},
		UrlPrefix: "/account",
		db:        dbOperater.NewAccountDbOperator(),
	}
	// home page
	var lhandler = listHandler{
		validate:  &baseValidate{},
		UrlPrefix: "/home",
		db:        dbOperater.NewListDboperater(),
	}

	// jobs
	var jobHandler = jobHandler{
		validate:   &baseValidate{},
		dbOperator: dbOperater.NewJobDbOperator(),
		UrlPrefix:  "/job",
	}

	var testh = TestHandler{
		UrlPrefix: "/test",
	}

	global := rg.NewGroupRouter(apphandler.UrlPrefix, router)
	{
		global.GET("/guidance", apphandler.AppGuidanceItems)
		global.GET("/advise/image", apphandler.AppAdvitiseImageURL)
		global.POST("/news", apphandler.News)

	}

	account := rg.NewGroupRouter(accoutHandler.UrlPrefix, router)
	{
		account.POST("/login/pwd", accoutHandler.LoginWithPassword)
		account.DELETE("/logout", accoutHandler.LogoutUser, middleware.AuthorizationVerify)
		account.GET("/anonymouse", accoutHandler.LoginWithAnonymous)
		account.PUT("/code/:phone", accoutHandler.SecurityCode)
		account.PUT("/phone/:phone", accoutHandler.ExistAccountCode)
		account.POST("/registry/pwd", accoutHandler.RegistryAccount)
		account.POST("/password", accoutHandler.ResetAccountPassword)
		account.POST("/login/code", accoutHandler.LoginWithCode)
		account.POST("/login/social", accoutHandler.LoginWithRelatedAccount)
		account.POST("/registry/social", accoutHandler.BindRelatedAccount)

	}

	homePage := rg.NewGroupRouter(lhandler.UrlPrefix, router)
	{
		homePage.GET("/banners", lhandler.bannerInfos)
		homePage.GET("/news", lhandler.latestNews)
		homePage.GET("/jobCategory", lhandler.jobCategories)
		homePage.GET("/latest", lhandler.jobTops)
		homePage.GET("/carrerTalks", lhandler.careerTalks)
		homePage.GET("/onlineApply", lhandler.onlineApply)
		homePage.POST("/jobs", lhandler.personalityJobs)
		homePage.GET("/recommand", lhandler.personalRecommand)

	}
	job := rg.NewGroupRouter(jobHandler.UrlPrefix, router)
	{

		job.POST("/kind/:kind", jobHandler.FindJobKind)
	}

	test := rg.NewGroupRouter(testh.UrlPrefix, router, middleware.AuthorizationVerify)
	{
		test.GET("/some/:demo", testh.FindSomething)
	}

}
