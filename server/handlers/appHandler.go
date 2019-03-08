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

// 选择的城市
func (app *appHandler) Citys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var citys = make(map[string][]string)
	citys = map[string][]string{
		"GHIJ":   []string{"桂林", "广州", "合肥", "呼和浩特", "海口", "杭州", "湖州", "济南", "济宁", "嘉兴", "江阴"},
		"热门城市":   []string{"全国", "西安", "深圳", "武汉", "长沙", "苏州", "南京", "上海", "成都", "广州", "杭州", "北京"},
		"ABCDEF": []string{"鞍山", "安阳", "保定", "北京", "成都", "重庆", "长沙", "常州", "大连", "东营", "德州", "佛山", "福州"},
		"KLMN":   []string{"特别长特别长特比长", "昆明", "昆山", "聊城", "廊坊", "洛阳", "连云港", "兰州", "绵阳", "宁波", "南京", "南宁", "南通"},
	}

	var c = struct {
		Citys map[string][]string `json:"citys"`
	}{
		Citys: citys,
	}
	app.JSON(w, c, http.StatusOK)
}

// 选择的行业领域 粗粒度
func (app *appHandler) BusinessField(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var fields = struct {
		Fields []string `json:"fields"`
	}{
		Fields: []string{"不限", "健康医疗", "生活服务", "旅游", "金融", "信息安全",
			"网络招聘", "互联网", "IT软件", "媒体", "公共会展", "机械制造", "游戏", "教育培训", "其他"},
	}

	app.JSON(w, fields, http.StatusOK)

}

// 细分的行业领域职位
func (app *appHandler) BusinessFieldJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var kind = struct {
		Fields map[string][]string `json:"fields"`
	}{
		Fields: map[string][]string{
			"IT互联网":  []string{"软件及系统开发", "算法/大数据", "智能硬件", "移动开发"},
			"电子电气":   []string{"电子/通信", "嵌入式", "电气工程"},
			"人事行政":   []string{"人事HR", "猎头", "行政"},
			"传媒设计":   []string{"广告", "编辑", "媒体", "视频后期"},
			"杀掉无多哇多": []string{"数据1", "数据2"},
		},
	}

	app.JSON(w, kind, http.StatusOK)
}

func (app *appHandler) CompanyType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var kind = struct {
		Type []string `json:"type"`
	}{
		Type: []string{"不限", "外资企业", "私营企业", "国有企业", "非营利企业", "其他"},
	}

	app.JSON(w, kind, http.StatusOK)
}

// 实习条件
func (app *appHandler) InternCondition(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var condition = struct {
		Condition map[string][]string `json:"condition"`
	}{
		Condition: map[string][]string{
			"每周实习天数": []string{"不限", "1天", "2天", "3天", "4天", "5天"},
			"实习日薪":   []string{"不限", "100以下", "100-200", "200以上"},
			"实习月数":   []string{"不限", "1月", "1-6月", "6月以上"},
			"学历":     []string{"不限", "大专", "本科", "硕士", "博士"},
			"是否转正":   []string{"不限", "提供转正"},
		},
	}

	app.JSON(w, condition, http.StatusOK)
}

// 城市和大学
func (app *appHandler) CitysCollege(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var citysCollege = struct {
		CityCollege map[string][]string `json:"city_college"`
	}{
		CityCollege: map[string][]string{
			"北京": []string{"清华大学", "北京大学", "北京邮电大学", "北京理工大学"},
			"上海": []string{"上海大学", "上海理工", "复旦大学", "交通大学", "徐家汇教会学院", "上海戏剧学院"},
			"成都": []string{"电子科技大学", "西南大学", "成都理工", "四川大学", "四川师范大学", "成都学院"},
			"沈阳": []string{"东北大学", "沈阳理工"},
			"全国": []string{},
		},
	}

	app.JSON(w, citysCollege, http.StatusOK)
}

// 职位举报信息列表

func (app *appHandler) jobWarns(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var warns = struct {
		Warns []string `json:"warns"`
	}{
		Warns: []string{"描述不符合", "虚假信息", "带我去的无群", "带我飞外国人", "核桃仁和投入和", "反而更让他很突然很突然", "其他"},
	}

	app.JSON(w, warns, http.StatusOK)
}
