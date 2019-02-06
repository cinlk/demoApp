package handlers

import (
	"demoApp/server/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type TestHandler struct {
	UrlPrefix string
}

func (t *TestHandler) FindSomething(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	utils.JsonResponse(w, map[string]string{
		"dqwd": "dwqdwq",
	}, http.StatusOK)
}
