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
	"reflect"
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

// TODO 验证 必须的参数
type jsonValidate struct {
	requiredTag string
}

func (j *jsonValidate) Validate(r *http.Request, obj interface{}) error {

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

	// bind check 的tag 检查
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		if tagv, ok := t.Field(i).Tag.Lookup(j.requiredTag); ok {
			if tagv != "required" {
				continue
			}
			switch v.Field(i).Kind() {
			case reflect.String:
				if v.Field(i).String() == "" {
					return &errorStatus.AppError{
						Code: http.StatusBadRequest,
						Err: errors.New(fmt.Sprintf("field with name %s is invalidate",
							t.Field(i).Name)),
					}
				}
			case reflect.Uint:
				if v.Field(i).Uint() == 0 {
					return &errorStatus.AppError{
						Code: http.StatusBadRequest,
						Err: errors.New(fmt.Sprintf("field with name %s is invalidate",
							t.Field(i).Name)),
					}
				}
			}

		}
	}

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

	// search
	var searchHandler = searchHandler{
		UrlPrefix: "/search",
		db:        dbOperater.NewSearchDboperator(),
		validate: &jsonValidate{
			requiredTag: "binding",
		},
	}
	//
	var recruitHandler = recruiteListHandle{
		UrlPrefix:  "/recruite",
		dbOperator: dbOperater.NewRecruiteDboperator(),
		validate: &jsonValidate{
			requiredTag: "binding",
		},
	}

	// message
	var messageHandler = messageHandler{
		urlPrefix:  "/message",
		dbOperator: dbOperater.NewMessageDbOperater(),
		validate: &jsonValidate{
			requiredTag: "binding",
		},
	}

	// forum
	var forumHandler = forumHandler{
		urlPrefix:  "/forum",
		dbOperator: dbOperater.NewForumDboperator(),
		validate: &jsonValidate{
			requiredTag: "binding",
		},
	}

	// person
	var personHandler = personHandler{
		UrlPrefix: "/person",
		db:        dbOperater.NewPersonDbOperator(),
		validate: &jsonValidate{
			requiredTag: "binding",
		},
	}
	var testh = TestHandler{
		UrlPrefix: "/test",
	}

	global := rg.NewGroupRouter(apphandler.UrlPrefix, router)
	{
		global.GET("/guidance", apphandler.AppGuidanceItems)
		global.GET("/advise/image", apphandler.AppAdvitiseImageURL)
		global.POST("/news", apphandler.News)
		global.GET("/citys", apphandler.Citys)
		global.GET("/business/field", apphandler.BusinessField)
		global.GET("/subBusiness/field", apphandler.BusinessFieldJob)
		global.GET("/company/type", apphandler.CompanyType)
		global.GET("/intern/condition", apphandler.InternCondition)
		global.GET("/city/college", apphandler.CitysCollege)
		//global.GET("/near/meetings", )
		global.GET("/jobs/warns", apphandler.jobWarns)
		global.GET("/forum/warns", apphandler.forumWarns)

		// 测试通知推送
		global.PUT("/systemNotify", messageHandler.systemNotifyMessage)
		global.PUT("/userNotify", messageHandler.SpecialNotifyMessage)

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
		account.POST("/recruite/info/:recruiterId", accoutHandler.RecuiterInfo)
		account.GET("/userinfo", accoutHandler.userInfo, middleware.AuthorizationVerify)
		//account.GET("/visitor/:userId", accoutHandler.CheckVisitor, middleware.AuthorizationVerify)
		//account.POST("/visitor/:userId", accoutHandler.CheckVisitor, middleware.AuthorizationVerify)

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

		homePage.POST("/near/meetings", lhandler.nearBy)
		homePage.POST("/near/company", lhandler.nearBy)

	}
	job := rg.NewGroupRouter(jobHandler.UrlPrefix, router)
	{

		job.POST("/kind/:kind", jobHandler.FindJobKind)
	}

	search := rg.NewGroupRouter(searchHandler.UrlPrefix, router)
	{
		search.GET("/word/:type", searchHandler.TopWords)
		search.POST("/similar", searchHandler.searchKeyword)
		search.POST("/online", searchHandler.searchOnlineApply)
		search.POST("/company", searchHandler.searchCompany)
		search.POST("/careerTalk", searchHandler.searchCarrerTalk)
		search.POST("/graduate", searchHandler.searchGraduateJobs)
		search.POST("/intern", searchHandler.searchInternJobs)

	}

	recruit := rg.NewGroupRouter(recruitHandler.UrlPrefix, router, middleware.AuthorizationVerify)
	{
		recruit.POST("/online", recruitHandler.onlineApplys)
		recruit.POST("/carreerTalk", recruitHandler.careerTalks)
		recruit.POST("/company", recruitHandler.companys)
		recruit.POST("/graduate", recruitHandler.graduatejobs)
		recruit.POST("/intern", recruitHandler.internjobs)
		recruit.GET("/online/:id", recruitHandler.findOnlineApply)
		recruit.PUT("/online/:onlineId/:positionId", recruitHandler.applyOnlineJob)
		recruit.POST("/online/collect", recruitHandler.collectOnlineApply)
		recruit.GET("/meeting/:id", recruitHandler.findCareerTalk)
		recruit.POST("/meeting/collect", recruitHandler.collecteCareerTalk)
		recruit.GET("/company/:id", recruitHandler.findCompany)
		recruit.POST("/company/collect", recruitHandler.collectCompany)
		recruit.GET("/graduate/:id", recruitHandler.findGraduate)
		recruit.PUT("/job/:type/:jobId", recruitHandler.applyJob)
		recruit.GET("/intern/:id", recruitHandler.findInternJob)
		recruit.POST("/job/collect", recruitHandler.collectJob)
		recruit.GET("/recruiter/:id", recruitHandler.recruiterCompanyAndJobs)
		recruit.POST("/tag/jobs", recruitHandler.companyTagJobs)
		recruit.POST("/company/recruit/meeting", recruitHandler.companyRecruitMeeting)

	}

	message := rg.NewGroupRouter(messageHandler.urlPrefix, router, middleware.AuthorizationVerify)
	{
		message.POST("/conversation", messageHandler.conversation)
		message.GET("/conversation/:userId/:jobId", messageHandler.conversation)
		message.GET("/talkWith/:userId", messageHandler.recruiterInfo)
		message.POST("/visitors", messageHandler.myVisitor)
		message.PUT("/visitor/status", messageHandler.visitorChecked)
		message.POST("/visitorTime/:userId", messageHandler.checkVisitorTime)
		message.GET("/newVisitor/:userId", messageHandler.CheckNewVisitor)
		message.GET("/newSystemMessage/:userId", messageHandler.HasNewSystemMessage)
		message.POST("/systemMessageTime/:userId", messageHandler.ReviewSystemMessage)
		message.GET("/newThumbUp/:userId", messageHandler.HasThumbUpMessage)
		message.POST("/thumbUpTime/:userId", messageHandler.ReviewThumbUpMessage)
		message.GET("/newForumReply/:userId", messageHandler.NewForumReply2Me)
		message.POST("/forumReplyTime/:userId", messageHandler.ReviewForumReply2Me)

	}

	forum := rg.NewGroupRouter(forumHandler.urlPrefix, router)
	{
		forum.POST("/articles", forumHandler.SectionArticles, middleware.FetchUserId)
		forum.POST("/new/article", forumHandler.NewArticle, middleware.AuthorizationVerify)
		forum.PUT("/article/count/:postId", forumHandler.ReadPostCount)
		// 获取所有子回复
		forum.POST("/article/replys", forumHandler.PostReply, middleware.FetchUserId)
		forum.PUT("/article/like", forumHandler.LikePost, middleware.AuthorizationVerify)
		forum.PUT("/reply/like", forumHandler.UserLikeReply, middleware.AuthorizationVerify)
		forum.PUT("/article/collect", forumHandler.CollectedPost, middleware.AuthorizationVerify)
		forum.POST("/article/alert", forumHandler.AlertPost, middleware.AuthorizationVerify)
		// 发布回复
		forum.POST("/article/reply", forumHandler.UserReplyPost, middleware.AuthorizationVerify)
		forum.DELETE("/article/:postId", forumHandler.RemovePost, middleware.AuthorizationVerify)
		forum.DELETE("/reply/:replyId", forumHandler.RemoveReply, middleware.AuthorizationVerify)
		forum.POST("/reply/alert", forumHandler.AlertReply, middleware.AuthorizationVerify)
		// 子回复
		forum.POST("/subreply", forumHandler.UserSubReplys, middleware.FetchUserId)
		forum.POST("/newSubreply", forumHandler.NewSubReply, middleware.AuthorizationVerify)
		forum.POST("/subReply/alert", forumHandler.AlertSubReply, middleware.AuthorizationVerify)
		forum.PUT("/subReply/like", forumHandler.UserLikeSubReply, middleware.AuthorizationVerify)
		forum.DELETE("/subReply/:subReplyId", forumHandler.RemoveMySubReply, middleware.AuthorizationVerify)

		// 搜索帖子
		forum.POST("/search", forumHandler.SearchForumPost, middleware.FetchUserId)

	}

	// 个人主页
	person := rg.NewGroupRouter(personHandler.UrlPrefix, router, middleware.AuthorizationVerify)
	{
		person.POST("/avatar", personHandler.updateAvatar)
		person.POST("/brief", personHandler.BriefInfos)
		person.GET("/delivery", personHandler.DeliveryJobList)
		person.GET("/delivery/history/:type/:jobId", personHandler.JobDeliveryHistoryStatus)
		person.GET("/onlineApplyId/:positionId", personHandler.FindOnlineApplyId)
		person.GET("/resumes", personHandler.MyResumes)
		person.POST("/textResume", personHandler.createTextResume)
		person.POST("/attacheResume", personHandler.createNewAttachResume)
		person.PUT("/primary/resume/:resumeId", personHandler.primaryResume)
		person.POST("/resume/name", personHandler.resumeName)
		person.DELETE("/resume/:type/:resumeId", personHandler.deleteResume)

		// 文本简历
		person.GET("/text/resume/:resumeId", personHandler.TextResumeInfo)
		person.PUT("/textResume/baseInfo/avatar/:resumeId", personHandler.BaseInfoAvatar)
		person.POST("/textResume/baseInfo/content", personHandler.BaseInfoContent)

		person.POST("/textResume/education", personHandler.NewEducationInfo)
		person.PUT("/textResume/education/:id", personHandler.UpdateEducationInfo)
		person.DELETE("/textResume/education/:resumeId/:id", personHandler.DeleteEducationInfo)


		person.POST("/textResume/work", personHandler.newWorkExperience)
		person.PUT("/textResume/work/:id", personHandler.updateWorkExperience)
		person.DELETE("/textResume/work/:resumeId/:id", personHandler.deleteWorkExperience)


		person.POST("/textResume/project", personHandler.newProjectExperience)
		person.PUT("/textResume/project/:id", personHandler.updateProjectExperience)
		person.DELETE("/textResume/project/:resumeId/:id", personHandler.deleteProjectExperience)


		person.POST("/textResume/college", personHandler.newCollegeActive)
		person.PUT("/textResume/college/:id", personHandler.updateCollegeActive)
		person.DELETE("/textResume/college/:resumeId/:id", personHandler.deleteCollegeActive)

		person.POST("/textResume/skill", personHandler.newResumeSkill)
		person.PUT("/textResume/skill/:id", personHandler.updateResumeSkill)
		person.DELETE("/textResume/skill/:resumeId/:id", personHandler.deleteResumeSkill)

		person.POST("/textResume/socialPractice", personHandler.newSocialPractice)
		person.PUT("/textResume/socialPractice/:id", personHandler.updateSocialPractice)
		person.DELETE("/textResume/socialPractice/:resumeId/:id", personHandler.deleteSocialPractice)

		person.POST("/textResume/other", personHandler.newResumeOther)
		person.PUT("/textResume/other/:id", personHandler.updateResumeOther)
		person.DELETE("/textResume/other/:resumeId/:id", personHandler.deleteResumeOther)

		//person.POST("/textResume/estimate", personHandler.newResumeEstimate)
		person.PUT("/textResume/estimate/:id", personHandler.updateResumeEstimate)

		person.GET("/attachResume/:resumeId", personHandler.attachResumeUrl)

		person.POST("/collect/jobs", personHandler.collectedJobs)
		person.POST("/collect/careerTalk", personHandler.collectedCareerTalk)
		person.POST("/collect/onlineApply", personHandler.collectedOnlineApply)
		person.POST("/collect/company", personHandler.collectedCompany)

		person.POST("/unCollect/jobs", personHandler.unSubscribeCollectedJobs)
		person.POST("/unCollect/company", personHandler.unSubScribeCollectedCompany)
		person.POST("/unCollect/onlineApply", personHandler.unSubScribeCollectedOnlineApply)
		person.POST("/unCollect/careerTalk", personHandler.unSubScribeCollectedCareerTalk)


	}

	test := rg.NewGroupRouter(testh.UrlPrefix, router, middleware.AuthorizationVerify)
	{
		test.GET("/some/:demo", testh.FindSomething)
	}

}
