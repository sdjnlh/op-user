package main

import (
	"flag"

	"code.letsit.cn/go/common/app"
	"code.letsit.cn/go/common/log"
	"code.letsit.cn/go/common/util"
	"code.letsit.cn/go/common/web"
	"code.letsit.cn/go/op-user/opu"
	userRest "code.letsit.cn/go/op-user/rest"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go/extra"
)

func main() {
	opu.Api.Mount(opu.Service)
	app.RegisterStarter(opu.Api)
	if err := app.Start(); err != nil {
		panic(err)
	}

	flag.Parse()
	log.Slog.Info("starting rest")
	router := gin.Default()
	router.Use(web.CorsHandler(opu.Api))
	//router.Use(session.InitSessionStore(db.Redis))
	userRest.RegisterAPIs(router)
	extra.SetNamingStrategy(util.LowerFirst)

	log.Slog.Fatal(router.Run(":" + opu.Api.Port))
	log.Slog.Info("rest started")
}
