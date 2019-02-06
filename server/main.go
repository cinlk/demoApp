package main

import (
	"demoApp/server/handlers"
	"demoApp/server/model"
	"flag"
	"github.com/gorilla/context"
	handlers2 "github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"goframework/cache"
	"goframework/config"
	"goframework/gLog"
	"net/http"
	"strings"
)

const (
	configFlag = "c"
)

var _ = flag.String(configFlag, "/etc/config.ini", "server configuration file")

func main() {

	// 初始化配置文件
	config.LoadConfig(configFlag)

	// 日志
	gLog.InitialLog(config.LogLevel, gLog.NewRotateFileHandler(config.LOGHandler.Key("logfile").Value(),
		config.LOGHandler.Key("unit").MustString("M"), config.LOGHandler.Key("maxByte").MustInt64(8),
		config.LOGHandler.Key("backup").MustInt(10)))

	// 初始化数据库
	model.CreateTables()
	defer model.CloseDB()

	// 初始化缓存
	cache.InitialCache()
	defer cache.CacheProxy.Close()

	// 加载rbac 权限控制
	config.LoadAccessControlPolicy(config.AccessControllCfg, config.AccessControllCsv)

	// router注册
	router := httprouter.New()
	router.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {
		gLog.LOG_ERRORF("url:%s  panic: %s \n", request.URL, i)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	handlers.RegisterRouter(router)
	// 中间件
	middlewares := handlers2.LoggingHandler(gLog.GetLogUtil(), router)

	// 热更新启动服务 TODO
	gLog.LOG_FATAL(http.ListenAndServe(
		strings.Join([]string{config.ServerIP, config.ServerPort}, ":"),
		context.ClearHandler(middlewares)))

}
