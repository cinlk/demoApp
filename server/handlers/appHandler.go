package handlers

import (
	"demoApp/server/model/dbOperater"
	"demoApp/server/utils"
	"demoApp/server/utils/qiniuStorage"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
)

type imageRes struct {
	ImageURL string `json:"image_url"`
}

type appHandler struct {
	baseHandler
	UrlPrefix  string
	dbOperator *dbOperater.AppDBoperator
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
