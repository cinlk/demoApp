package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"demoApp/server/utils/qiniuStorage"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)

type newQueryPage struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}

type imageRes struct {
	ImageURL string `json:"image_url"`
}

type appHandler struct {
	baseHandler
	UrlPrefix  string
	dbOperator *dbOperater.AppDBoperator
	validate   handlerValider
}

func (app *appHandler) AppGuidanceItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	res := app.dbOperator.AppGuidanceItems()
	app.JSON(w, res, http.StatusAccepted)
}

func (app *appHandler) AppAdvitiseImageURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// 七牛云获取 url
	imageURL := qiNiuStorage.PubluicBucketFile(utils.ADVIST_IMAGE_NAME)
	if imageURL == "" {
		app.ERROR(w, errors.New("empty image"), http.StatusNotFound)
		return
	}

	app.JSON(w, imageRes{
		ImageURL: imageURL,
	}, http.StatusOK)

}

// 新闻专栏
func (app *appHandler) News(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req newQueryPage
	err := app.validate.Validate(r, &req)
	if err != nil {
		app.ERROR(w, err, http.StatusBadRequest)
		return
	}

	res := app.dbOperator.News(req.Type, req.Offset)
	for i := 0; i < len(res); i++ {
		res[i].CreatedTime = res[i].Time.Unix()
	}

	app.JSON(w, res, http.StatusAccepted)
}
